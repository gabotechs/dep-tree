// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-expect-error
import { UnrealBloomPass } from 'three/examples/jsm/postprocessing/UnrealBloomPass.js'
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-expect-error
import { CSS2DObject, CSS2DRenderer } from 'three/examples/jsm/renderers/CSS2DRenderer.js'
import ForceGraph, { ForceGraphMethods, LinkObject, NodeObject } from "react-force-graph-3d";
import { useEffect, useRef, useState } from "react";
import { Leva, useControls } from "leva";

import { buildXGraph, XLink, XNode } from "./XGraph.ts";
import { HSVtoRGB } from "./@utils/HSVtoRGB.ts";
import { Graph } from "./types.ts";
import { Explorer } from "./Explorer/Explorer.tsx";

import './App.css'


class Data {
  static __INLINE_DATA = {} as Graph
}

// Uncomment this line for using the test data.
Data.__INLINE_DATA = { "nodes": [{ "id": 1334025261, "isEntrypoint": true, "fileName": "main.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "main.go"], "group": "main", "dirName": "./", "loc": 14, "size": 0, }, { "id": 3232333941, "isEntrypoint": false, "fileName": "root.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "cmd", "root.go"], "group": "cmd", "dirName": "cmd/", "loc": 247, "size": 10, }, { "id": 3022277916, "isEntrypoint": false, "fileName": "entropy.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "cmd", "entropy.go"], "group": "cmd", "dirName": "cmd/", "loc": 52, "size": 2, }, { "id": 1432468757, "isEntrypoint": false, "fileName": "tree.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "cmd", "tree.go"], "group": "cmd", "dirName": "cmd/", "loc": 69, "size": 2, }, { "id": 2893840283, "isEntrypoint": false, "fileName": "check.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "cmd", "check.go"], "group": "cmd", "dirName": "cmd/", "loc": 50, "size": 2, }, { "id": 806026951, "isEntrypoint": false, "fileName": "config.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "cmd", "config.go"], "group": "cmd", "dirName": "cmd/", "loc": 29, "size": 1, }, { "id": 29475502, "isEntrypoint": false, "fileName": "explain.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "cmd", "explain.go"], "group": "cmd", "dirName": "cmd/", "loc": 107, "size": 4, }, { "id": 679976714, "isEntrypoint": false, "fileName": "inarray.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "utils", "inarray.go"], "group": "utils", "dirName": "internal/utils/", "loc": 10, "size": 0, }, { "id": 2753948251, "isEntrypoint": false, "fileName": "config.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "config", "config.go"], "group": "config", "dirName": "internal/config/", "loc": 100, "size": 4, }, { "id": 1362904939, "isEntrypoint": false, "fileName": "language.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "language", "language.go"], "group": "language", "dirName": "internal/language/", "loc": 124, "size": 5, }, { "id": 648007272, "isEntrypoint": false, "fileName": "endswith.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "utils", "endswith.go"], "group": "utils", "dirName": "internal/utils/", "loc": 12, "size": 0, }, { "id": 3917048986, "isEntrypoint": false, "fileName": "language.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "js", "language.go"], "group": "js", "dirName": "internal/js/", "loc": 57, "size": 2, }, { "id": 1650987103, "isEntrypoint": false, "fileName": "language.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "rust", "language.go"], "group": "rust", "dirName": "internal/rust/", "loc": 33, "size": 1, }, { "id": 3967137051, "isEntrypoint": false, "fileName": "language.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "python", "language.go"], "group": "python", "dirName": "internal/python/", "loc": 51, "size": 2, }, { "id": 212004047, "isEntrypoint": false, "fileName": "language.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "go", "language.go"], "group": "golang", "dirName": "internal/go/", "loc": 73, "size": 2, }, { "id": 1895523041, "isEntrypoint": false, "fileName": "language.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "dummy", "language.go"], "group": "dummy", "dirName": "internal/dummy/", "loc": 63, "size": 2, }, { "id": 491731072, "isEntrypoint": false, "fileName": "file.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "utils", "file.go"], "group": "utils", "dirName": "internal/utils/", "loc": 23, "size": 0, }, { "id": 3177744305, "isEntrypoint": false, "fileName": "match.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "utils", "match.go"], "group": "utils", "dirName": "internal/utils/", "loc": 9, "size": 0, }, { "id": 3235280883, "isEntrypoint": false, "fileName": "node.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "graph", "node.go"], "group": "graph", "dirName": "internal/graph/", "loc": 78, "size": 3, }, { "id": 3672289200, "isEntrypoint": false, "fileName": "parser.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "language", "parser.go"], "group": "language", "dirName": "internal/language/", "loc": 142, "size": 5, }, { "id": 2141569406, "isEntrypoint": false, "fileName": "render.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "entropy", "render.go"], "group": "entropy", "dirName": "internal/entropy/", "loc": 56, "size": 2, }, { "id": 1331865351, "isEntrypoint": false, "fileName": "load.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "graph", "load.go"], "group": "graph", "dirName": "internal/graph/", "loc": 137, "size": 5, }, { "id": 356653727, "isEntrypoint": false, "fileName": "tree.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "tree", "tree.go"], "group": "tree", "dirName": "internal/tree/", "loc": 77, "size": 3, }, { "id": 341567957, "isEntrypoint": false, "fileName": "tui.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "tui", "tui.go"], "group": "tui", "dirName": "internal/tui/", "loc": 104, "size": 4, }, { "id": 1005564345, "isEntrypoint": false, "fileName": "check.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "check", "check.go"], "group": "check", "dirName": "internal/check/", "loc": 123, "size": 4, }, { "id": 182641601, "isEntrypoint": false, "fileName": "explain.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "explain", "explain.go"], "group": "explain", "dirName": "internal/explain/", "loc": 38, "size": 1, }, { "id": 801525785, "isEntrypoint": false, "fileName": "config.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "check", "config.go"], "group": "check", "dirName": "internal/check/", "loc": 35, "size": 1, }, { "id": 464063646, "isEntrypoint": false, "fileName": "config.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "js", "config.go"], "group": "js", "dirName": "internal/js/", "loc": 6, "size": 0, }, { "id": 3103188379, "isEntrypoint": false, "fileName": "config.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "rust", "config.go"], "group": "rust", "dirName": "internal/rust/", "loc": 3, "size": 0, }, { "id": 118294663, "isEntrypoint": false, "fileName": "config.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "python", "config.go"], "group": "python", "dirName": "internal/python/", "loc": 8, "size": 0, }, { "id": 1335279723, "isEntrypoint": false, "fileName": "config.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "go", "config.go"], "group": "golang", "dirName": "internal/go/", "loc": 3, "size": 0, }, { "id": 726487653, "isEntrypoint": false, "fileName": "package_json.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "js", "package_json.go"], "group": "js", "dirName": "internal/js/", "loc": 87, "size": 3, }, { "id": 1428314197, "isEntrypoint": false, "fileName": "grammar.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "js", "js_grammar", "grammar.go"], "group": "js_grammar", "dirName": "internal/js/js_grammar/", "loc": 67, "size": 2, }, { "id": 1560530462, "isEntrypoint": false, "fileName": "mod_tree.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "rust", "mod_tree.go"], "group": "rust", "dirName": "internal/rust/", "loc": 96, "size": 3, }, { "id": 2680775406, "isEntrypoint": false, "fileName": "cargo_toml.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "rust", "cargo_toml.go"], "group": "rust", "dirName": "internal/rust/", "loc": 123, "size": 4, }, { "id": 4099876531, "isEntrypoint": false, "fileName": "resolve.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "python", "resolve.go"], "group": "python", "dirName": "internal/python/", "loc": 157, "size": 6, }, { "id": 107011351, "isEntrypoint": false, "fileName": "grammar.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "python", "python_grammar", "grammar.go"], "group": "python_grammar", "dirName": "internal/python/python_grammar/", "loc": 65, "size": 2, }, { "id": 4112388782, "isEntrypoint": false, "fileName": "go_mod.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "go", "go_mod.go"], "group": "golang", "dirName": "internal/go/", "loc": 28, "size": 1, }, { "id": 2120775579, "isEntrypoint": false, "fileName": "package.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "go", "package.go"], "group": "golang", "dirName": "internal/go/", "loc": 91, "size": 3, }, { "id": 582030790, "isEntrypoint": false, "fileName": "find_closest_dir_with_root_file.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "utils", "find_closest_dir_with_root_file.go"], "group": "utils", "dirName": "internal/utils/", "loc": 33, "size": 1, }, { "id": 387927870, "isEntrypoint": false, "fileName": "parser.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "dummy", "parser.go"], "group": "dummy", "dirName": "internal/dummy/", "loc": 39, "size": 1, }, { "id": 4037712688, "isEntrypoint": false, "fileName": "cached.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "utils", "cached.go"], "group": "utils", "dirName": "internal/utils/", "loc": 80, "size": 3, }, { "id": 3222993554, "isEntrypoint": false, "fileName": "exports.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "language", "exports.go"], "group": "language", "dirName": "internal/language/", "loc": 111, "size": 4, }, { "id": 2428266738, "isEntrypoint": false, "fileName": "graph.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "entropy", "graph.go"], "group": "entropy", "dirName": "internal/entropy/", "loc": 161, "size": 6, }, { "id": 2958756130, "isEntrypoint": false, "fileName": "open.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "entropy", "open.go"], "group": "entropy", "dirName": "internal/entropy/", "loc": 20, "size": 0, }, { "id": 1216276501, "isEntrypoint": false, "fileName": "graph.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "graph", "graph.go"], "group": "graph", "dirName": "internal/graph/", "loc": 201, "size": 8, }, { "id": 1942096762, "isEntrypoint": false, "fileName": "cycles.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "graph", "cycles.go"], "group": "graph", "dirName": "internal/graph/", "loc": 91, "size": 3, }, { "id": 2485026308, "isEntrypoint": false, "fileName": "render.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "tui", "systems", "render.go"], "group": "systems", "dirName": "internal/tui/systems/", "loc": 157, "size": 6, }, { "id": 2562016124, "isEntrypoint": false, "fileName": "spatial.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "tui", "systems", "spatial.go"], "group": "systems", "dirName": "internal/tui/systems/", "loc": 73, "size": 2, }, { "id": 552151127, "isEntrypoint": false, "fileName": "vector.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "utils", "vector.go"], "group": "utils", "dirName": "internal/utils/", "loc": 20, "size": 0, }, { "id": 1339565699, "isEntrypoint": false, "fileName": "state.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "tui", "systems", "state.go"], "group": "systems", "dirName": "internal/tui/systems/", "loc": 16, "size": 0, }, { "id": 2241854762, "isEntrypoint": false, "fileName": "world.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "ecs", "world.go"], "group": "ecs", "dirName": "internal/ecs/", "loc": 37, "size": 1, }, { "id": 2567328503, "isEntrypoint": false, "fileName": "entity.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "ecs", "entity.go"], "group": "ecs", "dirName": "internal/ecs/", "loc": 31, "size": 1, }, { "id": 1094774888, "isEntrypoint": false, "fileName": "runtime.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "tui", "systems", "runtime.go"], "group": "systems", "dirName": "internal/tui/systems/", "loc": 94, "size": 3, }, { "id": 1440380694, "isEntrypoint": false, "fileName": "set.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "utils", "set.go"], "group": "utils", "dirName": "internal/utils/", "loc": 16, "size": 0, }, { "id": 2048363028, "isEntrypoint": false, "fileName": "resolve.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "js", "resolve.go"], "group": "js", "dirName": "internal/js/", "loc": 129, "size": 5, }, { "id": 3166841825, "isEntrypoint": false, "fileName": "import.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "js", "js_grammar", "import.go"], "group": "js_grammar", "dirName": "internal/js/js_grammar/", "loc": 35, "size": 1, }, { "id": 3431066080, "isEntrypoint": false, "fileName": "export.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "js", "js_grammar", "export.go"], "group": "js_grammar", "dirName": "internal/js/js_grammar/", "loc": 30, "size": 1, }, { "id": 2279247337, "isEntrypoint": false, "fileName": "unquote.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "utils", "unquote.go"], "group": "utils", "dirName": "internal/utils/", "loc": 40, "size": 1, }, { "id": 4065462691, "isEntrypoint": false, "fileName": "grammar.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "rust", "rust_grammar", "grammar.go"], "group": "rust_grammar", "dirName": "internal/rust/rust_grammar/", "loc": 69, "size": 2, }, { "id": 2476273619, "isEntrypoint": false, "fileName": "import.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "python", "python_grammar", "import.go"], "group": "python_grammar", "dirName": "internal/python/python_grammar/", "loc": 21, "size": 0, }, { "id": 2311248394, "isEntrypoint": false, "fileName": "export.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "python", "python_grammar", "export.go"], "group": "python_grammar", "dirName": "internal/python/python_grammar/", "loc": 22, "size": 0, }, { "id": 1720170331, "isEntrypoint": false, "fileName": "start_lexer.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "python", "python_grammar", "start_lexer.go"], "group": "python_grammar", "dirName": "internal/python/python_grammar/", "loc": 66, "size": 2, }, { "id": 1965161762, "isEntrypoint": false, "fileName": "callstack.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "utils", "callstack.go"], "group": "utils", "dirName": "internal/utils/", "loc": 57, "size": 2, }, { "id": 1375713270, "isEntrypoint": false, "fileName": "dirs.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "entropy", "dirs.go"], "group": "entropy", "dirName": "internal/entropy/", "loc": 195, "size": 7, }, { "id": 1759930456, "isEntrypoint": false, "fileName": "max.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "utils", "max.go"], "group": "utils", "dirName": "internal/utils/", "loc": 12, "size": 0, }, { "id": 2376376984, "isEntrypoint": false, "fileName": "edge.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "graph", "edge.go"], "group": "graph", "dirName": "internal/graph/", "loc": 25, "size": 1, }, { "id": 1076436185, "isEntrypoint": false, "fileName": "stack.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "board", "graphics", "stack.go"], "group": "graphics", "dirName": "internal/board/graphics/", "loc": 112, "size": 4, }, { "id": 944021367, "isEntrypoint": false, "fileName": "render.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "tree", "render.go"], "group": "tree", "dirName": "internal/tree/", "loc": 91, "size": 3, }, { "id": 3891966039, "isEntrypoint": false, "fileName": "clamp.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "utils", "clamp.go"], "group": "utils", "dirName": "internal/utils/", "loc": 12, "size": 0, }, { "id": 1089468033, "isEntrypoint": false, "fileName": "system.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "ecs", "system.go"], "group": "ecs", "dirName": "internal/ecs/", "loc": 43, "size": 1, }, { "id": 224600506, "isEntrypoint": false, "fileName": "workspaces.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "js", "workspaces.go"], "group": "js", "dirName": "internal/js/", "loc": 130, "size": 5, }, { "id": 2149662377, "isEntrypoint": false, "fileName": "tsconfig.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "js", "tsconfig.go"], "group": "js", "dirName": "internal/js/", "loc": 71, "size": 2, }, { "id": 21548004, "isEntrypoint": false, "fileName": "mod.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "rust", "rust_grammar", "mod.go"], "group": "rust_grammar", "dirName": "internal/rust/rust_grammar/", "loc": 23, "size": 0, }, { "id": 2714238133, "isEntrypoint": false, "fileName": "use.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "rust", "rust_grammar", "use.go"], "group": "rust_grammar", "dirName": "internal/rust/rust_grammar/", "loc": 61, "size": 2, }, { "id": 409742873, "isEntrypoint": false, "fileName": "pub.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "rust", "rust_grammar", "pub.go"], "group": "rust_grammar", "dirName": "internal/rust/rust_grammar/", "loc": 6, "size": 0, }, { "id": 4178430438, "isEntrypoint": false, "fileName": "append_front.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "utils", "append_front.go"], "group": "utils", "dirName": "internal/utils/", "loc": 10, "size": 0, }, { "id": 1372522740, "isEntrypoint": false, "fileName": "scale.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "utils", "scale.go"], "group": "utils", "dirName": "internal/utils/", "loc": 11, "size": 0, }, { "id": 971943911, "isEntrypoint": false, "fileName": "cell.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "board", "graphics", "cell.go"], "group": "graphics", "dirName": "internal/board/graphics/", "loc": 102, "size": 4, }, { "id": 3596385220, "isEntrypoint": false, "fileName": "lines.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "board", "graphics", "lines.go"], "group": "graphics", "dirName": "internal/board/graphics/", "loc": 85, "size": 3, }, { "id": 2548718478, "isEntrypoint": false, "fileName": "merge.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "utils", "merge.go"], "group": "utils", "dirName": "internal/utils/", "loc": 13, "size": 0, }, { "id": 3982654405, "isEntrypoint": false, "fileName": "board.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "board", "board.go"], "group": "board", "dirName": "internal/board/", "loc": 65, "size": 2, }, { "id": 2839729826, "isEntrypoint": false, "fileName": "block.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "board", "block.go"], "group": "board", "dirName": "internal/board/", "loc": 68, "size": 2, }, { "id": 3997750360, "isEntrypoint": false, "fileName": "connector.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "board", "connector.go"], "group": "board", "dirName": "internal/board/", "loc": 144, "size": 5, }, { "id": 23089748, "isEntrypoint": false, "fileName": "matrix.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "board", "graphics", "matrix.go"], "group": "graphics", "dirName": "internal/board/graphics/", "loc": 110, "size": 4, }, { "id": 1618455164, "isEntrypoint": false, "fileName": "prefixn.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "utils", "prefixn.go"], "group": "utils", "dirName": "internal/utils/", "loc": 10, "size": 0, }, { "id": 385340190, "isEntrypoint": false, "fileName": "trace.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "board", "graphics", "trace.go"], "group": "graphics", "dirName": "internal/board/graphics/", "loc": 93, "size": 3, }, { "id": 1284823955, "isEntrypoint": false, "fileName": "bool2int.go", "pathBuf": ["", "Users", "gabriel", "dep-tree", "dep-tree", "internal", "utils", "bool2int.go"], "group": "utils", "dirName": "internal/utils/", "loc": 9, "size": 0, }], "links": [{ "from": 1334025261, "to": 3232333941, "isCyclic": false }, { "from": 3232333941, "to": 3022277916, "isCyclic": false }, { "from": 3232333941, "to": 1432468757, "isCyclic": false }, { "from": 3232333941, "to": 2893840283, "isCyclic": false }, { "from": 3232333941, "to": 806026951, "isCyclic": false }, { "from": 3232333941, "to": 29475502, "isCyclic": false }, { "from": 3232333941, "to": 679976714, "isCyclic": false }, { "from": 3232333941, "to": 2753948251, "isCyclic": false }, { "from": 3232333941, "to": 1362904939, "isCyclic": false }, { "from": 3232333941, "to": 648007272, "isCyclic": false }, { "from": 3232333941, "to": 3917048986, "isCyclic": false }, { "from": 3232333941, "to": 1650987103, "isCyclic": false }, { "from": 3232333941, "to": 3967137051, "isCyclic": false }, { "from": 3232333941, "to": 212004047, "isCyclic": false }, { "from": 3232333941, "to": 1895523041, "isCyclic": false }, { "from": 3232333941, "to": 491731072, "isCyclic": false }, { "from": 3232333941, "to": 3177744305, "isCyclic": false }, { "from": 3232333941, "to": 3235280883, "isCyclic": false }, { "from": 3022277916, "to": 3672289200, "isCyclic": false }, { "from": 3022277916, "to": 2141569406, "isCyclic": false }, { "from": 3022277916, "to": 1331865351, "isCyclic": false }, { "from": 3022277916, "to": 1362904939, "isCyclic": false }, { "from": 1432468757, "to": 3672289200, "isCyclic": false }, { "from": 1432468757, "to": 356653727, "isCyclic": false }, { "from": 1432468757, "to": 1362904939, "isCyclic": false }, { "from": 1432468757, "to": 1331865351, "isCyclic": false }, { "from": 1432468757, "to": 341567957, "isCyclic": false }, { "from": 2893840283, "to": 2753948251, "isCyclic": false }, { "from": 2893840283, "to": 3672289200, "isCyclic": false }, { "from": 2893840283, "to": 1005564345, "isCyclic": false }, { "from": 2893840283, "to": 1362904939, "isCyclic": false }, { "from": 2893840283, "to": 1331865351, "isCyclic": false }, { "from": 806026951, "to": 2753948251, "isCyclic": false }, { "from": 29475502, "to": 3672289200, "isCyclic": false }, { "from": 29475502, "to": 182641601, "isCyclic": false }, { "from": 29475502, "to": 1362904939, "isCyclic": false }, { "from": 29475502, "to": 1331865351, "isCyclic": false }, { "from": 29475502, "to": 3235280883, "isCyclic": false }, { "from": 2753948251, "to": 801525785, "isCyclic": false }, { "from": 2753948251, "to": 464063646, "isCyclic": false }, { "from": 2753948251, "to": 3103188379, "isCyclic": false }, { "from": 2753948251, "to": 118294663, "isCyclic": false }, { "from": 2753948251, "to": 1335279723, "isCyclic": false }, { "from": 3917048986, "to": 464063646, "isCyclic": false }, { "from": 3917048986, "to": 726487653, "isCyclic": false }, { "from": 3917048986, "to": 1362904939, "isCyclic": false }, { "from": 3917048986, "to": 491731072, "isCyclic": false }, { "from": 3917048986, "to": 1428314197, "isCyclic": false }, { "from": 1650987103, "to": 3103188379, "isCyclic": false }, { "from": 1650987103, "to": 1560530462, "isCyclic": false }, { "from": 1650987103, "to": 2680775406, "isCyclic": false }, { "from": 1650987103, "to": 1362904939, "isCyclic": false }, { "from": 3967137051, "to": 118294663, "isCyclic": false }, { "from": 3967137051, "to": 4099876531, "isCyclic": false }, { "from": 3967137051, "to": 1362904939, "isCyclic": false }, { "from": 3967137051, "to": 107011351, "isCyclic": false }, { "from": 212004047, "to": 1335279723, "isCyclic": false }, { "from": 212004047, "to": 4112388782, "isCyclic": false }, { "from": 212004047, "to": 2120775579, "isCyclic": false }, { "from": 212004047, "to": 582030790, "isCyclic": false }, { "from": 212004047, "to": 1362904939, "isCyclic": false }, { "from": 1895523041, "to": 387927870, "isCyclic": false }, { "from": 1895523041, "to": 1362904939, "isCyclic": false }, { "from": 491731072, "to": 4037712688, "isCyclic": false }, { "from": 3235280883, "to": 4037712688, "isCyclic": false }, { "from": 3672289200, "to": 1362904939, "isCyclic": false }, { "from": 3672289200, "to": 3222993554, "isCyclic": false }, { "from": 3672289200, "to": 1331865351, "isCyclic": false }, { "from": 3672289200, "to": 3177744305, "isCyclic": false }, { "from": 3672289200, "to": 3235280883, "isCyclic": false }, { "from": 2141569406, "to": 2428266738, "isCyclic": false }, { "from": 2141569406, "to": 2958756130, "isCyclic": false }, { "from": 2141569406, "to": 1331865351, "isCyclic": false }, { "from": 2141569406, "to": 1362904939, "isCyclic": false }, { "from": 1331865351, "to": 3235280883, "isCyclic": false }, { "from": 1331865351, "to": 1216276501, "isCyclic": false }, { "from": 356653727, "to": 3235280883, "isCyclic": false }, { "from": 356653727, "to": 1216276501, "isCyclic": false }, { "from": 356653727, "to": 1331865351, "isCyclic": false }, { "from": 356653727, "to": 1942096762, "isCyclic": false }, { "from": 341567957, "to": 1331865351, "isCyclic": false }, { "from": 341567957, "to": 3235280883, "isCyclic": false }, { "from": 341567957, "to": 356653727, "isCyclic": false }, { "from": 341567957, "to": 2485026308, "isCyclic": false }, { "from": 341567957, "to": 2562016124, "isCyclic": false }, { "from": 341567957, "to": 552151127, "isCyclic": false }, { "from": 341567957, "to": 1339565699, "isCyclic": false }, { "from": 341567957, "to": 2241854762, "isCyclic": false }, { "from": 341567957, "to": 2567328503, "isCyclic": false }, { "from": 341567957, "to": 1094774888, "isCyclic": false }, { "from": 1005564345, "to": 801525785, "isCyclic": false }, { "from": 1005564345, "to": 1331865351, "isCyclic": false }, { "from": 1005564345, "to": 3235280883, "isCyclic": false }, { "from": 1005564345, "to": 1216276501, "isCyclic": false }, { "from": 1005564345, "to": 3177744305, "isCyclic": false }, { "from": 182641601, "to": 1331865351, "isCyclic": false }, { "from": 182641601, "to": 3235280883, "isCyclic": false }, { "from": 182641601, "to": 1216276501, "isCyclic": false }, { "from": 182641601, "to": 1440380694, "isCyclic": false }, { "from": 726487653, "to": 2048363028, "isCyclic": false }, { "from": 726487653, "to": 4037712688, "isCyclic": false }, { "from": 1428314197, "to": 3166841825, "isCyclic": false }, { "from": 1428314197, "to": 3431066080, "isCyclic": false }, { "from": 1428314197, "to": 2279247337, "isCyclic": false }, { "from": 1428314197, "to": 1362904939, "isCyclic": false }, { "from": 1560530462, "to": 4037712688, "isCyclic": false }, { "from": 1560530462, "to": 4065462691, "isCyclic": false }, { "from": 1560530462, "to": 491731072, "isCyclic": false }, { "from": 2680775406, "to": 1560530462, "isCyclic": false }, { "from": 2680775406, "to": 4037712688, "isCyclic": false }, { "from": 2680775406, "to": 491731072, "isCyclic": false }, { "from": 4099876531, "to": 4037712688, "isCyclic": false }, { "from": 4099876531, "to": 491731072, "isCyclic": false }, { "from": 4099876531, "to": 582030790, "isCyclic": false }, { "from": 107011351, "to": 2476273619, "isCyclic": false }, { "from": 107011351, "to": 2311248394, "isCyclic": false }, { "from": 107011351, "to": 1720170331, "isCyclic": false }, { "from": 107011351, "to": 1362904939, "isCyclic": false }, { "from": 4112388782, "to": 4037712688, "isCyclic": false }, { "from": 2120775579, "to": 4037712688, "isCyclic": false }, { "from": 582030790, "to": 491731072, "isCyclic": false }, { "from": 582030790, "to": 4037712688, "isCyclic": false }, { "from": 3222993554, "to": 1362904939, "isCyclic": false }, { "from": 3222993554, "to": 1965161762, "isCyclic": false }, { "from": 2428266738, "to": 1375713270, "isCyclic": false }, { "from": 2428266738, "to": 1331865351, "isCyclic": false }, { "from": 2428266738, "to": 1362904939, "isCyclic": false }, { "from": 2428266738, "to": 1216276501, "isCyclic": false }, { "from": 2428266738, "to": 3235280883, "isCyclic": false }, { "from": 2428266738, "to": 1759930456, "isCyclic": false }, { "from": 1216276501, "to": 3235280883, "isCyclic": false }, { "from": 1216276501, "to": 2376376984, "isCyclic": false }, { "from": 1942096762, "to": 1216276501, "isCyclic": false }, { "from": 1942096762, "to": 3235280883, "isCyclic": false }, { "from": 1942096762, "to": 1965161762, "isCyclic": false }, { "from": 2485026308, "to": 1339565699, "isCyclic": false }, { "from": 2485026308, "to": 2562016124, "isCyclic": false }, { "from": 2485026308, "to": 1076436185, "isCyclic": false }, { "from": 2485026308, "to": 944021367, "isCyclic": false }, { "from": 2485026308, "to": 552151127, "isCyclic": false }, { "from": 2485026308, "to": 3891966039, "isCyclic": false }, { "from": 2562016124, "to": 1339565699, "isCyclic": false }, { "from": 2562016124, "to": 552151127, "isCyclic": false }, { "from": 2562016124, "to": 3891966039, "isCyclic": false }, { "from": 1339565699, "to": 552151127, "isCyclic": false }, { "from": 2241854762, "to": 2567328503, "isCyclic": false }, { "from": 2241854762, "to": 1089468033, "isCyclic": false }, { "from": 1094774888, "to": 1339565699, "isCyclic": false }, { "from": 1094774888, "to": 2485026308, "isCyclic": false }, { "from": 2048363028, "to": 224600506, "isCyclic": false }, { "from": 2048363028, "to": 2149662377, "isCyclic": false }, { "from": 2048363028, "to": 491731072, "isCyclic": false }, { "from": 2048363028, "to": 4037712688, "isCyclic": false }, { "from": 4065462691, "to": 21548004, "isCyclic": false }, { "from": 4065462691, "to": 2714238133, "isCyclic": false }, { "from": 4065462691, "to": 409742873, "isCyclic": false }, { "from": 4065462691, "to": 1362904939, "isCyclic": false }, { "from": 1375713270, "to": 1362904939, "isCyclic": false }, { "from": 1375713270, "to": 4178430438, "isCyclic": false }, { "from": 1375713270, "to": 1372522740, "isCyclic": false }, { "from": 2376376984, "to": 3235280883, "isCyclic": false }, { "from": 1076436185, "to": 971943911, "isCyclic": false }, { "from": 1076436185, "to": 3596385220, "isCyclic": false }, { "from": 1076436185, "to": 2548718478, "isCyclic": false }, { "from": 944021367, "to": 356653727, "isCyclic": false }, { "from": 944021367, "to": 3982654405, "isCyclic": false }, { "from": 944021367, "to": 2839729826, "isCyclic": false }, { "from": 944021367, "to": 552151127, "isCyclic": false }, { "from": 224600506, "to": 491731072, "isCyclic": false }, { "from": 224600506, "to": 4037712688, "isCyclic": false }, { "from": 224600506, "to": 3177744305, "isCyclic": false }, { "from": 2714238133, "to": 21548004, "isCyclic": false }, { "from": 409742873, "to": 21548004, "isCyclic": false }, { "from": 971943911, "to": 2548718478, "isCyclic": false }, { "from": 3596385220, "to": 971943911, "isCyclic": false }, { "from": 3982654405, "to": 2839729826, "isCyclic": false }, { "from": 3982654405, "to": 3997750360, "isCyclic": false }, { "from": 3982654405, "to": 552151127, "isCyclic": false }, { "from": 3982654405, "to": 23089748, "isCyclic": false }, { "from": 3982654405, "to": 1076436185, "isCyclic": false }, { "from": 2839729826, "to": 552151127, "isCyclic": false }, { "from": 2839729826, "to": 23089748, "isCyclic": false }, { "from": 3997750360, "to": 2839729826, "isCyclic": false }, { "from": 3997750360, "to": 23089748, "isCyclic": false }, { "from": 3997750360, "to": 552151127, "isCyclic": false }, { "from": 3997750360, "to": 679976714, "isCyclic": false }, { "from": 3997750360, "to": 1618455164, "isCyclic": false }, { "from": 3997750360, "to": 385340190, "isCyclic": false }, { "from": 3997750360, "to": 1284823955, "isCyclic": false }, { "from": 23089748, "to": 1076436185, "isCyclic": false }, { "from": 23089748, "to": 552151127, "isCyclic": false }, { "from": 385340190, "to": 23089748, "isCyclic": false }, { "from": 385340190, "to": 971943911, "isCyclic": false }, { "from": 385340190, "to": 552151127, "isCyclic": false }, { "from": 385340190, "to": 2548718478, "isCyclic": false }, { "from": 385340190, "to": 1284823955, "isCyclic": false }, { "from": 3022277916, "to": 3232333941, "isCyclic": true }, { "from": 3222993554, "to": 3672289200, "isCyclic": true }, { "from": 1432468757, "to": 3232333941, "isCyclic": true }, { "from": 2839729826, "to": 3982654405, "isCyclic": true }, { "from": 3997750360, "to": 3982654405, "isCyclic": true }, { "from": 1089468033, "to": 2241854762, "isCyclic": true }, { "from": 2893840283, "to": 3232333941, "isCyclic": true }, { "from": 806026951, "to": 3232333941, "isCyclic": true }, { "from": 29475502, "to": 3232333941, "isCyclic": true }, { "from": 2048363028, "to": 3917048986, "isCyclic": true }, { "from": 224600506, "to": 726487653, "isCyclic": true }, { "from": 224600506, "to": 2048363028, "isCyclic": true }, { "from": 2149662377, "to": 2048363028, "isCyclic": true }, { "from": 2048363028, "to": 726487653, "isCyclic": true }, { "from": 4099876531, "to": 3967137051, "isCyclic": true }], "enableGui": false }

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
  DIR_LINK_STRENGTH_FACTOR: 1,
  PACKAGE_LINK_STRENGTH_FACTOR: 1,
  HIGHLIGHT_CYCLES: false
}

