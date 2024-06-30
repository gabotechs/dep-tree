export type FileLeaf<T = object> = {
  id: number
  pathBuf: string[]
  __parent?: FileTree<T>
  __index?: number
} & T

export interface NodeStats {
  /**
   * If the node is a subtree or a leaf.
   */
  kind: 'tree' | 'leaf'
  /**
   * The depth in the tree at which this node is located.
   */
  depth: number,
  /**
   * The index of the node among the nodes of the same kind that
   * share the same parent as the current node.
   */
  index: number,
  /**
   * The total amount of nodes of the same kind that share the
   * same parent as the current node.
   */
  total: number
}

export class FileTree<T = object> {
  static readonly ROOT_NAME = '__dep_tree_root__'

  name: string
  __data?: T
  __parent?: FileTree<T>
  __index?: number
  subTrees: Map<string, FileTree<T>> = new Map();
  leafs: Map<string, FileLeaf<T>> = new Map();

  protected constructor (name: string) {
    this.name = name;
  }

  static root<T = object> (): FileTree<T> {
    return new FileTree<T>(this.ROOT_NAME)
  }

  /**
   * Pushes a node into the tree based on its path buffer.
   */
  pushNode (node: FileLeaf<T>, i = 0) {
    if (i >= node.pathBuf.length) return

    const tree = this.subTrees.get(node.pathBuf[i])
    if (tree instanceof FileTree) {
      // Already a FileTree node in this location.
      tree.pushNode(node, i + 1)
    } else if (tree !== undefined) {
      // Already a `T` node in this location, which was supposed to be a leaf.
      throw new Error(`Cannot push node ${node.id} with pathBuf ${node.pathBuf} into the tree, there's already a leaf node at ${node.pathBuf.slice(0, i)}`)
    } else if (i === node.pathBuf.length - 1) {
      // Exhausted all pathBuf elements, so add a leaf node.
      node.__parent = this
      node.__index = this.leafs.size
      this.leafs.set(node.pathBuf[i], node)
    } else {
      // New FileTree element.
      const tree = new FileTree<T>(node.pathBuf[i])
      tree.__parent = this
      tree.__index = this.subTrees.size
      this.subTrees.set(node.pathBuf[i], tree)
      tree.pushNode(node, i + 1)
    }
  }

  /**
   * Squashes all the single dir nestings into one, for example:
   *
   * a
   *  b
   *   c
   *    d.ts
   *    e.ts
   *
   * will be squashed into
   *
   * a/b/c
   *  d.ts
   *  e.ts
   */
  squash (): void {
    if (this.subTrees.size === 1 && this.leafs.size === 0) {
      const child = [...this.subTrees.values()][0]
      const oldName = this.name
      this.name += '/' + child.name
      this.subTrees = child.subTrees
      this.leafs = child.leafs
      this.__parent?.subTrees.delete(oldName)
      this.__parent?.subTrees.set(this.name, this)
      this.squash()
      for (const tree of this.subTrees.values()) {
        tree.__parent = this
      }
      for (const leaf of this.leafs.values()) {
        leaf.__parent = this
      }
    } else {
      for (const child of this.subTrees.values()) {
        child.squash()
      }
    }
  }

  /**
   * A Map in JS will be iterated in the order of insertion. This function
   * orders the internal Maps alphabetically.
   */
  order() {
    const subTrees = [...this.subTrees.entries()]
    this.subTrees.clear()
    subTrees.sort(([a], [b]) => a > b ? 1 : -1)
    let i = 0
    for (const [k, v] of subTrees) {
      v.__index = i
      v.order()
      this.subTrees.set(k, v)
      i++
    }

    const leafs = [...this.leafs.entries()]
    this.leafs.clear()
    leafs.sort(([a], [b]) => a > b ? 1 : -1)
    i = 0
    for (const [k, v] of leafs) {
      v.__index = i
      this.leafs.set(k, v)
      i++
    }
  }

  /**
   * Retrieves the parent {@link FileTree} of a leaf node.
   */
  static parentTree<T> (leaf: FileLeaf<T>): FileTree<T> {
    if (leaf.__parent === undefined) {
      throw new Error(`Node ${leaf.id} with pathBuf ${leaf.pathBuf} does not have a parent, maybe it was never added to the FileTree?`)
    }
    return leaf.__parent
  }

  /**
   * Retrieves the parent dirs to which this node belongs, for example:
   *
   * a
   *  b
   *   c
   *    node.ts
   *
   * will return:
   *
   * ['a', 'b', 'c']
   */
  static parentFolders<T> (node: FileLeaf<T>): string[] {
    const parents: string[] = []
    let curr = node.__parent
    while (curr !== undefined && curr.name !== this.ROOT_NAME) {
      parents.push(curr.name)
      curr = curr.__parent
    }
    parents.reverse()
    return parents
  }

  /**
   * Return some stats about the position of the node in the graph, see {@link NodeStats}
   * @param node
   */
  static stats<T> (node: FileLeaf<T> | FileTree<T>): NodeStats {
    let depth = 0
    let parent = node.__parent
    while (parent !== undefined) {
      depth++
      parent = parent.__parent
    }

    return {
      kind: node instanceof FileTree ? 'tree' : 'leaf',
      depth,
      // NOTE: if node.__index and node.__parent are undefined, then `node` is the root node
      index: node.__index ?? 0,
      total: node.__parent?.subTrees.size ?? 1
    }
  }

  stats (): NodeStats {
    return FileTree.stats(this)
  }

  /**
   * Iterates all the leaf nodes.
   */
  * iterLeafs (): Generator<FileLeaf<T>> {
    for (const tree of this.subTrees.values()) {
      yield * tree.iterLeafs()
    }
    for (const leaf of this.leafs.values()) {
      yield leaf
    }
  }

  /**
   * Renders the tree in a human-readable format.
   */
  render (
    {
      renderTree = () => '',
      renderLeaf = leaf => ` -> ${leaf.id}`
    }: {
      renderLeaf?: (leaf: FileLeaf<T>) => string,
      renderTree?: (tree: FileTree<T>) => string,
    } = {}
  ): string {
    function render (node: FileTree<T>, indent: number): string[] {
      let line = `${' '.repeat(indent)}${node.name}`
      line += renderTree(node)
      const lines = [line]
      const newIndent = indent + 1
      for (const child of node.subTrees.values()) {
        lines.push(...render(child, newIndent))
      }
      for (const [name, child] of node.leafs.entries()) {
        let line = `${' '.repeat(newIndent)}${name}`
        line += renderLeaf(child)
        lines.push(line)
      }
      return lines
    }

    return render(this, 0).join('\n')
  }

  toString (): string {
    return this.render()
  }
}
