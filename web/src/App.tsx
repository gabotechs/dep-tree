import { useEffect, useRef, useState } from "react";
import ForceGraph, { ForceGraphMethods, LinkObject, NodeObject } from "react-force-graph-3d";

// @ts-ignore
import { UnrealBloomPass } from 'three/examples/jsm/postprocessing/UnrealBloomPass.js'
// @ts-ignore
import { CSS2DObject, CSS2DRenderer } from 'three/examples/jsm/renderers/CSS2DRenderer.js'
import { Data } from "./data.ts";
import { useData, XLink, XNode } from "./useData.ts";
import './App.css'

const SETTINGS = {
  DEFAULT_DISTANCE: 400,
  NODE_RESOLUTION: 16,
  LINK_HIGHLIGHT_WIDTH: 2,
  BLOOM_PASS_STRENGTH: .5,
  BLOOM_PASS_RADIUS: 0.5,
  BLOOM_PASS_THRESHOLD: 0.1,
  DOUBLE_CLICK_INTERVAL: 350,

  NODE_ALPHA: 1,
  UNSELECTED_NODE_ALPHA: 0.1,
  LINK_ALPHA: 0.3,
  UNSELECTED_LINK_ALPHA: 0.1,

  LINK_DISTANCE: 30, // https://github.com/vasturiano/d3-force-3d?tab=readme-ov-file#link_distance
  FILE_NODE_REPULSION_FORCE: 30, // https://github.com/vasturiano/d3-force-3d?tab=readme-ov-file#manyBody_strength
  DIR_NODE_REPULSION_FORCE: 50,
  PACKAGE_NODE_REPULSION_FORCE: 50,
  FILE_LINK_STRENGTH_FACTOR: 1,
  DIR_LINK_STRENGTH_FACTOR: 3,
  PACKAGE_LINK_STRENGTH_FACTOR: 1,
  HIGHLIGHT_CYCLES: false
}

