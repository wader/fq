# Implementation details

- fq uses a gojq fork that can be found at https://github.com/wader/gojq/tree/fq (the "fq" branch)
- cli readline uses raw mode so blocks ctrl-c to become a SIGINT
- TODO: `scope` and `scopedump` functions used to implement REPL completion
- TODO: Custom object interface used to traverse fq's field tree and to allowing a terse syntax for comparing and working with fields, accessing child fields and special properties like `_range`.

## Decoder implementation help

- Main goal in the end is to produce a tree that is user-friendly and easy to work with.
So there are always excepts to these rules and sometimes it might be better to let the
decoder code be a bit ugly over producing a tree that is hard to understand.

- Try use same names, symbols, constant number base etc as in specification

- TODO: Decode only what you know. If possible let "parent" decide what to do with unknown bits by using `*Decode*Len/Range/Limit`  funcitions

- Use sub decoders if possible for frames, metadata etc. You can pass data between them. Also makes it possible to call them separately

- Try to no decode to much as one field or value. A length encoded int could be two fields, a flags byte can be struct with bit fields.

- Try to have add symbols for all named constants.

## Debug

Send `log` package output and stderr to a file that can be `tail -f`:ed:
```sh
LOGFILE=/tmp/log go run fq.go ... 2>>/tmp/log
```

gojq execution debug:
```sh
GOJQ_DEBUG=1 go run -tags debug fq.go ...
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
