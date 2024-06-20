import { CSSProperties, HTMLProps, useMemo } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faFile, faFolder } from "@fortawesome/free-solid-svg-icons";
import { faGolang, faJs, faPython, faRust, IconDefinition } from "@fortawesome/free-brands-svg-icons";

import { FileTree } from "../FileTree.ts";
import { XNode } from "../XGraph.ts";
import { FolderState } from "./FolderState.ts";
import { useForceUpdate } from "../@utils/useForceUpdate.ts";

export interface ExplorerProps {
  className?: string;
  fileTree: FileTree<XNode>
  selected?: number
  highlighted?: Set<number>
}

export function Explorer ({ className = '', fileTree }: ExplorerProps) {
  const folderState = useMemo(() => {
    const folderState = FolderState.fromFileTree(fileTree)
    folderState.name = folderState.name.replace(FileTree.ROOT_NAME + "/", "")
    folderState.unfold()
    return folderState
  }, [fileTree])

  return <div className={`${className} flex flex-col overflow-y-auto`}>
    <ExplorerFolder folderState={folderState}/>
  </div>
}

interface ExplorerLevelProps {
  folderState: FolderState<XNode>
}

function ExplorerFolder (props: ExplorerLevelProps & HTMLProps<HTMLDivElement>) {
  const { folderState } = props
  const folderStyle: CSSProperties = { color: folderState.color, cursor: 'pointer' }
  const fileStyle: CSSProperties = { color: folderState.color, marginLeft: 16 }

  const forceUpdate = useForceUpdate()

  function unfold () {
    folderState.unfold()
    forceUpdate()
  }

  function collapse () {
    folderState.collapseAll()
    forceUpdate()
  }

  if (folderState.folded) {
    return <Folder name={folderState.name} style={{ ...folderStyle, ...props.style }} onClick={unfold}/>
  }

  return <div className="flex flex-col" {...props}>
    <Folder name={folderState.name} style={folderStyle} onClick={collapse}/>
    {[...folderState.folders.values()].map(folder =>
      <ExplorerFolder
        style={{ marginLeft: 16 }}
        key={folder.name}
        folderState={folder}
      />
    )}
    {[...folderState.files.values()].map(file =>
      <File key={file.id} name={file.fileName} style={fileStyle}/>
    )}
  </div>
}


function Folder ({ name, ...rest }: { name: string } & HTMLProps<HTMLDivElement>) {
  return <div className={'flex flex-row items-center '} {...rest}>
    <FontAwesomeIcon icon={faFolder} color={rest.style?.color}/>
    <span className={'text-white ml-2'}>{name}</span>
  </div>
}

function File ({ name, ...rest }: { name: string } & HTMLProps<HTMLDivElement>) {
  const ext = name.split('.').slice(-1)[0]
  return <div className='flex flex-row items-center' {...rest}>
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

