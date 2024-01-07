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
    Visualize the <strong>entropy</strong> of a code base with a 3d force-directed graph. 
</p>

<p align="center">
    The more decoupled and modular a code base is, the more spread the graph will look like.
</p>

<p align="center">
    <img width="819" src="docs/demo.gif" alt="File structure">
</p>

<p align="center">
    Ensure your code base decoupling by creating your own rules and enforcing them with <code>dep-tree check</code>
</p>

## Checkout the entropy graph of well-known projects
- [typescript](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Fmicrosoft%2FTypeScript&entrypoint=src%2Ftypescript%2Ftypescript.ts)
- [react](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Ffacebook%2Freact&entrypoint=packages%2Freact-dom%2Findex.js)
- [svelte](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Fsveltejs%2Fsvelte&entrypoint=packages%2Fsvelte%2Fsrc%2Fcompiler%2Findex.js)
- [vuejs](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Fvuejs%2Fvue&entrypoint=src%2Fcore%2Findex.ts)
- [angular](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Fangular%2Fangular&entrypoint=packages%2Fcompiler%2Findex.ts)
- [storybook](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Fstorybookjs%2Fstorybook&entrypoint=code%2Fui%2Fblocks%2Fsrc%2Findex.ts)
- [three.js](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Fmrdoob%2Fthree.js&entrypoint=src%2FThree.js)
- [expressjs](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Fexpressjs%2Fexpress&entrypoint=lib%2Fexpress.js)
- [material-ui](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Fmui%2Fmaterial-ui&entrypoint=packages%2Fmui-material%2Fsrc%2Findex.js)
- [eslint](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Feslint%2Feslint&entrypoint=lib%2Fcli.js)
- [prettier](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Fprettier%2Fprettier&entrypoint=index.js)
- [langchain](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Flangchain-ai%2Flangchain&entrypoint=libs%2Flangchain%2Flangchain%2F__init__.py)
- [pytorch](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Fpytorch%2Fpytorch&entrypoint=torch%2Fnn%2F__init__.py)
- [tensorflow](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Ftensorflow%2Ftensorflow&entrypoint=tensorflow%2Fpython%2Fkeras%2Fmodels.py)
- [fastapi](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Ftiangolo%2Ffastapi&entrypoint=fastapi%2F__init__.py)
- [numpy](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Fnumpy%2Fnumpy&entrypoint=numpy%2F__init__.py)
- [scikit-learn](https://dep-tree-explorer.vercel.app/api?repo=https%3A%2F%2Fgithub.com%2Fscikit-learn%2Fscikit-learn&entrypoint=sklearn%2F__init__.py)

## Install

On Mac and Linux, it can be installed using brew:
```shell
brew install gabotechs/taps/dep-tree
```

Alternatively, on any platform including Windows it can be installed with `pip`...
```shell
pip install python-dep-tree
```

...or `npm`:
```shell
npm install @dep-tree/cli
```

## Supported languages

<div>
    <img height="40px" src="docs/js-logo.png">
    <img height="40px" src="docs/ts-logo.png">
    <img height="40px" src="docs/python-logo.png">
    <img height="40px" src="docs/rust-logo.png">
</div>

## About Dep Tree

`dep-tree` is a cli tool for visualizing the complexity of a code base, and creating
rules for ensuring its loosely coupling.

It works with files, meaning that each file is a node in the dependency tree:
- It starts from an entrypoint, which is usually the main executable file in a
program or the file that exposes the contents of a library (like `package/main.py`, `src/index.ts`, `src/lib.rs`...).
- It makes a parent node out of the root file, and one child node for each imported file.

> [!NOTE]
> it only takes into account local files, not files imported from external libraries.

- That process is repeated recursively with all the imported files, until the file dependency
graph is formed.
- If rendering the **code base entropy**, the nodes will be rendered using a 3d force-directed graph
in the browser.
- If rendering the **dependency tree** in the terminal, the nodes will be placed in a human-readable
way, and users can navigate through the graph using the keyboard.
- If validating the **dependency rules** in a CI system, it will check that the dependencies between files
match some boundaries declared in a `.dep-tree.yml` file.

## Usage

### Entropy

Choose the file that will act as the root of the dependency graph (for example `src/index.ts`), and run:

```shell
dep-tree entropy src/index.ts
```

It will open a browser window and will render your file dependency graph using a 3d force-directed graph.

The spheres (files) will be placed mimicking some attraction/repulsion forces. Some parts of your code
base will tend to gravitate together if they are tightly coupled, and will tend to be separated if
they are loosely coupled.

The 3d graph for a clean code base will have groups of nodes clustered together and clearly separated
from other clusters:

<img height="200px" src="docs/decoupled-code-base.png">

The 3d graph for a tightly coupled code base will have all the nodes grouped together with no
clustering and no clear separation between them:

<img height="200px" src="docs/coupled-code-base.png">

### Tree

Choose the file that will act as the root of the dependency graph (for example `my-file.py`), and run:

```shell
dep-tree tree my-file.py
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

### Check

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

# Whether to unwrap re-exports to the target file or not.
# Imagine that you have the following setup:
#
#  src/index.ts     -> import { foo } from './foo'
#  src/foo/index.ts -> export { bar as foo } from './bar'
#  src/foo/bar.ts   -> export function bar() {}
#
# If `unwrapExports` is true, a dependency will be created from
# `src/index.ts` to `src/foo/bar.ts`, and the middle file `src/foo/index.ts`
# will be ignored, as it's just there for re-exporting the `bar` symbol,
# which is actually declared on `src/foo/bar.ts`
#
# If `unwrapExports` is false, re-exported symbols will not be traced back
# to where they are declared, and instead two dependencies will be created:
# - from `src/index.ts` to `src/foo/index.ts`
# - from `src/foo/index.ts` to `src/foo/bar.ts`
#
# Entropy visualization tends to lead to better results if this is set to `false`,
# but CLI rendering is slightly better with this set to `true`.
unwrapExports: false

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
  # Whether to take package.json workspaces into account while resolving paths
  # or not. You might want to disable it if you only want to analyze one workspace
  # in a monorepo.
  workspaces: true
  # Whether to follow tsconfig.json paths or not. You will typically want to
  # enable this, but for some monorepo setups, it might be better to leave this off
  # if you want to analyze only one package.
  tsConfigPaths: true

# Python specific settings.
python:
  # Whether to take into account conditional imports as dependencies between files or not.
  # A conditional import is an `import` statement that is wrapped inside an `if` block or
  # a function, for example:
  #
  # if SHOULD_IMPORT:
  #     from foo import *
  #
  # by default these statements introduce a dependency between importing and imported file,
  # but depending on your use case you might want to disable it.
  excludeConditionalImports: false

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

