<p align="center">
    <img height="96" src="./docs/dep-tree.svg"/>
    <img height="74" src="./docs/dep-tree-name.svg"/>
</p>

<br/>

<p align="center">
    <img src="https://coveralls.io/repos/github/gabotechs/dep-tree/badge.svg?branch=main"/>
    <img src="https://goreportcard.com/badge/github.com/gabotechs/dep-tree"/>
    <img src="https://img.shields.io/github/v/release/gabotechs/dep-tree?color=%e535abff"/>
</p>

<p align="center">
    Dep Tree is a tool for helping developers maintain their code bases clean and decoupled. It 
    allows rendering the "project's entropy" using a 3D force directed graph, where each file is a
    node, and each dependency between two files is an edge.
</p>
<p align="center">
    The more decoupled a code base is, the more spread the 3d graph will look like.
</p>

<p align="center">
    <img height="430" src="docs/demo.gif" alt="File structure">
</p>

<p align="center">
    Additionally, it enables the specification of prohibited dependencies,
    ensuring your CI validates the integrity of your dependency graph, keeping your code base clean.
</p>

## Checkout the entropy graph of well-known libraries
- [three.js](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Fmrdoob%2Fthree.js&entrypoint=src%2FThree.js)
- [langchain](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Flangchain-ai%2Flangchain&entrypoint=libs%2Flangchain%2Flangchain%2F__init__.py)
- [vuejs](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Fvuejs%2Fvue&entrypoint=src%2Fcore%2Findex.ts)
- [pytorch](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Fpytorch%2Fpytorch&entrypoint=torch%2Fnn%2F__init__.py)
- [tensorflow](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Ftensorflow%2Ftensorflow&entrypoint=tensorflow%2Fpython%2Fkeras%2Fmodels.py)
- [storybook](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Fstorybookjs%2Fstorybook&entrypoint=code%2Fui%2Fblocks%2Fsrc%2Findex.ts) (ui module)

## Install

On Mac and Linux, it can be installed using brew:
```shell
brew install gabotechs/taps/dep-tree
```

Alternatively, on any platform including Windows it can be installed with `pip`:
```shell
pip install python-dep-tree
```

There is also a node wrapper that can be installed with:
```shell
npm install @dep-tree/cli
```

## About Dep Tree

`dep-tree` is a cli tool that allows users to visualize the complexity of a code base, and create
rules for ensuring its loosely coupling.

It works with files, meaning that each file is a node in the dependency tree:
- It starts from an entrypoint, which is usually the main executable file in a
program or the file that exposes the contents of a library (like `package/main.py`).
- It reads its import statements, it makes a parent node out of the main file,
and one child node for each imported file.
> **NOTE**: it only takes into account local files, not files imported from external libraries.
- That process is repeated recursively with all the child files, until the file dependency
tree is formed.
- If rendering the code base entropy, the nodes will be rendered in a 3d force directed graph
in the browser.
- If rendering the dependency tree in the terminal, the nodes will be placed in a human-readable
way, and users can navigate through the graph using the keyboard.
- If validating the dependency tree in a CI system, it will check that the dependencies between files
match some boundaries declared in a `.dep-tree.yml` file.

## Usage

### Entropy visualization

Choose the file that will act as the root of the dependency tree (for example `src/index.ts`), and run:

```shell
dep-tree entropy src/index.ts
```

It will open a browser window and will render your file dependency graph. You will see a lot of spheres
and lines connecting them. Each sphere is a file in your code base, and each line indicates a dependency
between two files.

The spheres will be placed mimicking some attraction/repulsion forces, that way parts of your code
base will tend to gravitate together if they are tightly coupled, and will tend to be separated if
they are loosely coupled.

### CLI tree visualization

Choose the file that will act as the root of the dependency tree (for example `my-file.py`), and run:

```shell
dep-tree render my-file.py
```

You can see the controls for navigating through the graph pressing `h` at any time:

```
j      -> move one step down
k      -> move one step up
Ctrl d -> move half page down
Ctrl u -> move half page up
Enter  -> select the current node as the root node
q      -> navigate backwards on selected nodes or quit
h      -> show this help section
```

### Dependency linting

The dependency linting can be executed with:

```shell
dep-tree check
```

This is specially useful for CI systems, for ensuring that parts of an application that
should not be coupled remain decoupled as the project evolves.

These are the parameters that can be configured in the `.dep-tree.yml` file:

### `entrypoints`: 
List of entrypoints that will act as root nodes for evaluating multiple
dependency trees. Some applications might expose more than one entrypoint, for that reason,
this parameter is a list. The most typical thing is that there is only one entrypoint.

### `allow`:
Map from glob pattern to list of glob patterns that define, using a "white list"
logic, what files can depend on what other files. For example:
```yml
allow:
  "src/products/**":
    - "src/products/**"
    - "src/common/**"
```
In the example above, any file under the `src/products` folder has the restriction of only
being able to import files that live either in the same `src/products` folder, or in the
`src/common` folder.

### `deny`: 
Map from glob pattern to list of glob patterns that define, using a "black list"
logic, what dependencies are forbidden. For example:

```yml
deny:
  "api/routes.py":
    - "adapters/**"
```

In the example above, the file `api/routes.py` can import from anywhere but the `adapters` folder.

