import { FileTree } from "../FileTree.ts";
import { HSVtoRGB } from "../@utils/HSVtoRGB.ts";
import { ColoredFileLeaf } from "../color.ts";

export type Message = ['collapsed'] | ['expanded'] | ['tagged'] | ['untagged'] | ['fileTagged', string]

export class FolderState<F> {
  /** the name of the folder */
  name: string = ''
  /** the subfolders that live in this folder */
  folders: Map<string, FolderState<F>> = new Map();
  /** the files that live in this folder */
  files: Map<string, F> = new Map();
  /** whether the folder is collapsed or not */
  isCollapsed: boolean = true
  /** the color of the folder */
  color: string = 'white'
  /** tags for this folder */
  tags: Record<string, string> = {}
  /** tags for each file, the first level is a file */
  fileTags: Record<string, Record<string, string>> = {}
  /** record for accessing nested tags in O(1) time */
  private taggedFolders: Record<string, Set<string>> = {}


  private listeners: Record<string, (message: Message) => void> = {}

  registerListener (name: string, f: (message: Message) => void): void {
    this.listeners[name] = f;
  }

  private notifyListeners (message: Message) {
    for (const l of Object.values(this.listeners)) {
      l(message)
    }
  }

  expandAll () {
    for (const folder of this.folders.values()) folder.expandAll()
    this.expand()
  }

  expand () {
    this.isCollapsed = false
    this.notifyListeners(['expanded'])
  }

  expandRecursive (name: string[], i = 0) {
    this.expand()
    if (i >= name.length) return
    const folder = this.folders.get(name[i])
    if (folder !== undefined) {
      folder.expandRecursive(name, i + 1)
    }
  }

  collapseAll () {
    for (const folder of this.folders.values()) folder.collapse()
    this.collapse()
  }

  collapse () {
    this.isCollapsed = true
    this.notifyListeners(['collapsed'])
  }

  collapseRecursive (name: string[], i = 0) {
    if (i >= name.length) return
    const folder = this.folders.get(name[i])
    if (folder !== undefined) {
      if (i === name.length - 1) {
        folder.collapseAll()
      } else {
        folder.expandRecursive(name, i + 1)
      }
    }
  }

  untag (tag: string) {
    delete this.tags[tag]
    this.notifyListeners(['untagged'])
  }

  untagAll (tag: string) {
    for (const folder of this.taggedFolders[tag]?.values() ?? []) {
      this.folders.get(folder)?.untagAll(tag)
    }
    delete this.taggedFolders[tag]
    for (const tags of Object.values(this.fileTags)) {
      delete tags[tag]
    }
    this.untag(tag)
  }

  tag (tag: string, value: string) {
    this.tags[tag] = value
    this.notifyListeners(['tagged'])
  }

  tagFile (file: string, tag: string, value: string) {
    if (this.files.has(file)) {
      this.fileTags[file] ??= {}
      this.fileTags[file][tag] = value
      this.notifyListeners(['fileTagged', file])
    }
  }

  tagRecursive (name: string[], tag: string, value: string, i = 0) {
    if (i >= name.length) {
      // nothing
    } else if (i === name.length - 1) {
      if (this.files.has(name[i])) {
        this.tagFile(name[i], tag, value)
      } else if (this.folders.has(name[i])) {
        this.taggedFolders[tag] ??= new Set<string>()
        this.taggedFolders[tag].add(name[i])
        this.folders.get(name[i])!.tag(tag, value)
      }
    } else {
      if (this.folders.has(name[i])) {
        this.taggedFolders[tag] ??= new Set<string>()
        this.taggedFolders[tag].add(name[i])
        const folder = this.folders.get(name[i])!
        folder.tag(tag, value)
        folder.tagRecursive(name, tag, value, i + 1)
      }
    }
  }

  static fromFileTree<T extends ColoredFileLeaf> (tree: FileTree<T>): FolderState<T> {
    const folderState = new FolderState<T>()
    const folderKeys = [...tree.subTrees.keys()].sort()
    for (const folder of folderKeys) {
      folderState.folders.set(folder, FolderState.fromFileTree(tree.subTrees.get(folder)!))
    }

    const fileKeys = [...tree.leafs.keys()].sort()
    for (const file of fileKeys) {
      folderState.files.set(file, tree.leafs.get(file)!)
    }
    const { h, s, v } = tree.__data?.__color ?? { h: 0, s: 0, v: 1 }
    const [r, g, b] = HSVtoRGB(h, s, v)
    folderState.color = `rgb(${r}, ${g}, ${b})`
    folderState.name = tree.name
    return folderState
  }

  render (lvl = 0): string {
    if (this.isCollapsed) return ''
    const stringBuilder: string[] = []
    for (const [name, folder] of this.folders.entries()) {
      stringBuilder.push(' '.repeat(lvl) + '> ' + name + " " + JSON.stringify(folder.tags))
      stringBuilder.push(folder.render(lvl + 1))
    }
    for (const name of this.files.keys()) {
      stringBuilder.push(' '.repeat(lvl) + name +  " " +JSON.stringify(this.fileTags[name]))
    }
    return stringBuilder.filter(_ => _.length > 0).join('\n')
  }

  toString () {
    return this.render()
  }
}
