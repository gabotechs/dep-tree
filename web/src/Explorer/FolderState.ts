import { FileTree } from "../FileTree.ts";
import { HSVtoRGB } from "../@utils/HSVtoRGB.ts";
import { ColoredFileLeaf } from "../color.ts";

export class FolderState<F> {
  name: string = ''
  folders: Map<string, FolderState<F>> = new Map();
  files: Map<string, F> = new Map();
  folded: boolean = true
  color: string = 'white'

  private modifyAll(value: boolean) {
    for (const folder of this.folders.values()) {
      folder.modifyAll(value)
    }
    this.modify(value)
  }

  private modify(value: boolean) {
    this.folded = value
  }

  unfoldAll() {
    this.modifyAll(false)
  }

  unfold() {
    this.modify(false)
  }

  unfoldByName(name: string[], i=0) {
    if (i >= name.length) return
    const folder = this.folders.get(name[i])
    if (folder !== undefined) {
      folder.modify(false)
      folder.unfoldByName(name, i+1)
    }
  }

  collapseAll() {
    this.modifyAll(true)
  }

  collapse() {
    this.modify(true)
  }

  collapseByName(name: string[], i=0) {
    if (i >= name.length) return
    const folder = this.folders.get(name[i])
    if (folder !== undefined) {
      if (i === name.length-1) {
        folder.collapseAll()
      } else {
        folder.unfoldByName(name, i+1)
      }
    }
  }

  static fromFileTree<T extends ColoredFileLeaf>(tree: FileTree<T>): FolderState<T> {
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

  render(lvl=0): string {
    if (this.folded) return ''
    const stringBuilder: string[] = []
    for (const [name, folder] of this.folders.entries()) {
      stringBuilder.push(' '.repeat(lvl) + '> ' + name)
      stringBuilder.push(folder.render(lvl+1))
    }
    for (const name of this.files.keys()) {
      stringBuilder.push(' '.repeat(lvl) + name)
    }
    return stringBuilder.filter(_ => _.length > 0).join('\n')
  }

  toString() {
    return this.render()
  }
}
