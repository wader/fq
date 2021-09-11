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
LOGFILE=/tmp/log go run main.go ... 2>>/tmp/log
```

gojq execution debug:
```sh
GOJQ_DEBUG=1 go run -tags debug main.go ...
```
