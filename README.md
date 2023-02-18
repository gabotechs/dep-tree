<p align="center">
    <img height="96" src="./docs/dep-tree.svg"/>
    <img height="74" src="./docs/dep-tree-name.svg"/>
</p>

<br/>

<p align="center">
    <img src="https://coveralls.io/repos/github/gabotechs/dep-tree/badge.svg?branch=main">
    <img src="https://img.shields.io/github/v/release/gabotechs/dep-tree?color=%e535abff">
</p>


<p align="center">
    Render your project's dependency tree in the terminal and/or validate it against your rules.
</p>

<p align="center">
    <img width="440" src="docs/demo.gif" alt="Dependency tree render">
</p>

## Install

There is a node wrapper that can be installed with:

```shell
npm install @dep-tree/cli
# or
yarn add @dep-tree/cli
```

Installing the standalone precompiled binary can be done using [brew](https://brew.sh/index_es):
```shell
brew install gabotechs/taps/dep-tree
```

## Usage

With dep-tree you can either render an interactive dependency tree in your terminal, or check
that your project's dependency graph matches some user defined rules.

### Render

Choose the file that will act as the root of the dependency tree and run:

```shell
dep-tree render my-file.js
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

### Dependency check

Create a configuration file `.dep-tree.yml` with some rules

```yml
entrypoints:
  - src/index.ts
allow:
  "src/utils/**/*.ts":
    - "src/utils/**/*.ts"  # The files in src/utils can only depend on other utils
deny:
  "src/ports/**/*.ts":
    - "**"  # A port cannot have any dependency
```

and check that your project matches those rules:

```shell
dep-tree check
```

## Supported languages

- JavaScript/TypeScript
- Python (coming soon...)
- Rust (coming soon...)
- Golang (coming soon...)
