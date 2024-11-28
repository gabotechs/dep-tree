import { describe, expect, test } from 'vitest'
import { FileLeaf, FileTree, NodeStats } from './FileTree'

describe("FileTree", () => {
  it(
    'non-squash',
    {
      nodes: [
        ['foo', 'bar', 'a.ts'],
        ['foo', 'bar', 'b.ts'],
        ['foo', 'baz', 'c.ts'],
        ['foo', 'd.ts'],
        ['a', 'b', 'c', 'd', 'e.ts'],
        ['f.ts']
      ],
      squash: false
    },
    {
      render: `\
__dep_tree_root__
 a
  b
   c
    d
     e.ts -> 4
 foo
  bar
   a.ts -> 0
   b.ts -> 1
  baz
   c.ts -> 2
  d.ts -> 3
 f.ts -> 5`,
      parents: [
        ['a', 'b', 'c', 'd'],
        ['foo', 'bar'],
        ['foo', 'bar'],
        ['foo', 'baz'],
        ['foo'],
        []
      ],
      parentStats: [
        { kind: 'tree', depth: 4, total: 1, index: 0 },
        { kind: 'tree', depth: 2, total: 2, index: 0 },
        { kind: 'tree', depth: 2, total: 2, index: 0 },
        { kind: 'tree', depth: 2, total: 2, index: 1 },
        { kind: 'tree', depth: 1, total: 2, index: 1 },
        { kind: 'tree', depth: 0, total: 1, index: 0 },
      ],
      leafs: [4, 0, 1, 2, 3, 5]
    }
  )

  it(
    'squash',
    {
      nodes: [
        ['a', 'b', 'c', 'd', 'e.ts'],
        ['foo', 'bar', 'a.ts'],
        ['foo', 'bar', 'b.ts'],
        ['foo', 'baz', 'c.ts'],
        ['foo', 'd.ts'],
        ['f.ts'],
        ['foo', 'bar', 'a', 'b', 'g.ts']
      ],
      squash: true
    },
    {
      render: `\
__dep_tree_root__
 a/b/c/d
  e.ts -> 0
 foo
  bar
   a/b
    g.ts -> 6
   a.ts -> 1
   b.ts -> 2
  baz
   c.ts -> 3
  d.ts -> 4
 f.ts -> 5`,
      parents: [
        ['a/b/c/d'],
        ['foo', 'bar', 'a/b'],
        ['foo', 'bar'],
        ['foo', 'bar'],
        ['foo', 'baz'],
        ['foo'],
        [],
      ],
      parentStats: [
        { kind: 'tree', depth: 1, total: 2, index: 0 },
        { kind: 'tree', depth: 3, total: 1, index: 0 },
        { kind: 'tree', depth: 2, total: 2, index: 0 },
        { kind: 'tree', depth: 2, total: 2, index: 0 },
        { kind: 'tree', depth: 2, total: 2, index: 1 },
        { kind: 'tree', depth: 1, total: 2, index: 1 },
        { kind: 'tree', depth: 0, total: 1, index: 0 },
      ],
      leafs: [0, 6, 1, 2, 3, 4, 5]
    }
  )
})

function it (
  name: string,
  input: { nodes: string[][], squash: boolean },
  expected: { render: string, parents: string[][], parentStats: NodeStats[], leafs: number[] }
): void {
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

    // ensure tree integrity
    ensureChildrenParents(fileTree)
    ensureFolderNamesMatchKeys(fileTree)

    // expect the serialization to be equal
    expect(fileTree.toString()).to.equal(expected.render)

    // check the parent dirs for each node
    const parents: string[][] = []
    for (const n of fileTree.iterLeafs()) {
      parents.push(FileTree.parentFolders(n))
    }
    expect(parents).to.deep.equal(expected.parents)

    // check the stats of each node
    const stats = []
    for (const n of fileTree.iterLeafs()) {
      stats.push(FileTree.stats(FileTree.parentTree(n)))
    }
    expect(stats).to.deep.equal(expected.parentStats)

    // check the leafs
    expect([...fileTree.iterLeafs()].map(_ => _.id)).to.deep.equal(expected.leafs)
  })
}

function ensureChildrenParents (tree: FileTree): void {
  for (const [name, child] of tree.subTrees.entries()) {
    if (child.__parent !== tree) throw new Error(`${name} is a child of ${tree.name}, but ${name}'s parent is ${child.__parent?.name}`)
    ensureChildrenParents(child)
  }
  for (const [name, child] of tree.leafs.entries()) {
    if (child.__parent !== tree) throw new Error(`${name} is a child of ${tree.name}, but ${name}'s parent is ${child.__parent?.name}`)
  }
}

function ensureFolderNamesMatchKeys (tree: FileTree): void {
  for (const [name, child] of tree.subTrees.entries()) {
    if (child.name !== name) throw new Error(`sub-tree ${tree.name} has a child named ${child.name} under the key ${name}. They key should match the child's name`)
    ensureFolderNamesMatchKeys(child)
  }
}
