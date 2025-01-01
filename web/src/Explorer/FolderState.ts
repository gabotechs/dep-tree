import { FileTree } from "../FileTree.ts";
import { HSVtoRGB } from "../@utils/HSVtoRGB.ts";
import { ColoredFileLeaf } from "../color.ts";

export type Message =
  ['tagged', string, string] |
  ['untagged', string] |
  ['fileTagged', string, string, string] |
  ['fileUntagged', string, string]

export class FolderState<F> {

  /** the name of the folder */
  name: string = ''
  /** the subfolders that live in this folder */
  folders: Map<string, FolderState<F>> = new Map();
  /** the files that live in this folder */
  files: Map<string, F> = new Map();
  /** the color of the folder */
  color: string = 'white'
  /** tags for this folder */
  private _tags: Record<string, string> = {}
  get tags (): Record<string, string> {
    return this._tags;
  }

  /** tags for each file, the first level is a file */
  private _fileTags: Record<string, Record<string, string>> = {}
  get fileTags (): Record<string, Record<string, string>> {
    return this._fileTags;
  }

  /**
   * record for accessing nested tags in O(1) time.
   *
   * This is a record from tag -> folder set
   */
  private taggedFolders: Record<string, Set<string>> = {}
  private parent?: FolderState<F>

  private listeners: Record<string, (message: Message) => void> = {}

  registerListener (name: string, f: (message: Message) => void): void {
    this.listeners[name] = f;
  }

  private notifyListeners (message: Message) {
    for (const l of Object.values(this.listeners)) {
      l(message)
    }
  }

  tagAllFolders (tag: string, value: string) {
    for (const folder of this.folders.values()) folder.tagAllFolders(tag, value)
    this.tag(tag, value)
  }

  tagAll (tag: string, value: string) {
    for (const folder of this.folders.values()) folder.tagAll(tag, value)
    this.tag(tag, value)
    for (const file of this.files.keys()) {
      this.tagFile(file, tag, value)
    }
  }

  untagAllFolders (tag: string) {
    for (const folder of this.taggedFolders[tag]?.values() ?? []) {
      this.folders.get(folder)?.untagAllFolders(tag)
    }
    this.untag(tag)
  }

  untagAllFoldersFrom (name: string[], tag: string) {
    // eslint-disable-next-line @typescript-eslint/no-this-alias
    let folder: FolderState<F> | undefined = this
    for (const n of name) {
      folder = folder.folders.get(n)
      if (!folder) return
    }
    folder.untagAllFolders(tag)
  }

  untag (tag: string) {
    delete this._tags[tag]
    if (this.parent !== undefined) {
      this.parent.taggedFolders[tag]?.delete(this.name)
    }
    this.notifyListeners(['untagged', tag])
  }

  untagAll (tag: string) {
    for (const folder of this.taggedFolders[tag]?.values() ?? []) {
      this.folders.get(folder)?.untagAll(tag)
    }
    for (const tags of Object.values(this._fileTags)) {
      delete tags[tag]
    }
    this.untag(tag)
  }

  tag (tag: string, value: string) {
    this._tags[tag] = value
    if (this.parent !== undefined) {
      this.parent.taggedFolders[tag] ??= new Set<string>()
      this.parent.taggedFolders[tag].add(this.name)
    }
    this.notifyListeners(['tagged', tag, value])
  }

  tagFile (file: string, tag: string, value: string) {
    if (this.files.has(file)) {
      this._fileTags[file] ??= {}
      this._fileTags[file][tag] = value
      this.notifyListeners(['fileTagged', file, tag, value])
    }
  }

  tagRecursive (name: string[], tag: string, value: string) {
    if (name.length === 0) return
    this.tag(tag, value)
    // eslint-disable-next-line @typescript-eslint/no-this-alias
    let folder: FolderState<F> | undefined = this
    for (const n of name.slice(0, name.length - 1)) {
      folder = folder.folders.get(n)
      if (!folder) return
      folder.tag(tag, value)
    }
    const last = name[name.length - 1]

    if (folder.files.has(last)) {
      folder.tagFile(last, tag, value)
    } else if (folder.folders.has(last)) {
      folder.folders.get(last)!.tag(tag, value)
    }
  }

  * allFiles (): Generator<F> {
    for (const folder of this.folders.values()) {
      yield  * folder.allFiles()
    }
    for (const file of this.files.values()) {
      yield file
    }
  }

  static fromFileTree<T extends ColoredFileLeaf> (tree: FileTree<T>): FolderState<T> {
    const folderState = new FolderState<T>()
    for (const folder of tree.subTrees.keys()) {
      const childFolderState = FolderState.fromFileTree(tree.subTrees.get(folder)!)
      childFolderState.parent = folderState
      folderState.folders.set(folder, childFolderState)
    }

    for (const file of tree.leafs.keys()) {
      folderState.files.set(file, tree.leafs.get(file)!)
    }
    const { h, s, v } = tree.__data?.__color ?? { h: 0, s: 0, v: 1 }
    const [r, g, b] = HSVtoRGB(h, s, v)
    folderState.color = `rgb(${r}, ${g}, ${b})`
    folderState.name = tree.name
    return folderState
  }

  render (lvl = 0): string {
    const stringBuilder: string[] = []
    for (const [name, folder] of this.folders.entries()) {
      stringBuilder.push(' '.repeat(lvl) + '> ' + name + " " + JSON.stringify(folder._tags))
      stringBuilder.push(folder.render(lvl + 1))
    }
    for (const name of this.files.keys()) {
      stringBuilder.push(' '.repeat(lvl) + name + " " + JSON.stringify(this._fileTags[name]))
    }
    return stringBuilder.filter(_ => _.length > 0).join('\n')
  }

  toString () {
    return this.render()
  }
}
