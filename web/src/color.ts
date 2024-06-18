import { FileLeaf, FileTree } from "./FileTree.ts";
import { scale } from "./@utils/scale.ts";

export interface ColoredFileLeaf {
  __color?: { h: number, s: number, v: number }
}

export function color<T extends object>(tree: FileTree<T>): FileTree<T & ColoredFileLeaf> {
  for (const leaf of tree.iterLeafs()) {
    colorNode(leaf)
  }
  return tree
}

function colorNode<T extends object> (node: FileTree<T> | FileLeaf<T>): { h: number, s: number, v: number } {
  if (node.__parent === undefined) {
    // This is the root node, we want it white.
    const color = { h: 0, s: 0, v: 1 }
    const n = (node as FileTree<ColoredFileLeaf>)
    n.__data ??= {}
    n.__data.__color = color
    return color
  }

  const { h, s, v } = colorNode(node.__parent)
  if (node instanceof FileTree) {
    // this node is a tree, need to accumulate colors.
    const stats = node.stats()

    const nh = (h + 360 * stats.index / stats.total) % 360

    let ns = s === 0 ? 1 : s
    ns = scale(ns - .2, 0, 1, .2, .9)

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
// func (d *DirTree) ColorForDir(dirs []string, format colorFormat) []float64 {
//   node := d.inner()
//   h, s, v := 0., 0., 1.
//   depth := 0
//   for depth < len(dirs) {
//     el, ok := node.Get(dirs[depth])
//     if !ok {
//       return []float64{0, 0, 0}
//     }
//
//     // It might happen that all the nodes have some common folders, like src/,
//     // so if literally all of them have the same common folders, we do not want to take
//     // them into account for reducing the saturation, as they will appear very faded.
//     if node.Len() > 1 {
//       h = float64(int(h+360*float64(el.index)/float64(node.Len())) % 360)
//       if s == 0 {
//         s = 1
//       }
//       s -= .2
//       s = utils.Scale(s, 0, 1, .2, .9)
//     }
//
//     depth += 1
//     node = el.entry.inner()
//   }
//   if format == RGB {
//     r, g, b := HSVToRGB(h, s, v)
//     return []float64{float64(r), float64(g), float64(b)}
//   } else {
//     return []float64{h, s, v}
//   }
// }
