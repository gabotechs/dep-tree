import { CSSProperties, HTMLProps, useEffect, useMemo } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faFile, faFolder } from "@fortawesome/free-solid-svg-icons";
import { faGolang, faJs, faPython, faRust, IconDefinition } from "@fortawesome/free-brands-svg-icons";

import { FileTree } from "../FileTree.ts";
import { XNode } from "../XGraph.ts";
import { FolderState } from "./FolderState.ts";
import { useForceUpdate } from "../@utils/useForceUpdate.ts";

const ID_PREFIX = '__explorer_'

enum TAGS {
  SELECTED = 's',
  HIGHLIGHTED = 'h',
  EXPANDED = 'x'
}

enum VALUES {
  IN = 'in',
  OUT = 'Out',
  YES = 'y'
}

enum COLORS {
  DIR_SELECTED = '#ffffff11',
  FILE_SELECTED = 'rgba(31,194,218,0.34)',
  FILE_IN_HIGHLIGHTED = 'rgba(31,218,47,0.16)',
  FILE_OUT_HIGHLIGHTED = 'rgba(218,187,31,0.16)'
}

export interface ExplorerProps {
  className?: string;
  fileTree: FileTree<XNode>
  selected?: XNode
  highlighted?: Set<XNode>
  onSelectNode?: (x: XNode) => void
}

export function Explorer (
  {
    className = '',
    fileTree,
    onSelectNode,
    highlighted,
    selected
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
      for (const node of highlighted?.values() ?? []) {
        const names = FileTree.parentFolders(node)
        names.shift()
        folderState.tagRecursive(names, TAGS.EXPANDED, VALUES.YES)
        names.push(node.fileName)
        folderState.tagRecursive(names, TAGS.HIGHLIGHTED, outLinks?.has(node.id) ? VALUES.IN : VALUES.OUT)
      }
    }
  }, [folderState, highlighted, selected])

  return (
    <div
      className={`${className} flex flex-col overflow-y-scroll pb-8 pt-1 scrollbar-thin scrollbar-transparent`}
      dir={'rtl'}
    >
      <ExplorerFolder folderState={folderState} onSelectNode={onSelectNode} dir={'ltr'}/>
    </div>
  )
}

interface ExplorerFolderProps {
  onSelectNode?: (x: XNode) => void
  folderState: FolderState<XNode>
}

function ExplorerFolder (
  {
    folderState,
    onSelectNode,
    style,
    ...props
  }: ExplorerFolderProps & HTMLProps<HTMLDivElement>) {
  const forceUpdate = useForceUpdate()

  useEffect(() => folderState.registerListener('update', forceUpdate), [folderState, forceUpdate]);

  if (!folderState.tags[TAGS.EXPANDED]) {
    return <Folder
      name={folderState.name}
      tags={folderState.tags}
      logoColor={folderState.color}
      style={style}
      onClick={() => folderState.tag(TAGS.EXPANDED, VALUES.YES)} dir={props.dir}
      {...props}
    />
  }

  return <div className="flex flex-col" style={style} {...props}>
    <Folder
      name={folderState.name}
      logoColor={folderState.color}
      tags={folderState.tags}
      onClick={() => folderState.untagAllFolders(TAGS.EXPANDED)}
    />
    {[...folderState.folders.values()].map(folder =>
      <ExplorerFolder
        style={{ marginLeft: 16 }}
        key={folder.name}
        folderState={folder}
        onSelectNode={onSelectNode}
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
  if (tags[TAGS.SELECTED] === VALUES.YES) backgroundColor = COLORS.DIR_SELECTED
  return <div
    className={'flex flex-row items-center cursor-pointer'}
    style={{ backgroundColor, ...style }}
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
  if (tags?.[TAGS.HIGHLIGHTED] === VALUES.OUT) backgroundColor = COLORS.FILE_OUT_HIGHLIGHTED
  if (tags?.[TAGS.HIGHLIGHTED] === VALUES.IN) backgroundColor = COLORS.FILE_IN_HIGHLIGHTED
  if (tags?.[TAGS.SELECTED] === VALUES.YES) backgroundColor = COLORS.FILE_SELECTED
  return <div
    className='flex flex-row items-center cursor-pointer'
    style={{ backgroundColor, ...style }}
    {...props}
  >
    <FontAwesomeIcon icon={FA_MAP[ext] ?? faFile} color={logoColor}/>
    <span className={'text-white ml-2'}>{name}</span>
  </div>
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

