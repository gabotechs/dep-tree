import { describe, expect, test } from "vitest";
import { FileLeaf, FileTree } from "./FileTree.ts";
import { color } from "./color.ts";

describe('ColoredFileTree', () => {
  it(
    'computes colors',
    {
      nodes: [
        ['foo', 'bar', 'a.ts'],
        ['foo', 'bar', 'b.ts'],
        ['foo', 'baz', 'c.ts'],
        ['foo', 'd.ts'],
        ['a', 'b', 'c', 'd', 'e.ts'],
        ['f.ts']
      ],
    },
    {
      render: `\
__dep_tree_root__ (0, 0.00, 1)
 foo (0, 0.76, 1)
  bar (0, 0.59, 1)
   a.ts -> 0 (0, 0.59, 1)
   b.ts -> 1 (0, 0.59, 1)
  baz (180, 0.59, 1)
   c.ts -> 2 (180, 0.59, 1)
  d.ts -> 3 (0, 0.76, 1)
 a (180, 0.76, 1)
  b (180, 0.59, 1)
   c (180, 0.47, 1)
    d (180, 0.39, 1)
     e.ts -> 4 (180, 0.39, 1)
 f.ts -> 5 (0, 0.00, 1)`,
    }
  )
})

function it (
  name: string,
  input: { nodes: string[][] },
  expected: { render: string }
): void {
  let id = 0

  function newNode (pathBuf: string[]): FileLeaf {
    return { pathBuf, id: id++ }
  }

  test(name, () => {
    const fileTree = FileTree.root<object>()
    for (const node of input.nodes) {
      fileTree.pushNode(newNode(node))
    }
    const coloredFileTree = color(fileTree)

    expect(coloredFileTree.render({
      renderLeaf: leaf => {
        const { h, s, v } = leaf.__color ?? {}
        return ` -> ${leaf.id} (${h}, ${s?.toFixed(2)}, ${v})`;
      },
      renderTree: tree => {
        const { h, s, v } = tree.__data?.__color ?? {}
        return ` (${h}, ${s?.toFixed(2)}, ${v})`;
      }
    })).to.equal(expected.render)
  })
}

