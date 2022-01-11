## Implement a decoder

### Steps to add new decoder

- Create a directory  `format/<name>`
- Copy some similar decoder, `format/format/bson.go` is quite small, to `format/<name>/<name>.go`
- Cleanup and fill in the register struct, rename `format.BSON` and add it
to `format/fromat.go` and don't forget to change the string constant.
- Add an import to `format/all/all.go`

### Some general tips

- Main goal is to produce a tree structure that is user-friendly and easy to work with.
Prefer a nice and easy query tree structure over nice decoder implementation.
- Use same names, symbols, constant number bases etc as in specification.
But maybe in lowercase to be jq/JSON-ish.
- Decode only ranges you know what they are. If possible let "parent" decide what to do with unknown
bits by using `*Decode*Len/Range/Limit` functions. fq will also automatically add "unknown" fields if
it finds gaps.
- Try to not decode too much as one value.
A length encoded int could be two fields, but maybe a length prefixed string should be one.
Flags can be struct with bit-fields.
- Map as many value as possible to more symbolic values.
- Endian is inherited inside one format decoder, defaults to big endian for new format decoder
- Make sure zero length or no frames found etc fails decoding
- If format is in the probe group make sure to validate input to make it non-ambiguous with other decoders
- Try keep decoder code as declarative as possible
- Split into multiple sub formats if possible. Makes it possible to use them separately.
- Validate/Assert
- Error/Fatal/panic
- Is format probeable or not
- Can new formats be added to other formats
- Does the new format include existing formats

### Development tips

I ususally use `-d <format>` and `v` while developing, that way you will get a decode tree
even if it fails. `v` gives verbose output and also includes stacktrace.

```sh
go run fq.go -d <format> v file
```

If the format is inside some other format it can be handy to first extract the bits and run
the decode directly. For example if working a `aac_frame` decoder issue:

```sh
fq '.tracks[0].samples[1234] | tobytes' file.mp4 > aac_frame_1234
fq -d aac_frame v aac_frame_1234
```

Sometimes nested decoding fails then maybe a good way is to change the parent decoder to
use `d.RawLen()` etc instead of `d.FormatLen()` etc temporary to extract the bits. Hopefully
there will be some option to do this in the future.

When researching or investinging something I can recommend to use `watchexec`, `modd` etc to
make things more comfortable. Also using vscode/delve for debugging should work fine once
launch `args` are setup etc.

```
watchexec "go run fq.go -d aac_frame v aac_frame"
```

Some different ways to run tests:
```sh
# run all tests
make test
# run all go tests
go test ./...
# run all tests for one format
go test -run TestFQTests/mp4 ./format/
# write all actual outputs
make actual
# write actual output for specific tests
WRITE_ACTUAL=1 go run -run ...
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
LOGFILE=/tmp/log go run fq.go ... 2>>/tmp/log
```

gojq execution debug:
```sh
GOJQ_DEBUG=1 go run -tags debug fq.go ...
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
bitio.Reader (implements bitio.Bit* interfaces)
^
| (bitio.Bit* interfaces)
|
bitio.Buffer convenience wrapper to read bytes from bit reader, create section readers etc
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
docker --context 2016-box run --rm -ti -v "C:${PWD//\//\\}:C:${PWD//\//\\}" -w "$PWD" golang:1.17.5-windowsservercore-ltsc2016
```

# Implementation details

- fq uses a gojq fork that can be found at https://github.com/wader/gojq/tree/fq (the "fq" branch)
- fq uses a readline fork that can be found at https://github.com/wader/readline/tree/fq (the "fq" branch)
- cli readline uses raw mode so blocks ctrl-c to become a SIGINT

## Release process

Run and follow instructions:
```
make release=1.2.3
```
