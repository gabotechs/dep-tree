import { Graph, Node, Link } from "./types.ts";
import { useMemo } from "react";

export interface XLink extends Link {}

export interface XNode extends Node {
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
  const [data, nodes]: [XGraph, Record<number, XNode>] = useMemo(() => {
    const nodes: Record<number, XNode> = {}

    graph.nodes.forEach(node => {
      nodes[node.id] = node
    })

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
    return [graph, nodes]
  }, [graph])

  return { data, nodes }
}
