import { Graph, Node, Link } from "./types.ts";
import { useMemo } from "react";
import { FileTree } from "./FileTree.ts";
import { ColoredFileLeaf, color } from "./color.ts";

export interface XLink extends Link {}

export interface XNode extends Node, ColoredFileLeaf {
  neighbors?: XNode[]
  links?: XLink[]
  x?: number
  y?: number
  z?: number
}

export interface XGraph extends Graph {
  nodes: XNode[];
  links: XLink[];
}

export function useData(graph: Graph) {
  const [data, nodes, fileTree] = useMemo(() => {
    const nodes: Record<number, XNode> = {}
    const fileTree = FileTree.root<XNode>()

    graph.nodes.forEach(node => {
      fileTree.pushNode(node)
      nodes[node.id] = node
    })

    fileTree.squash()
    const coloredFileTree = color(fileTree)

    // cross-link node objects
    graph.links.forEach(link => {
      const a = nodes[link.from];
      const b = nodes[link.to];
      !a.neighbors && (a.neighbors = []);
      !b.neighbors && (b.neighbors = []);
      a.neighbors.push(b);
      b.neighbors.push(a);

      !a.links && (a.links = []);
      !b.links && (b.links = []);
      a.links.push(link);
      b.links.push(link);
    });
    return [graph as XGraph, nodes, coloredFileTree]
  }, [graph])

  return { data, nodes, fileTree }
}
