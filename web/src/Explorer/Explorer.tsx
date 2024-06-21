import { HTMLProps, useEffect, useMemo } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faFile, faFolder } from "@fortawesome/free-solid-svg-icons";
import { faGolang, faJs, faPython, faRust, IconDefinition } from "@fortawesome/free-brands-svg-icons";

import { FileTree } from "../FileTree.ts";
import { XNode } from "../XGraph.ts";
import { FolderState } from "./FolderState.ts";
import { useForceUpdate } from "../@utils/useForceUpdate.ts";

const ID_PREFIX = '__explorer_'
const SELECTED_TAG = 's'
const HIGHLIGHTED_TAG = 'h'
const IN_TAG_VALUE = 'in'
const OUT_TAG_VALUE = 'out'
const DIR_SELECTED_COLOR = '#ffffff11'
const FILE_SELECTED_COLOR = 'rgba(31,194,218,0.34)'
const FILE_IN_HIGHLIGHTED_COLOR = 'rgba(31,218,47,0.16)'
const FILE_OUT_HIGHLIGHTED_COLOR = 'rgba(218,187,31,0.16)'

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
    folderState.expand()

    return folderState
  }, [fileTree])

  useEffect(() => {
    folderState.untagAll(SELECTED_TAG)
    if (selected !== undefined) {
      const names = FileTree.parentFolders(selected)
      names.shift()
      folderState.expandRecursive(names)
      names.push(selected.fileName)
      folderState.tagRecursive(names, SELECTED_TAG, 'true')
      setTimeout(
        () => document.getElementById(ID_PREFIX + selected.id.toString())
          ?.scrollIntoView({ behavior: "smooth", block: 'nearest' }),
        50
      )
    }

    folderState.untagAll(HIGHLIGHTED_TAG)
    if (highlighted !== undefined && highlighted.size > 0) {
      const outLinks = selected?.links?.reduce((acc, link) => (acc.add(link.to)), new Set<number>())
      for (const node of highlighted?.values() ?? []) {
        const names = FileTree.parentFolders(node)
        names.shift()
        folderState.expandRecursive(names)
        names.push(node.fileName)
        folderState.tagRecursive(names, HIGHLIGHTED_TAG, outLinks?.has(node.id) ? IN_TAG_VALUE : OUT_TAG_VALUE)
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

function ExplorerFolder (props: ExplorerFolderProps & HTMLProps<HTMLDivElement>) {
  const { folderState, onSelectNode } = props
  const forceUpdate = useForceUpdate()

  useEffect(() => folderState.registerListener('update', forceUpdate), [folderState, forceUpdate]);

  if (folderState.isCollapsed) {
    return <Folder
      name={folderState.name}
      style={{
        color: folderState.color,
        backgroundColor: folderState.tags[SELECTED_TAG] ? DIR_SELECTED_COLOR : undefined,
        ...props.style
      }}
      onClick={() => folderState.expand()} dir={props.dir}
    />
  }

  return <div className="flex flex-col" {...props}>
    <Folder
      name={folderState.name}
      style={{
        color: folderState.color,
        backgroundColor: folderState.tags[SELECTED_TAG] ? DIR_SELECTED_COLOR : undefined,
      }}
      onClick={() => folderState.collapseAll()}
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
        style={{
          color: folderState.color,
          marginLeft: 16,
          backgroundColor:
            folderState.fileTags[file.fileName]?.[SELECTED_TAG]
              ? FILE_SELECTED_COLOR
              : folderState.fileTags[file.fileName]?.[HIGHLIGHTED_TAG] === IN_TAG_VALUE
                ? FILE_IN_HIGHLIGHTED_COLOR
                : folderState.fileTags[file.fileName]?.[HIGHLIGHTED_TAG] === OUT_TAG_VALUE
                  ? FILE_OUT_HIGHLIGHTED_COLOR
                  : undefined
        }}
        onClick={() => onSelectNode?.(file)}
      />
    )}
  </div>
}


function Folder ({ name, ...rest }: { name: string } & HTMLProps<HTMLDivElement>) {
  return <div className={'flex flex-row items-center cursor-pointer'} {...rest}>
    <FontAwesomeIcon icon={faFolder} color={rest.style?.color}/>
    <span className={'text-white ml-2'}>{name}</span>
  </div>
}

function File ({ name, ...rest }: { name: string } & HTMLProps<HTMLDivElement>) {
  const ext = name.split('.').slice(-1)[0]
  return <div className='flex flex-row items-center cursor-pointer' {...rest}>
    <FontAwesomeIcon icon={FA_MAP[ext] ?? faFile} color={rest.style?.color}/>
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

