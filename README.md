# StringerGen

stringergen is a command-line tool that generates `String` methods for structs in Go source files. It supports two modes of operation: **source mode** and **recursive mode**.

Explore the motivations and design choices behind StringerGen in my blog post [here](https://chasemao.com/article/exploring-go-stringer-usage/).

## Installation

Once you have [installed Go](https://go.dev/doc/install#releases), install the `stringergen` tool.

**Note**: If you have not done so already be sure to add `$GOPATH/bin` to your
`PATH`.

To get the latest released version use:

### Go version < 1.16

```sh
GO111MODULE=on go get github.com/chasemao/stringergen @v1.0.0
```

### Go 1.16+

```sh
go install github.com/chasemao/stringergen @v1.0.0
```

## Running stringergn

### Source Mode

In source mode, StringerGen generates `String` methods for structs from a specified source file.

**Usage:**
- Enable source mode with the `-source` flag.
- Optionally, use the `-destination` flag to specify an output file. If not set, the output will be written to stdout.

**Example:**
```sh
stringergen -source=foo.go -destination=bar.go
```

### Recursive Mode

In recursive mode, StringerGen generates `String` methods for all structs in files within the specified directory and its subdirectories.

**Usage:**

Enable recursive mode with the `-recursive` flag.
Optionally, use the `-save` flag to save the output to files named `xx_stringer.go` for `xx.go` files. If not set, the output will be written to stdout.
Use the `-skipdir` flag to specify directories to skip.

**Example:**

```sh
stringergen -recursive=/path/to/directory/ -save -skipdir=/path/to/directory/skip1/,/path/to/directory/skip2/
```

## Output

stringergen use `methol` flagsto determine method for the String method generation. Supported values: json, jsoniter, fmt; defaults to json.

There is a [benchmark result](./benchmark/README.md) on the performace of different method.

Below is examples of generated string method.

### json

```go
package main

import (
        "encoding/json"
)

// String Used in fmt to generate string
func (o *output) String() string {
        v, _ := json.Marshal(o)
        return string(v)
}
```

### jsoniter

```go
package main

import (
        jsoniter "github.com/json-iterator/go"
)

// String Used in fmt to generate string
func (o *output) String() string {
        v, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(o)
        return string(v)
}
```

### fmt

```go
package main

import (
        "fmt"
)

// String Used in fmt to generate string
func (o *output) String() string {
        return fmt.Sprintf("%+v", *o)
}
```


## Flags

* `-destination string`

(source mode) Output file; defaults to stdout, used in source mode.

* `-exclude string`

Regular expression patterns for struct names to exclude from generation, separated by commas (without quotation marks); defaults to none.

* `-method string`

Method for the String method generation. Supported values: json, jsoniter, fmt; defaults to json.

* `-recursive string`

(recursive mode) Input directory, will handle all files recursively.

* `-save`
(recursive mode) Write to file like `xx_stringer.go` for `xx.go`, used in recursive mode.


* `-skipdir string`

(recursive mode) Name of directory to skip, not to generate String methods within these directories; defaults to none.

* `-source string`

(source mode) Input Go source file.

* `-v`

Output detailed information.

* `-version`

Print version information.

## Examples

### Source Mode Example

Generate `String` methods for structs in `foo.go` and output to `bar.go`:

```sh
stringergen -source=foo.go -destination=bar.go
```

### Recursive Mode Example

Generate `String` methods for all structs in the directory `/path/to/directory/` and its subdirectories, saving the output to `xx_stringer.go` files, and skipping the directories `/path/to/directory/skip1/` and `/path/to/directory/skip2/`:

```sh
stringergen -recursive=/path/to/directory/ -save -skipdir=/path/to/directory/skip1/,/path/to/directory/skip2/
```

## Notes

* If the `-destination` flag is not set in source mode, the output will be written to stdout.
* If the `-save` flag is not set in recursive mode, the output will be written to stdout.
* Use the `-exclude` flag to provide regular expression patterns for struct names to exclude from generation.
* Use the `-method` flag to choose the method for the `String` method generation (`json`, `jsoniter`, `fmt`).

## Version

To print the version information, use the `-version` flag:

```sh
stringergen -version
```

## Verbose Output

To enable detailed output, use the `-v` flag:

```sh
stringergen -v
```

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

## Contributing
Contributions are welcome!