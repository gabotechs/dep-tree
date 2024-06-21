import { describe, expect, test } from "vitest";
import { FileLeaf, FileTree } from "../FileTree.ts";
import { FolderState, Message } from "./FolderState.ts";

describe('FolderState', () => {
  it(
    'should untag all',
    {
      nodes: [
        ['foo', 'bar', 'a.ts'],
      ],
      modify: folderState => {
        folderState.expandAll()
        folderState.tagRecursive(['foo', 'bar', 'a.ts'], 'selected', 'true')
        folderState.untagAll('selected')
      }
    },
    {
      render: `\
> foo {}
 > bar {}
  a.ts {}`,
      events: [
        ['expanded'],
        ['untagged'],
      ]
    }
  );

  it(
    'Multiple operations',
    {
      nodes: [
        ['foo', 'bar', 'a.ts'],
        ['foo', 'bar', 'b.ts'],
        ['foo', 'baz', 'c.ts'],
        ['foo', 'd.ts'],
        ['a', 'b', 'c', 'd', 'e.ts'],
        ['f.ts'],
        ['a', 'b', 'c1', 'd1', 'g.ts'],
      ],
      modify: folderState => {
        folderState.collapseAll()
        folderState.expandAll()
        folderState.collapseRecursive(['a', 'b', 'c'])
        folderState.folders.get('foo')!.folders.get('baz')!.collapse()
        folderState.tagRecursive(['a', 'b', 'c', 'd', 'e.ts'], 'selected', 'true')
        folderState.untagAll('selected')
        folderState.tagRecursive(['foo', 'bar', 'a.ts'], 'selected', 'true')
        folderState.tagRecursive(['f.ts'], 'selected', 'true')
        folderState.folders.get('a')!.folders.get('b')!.folders.get('c1')!.collapseAll()
        folderState.folders.get('a')!.folders.get('b')!.folders.get('c1')!.expand()
      }
    },
    {
      render: `\
> a {}
 > b {}
  > c {}
   > d {}
    e.ts {}
  > c1 {}
   > d1 {}
> foo {"selected":"true"}
 > bar {"selected":"true"}
  a.ts {"selected":"true"}
  b.ts undefined
 > baz {}
 d.ts undefined
f.ts {"selected":"true"}`,
      events: [
        ['collapsed'],
        ['expanded'],
        ["untagged"],
        ["fileTagged", "f.ts"],
      ]
    }
  )
})

function it (
  name: string,
  input: {
    nodes: string[][],
    modify: (folderState: FolderState<object>) => void
  },
  expected: {
    render: string,
    events: Message[]
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
    const evs: Message[] = []
    const folderState = FolderState.fromFileTree(fileTree)
    folderState.registerListener('test', (m) => evs.push(m))
    input.modify(folderState)

    expect(folderState.toString()).to.equal(expected.render)
    expect(evs).to.deep.equal(expected.events)
  })
}
