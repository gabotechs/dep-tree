import React, { CSSProperties, HTMLProps, ReactNode, useEffect, useMemo, useRef, useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faFile, faFolder } from "@fortawesome/free-solid-svg-icons";
import { faGolang, faJs, faPython, faRust, IconDefinition } from "@fortawesome/free-brands-svg-icons";

import { FileTree } from "../FileTree.ts";
import { XNode } from "../XGraph.ts";
import { FolderState } from "./FolderState.ts";
import { useForceUpdate } from "../@utils/useForceUpdate.ts";
import './Explorer.css'

const ID_PREFIX = '__explorer_'

enum TAGS {
  SELECTED = 's',
  HIGHLIGHTED = 'h',
  EXPANDED = 'x',
  IGNORE = 'i'
}

enum VALUES {
  IN = 'in',
  OUT = 'out',
  BOTH = 'both',
  YES = 'y'
}

enum COLORS {
  DIR_SELECTED = '#ffffff11',
  FILE_SELECTED = 'rgba(31,194,218,0.5)',
  FILE_IN_HIGHLIGHTED = 'rgba(31,218,47,0.16)',
  FILE_OUT_HIGHLIGHTED = 'rgba(238,255,2,0.16)',
  FILE_BOTH_HIGHLIGHTED = 'rgba(231,116,54,0.16)'
}

const FA_MAP: Record<string, IconDefinition> = {
  // python
  py: faPython,
  pyi: faPython,
  pyx: faPython,
  // golang
  go: faGolang,
  // rust
  rs: faRust,
  // js/ts
  js: faJs,
  cjs: faJs,
  mjs: faJs,
  jsx: faJs, // TODO
  ts: faJs, // TODO
  tsx: faJs, // TODO
  'd.ts': faJs, // TODO
}

export interface ExplorerProps {
  className?: string;
  fileTree: FileTree<XNode>
  selected?: XNode
  highlighted?: Set<XNode>
  onSelectNode?: (x: XNode) => void
  onNodesMutated?: () => void
}

export function Explorer (
  {
    className = '',
    fileTree,
    onSelectNode,
    highlighted,
    selected,
    onNodesMutated
  }: ExplorerProps
) {
  const folderState = useMemo(() => {
    const folderState = FolderState.fromFileTree(fileTree)
    folderState.name = folderState.name.replace(FileTree.ROOT_NAME + "/", "")
    folderState.tag(TAGS.EXPANDED, VALUES.YES)

    return folderState
  }, [fileTree])

  useEffect(() => {
    folderState.untagAll(TAGS.SELECTED)
    if (selected !== undefined) {
      const names = FileTree.parentFolders(selected)
      names.shift()
      folderState.tagRecursive(names, TAGS.EXPANDED, VALUES.YES)
      names.push(selected.fileName)
      folderState.tagRecursive(names, TAGS.SELECTED, VALUES.YES)
      setTimeout(
        () => document.getElementById(ID_PREFIX + selected.id.toString())
          ?.scrollIntoView({ behavior: "smooth", block: 'nearest' }),
        50
      )
    }

    folderState.untagAll(TAGS.HIGHLIGHTED)
    if (highlighted !== undefined && highlighted.size > 0) {
      const outLinks = selected?.links?.reduce((acc, link) => (acc.add(link.to)), new Set<number>())
      const inLinks = selected?.links?.reduce((acc, link) => (acc.add(link.from)), new Set<number>())
      for (const node of highlighted?.values() ?? []) {
        const names = FileTree.parentFolders(node)
        names.shift()
        folderState.tagRecursive(names, TAGS.EXPANDED, VALUES.YES)
        names.push(node.fileName)
        if (outLinks?.has(node.id) && inLinks?.has(node.id)) {
          folderState.tagRecursive(names, TAGS.HIGHLIGHTED, VALUES.BOTH)
        } else if (outLinks?.has(node.id)) {
          folderState.tagRecursive(names, TAGS.HIGHLIGHTED, VALUES.IN)
        } else if (inLinks?.has(node.id)) {
          folderState.tagRecursive(names, TAGS.HIGHLIGHTED, VALUES.OUT)
        }
      }
    }
  }, [folderState, highlighted, selected])

  const [contextMenuProps, setContextMenuProps] = useState<ContextMenuProps>()

  return (
    <div
      className={`${className} flex flex-col overflow-y-scroll pb-8 pt-1 scrollbar-thin scrollbar-transparent`}
      dir={'rtl'}
    >
      {contextMenuProps && <ContextMenu
        {...contextMenuProps}
        onClose={() => setContextMenuProps(undefined)}
        onNodesMutated={onNodesMutated}
      />}
      <ExplorerFolder
        folderState={folderState}
        onSelectNode={onSelectNode}
        onRightClick={setContextMenuProps}
        dir={'ltr'}
      />
    </div>
  )
}

