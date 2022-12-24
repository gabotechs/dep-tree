# Dep Tree

[![Coverage Status](https://coveralls.io/repos/github/gabotechs/dep-tree/badge.svg?branch=main)](https://coveralls.io/github/gabotechs/dep-tree?branch=main)
![](https://img.shields.io/github/v/release/gabotechs/dep-tree?color=%e535abff)

Render your project's dependency tree in the terminal.

## Install

Using brew

```shell
brew install gabotechs/taps/dep-tree
```

## Usage

Choose the file that will act as the root of the dependency tree and run:

```shell
dep-tree my-file.js
```

The dependency tree will be formed recursively based on the imports declared
in that file.

## Supported languages

- JavaScript/TypeScript
- Python (coming soon...)
- Rust (coming soon...)
- Golang (coming soon...)
