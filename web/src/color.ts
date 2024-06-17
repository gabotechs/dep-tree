import { FileLeaf, FileTree } from "./FileTree.ts";
import { scale } from "./@utils/scale.ts";

export interface ColoredFileLeaf {
  __color?: { h: number, s: number, v: number }
}

export function color<T extends {}>(tree: FileTree<T>): FileTree<T & ColoredFileLeaf> {
  for (const leaf of tree.iterLeafs()) {
    colorNode(leaf)
  }
  return tree
}

function colorNode<T extends {}> (node: FileTree<T> | FileLeaf<T>): { h: number, s: number, v: number } {
  if (node.__parent === undefined) {
    // This is the root node, we want it white.
    const color = { h: 0, s: 0, v: 1 }
    const n = (node as FileTree<ColoredFileLeaf>)
    n.__data ??= {}
    n.__data.__color = color
    return color
  }

  let { h, s, v } = colorNode(node.__parent)
  if (node instanceof FileTree) {
    // this node is a tree, need to accumulate colors.
    const stats = node.stats()

    const nh = (h + 360 * stats.index / stats.total) % 360

    if (s === 0) {
      s = 1
    }
    s -= .2
    const ns = scale(s, 0, 1, .2, .9)

    const color = { h: nh, s: ns, v }
    const n = (node as FileTree<ColoredFileLeaf>)
    n.__data ??= {}
    n.__data.__color = color
    return color
  } else {
    // this is a leaf, just show the parent tree color.
    const color = { h, s, v }
    ;(node as ColoredFileLeaf).__color = color
    return color
  }
}