function App () {
  const [highlightNodes, setHighlightNodes] = useState(new Set<XNode>())
  const [highlightLinks, setHighlightLinks] = useState(new Set<XLink>())
  const [selectedNode, setSelectedNode] = useState<XNode>()

  const graph = useRef<ForceGraphMethods<NodeObject<XNode>, LinkObject<XNode, XLink>>>();
  const { data, nodes } = useData(Data.__INLINE_DATA)

  const lastBackgroundClick = useRef(0);

  function backgroundClick () {
    const now = new Date().getTime()
    if (selectedNode) {
      selectNode(undefined)
    } else {
      if (now - lastBackgroundClick.current < SETTINGS.DOUBLE_CLICK_INTERVAL) {
        graph.current?.zoomToFit(SETTINGS.DEFAULT_DISTANCE)
      }
    }
    lastBackgroundClick.current = now
  }

  function colorNode (node: XNode) {
    const [r, g, b] = node.isEntrypoint || node.color === undefined ? [255, 255, 255] : node.color
    let alpha = SETTINGS.NODE_ALPHA
    if (highlightNodes.size > 0 && !highlightNodes.has(node)) alpha = SETTINGS.UNSELECTED_NODE_ALPHA
    return `rgba(${r}, ${g}, ${b}, ${alpha})`;
  }

  function colorLink (link: XLink) {
    let alpha = SETTINGS.LINK_ALPHA
    if (highlightLinks.size > 0 && !highlightLinks.has(link)) alpha = SETTINGS.UNSELECTED_LINK_ALPHA
    if (link.isCyclic && SETTINGS.HIGHLIGHT_CYCLES) return `indianred`;
    return `rgba(255, 255, 255, ${alpha})`;
  }

  function nodeThreeObject (node: XNode) {
    const id = node.id.toString()
    let nodeEl = document.getElementById(id)
    if (!nodeEl) {
      nodeEl = document.createElement('div')
      nodeEl.id = id
      nodeEl.className = 'nodeLabel nodeLabelSelected'
      nodeEl.textContent = node.dirName + node.fileName
      nodeEl.style.color = colorNode(node)
    }
    if (highlightNodes.has(node)) {
      return new CSS2DObject(nodeEl)
    } else {
      return undefined
    }
  }

  function selectNode (node?: XNode) {
    if (node === selectedNode) {
      node = undefined
    }
    setSelectedNode(node)

    const newHighlightNodes = new Set<XNode>()
    const newHighlightLinks = new Set<XLink>()

    if (node !== undefined) {
      newHighlightNodes.add(node);
      node.neighbors?.forEach(neighbor => newHighlightNodes.add(neighbor));
      node.links?.forEach(link => newHighlightLinks.add(link));
    }

    setHighlightNodes(newHighlightNodes)
    setHighlightLinks(newHighlightLinks)
  }

  function centerOnNode (node: XNode) {
    const distance = SETTINGS.DEFAULT_DISTANCE;
    const { x = 1, y = 1, z = 1 } = node
    const distRatio = 1 + distance / Math.hypot(x, y, z);

    const newPos = node.x || node.y || node.z
      ? { x: x * distRatio, y: y * distRatio, z: z * distRatio }
      : { x: 0, y: 0, z: distance }; // special case if node is in (0,0,0)

    graph.current?.cameraPosition(newPos, { x, y, z }, 1000)
  }

  useEffect(() => {
    const g = graph.current
    if (g === undefined)  return

    const bloomPass = new UnrealBloomPass(
      undefined, // resolution
      SETTINGS.BLOOM_PASS_STRENGTH, // strength
      SETTINGS.BLOOM_PASS_RADIUS, // radius
      SETTINGS.BLOOM_PASS_THRESHOLD // threshold
    );

    g.postProcessingComposer().addPass(bloomPass)

    g.d3Force('link')
      ?.distance(() => SETTINGS.LINK_DISTANCE)
      .strength((link: XLink) => {
        let f = SETTINGS.FILE_LINK_STRENGTH_FACTOR
        if (link.isDir) f = SETTINGS.DIR_LINK_STRENGTH_FACTOR
        if (link.isPackage) f = SETTINGS.PACKAGE_LINK_STRENGTH_FACTOR
        return f / Math.min(nodes[link.from].neighbors?.length ?? 1, nodes[link.to].neighbors?.length ?? 1);
      })

    g.d3Force('charge')
      ?.strength((node: XNode) => {
        let f = SETTINGS.FILE_NODE_REPULSION_FORCE
        if (node.isDir) f = SETTINGS.DIR_NODE_REPULSION_FORCE
        if (node.isPackage) f = SETTINGS.PACKAGE_NODE_REPULSION_FORCE
        return -f
      })

    setTimeout(() => g.zoomToFit(SETTINGS.DEFAULT_DISTANCE), 1000)
  }, [])

  return (
    <ForceGraph
      ref={graph}
      extraRenderers={[new CSS2DRenderer()]}
      graphData={data}
      backgroundColor={'#000003'}
      nodeResolution={SETTINGS.NODE_RESOLUTION}
      onBackgroundClick={backgroundClick}
      nodeLabel={({ fileName, dirName, group, loc }) => selectedNode ? '' : `
        <div class="nodeLabel">
            <span>${dirName}<span style="font-weight: bold">${fileName}</span></span>
            <span>${group != null ? `package: ${group}` : ''}</span>
            <span>LOC: ${loc}</span>
        </div>`}
      nodeThreeObject={nodeThreeObject}
      nodeThreeObjectExtend={true}
      nodeVal={'size' satisfies keyof XNode}
      nodeVisibility={node => !node.isDir && !node.isPackage}
      nodeColor={colorNode}
      nodeOpacity={1}
      onNodeClick={node => {
        selectNode(node)
        centerOnNode(node)
      }}
      linkDirectionalArrowLength={4}
      linkDirectionalArrowRelPos={1}
      linkColor={colorLink}
      linkDirectionalArrowColor={colorLink}
      linkSource={'from' satisfies keyof XLink}
      linkTarget={'to' satisfies keyof XLink}
      linkVisibility={link => !link.isDir && !link.isPackage}
      linkWidth={link => highlightLinks.has(link) ? SETTINGS.LINK_HIGHLIGHT_WIDTH : 0.25}
      linkDirectionalParticles={link => highlightLinks.has(link) ? 2 : 0}
      linkDirectionalParticleWidth={SETTINGS.LINK_HIGHLIGHT_WIDTH}
    />
  )
}

export default App
