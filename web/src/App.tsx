import { useEffect, useRef, useState } from "react";
import ForceGraph, { ForceGraphMethods, LinkObject, NodeObject } from "react-force-graph-3d";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-expect-error
import { UnrealBloomPass } from 'three/examples/jsm/postprocessing/UnrealBloomPass.js'
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-expect-error
import { CSS2DObject, CSS2DRenderer } from 'three/examples/jsm/renderers/CSS2DRenderer.js'
import { Data } from "./data.ts";
import { useData, XLink, XNode } from "./useData.ts";
import './App.css'
import { Leva, useControls } from "leva";
import { HSVtoRGB } from "./@utils/HSVtoRGB.ts";


const DEFAULT_SETTINGS = {
  DEFAULT_DISTANCE: 400,
  NODE_RESOLUTION: 16,
  LINK_WIDTH: 0.5,
  LINK_HIGHLIGHT_WIDTH: 2,
  BLOOM_PASS_STRENGTH: 1,
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
  const SETTINGS = useControls(DEFAULT_SETTINGS)

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
    let alpha = SETTINGS.NODE_ALPHA
    if (highlightNodes.size > 0 && !highlightNodes.has(node)) alpha = SETTINGS.UNSELECTED_NODE_ALPHA
    const { h, s, v } = (node.isEntrypoint || node.__color === undefined) ? { h: 0, s: 0, v: 1 } : node.__color
    const [ r, g, b ] = HSVtoRGB(h, s, v)
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

    // A node was selected. Highlight it and its neighbors.
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

    graph.current?.cameraPosition({ x: x * distRatio, y: y * distRatio, z: z * distRatio }, { x, y, z }, 1000)
  }

  useEffect(() => {
    graph.current?.postProcessingComposer().reset()
    graph.current?.postProcessingComposer().addPass(new UnrealBloomPass(
      undefined, // resolution
      SETTINGS.BLOOM_PASS_STRENGTH, // strength
      SETTINGS.BLOOM_PASS_RADIUS, // radius
      SETTINGS.BLOOM_PASS_THRESHOLD // threshold
    ))
  }, [SETTINGS.BLOOM_PASS_RADIUS, SETTINGS.BLOOM_PASS_STRENGTH, SETTINGS.BLOOM_PASS_THRESHOLD])

  useEffect(() => {
    graph.current?.d3Force('link')
      ?.distance(() => SETTINGS.LINK_DISTANCE)
      .strength((link: XLink) => {
        let f = SETTINGS.FILE_LINK_STRENGTH_FACTOR
        if (link.isDir) f = SETTINGS.DIR_LINK_STRENGTH_FACTOR
        if (link.isPackage) f = SETTINGS.PACKAGE_LINK_STRENGTH_FACTOR
        return f / Math.min(nodes[link.from].neighbors?.length ?? 1, nodes[link.to].neighbors?.length ?? 1);
      })
    graph.current?.d3ReheatSimulation()
  }, [SETTINGS.DIR_LINK_STRENGTH_FACTOR, SETTINGS.FILE_LINK_STRENGTH_FACTOR, SETTINGS.LINK_DISTANCE, SETTINGS.PACKAGE_LINK_STRENGTH_FACTOR, nodes])

  useEffect(() => {
    graph.current?.d3Force('charge')
      ?.strength((node: XNode) => {
        let f = SETTINGS.FILE_NODE_REPULSION_FORCE
        if (node.isDir) f = SETTINGS.DIR_NODE_REPULSION_FORCE
        if (node.isPackage) f = SETTINGS.PACKAGE_NODE_REPULSION_FORCE
        return -f
      })
    graph.current?.d3ReheatSimulation()
  }, [SETTINGS.DIR_NODE_REPULSION_FORCE, SETTINGS.FILE_NODE_REPULSION_FORCE, SETTINGS.PACKAGE_NODE_REPULSION_FORCE])

  useEffect(() => {
    setTimeout(() => graph.current?.zoomToFit(SETTINGS.DEFAULT_DISTANCE), 1000)
  }, [SETTINGS.DEFAULT_DISTANCE]);

  return (
    <>
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
        linkWidth={link => highlightLinks.has(link) ? SETTINGS.LINK_HIGHLIGHT_WIDTH : SETTINGS.LINK_WIDTH}
        linkDirectionalParticles={link => highlightLinks.has(link) ? 2 : 0}
        linkDirectionalParticleWidth={SETTINGS.LINK_HIGHLIGHT_WIDTH}
      />
      <Leva hidden={!data.enableGui}/>
    </>
  )
}

export default App

