import { Graph, Link, Node } from "./types.ts";
import { useMemo } from "react";
import { FileTree } from "./FileTree.ts";
import { color, ColoredFileLeaf } from "./color.ts";
import { hashInt } from "./@utils/hashInt.ts";

export interface XLink extends Link {
  isDir?: boolean
  isPackage?: boolean
}

export interface XNode extends Node, ColoredFileLeaf {
  neighbors?: XNode[]
  links?: XLink[]
  x?: number
  y?: number
  z?: number
  isDir?: boolean
  isPackage?: boolean
}

export interface XGraph extends Graph {
  nodes: XNode[];
  links: XLink[];
}

/**
 * Processes the graph performing some operations like:
 * - storing the nodes in a Record format
 * - adding Folder and Package nodes
 * - coloring the nodes
 *
 * WARNING: this function mutates the provided graph for performance reasons
 */
export function useData (graph: Graph) {
  const [data, nodes] = useMemo(() => {

    // Build a record of nodes. This will be useful across the application for accessing
    // a node given its id in O(1) time.
    const nodes: Record<number, XNode> = {}
    graph.nodes.forEach(node => (nodes[node.id] = node))

    // Build the FileTree with all the nodes.
    const fileTree = FileTree.root<XNode>()
    graph.nodes.forEach(node => fileTree.pushNode(node))
    fileTree.squash()
    color(fileTree)

    // Check if there are more than one groups in all the nodes.
    const groups = new Set<string>()
    for (const node of Object.values(nodes)) {
      groups.add(node.group ?? '')
      if (groups.size > 1) break
    }
    const needsToAddGroups = groups.size > 1

    // Create virtual nodes representing folders and packages, so that additional forces are applied
    // between nodes withing the same folder/package. This will artificially concentrate related files
    // together, and will stabilize the visualization.
    for (const node of Object.values(nodes)) {
      if (needsToAddGroups) {
        const groupId = hashInt(`__dep_tree_group__${node.group}`)
        if (!(groupId in nodes)) {
          nodes[groupId] = newGroupNode(groupId)
          graph.nodes.push(nodes[groupId])
        }
        graph.links.push(newGroupLink(node.id, groupId))
      }

      for (const folderName of FileTree.parentFolders(node)) {
        if (folderName.startsWith(FileTree.ROOT_NAME)) continue
        const folderId = hashInt(`__dep_tree_folder__${folderName}`)
        if (!(folderId in nodes)) {
          nodes[folderId] = newFolderNode(folderId)
          graph.nodes.push(nodes[folderId])
        }
        graph.links.push(newFolderLink(node.id, folderId))
      }
    }

    // Cross-link nodes and edges. Needed for the rendering library to work.
    graph.links.forEach(link => {
      const a = nodes[link.from];
      const b = nodes[link.to];
      a.neighbors ??= [];
      b.neighbors ??= [];
      a.neighbors.push(b);
      b.neighbors.push(a);

      a.links ??= [];
      b.links ??= [];
      a.links.push(link);
      b.links.push(link);
    });
    return [graph as XGraph, nodes]
  }, [graph])

  return { data, nodes }
}

function newFolderNode (id: number): XNode {
  return {
    id,
    isDir: true,
    // bellow are just defaults
    dirName: "",
    fileName: "",
    isEntrypoint: false,
    isPackage: false,
    loc: 0,
    pathBuf: [],
    size: 0
  }
}

function newGroupNode (id: number): XNode {
  return {
    id,
    isPackage: true,
    // bellow are just defaults
    dirName: "",
    fileName: "",
    isEntrypoint: false,
    isDir: false,
    loc: 0,
    pathBuf: [],
    size: 0
  }
}

function newFolderLink (nodeId: number, folderId: number): XLink {
  return {
    from: nodeId,
    to: folderId,
    isDir: true,
    // below are just defaults
    isCyclic: false,
    isPackage: false,
  }
}

function newGroupLink (nodeId: number, folderId: number): XLink {
  return {
    from: nodeId,
    to: folderId,
    isPackage: true,
    // below are just defaults
    isCyclic: false,
    isDir: false,
  }
}
