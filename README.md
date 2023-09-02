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
    Dep Tree is a command line tool for visualizing your project's 
    import structure with an interactive graph.
</p>
<p align="center">
    Furthermore, it enables the specification of prohibited dependencies,
    ensuring your CI validates the integrity of your dependency graph.
</p>

<table align="center">
    <thead>
        <tr>
            <th>
                Given your project's file structure...
            </th>
            <th>
                ...it renders an interactive file import graph
            </th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>

<p align="center">
    <img height="600px" src="docs/tree.png" alt="File structure">
</p>

</td><td>

<p align="center">
    <img height="600px" src="docs/demo.gif" alt="Dependency tree render">
</p>

</td></tr></tbody></table>

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

## Dep Tree

`dep-tree` is a cli tool that allows users to render their file dependency tree in the terminal, or
check that it matches some dependency rules in CI systems.

It works with files, meaning that each file is a node in the dependency tree:
- It starts from an entrypoint, which is usually the main executable file in a
program or the file that exposes the contents of a library (like `package/main.py`).
- It reads its import statements, it makes a parent node out of the main file,
and one child node for each imported file.
> **NOTE**: it only takes into account local files, not files imported from external libraries.
- That process is repeated recursively with all the child files, until the file dependency
tree is formed.
- If rendering the dependency tree in the terminal, the nodes will be placed in a human-readable
way, and users can navigate through the graph using the keyboard.
- If validating the dependency tree in a CI system, it will check that the dependencies between files
match some boundaries declared in a `.dep-tree.yml` file.

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
npm install --save-dev @dep-tree/cli
```

## Usage
### Render

Choose the file that will act as the root of the dependency tree (for example `package/main.py`), and run:

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
Create a configuration file `.dep-tree.yml` with some rules:

```yml
entrypoints:
  - src/main.py
allow:
  "src/utils/**/*.py":
    - "src/utils/**/*.py"  # The files in src/utils can only depend on other utils
deny:
  "src/ports/**/*.py":
    - "**"  # A port cannot have any dependency
  "src/user/**":
    - "src/products/**" # The users domain cannot be related to the products domain
```

And check that your project matches those rules:

```shell
dep-tree check
```

## Supported languages

- Python
- JavaScript/TypeScript (es imports/exports)
- Rust (beta)
- Golang (coming soon...)
