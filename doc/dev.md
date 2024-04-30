## Implement a decoder

### Steps to add new decoder

- Create a directory  `format/<name>`
- Copy some similar decoder, `format/format/bson.go` is quite small, to `format/<name>/<name>.go`
- Cleanup and fill in the register struct, rename `format.BSON` and add it
to `format/format.go` and don't forget to change the string constant.
- Add an import to `format/all/all.go`

### Some general tips

- Main goal is to produce a tree structure that is user-friendly and easy to work with.
Prefer a nice and easy tree structure over nice decoder implementation.
- Use same names, symbols, constant number bases etc as in specification.
But maybe in lowercase to be jq/JSON-ish.
- Decode only ranges you know what they are. If possible let "parent" decide what to do with unknown gaps
bits by using `*Decode*Len/Range/Limit` functions. fq will also automatically add "gap" fields if
it finds gaps.
- If you have decode helpers functions that decode a bunch of fields etc it is usually nice to make it only decode fields, not seek or add it's own "containing" struct. That way the function will be easier to reuse and only do one thing. Ex the helper `func decodeHeader(d *decode.D)` can then be use as `d.FieldStruct("header", decodeHeader)`, `d.SeekRel(1234, decodeHeader)` or `d.SeekRel(1234, func(d *decode.D) { d.FieldStruct("header, decodeHeader") }`
- Try to not decode too much as one value.
A length encoded int could be two fields, but maybe a length prefixed string should be one.
Flags can be struct with bit-fields.
- Map as many value as possible to symbolic values.
- Endian is inherited inside one format decoder, defaults to big endian for new format decoder
- Make sure zero length or no frames/packets etc fails decoding
- If format is in the probe group make sure to validate input to make it non-ambiguous with other decoders
- Try keep decoder code "declarative" if possible
- Split into multiple sub formats if possible. Makes it possible to use them separately.
- Validate/Assert
- Error/Fatal/panic
- Can new formats be added to other formats?
- Does the new format include existing formats?

### Checklist

- Commits:
  - Use commit messages with a context prefix to make it easier to find and understand, ex:<br>
  `mp3: Validate sync correctly`
- Tests:
  - If possible use a pair of `testdata/file` and `testdata/file.fqtest` where `file.fqtest` is `$ fq dv file` or `$ fq 'dv,torepr' file` if there is `torepr` support.
  - If `dv` produces a lof of output maybe use `dv({array_truncate: 50})` etc
  - Run `go test ./format -run TestFormats/<name>` to test expected output.
  - Run `go test ./format -run TestFormats/<name> -update` to update current output as expected output.
- If you have format specific documentation:
  - Put it in `format/*/<name>.md` and use `//go:embed <name>.md`/`interp.RegisterFS(..)` to embed/register it.
  - Use simple markdown, just sections (depth starts at 3, `### Section`), paragraphs, lists and links.
  - No heading section is needs with format name, will be added by `make doc` and fq cli help system.
  - Add a `testdata/<name>_help.fqtest` with just `$ fq -h <name>` to test CLI help.
  - If in doubt look at `mp4.md`/`mp4.go` etc.
  - Run `make README.md doc/formats.md` to update md files.
- Run linter `make lint`
- Run fuzzer `make fuzz GROUP=<name>`, see usage in Makefile

### Decoder API

`*decode.D` reader methods use this name convention:

`<Field>?(<reader<length>?>|<type>Fn>)(...[, scalar.Mapper...]) <type>`

- If it starts with `Field` a field will be added and first argument will be name of field. If not it will just read.
- `<try>?<reader<length>?>|<try>?<type>Fn>` a reader or a reader function
  - `<try>?` If prefixed with `Try` function return error instead of panic on error.
  - `<reader<length>?>` Read bits using some decoder.
    - `U16` unsigned 16 bit integer.
    - `UTF8` UTF8 with byte length as argument.
  - `<type>Fn>` read using a `func(d *decode.D) <type>`  function.
    - This can be used to implement own custom readers.

All `Field` functions takes a var args of `scalar.Mapper`:s that will be applied after reading.

`<type>` are these types:

| `<type>` | Go type | jq type |
| -------- | ------- | ------- |
| U        | uint64  | number  |
| S        | int64   | number  |
| F        | float64 | number  |
| Str      | string  | string  |
| Bool     | bool    | boolean |
| Nil      | nil     | null    |


TODO: there are some more (BitBuf etc, should be renamed)

To add a struct or array use `d.FieldStruct(...)` and `d.FieldArray(...)`.

TODO: nested formats, buffers, own decoders, scalar mappers

TODO: seeking, framed/limited/range decode

For example this decoder:

```go
// read 4 byte UTF8 string and add it as "magic", return a string
d.FieldUTF8("magic", 4)
// create a new struct and add it as "headers", returns a *decode.D
d.FieldStruct("headers", func(d *decode.D) {
    // read 8 bit unsigned integer, map it and add it as "type", returns a uint64
    d.FieldU8("type", scalar.UintMapSymStr{
        1: "start",
        // ...
    })
})
```