const UNREAL_BLOOM_PASS = new UnrealBloomPass()

const { xGraph: X_GRAPH, nodes: NODES, fileTree: FILE_TREE } = buildXGraph(Data.__INLINE_DATA)

function App () {
  const [highlightNodes, setHighlightNodes] = useState(new Set<XNode>())
  const [highlightLinks, setHighlightLinks] = useState(new Set<XLink>())
  const [selectedNode, setSelectedNode] = useState<XNode>()

  const graph = useRef<ForceGraphMethods<NodeObject<XNode>, LinkObject<XNode, XLink>>>();
  const settings = useControls(DEFAULT_SETTINGS)

  const lastBackgroundClick = useRef(0);

  function backgroundClick () {
    const now = new Date().getTime()
    if (selectedNode) {
      selectNode(undefined)
    } else {
      if (now - lastBackgroundClick.current < settings.DOUBLE_CLICK_INTERVAL) {
        graph.current?.zoomToFit(settings.DEFAULT_DISTANCE)
      }
    }
    lastBackgroundClick.current = now
  }

  function colorNode (node: XNode) {
    let alpha = settings.NODE_ALPHA
    if (highlightNodes.size > 0 && !highlightNodes.has(node)) alpha = settings.UNSELECTED_NODE_ALPHA
    const { h, s, v } = (node.isEntrypoint || node.__color === undefined) ? { h: 0, s: 0, v: 1 } : node.__color
    const [r, g, b] = HSVtoRGB(h, s, v)
    return `rgba(${r}, ${g}, ${b}, ${alpha})`;
  }

  function colorLink (link: XLink) {
    let alpha = settings.LINK_ALPHA
    if (highlightLinks.size > 0 && !highlightLinks.has(link)) alpha = settings.UNSELECTED_LINK_ALPHA
    if (link.isCyclic && settings.HIGHLIGHT_CYCLES) return `indianred`;
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
    const distance = settings.DEFAULT_DISTANCE;
    const { x = 1, y = 1, z = 1 } = node
    const distRatio = 1 + distance / Math.hypot(x, y, z);

    graph.current?.cameraPosition({ x: x * distRatio, y: y * distRatio, z: z * distRatio }, { x, y, z }, 1000)
  }

  function nodeClick(node: XNode) {
    selectNode(node)
    centerOnNode(node)
  }

  useEffect(() => {
    graph.current?.postProcessingComposer().removePass(UNREAL_BLOOM_PASS)
    UNREAL_BLOOM_PASS.strength = settings.BLOOM_PASS_STRENGTH
    UNREAL_BLOOM_PASS.radius = settings.BLOOM_PASS_RADIUS
    UNREAL_BLOOM_PASS.threshold = settings.BLOOM_PASS_THRESHOLD
    graph.current?.postProcessingComposer().addPass(UNREAL_BLOOM_PASS)
  }, [settings.BLOOM_PASS_RADIUS, settings.BLOOM_PASS_STRENGTH, settings.BLOOM_PASS_THRESHOLD])

  useEffect(() => {
    graph.current?.d3Force('link')
      ?.distance(() => settings.LINK_DISTANCE)
      .strength((link: XLink) => {
        let f = settings.FILE_LINK_STRENGTH_FACTOR
        if (link.isDir) f = settings.DIR_LINK_STRENGTH_FACTOR
        if (link.isPackage) f = settings.PACKAGE_LINK_STRENGTH_FACTOR
        return f / Math.min(NODES[link.from].neighbors?.length ?? 1, NODES[link.to].neighbors?.length ?? 1);
      })
    graph.current?.d3ReheatSimulation()
  }, [settings.DIR_LINK_STRENGTH_FACTOR, settings.FILE_LINK_STRENGTH_FACTOR, settings.LINK_DISTANCE, settings.PACKAGE_LINK_STRENGTH_FACTOR])

  useEffect(() => {
    graph.current?.d3Force('charge')
      ?.strength((node: XNode) => {
        let f = settings.FILE_NODE_REPULSION_FORCE
        if (node.isDir) f = settings.DIR_NODE_REPULSION_FORCE
        if (node.isPackage) f = settings.PACKAGE_NODE_REPULSION_FORCE
        return -f
      })
    graph.current?.d3ReheatSimulation()
  }, [settings.DIR_NODE_REPULSION_FORCE, settings.FILE_NODE_REPULSION_FORCE, settings.PACKAGE_NODE_REPULSION_FORCE])

  useEffect(() => {
    setTimeout(() => graph.current?.zoomToFit(settings.DEFAULT_DISTANCE), 1000)
  }, [settings.DEFAULT_DISTANCE]);

  return (
    <>
      <ForceGraph
        ref={graph}
        extraRenderers={[new CSS2DRenderer()]}
        graphData={X_GRAPH}
        backgroundColor={'#000003'}
        nodeResolution={settings.NODE_RESOLUTION}
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
        onNodeClick={nodeClick}
        linkDirectionalArrowLength={4}
        linkDirectionalArrowRelPos={1}
        linkColor={colorLink}
        linkDirectionalArrowColor={colorLink}
        linkSource={'from' satisfies keyof XLink}
        linkTarget={'to' satisfies keyof XLink}
        linkVisibility={link => !link.isDir && !link.isPackage}
        linkWidth={link => highlightLinks.has(link) ? settings.LINK_HIGHLIGHT_WIDTH : settings.LINK_WIDTH}
        linkDirectionalParticles={link => highlightLinks.has(link) ? 2 : 0}
        linkDirectionalParticleWidth={settings.LINK_HIGHLIGHT_WIDTH}
      />
      <Explorer
        className={'fixed top-1 left-1 max-h-full bg-transparent'}
        fileTree={FILE_TREE}
        onSelectNode={nodeClick}
      />
      <Leva hidden={!X_GRAPH.enableGui}/>
    </>
  )
}

export default App
