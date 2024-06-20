import { describe, expect, test } from "vitest";
import { FileLeaf, FileTree } from "../FileTree.ts";
import { FolderState } from "./FolderState.ts";

describe('FolderState', () => {
  it(
    'name',
    {
      nodes: [
        ['foo', 'bar', 'a.ts'],
        ['foo', 'bar', 'b.ts'],
        ['foo', 'baz', 'c.ts'],
        ['foo', 'd.ts'],
        ['a', 'b', 'c', 'd', 'e.ts'],
        ['f.ts']
      ],
      modify: folderState => folderState.unfoldAll().collapseByName(['a', 'b', 'c'])
    },
    {
      render: `\
> a
 > b
  > c
> foo
 > bar
  a.ts
  b.ts
 > baz
  c.ts
 d.ts
f.ts`
    }
  )
})

function it(
  name: string,
  input: {
    nodes: string[][],
    modify: (folderState: FolderState<object>) => void
  },
  expected: {
    render: string
  }
) {
  let id = 0

  function newNode (pathBuf: string[]): FileLeaf {
    return { pathBuf, id: id++ }
  }
  test(name, () => {
    const fileTree = FileTree.root()
    for (const node of input.nodes) {
      fileTree.pushNode(newNode(node))
    }
    const folderState = FolderState.fromFileTree(fileTree)
    input.modify(folderState)
    expect(folderState.toString()).to.equal(expected.render)
  })
}