export interface ContextMenuProps {
  x: number
  y: number
  folderState: FolderState<XNode>
  onClose?: () => void
  onNodesMutated?: () => void
}

function ContextMenu (
  {
    onClose,
    x,
    y,
    folderState,
    onNodesMutated
  }: ContextMenuProps) {
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleOutsideClick (event: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        onClose?.();
      }
    }

    document.addEventListener('mousedown', handleOutsideClick);
    return () => document.removeEventListener('mousedown', handleOutsideClick)
  }, [onClose]);

  function ignore () {
    for (const node of folderState.allFiles()) {
      node.ignore = true
      node.links?.forEach(link => (link.ignore = true))
    }
    folderState.tagAll(TAGS.IGNORE, VALUES.YES)
    onClose?.()
    onNodesMutated?.()
  }

  function unIgnore () {
    for (const node of folderState.allFiles()) {
      node.ignore = false
      node.links?.forEach(link => (link.ignore = false))
    }
    folderState.untagAll(TAGS.IGNORE)
    onClose?.()
    onNodesMutated?.()
  }

  return (
    <div
      ref={menuRef}
      className="absolute bg-white border border-gray-300 z-50 rounded shadow-md w-[100px]"
      style={{ top: y, left: x }}
    >
      <ul className="space-y-2">
        <li
          className="cursor-pointer px-2 py-1 hover:bg-gray-200 text-center"
          onClick={folderState.tags[TAGS.IGNORE] ? unIgnore : ignore}
        >
          {folderState.tags[TAGS.IGNORE] ? 'un-ignore' : 'ignore'}
        </li>
      </ul>
    </div>
  );
}

interface ExplorerFolderProps {
  onSelectNode?: (x: XNode) => void
  onRightClick?: (props: ContextMenuProps) => void
  folderState: FolderState<XNode>
}

function ExplorerFolder (
  {
    folderState,
    onSelectNode,
    onRightClick,
    style,
    ...props
  }: ExplorerFolderProps & HTMLProps<HTMLDivElement>) {
  const [, forceUpdate] = useForceUpdate()

  useEffect(() => folderState.registerListener('update', forceUpdate), [folderState, forceUpdate]);

  function onContextMenu (e: React.MouseEvent<HTMLDivElement>) {
    e.preventDefault()
    onRightClick?.({
      folderState,
      x: e.clientX,
      y: e.clientY,
    })
  }

  if (!folderState.tags[TAGS.EXPANDED]) {
    return <Folder
      name={folderState.name}
      tags={folderState.tags}
      logoColor={folderState.color}
      style={style}
      onClick={() => folderState.tag(TAGS.EXPANDED, VALUES.YES)} dir={props.dir}
      onContextMenu={onContextMenu}
      {...props}
    />
  }

  return <div className="flex flex-col" style={style} {...props}>
    <Folder
      name={folderState.name}
      logoColor={folderState.color}
      tags={folderState.tags}
      onClick={() => folderState.untagAllFolders(TAGS.EXPANDED)}
      onContextMenu={onContextMenu}
    />
    {[...folderState.folders.values()].map(folder =>
      <ExplorerFolder
        style={{ marginLeft: 16 }}
        key={folder.name}
        folderState={folder}
        onSelectNode={onSelectNode}
        onRightClick={onRightClick}
      />
    )}
    {[...folderState.files.values()].map(file =>
      <File
        id={`${ID_PREFIX}${file.id}`}
        key={file.id}
        name={file.fileName}
        logoColor={folderState.color}
        tags={folderState.fileTags[file.fileName]}
        style={{ marginLeft: 16 }}
        onClick={() => onSelectNode?.(file)}
      />
    )}
  </div>
}


