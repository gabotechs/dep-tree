export type FileLeaf<T = {}> = {
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

export class FileTree<T = {}> {
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

  static root<T> (): FileTree<T> {
    return new FileTree<T>(this.ROOT_NAME)
  }

  /**
   * Pushes a node into the tree based on its path buffer.
   */
  pushNode (node: FileLeaf<any>, i = 0) {
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
      this.name += '/' + child.name
      this.subTrees = child.subTrees
      this.leafs = child.leafs
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
   * ['a', 'a/b', 'a/b/c']
   */
  static parents<T> (node: FileLeaf<T>): string[] {
    const parents: string[] = []
    let curr = node.__parent
    while (curr !== undefined && curr.name !== this.ROOT_NAME) {
      parents.push(curr.name)
      curr = curr.__parent
    }
    parents.reverse()
    for (let i = 1; i < parents.length; i++) {
      parents[i] = parents[i - 1] + '/' + parents[i]
    }
    return parents
  }

  /**
   * Return some stats about the position of the node in the graph, see {@link NodeStats}
   * @param node
   */
  static stats<T> (node: FileLeaf<T> | FileTree<T>): NodeStats {
    let depth = 0 // starts with -1 because there's always going to be at least 1 parent.
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

  stats(): NodeStats {
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
    opts: {
      renderLeaf: (leaf: FileLeaf<T>) => string,
      renderTree: (tree: FileTree<T>) => string,
    } = {
      renderTree: () => '',
      renderLeaf: leaf => ` -> ${leaf.id}`
    }
  ): string {
    function render (node: FileTree<T>, indent: number): string[] {
      let line = `${' '.repeat(indent)}${node.name}`
      line += opts.renderTree(node)
      const lines = [line]
      const newIndent = indent + 1
      for (const [_, child] of node.subTrees.entries()) {
        lines.push(...render(child, newIndent))
      }
      for (const [name, child] of node.leafs.entries()) {
        let line = `${' '.repeat(newIndent)}${name}`
        line += opts.renderLeaf(child)
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