### `allowCircularDependencies`:

Boolean parameter that defines whether circular dependencies are allowed or not. By default
they are not allowed.

```yml
allowCircularDependencies: true
```

### `aliases`:
Map from string to glob pattern that gathers utility groups of glob patterns that
can be reused in the `deny` and `allow` fields. For example:

```yml
aliases:
  "common-stuff":
    - "src/common/**"
    - "src/utils/**"
    - "src/helpers/**"
allow:
  "src/function.py":
    - "common-stuff"
    - "src/class.py"
```
is the same as saying:

```yml
allow:
  "src/function.py":
    - "src/common/**"
    - "src/utils/**"
    - "src/helpers/**"
    - "src/class.py"
```

### Example configuration file
Dep Tree by default will read the configuration file in `.dep-tree.yml`, which is expected to be a file
that contains the following settings:

```yml
# Files that should be completely ignored by dep-tree. It's fine to ignore
# some big files that everyone depends on and that don't add
# value to the visualization, like auto generated code.
exclude:
  - 'some-glob-pattern/**/*.ts'

# Whether to follow re-exports to the target file or not.
# Imagine that you have the following setup:
#
#  src/index.ts     -> import { foo } from './foo'
#  src/foo/index.ts -> export { bar as foo } from './bar'
#  src/foo/bar.ts   -> export function bar() {}
#
# If `followReExports` is true, a dependency will be created from
# `src/index.ts` to `src/foo/bar.ts`, and the middle file `src/foo/index.ts`
# will be ignored, as it's just there for re-exporting the `bar` symbol,
# which is actually declared on `src/foo/bar.ts`
#
# If `followReExports` is false, re-exported symbols will not be traced back
# to where they are declared, and instead two dependencies will be created:
# - from `src/index.ts` to `src/foo/index.ts`
# - from `src/foo/index.ts` to `src/foo/bar.ts`
#
# Entropy visualization tends to lead to better results if this is set to `false`,
# but CLI rendering is slightly better with this set to `true`.
followReExports: false

# Check configuration for the `dep-tree check` command. Dep Tree will check for dependency
# violation rules declared here, and fail if there is at least one unsatisfied rule.
check:
  # These are the entrypoints to your application. Dependencies will be checked with
  # these files as root nodes. Typically, an application has only one entrypoint, which
  # is the executable file (`src/index.ts`, `main.py`, `src/lib.rs`, ...), but here
  # you can declare as many as you want.
  entrypoints:
    - src/index.ts

  # Whether to allow circular dependencies or not. Languages typically allow
  # having circular dependencies, but that has an impact in execution path
  # traceability, so you might want to disallow it.
  allowCircularDependencies: false

  # map from glob pattern to array of glob patterns that determines the exclusive allowed
  # dependencies that a file matching a key glob pattern might have. If file that
  # matches a key glob pattern depends on another file that does not match any of
  # the glob patterns declared in the values array, the check will fail.
  allow:
    # example: any file in `src/products` can only depend on files that are also
    # in the `src/products` folder or in the `src/helpers` folder.
    'src/products/**':
      - 'src/products/**'
      - 'src/helpers/**'

  # map from glob pattern to array of glob patterns that determines forbidden
  # dependencies. If a file that matches a key glob pattern depends on another
  # file that matches at least one of the glob patterns declared in the values
  # array, the check will fail.
  deny:
    # example: files inside `src/products` cannot depend on files inside `src/users`,
    # as they are supposed to belong to different domains.
    'src/products/**':
      - 'src/users/**'

  # typically, in a project, there is a set of files that are always good to depend
  # on, because they are supposed to be common helpers, or parts that are actually
  # designed to be widely depended on. This allows you to create aliases to group
  # of files that are meant to be widely depended on, so that you can reference
  # them afterward in the `allow` or `deny` sections.
  aliases:
    # example: this 'common' entry will be available in the other sections:
    #
    # check:
    #   allow:
    #     'src/products/**':
    #       - 'common'
    'common':
      - 'src/helpers/**'
      - 'src/utils/**'
      - 'src/generated/**'

# JavaScript and TypeScript specific settings.
js:
  # Whether to follow tsconfig.json paths or not. You will typically want to
  # enable this, but for some monorepo setups, it might be better to leave this off.
  followTsConfigPaths: true

# Python specific settings.
python:
# None available at the moment.

# Rust specific settings.
rust:
# None available at the moment.
```

## Motivation

As codebases expand and teams grow, complexity inevitably creeps in.
While maintaining a cohesive and organized structure is key to
a project's scalability and maintainability,
the current developer toolbox often falls short in one critical
area: file structure and dependency management.

Luckily, the community has come up with very useful tools
for keeping our projects in check:
- **Type checkers** ensure correct interactions between code segments.
- **Linters** elevate code quality and maintain a consistent style.
- **Formatters** guarantee a uniform code format throughout.
- But what about file structure and file dependency management...

Dep Tree is a dedicated tool addressing this very challenge,
it aids developers in preserving a project's structural integrity
throughout its lifecycle. And with integration capabilities in CI systems,
the tool ensures that this architectural "harmony" remains undisturbed.


## Supported languages

- Python
- JavaScript/TypeScript (es imports/exports)
- Rust (beta)
- Golang (coming soon...)