function Folder (
  {
    name,
    logoColor,
    tags,
    style,
    ...props
  }: {
    name: string
    logoColor: string,
    tags: Record<string, string>
  } & HTMLProps<HTMLDivElement>) {
  let backgroundColor: CSSProperties['color'] = undefined
  let opacity = 1
  if (tags[TAGS.SELECTED] === VALUES.YES) backgroundColor = COLORS.DIR_SELECTED
  if (tags[TAGS.IGNORE] === VALUES.YES) opacity = 0.2
  return <div
    className={'flex flex-row items-center cursor-pointer'}
    style={{ backgroundColor, opacity, ...style }}
    {...props}
  >
    <FontAwesomeIcon icon={faFolder} color={logoColor}/>
    <span className={'text-white ml-2'}>{name}</span>
  </div>
}

function File (
  {
    name,
    logoColor,
    style,
    tags,
    ...props
  }: {
    name: string,
    logoColor: string,
    tags?: Record<string, string>
  } & HTMLProps<HTMLDivElement>) {
  const ext = useMemo(() => name.split('.').slice(-1)[0], [name])
  let backgroundColor: CSSProperties['color'] = undefined
  let animation: ReactNode = null
  let opacity = 1
  if (tags?.[TAGS.HIGHLIGHTED] === VALUES.OUT) {
    backgroundColor = COLORS.FILE_OUT_HIGHLIGHTED
    animation = <AnimatedDot dir={'ltr'}/>
  }
  if (tags?.[TAGS.HIGHLIGHTED] === VALUES.IN) {
    backgroundColor = COLORS.FILE_IN_HIGHLIGHTED
    animation = <AnimatedDot dir={'rtl'}/>
  }
  if (tags?.[TAGS.HIGHLIGHTED] === VALUES.BOTH) {
    backgroundColor = COLORS.FILE_BOTH_HIGHLIGHTED
    animation = <AnimatedDot dir={'both'}/>
  }
  if (tags?.[TAGS.SELECTED] === VALUES.YES) {
    backgroundColor = COLORS.FILE_SELECTED
    animation = null
  }
  if (tags?.[TAGS.IGNORE] === VALUES.YES) {
    opacity = 0.2
  }
  return <div
    className='flex flex-row items-center cursor-pointer'
    style={{ backgroundColor, opacity, ...style }}
    {...props}
  >
    <FontAwesomeIcon icon={FA_MAP[ext] ?? faFile} color={logoColor}/>
    <span className={'text-white mx-2'}>{name}</span>
    <div className={'flex-1 min-w-12'}>
      {animation}
    </div>
  </div>
}


function AnimatedDot ({ dir }: { dir: 'ltr' | 'rtl' | 'both' }) {
  return (
    <div className="relative w-full h-2 overflow-hidden">
      {dir === 'ltr'
        ? <div className={`absolute h-full w-2 bg-gray-500 rounded-full animate-left-right`}/>
        : dir === 'rtl'
          ? <div className={`absolute h-full w-2 bg-gray-500 rounded-full animate-right-left`}/>
          : <>
            <div className={`absolute h-full w-2 bg-gray-500 rounded-full animate-left-right`}/>
            <div className={`absolute h-full w-2 bg-gray-500 rounded-full animate-right-left`}/>
          </>
      }
    </div>
  )
}
