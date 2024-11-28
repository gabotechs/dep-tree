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
        folderState.tagAllFolders('expanded', 'true')
        folderState.tagRecursive(['foo', 'bar', 'a.ts'], 'selected', 'true')
        folderState.untagAll('selected')
      }
    },
    {
      render: `\
> foo {"expanded":"true"}
 > bar {"expanded":"true"}
  a.ts {}`,
      events: [
        ["tagged", "expanded", "true"],
        ['tagged', 'selected', 'true'],
        ["untagged", "selected"],
      ]
    }
  );

  it(
    'should work with squashed folders',
    {
      nodes: [
        ['foo', 'bar', 'a.ts'],
        ['foo', 'bar', 'b.ts'],
        ['bar', 'c.ts'],
        ['d.ts']
      ],
      squash: true,
      modify: folderState => {
        folderState.tagAllFolders('expanded', 'true')
      }
    },
    {
      render: `\
> bar {"expanded":"true"}
 c.ts undefined
> foo/bar {"expanded":"true"}
 a.ts undefined
 b.ts undefined
d.ts undefined`,
      events: [
        ['tagged', 'expanded', 'true']
      ]
    }
  )

  it(
    'Untag does not remove all parent tag folders',
    {
      nodes: [
        ['foo', 'bar', 'a.ts'],
        ['foo', 'baz', 'b.ts'],
        ['_']
      ],
      squash: true,
      modify: folderState => {
        folderState.tagAllFolders('expanded', 'true')
        folderState.folders.get('foo')!.folders.get('bar')!.untag('expanded')
        folderState.untagAllFolders('expanded')
      }
    },
    {
      render: `\
> foo {}
 > bar {}
  a.ts undefined
 > baz {}
  b.ts undefined
_ undefined`,
      events: [
        ['tagged', 'expanded', 'true'],
        ['untagged', 'expanded']
      ]
    }
  )

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
        folderState.tagAllFolders('expanded', 'true')
        folderState.untagAllFoldersFrom(['a', 'b', 'c'], 'expanded')
        folderState.folders.get('foo')!.folders.get('baz')!.untag('expanded')
        folderState.tagRecursive(['a', 'b', 'c', 'd', 'e.ts'], 'selected', 'true')
        folderState.untagAll('selected')
        folderState.tagRecursive(['foo', 'bar', 'a.ts'], 'selected', 'true')
        folderState.tagRecursive(['f.ts'], 'selected', 'true')
        folderState.folders.get('a')!.folders.get('b')!.folders.get('c1')!.untagAll('expanded')
        folderState.folders.get('a')!.folders.get('b')!.folders.get('c1')!.tag('expanded', 'true')
      }
    },
    {
      render: `\
> a {"expanded":"true"}
 > b {"expanded":"true"}
  > c {}
   > d {}
    e.ts {}
  > c1 {"expanded":"true"}
   > d1 {}
    g.ts undefined
> foo {"expanded":"true","selected":"true"}
 > bar {"expanded":"true","selected":"true"}
  a.ts {"selected":"true"}
  b.ts undefined
 > baz {}
  c.ts undefined
 d.ts undefined
f.ts {"selected":"true"}`,
      events: [
        ["tagged", "expanded", "true"],
        ["tagged", "selected", "true"],
        ["untagged", "selected"],
        ["tagged", "selected", "true"],
        ["tagged", "selected", "true"],
        ["fileTagged", "f.ts", "selected", "true"],
      ]
    }
  )
})

function it (
  name: string,
  input: {
    nodes: string[][],
    modify: (folderState: FolderState<object>) => void,
    squash?: boolean
  },
  expected: {
    render: string,
    events: Message[],
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
    if (input.squash) fileTree.squash()
    fileTree.order()
    const evs: Message[] = []
    const folderState = FolderState.fromFileTree(fileTree)
    folderState.registerListener('test', (m) => evs.push(m))
    input.modify(folderState)

    expect(folderState.toString()).to.equal(expected.render)
    expect(evs).to.deep.equal(expected.events)
  })
}
