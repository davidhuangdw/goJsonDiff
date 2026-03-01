# goJsonDiff

[![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A Go library for computing and visualizing differences between two JSON values with human-readable output formats.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Output Examples](#output-examples)
- [CLI Tool](#cli-tool)
- [Documentation](#documentation)

## Features

- **Multiple Output Formats**
  - **HTML**: Rich visual diff with color-coded changes
  - **Console**: Colored terminal output for command-line viewing
  - **JsonPatch**: Standard [RFC 6902](https://www.rfc-editor.org/rfc/rfc6902.html) compliant patches

- **Efficient Array Comparison**
  - Uses LCS (Longest Common Subsequence) algorithm on item hashes
  - Detects array item moves, additions, and removals efficiently
  - Handles complex nested structures

## Installation

```bash
go get github.com/davidhuangdw/goJsonDiff
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/davidhuangdw/goJsonDiff"
)

func main() {
    fromJson := `{"name": "Pluto", "category": "planet"}`
    toJson := `{"name": "Pluto", "category": "dwarf planet"}`

    // Compute the delta
    delta, err := goJsonDiff.DiffJsonStr(fromJson, toJson)
    if err != nil {
        panic(err)
    }

    // Render to console
    result, err := goJsonDiff.RenderConsole(delta)
    if err != nil {
        panic(err)
    }

    fmt.Println(result)
}
```

## Usage

For more examples, see [/docs/examples.go](/docs/examples.go) and [/cmd/diff.go](/cmd/diff.go)

```go
import (
    "github.com/davidhuangdw/goJsonDiff"
)

func example() {
    var fromJson, toJson string

    // Step 1: Compute delta between two JSON strings
    delta, err := goJsonDiff.DiffJsonStr(fromJson, toJson)
    if err != nil {
        // Handle error
        return
    }

    // Step 2: Render delta to desired format

    // Console format (colored text)
    consoleOutput, err := goJsonDiff.RenderConsole(delta)
    if err != nil {
        // Handle error
        return
    }

    // HTML format
    htmlOutput, err := goJsonDiff.RenderHtml(delta)
    if err != nil {
        // Handle error
        return
    }

    // JsonPatch format (RFC 6902)
    patchOutput, err := goJsonDiff.RenderJsonPatch(delta)
    if err != nil {
        // Handle error
        return
    }

    // Use the output
    println(consoleOutput)
}
```

## Output Examples

For the following two JSON inputs:

<table>
<tr>
<th>From JSON</th>
<th>To JSON</th>
</tr>

<tr>
<td>
<pre>
{
  "name": "Pluto",
  "number": 1.0,
  "category": "planet",
  "to_replace": "foo",
  "to_remove": {"a": 1},
  "composition": [
    {"x": 1, "y": 2},
    "methane",
    null,
    "nitrogen",
    "hi",
    "bar",
    {"a": 1, "b": 2, "c": 3}
  ],
  "deep": {
    "rp": {"x": 100},
    "kp": {
      "rm": 2,
      "x": 1,
      "sub_ign": [1, 2]
    }
  }
}
</pre>
</td>

<td>
<pre>
{
  "name": "Pluto",
  "number": 1.0,
  "category": "dwarf planet",
  "to_replace": 999,
  "composition": [
    {"x": 1, "+": 9, "y": 3},
    "foo",
    "methane",
    null,
    "nitrogen",
    {"a": 1, "b": 2, "c": 3},
    1
  ],
  "deep": {
    "rp": 9,
    "ad": null,
    "kp": {
      "x": 1,
      "sub_ign": 999
    }
  }
}
</pre>
</td>
</tr>
</table>

### Console Format

<img src="/docs/demo/result/console.png" alt="console demo" width="800px">

### HTML Format

View the full example: [result.html](/docs/demo/result/result.html)

<img src="/docs/demo/result/html.png" alt="html demo" width="300px">

### JsonPatch Format

View the full example: [result.jsonpatch](/docs/demo/result/result.jsonpatch)

```json
[
    {
        "op": "replace",
        "path": "/category",
        "value": "dwarf planet"
    },
    {
        "op": "replace",
        "path": "/to_replace",
        "value": 999
    },
    {
        "op": "remove",
        "path": "/to_remove",
        "value": {"a": 1}
    },
    {
        "op": "add",
        "path": "/composition/0/+",
        "value": 9
    },
    {
        "op": "replace",
        "path": "/composition/0/y",
        "value": 3
    }
    // ... (truncated for brevity, see full output in docs/demo/result/result.jsonpatch)
]
```

## CLI Tool

A command-line tool built on goJsonDiff for comparing JSON files directly from the terminal.

**Source**: [/cmd/diff.go](/cmd/diff.go)

### Building the CLI

```bash
go build cmd/diff.go
```

### Usage

```bash
# Show help
./diff -h

# Console output with debug info
./diff -debug docs/demo/from.json docs/demo/to.json

# HTML output
./diff -debug -f=html docs/demo/from.json docs/demo/to.json > tmp/result.html

# JsonPatch output
./diff -debug -f=patch docs/demo/from.json docs/demo/to.json > tmp/result.jsonpatch
```

### Options

```
Usage: ./diff [OPTIONS] {from_json_file_name} {to_json_file_name}

  -debug
        Show timing and extra debug information
  -f string
        Output format: console, html, patch (default "console")
```
## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