will produce something like this:

```go
*decode.Value{
    Parent: nil,
    V: *decode.Compound{
        IsArray: false, // is struct
        Children: []*decode.Value{
            *decode.Value{
                Name: "magic",
                V: scalar.Str{
                    Actual: "abcd", // read and set by UTF8 reader
                },
                Range: ranges.Range{Start: 0, Len: 32},
            },
            *decode.Value{
                Parent: &... // ref parent *decode.Value>,
                Name: "headers",
                V: *decode.Compound{
                    IsArray: false, // is struct
                    Children: []*decode.Value{
                        *decode.Value{
                            Name: "type",
                            V: scalar.Uint{
                                Actual: uint64(1), // read and set by U8 reader
                                Sym: "start", // set by UintMapSymStr scalar.Mapper
                            },
                            Range: ranges.Range{Start: 32, Len: 8},
                        },
                    },
                },
                Range: ranges.Range{Start: 32, Len: 8},
            },
        },
    },
    Range: ranges.Range{Start: 0, Len: 40},
}
```

and will look like this in jq/JSON:

```json
{
    "magic": "abcd",
    "headers": {
        "type": "start"
    }
}
```

#### `*decode.D` type

This is the main type used during decoding. It keeps track of:

- A current array or struct [`*decode.Value`](#decodevalue-type) where fields will be added.
- Current bit reader
- Current default endian
- Decode options

New [`*decode.D`](#decoded-type) are created during decoding when `d.FieldStruct` etc is used. It is also a kitchen sink of all kind functions for reading various standard number and string encodings etc.

Decoder authors do not have to create them.

#### `*decode.Value` type

Is what [`*decode.D`](#decoded-type) produces and it used to represent the decoded structure. Can be array, struct, number, string etc. It is the underlying type used by `interp.DecodeValue` that implements `gojq.JQValue` to expose it as various jq types, which in turn is used to produce JSON.

It stores:
- Parent [`*decode.Value`](#decodevalue-type) unless it's a root.
- A decoded value, a [`scalar.S`](#scalars-type) or [`*decode.Compound`](#decodecompound-type) (struct or array)
- Name in parent struct or array. If parent is a struct the name is unique.
- Index in parent array. Not used if parent is a struct.
- A bit range. Also struct and array have a range that is the min/max range of its children.
- A bit reader where the bit range can be read from.

Decoder authors will probably not have to create them.

#### `scalar.S` type

Keeps track of
- Actual value. Decoded value represented using a go type like `uint64`, `string` etc. For example a value reader by a utf8 or utf16 reader both will ends up as a `string`.
- Symbolic value. Optional symbolic representation of the actual value. For example a `scalar.UintMapSymStr` would map an actual `uint64` to a symbolic `string`.
- String description of the value.
- Number representation

The `scalar` package has `scalar.Mapper` implementations for all types to map actual to whole [`scalar.S`](#scalars-type) value `scalar.<type>ToScalar` or to just to set symbolic value `scalar.<type>ToSym<type>`. There is also mappers to just set values or to change number representations `scalar.Hex`/`scalar.SymHex` etc.

Decoder authors will probably not have to create them. But you might implement your own `scalar.Mapper` to modify them.

#### `*decode.Compound` type

Used to store struct or array of [`*decode.Value`](#decodevalue-type).

Decoder authors do not have to create them.

## Development tips

I usually use `-d <format>` and `dv` while developing, that way you will get a decode tree
even if it fails. `dv` gives verbose output and also includes stacktrace.

```sh
go run . -d <format> dv file
```

If the format is inside some other format it can be handy to first extract the bits and run
the decode directly. For example if working a `aac_frame` decoder issue:

```sh
fq '.tracks[0].samples[1234] | tobytes' file.mp4 > aac_frame_1234
fq -d aac_frame dv aac_frame_1234
```

Sometimes nested decoding fails then maybe a good way is to change the parent decoder to
use `d.RawLen()` etc instead of `d.FormatLen()` etc temporary to extract the bits. Hopefully
there will be some option to do this in the future.

When researching or investinging something I can recommend to use `watchexec`, `modd` etc to
make things more comfortable. Also using vscode/delve for debugging should work fine once
launch `args` are setup etc.

```
watchexec "go run . -d aac_frame dv aac_frame"
```

Some different ways to run tests:
```sh
# run all tests
make test
# run all go tests
go test ./...
# run all tests for one format
go test -run TestFormats/mp4 ./format/
# update all expected outputs for tests
go test ./pkg/interp ./format -update
# update actual output for specific tests
go run ./format -run TestFormats/elf -update
# color diff
DIFF_COLOR=1 go test ...
```

To lint source use:
```
make lint
```

Generate documentation. Requires [FFmpeg](https://github.com/FFmpeg/FFmpeg) and [Graphviz](https://gitlab.com/graphviz/graphviz):
```sh
make doc
```

TODO: `make fuzz`

## Debug

Split debug and normal output even when using repl:

Write `log` package output and stderr to a file that can be `tail -f`:ed in another terminal:
```sh
LOGFILE=/tmp/log go run . ... 2>>/tmp/log
```

gojq execution debug:
```sh
GOJQ_DEBUG=1 go run -tags debug . ...
```

Memory and CPU profile (will open a browser):
```sh
make memprof ARGS=". file"
make cpuprof ARGS=". test.mp3"
```

## From start to decoded value

```
main:main()
    cli.Main(default registry)
        interp.New(registry, std os interp implementation)
        interp.(*Interp).Main()
            interp.jq _main/0:
                args.jq _args_parse/2
                populate filenames for input/0
                interp.jq inputs/0
                    foreach valid input/0 output
                        interp.jq open
                            funcs.go _open
                        interp.jq decode
                            funcs.go _decode
                                decode.go Decode(...)
                                    ...
                        interp.jq eval expr
                            funcs.go _eval
                        interp.jq display
                            funcs.go _display
                                for interp.(decodeValueBase).Display()
                                    dump.go
                                        print tree
                                empty output
```

## bitio and other io packages

```
*os.File, *bytes.Buffer
^
ctxreadseeker.Reader defers blocking io operations to a goroutine to make them cancellable
^
progressreadseeker.Reader approximates how much of a file has been read
^
aheadreadseeker.Reader does readahead caching
^
| (io.ReadSeeker interface)
|
bitio.IOBitReader (implements bitio.Bit* interfaces)
SectionBitReader
MultiBitReader
```

## jq oddities

```
jq -n '[1,2,3,4] | .[null:], .[null:2], .[2:null], .[:null]'
```

## Setup docker desktop with golang windows container

```sh
git clone https://github.com/StefanScherer/windows-docker-machine.git
cd windows-docker-machine
vagrant up 2016-box
cd ../fq
docker --context 2016-box run --rm -ti -v "C:${PWD//\//\\}:C:${PWD//\//\\}" -w "$PWD" golang:1.18-windowsservercore-ltsc2016
```

## Implementation details

- fq uses a gojq fork that can be found at https://github.com/wader/gojq/tree/fq (the "fq" branch)
- cli readline uses raw mode to blocks ctrl-c to become a SIGINT

## Dependencies and source origins

- [gojq](https://github.com/itchyny/gojq) fork that can be found at https://github.com/wader/gojq/tree/fq<br>
Issues and PR:s related to fq:<br>
[#43](https://github.com/itchyny/gojq/issues/43) Support for functions written in go when used as a library<br>
[#46](https://github.com/itchyny/gojq/pull/46) Support custom internal functions<br>
[#56](https://github.com/itchyny/gojq/issues/56) String format query with no operator using %#v or %#+v panics
[#65](https://github.com/itchyny/gojq/issues/65) Try-catch with custom function<br>
[#67](https://github.com/itchyny/gojq/pull/67) Add custom iterator function support which enables implementing a REPL in jq<br>
[#81](https://github.com/itchyny/gojq/issues/81) path/1 behaviour and path expression question<br>
[#86](https://github.com/itchyny/gojq/issues/86) ER: basic TCO
[#109](https://github.com/itchyny/gojq/issues/109) jq halt_error behaviour difference<br>
[#113](https://github.com/itchyny/gojq/issues/113) error/0 and error/1 behavior difference<br>
[#117](https://github.com/itchyny/gojq/issues/117) Negative number modulus *big.Int behaves differently to int<br>
[#118](https://github.com/itchyny/gojq/issues/118) Regression introduced by "remove fork analysis from tail call optimization (ref #86)"<br>
[#122](https://github.com/itchyny/gojq/issues/122) Slow performance for large error values that ends up using typeErrorPreview()<br>
[#125](https://github.com/itchyny/gojq/pull/125) improve performance of join by make it internal<br>
[#141](https://github.com/itchyny/gojq/issues/141) Empty array flatten regression since "improve flatten performance by reducing copy"

- [gopacket](https://github.com/gopacket/gopacket) for TCP and IPv4 reassembly
- [mapstructure](https://github.com/mitchellh/mapstructure) for convenient JSON/map conversion
- [go-difflib](https://github.com/pmezard/go-difflib) for diff tests
- [golang.org/x/text](https://pkg.go.dev/golang.org/x/text) for text encoding conversions
- [float16.go](https://android.googlesource.com/platform/tools/gpu/+/gradle_2.0.0/binary/float16.go) to convert bits into 16-bit floats

## Release process

Run and follow instructions:
```
make release VERSION=1.2.3
```

Commits since release
```
git log --no-decorate --no-merges --oneline v0.0.4..wader/master | sort -t " " -k 2 | sed 's/\(.*\)/* \1/'
```
