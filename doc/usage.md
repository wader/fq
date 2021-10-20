## Basic usage

fq tries to behave the same way as jq as much as possible, so you can do:
```
fq . file.mp3
fq < file.mp3
fq . < file.mp3
fq . *.png *.jpg
fq '.frames[0]' file.mp3
```

## Interactive REPL

fq has an interactive [REPL](https://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop)
with auto completion and nested REPL support:

```
# start REPL with null input
fq -i
# same as
fq -ni
# start REPL with one file as input
fq -i . file.mp3
```

In the REPL you will see a prompt indicating current input and you can type jq expression to evaluate.

```
# basic arithmetics
mp3> 1+1
2
# "." is the identity function, returns current input, the mp3 file.
mp3> .
# access the first frame in the mp3 file
mp3> .frames[0]
# start a new nested REPl with first frame as input
mp3> .frames[0] | repl
# prompt shows "path" to current input and that it's an mp3_frame.
# do Ctrl-D to exit REPL
> .frames[0] mp3_frame> ^D
# do Ctrl-D to exit to shell
mp3> ^D
$
```

Use Ctrl-D to exits, Ctrl-C to interrupt current evaluation.

## The jq langauge

fq is based on the [jq language](https://stedolan.github.io/jq/) and for basic usage its syntax
is similar to how object and array access looks in JavaScript or JSON path, `.food[10]` etc.

To get the most out of fq it's recommended to learn more about jq, here are some good starting points:

- [jq manual](https://stedolan.github.io/jq/manual/)
- jq wiki pages
[jq Language Description](https://github.com/stedolan/jq/wiki/jq-Language-Description),
[jq wiki page Cookbook](https://github.com/stedolan/jq/wiki/Cookbook),
[FAQ](https://github.com/stedolan/jq/wiki/FAQ) and
[Pitfalls](https://github.com/stedolan/jq/wiki/How-to:-Avoid-Pitfalls)

The most common beginner gotcha is probably jq's use of `;` and `,`. jq uses `;` as argument separator
and `,` as output separator.
To call a function `f` with two arguments use `f(1; 2)`. If you do `f(1, 2)` you pass a single
argument `1, 2` (a lambda expression that output `1` and then output `2`) to `f`.

## Support formats

See [formats](formats.md)

## Arguments

TODO: examples, stdin/stdout

<pre sh>
$ fq -h 
fq - jq for files
Tool, language and format decoders for exploring binary data.
For more information see https://github.com/wader/fq

Usage: fq [OPTIONS] [--] [EXPR] [FILE...]
--arg NAME VALUE         Set variable $NAME to string VALUE
--argjson NAME JSON      Set variable $NAME to JSON
--color-output,-C        Force color output
--compact-output,-c      Compact output
--decode,-d NAME         Decode format (probe)
--decode-file NAME PATH  Set variable $NAME to decode of file
--formats                Show supported formats
--from-file,-f PATH      Read EXPR from file
--help,-h                Show help
--include-path,-L PATH   Include search path
--join-output,-j         No newline between outputs
--monochrome-output,-M   Force monochrome output
--null-input,-n          Null input (use input/0 and inputs/0 to read input)
--null-output,-0         Null byte between outputs
--option,-o KEY=VALUE    Set option, eg: color=true (use options/0 to see all options)
--raw-input,-R           Read raw input strings (don't decode)
--raw-output,-r          Raw string output (without quotes)
--rawfile NAME PATH      Set variable $NAME to string content of file
--repl,-i                Interactive REPL
--slurp,-s               Read (slurp) all inputs into an array
--version,-v             Show version
</pre>

## Use as script interpreter

fq can be used as a scrip interpreter:

`mp3_duration.jq`:
```jq
#!/usr/bin/env fq -d mp3 -rf
[.frames[].header | .sample_count / .sample_rate] | add
```

## Differences to jq

- [gojq's differences to jq](https://github.com/itchyny/gojq#difference-to-jq),
notable is support for arbitrary-precision integers.
- Supports hexdecimal `0xab`, octal `0o77` and binary `0b101` integer literals.
- Has bitwise operations, `band`, `bor`, `bxor`, `bsl`, `bsr`, `bnot`.
- Try include `include "file?";` that don't fail if file is missing.
- Some values can act as a object with keys even when it's an array, number etc.
- There can be keys hidden from `keys` and `[]`. Used for, `_format`, `_bytes` etc.
- Some values are readonly and can't be updated.

## Functions

- All standard library functions from jq
- Adds a few new general functions:
  - `streaks/0`, `streaks_by/1` like `group` but groups streaks based on condition.
  - `count`, `count_by/1` like `group` but counts groups lengths.
  - `debug/1` like `debug/0` but uses arg to produce debug message. `{a: 123} | debug({a}) | ...`.
  - `path_to_expr` from `["key", 1]` to `".key[1]"`.
  - `expr_to_path` from `".key[1]"` to `["key", 1]`.
  - `diff/2` produce diff object between two values.
  - `delta`, `delta_by/1`, array with difference between all consecutive pairs.
  - `chunk/1`, split array or string into even chunks
- Adds some decode value specific functions:
  - `root/0` return tree root for value
  - `buffer_root/0` return root value of buffer for value
  - `format_root/0` return root value of format for value
  - `parent/0` return parent value
  - `parents/0` output parents of value
  - All `match` and `grep` functions take 1 or 2 arguments. First is a scalar to match, where a string is
  treated as a regexp. A buffer scalar will be matches exact bytes. Second argument are regexp
  flags with addition that "b" will treat each byte in the input buffer as a code point, this
  makes it possible to match exact bytes, ex: `match("\u00ff"; b")` will match the byte `0xff` and not
  the UTF-8 encoded codepoint for 255.
    - `match/1`, `match/2` overloaded to support buffers. Match in buffer and output match buffers
    - `grep/1`, `grep/2` recursively match value and buffer
    - `vgrep/1`, `vgrep/2` recursively match value
    - `bgrep/1`, `bgrep/2` recursively match buffer
    - `fgrep/1`, `fgrep/2` recursively match field name
  - Buffers:
    - `tobits` - Transform input into a bits buffer not preserving source range, will start at zero.
    - `tobitsrange` - Transform input into a bits buffer preserving source range if possible.
    - `tobytes` - Transform input into a bytes buffer not preserving source range, will start at zero.
    - `tobytesrange` - Transform input into a byte buffer preserving source range if possible.
    - `buffer[start:end]`, `buffer[:end]`, `buffer[start:]` - Create a sub buffer from start to end in buffer units preserving source range.
- `open` open file for reading
- `probe` or `decode` probe format and decode
- `mp3`, `matroska`, ..., `<name>`, `decode([name])` force decode as format
- `d`/`display` display value and truncate long arrays
- `f`/`full` display value and don't truncate arrays
- `v`/`verbose` display value verbosely and don't truncate array
- `p`/`preview` show preview of field tree
- `hd`/`hexdump` hexdump value
- `repl` nested REPL, must be last in a pipeline. `1 | repl`, can "slurp" multiple outputs `1, 2, 3 | repl`.

## Decoded values (TODO: better name?)

When you decode something you will get a decode value. A decode values work like
normal jq values but has special abilities and is used to represent a tree structure of the decoded
binary data. Each value always has a name, type and a bit range.

A value has these special keys:

- `_name` name of value
- `_value` jq value of value
- `_start` bit range start
- `_stop` bit range stop
- `_len` bit range length (TODO: rename)
- `_bits` bits in range as a binary
- `_bytes` bits in range as binary using byte units
- `_path` jq path to value
- `_unknown` value is un-decoded gap
- `_symbol` symbolic string representation of value (optional)
- `_description` longer description of value (optional)
- `_format` name of decoded format (optional)
- `_error` error message (optional)

- TODO: unknown gaps

## Binary and IO lists

- TODO: similar to erlang io lists, [], binary, string (utf8) and numbers

## Configuration

To add own functions you can use `init.fq` that will be read from
- `$HOME/Library/Application Support/fq/init.jq` on macOS
- `$HOME/.config/fq/init.jq` on Linux, BSD etc
- `%AppData%\fq\init.jq` on Windows (TODO: not tested)

## Own decoders and use as library

TODO


## Known issues and useful tricks

### Run interactive mode with no input
```sh
fq -i
null>
```

### `select` fails with `expected an ... but got: ...`

Try add `select(...)?` to catch and ignore type errors in the select expression.

### Manual decode

Sometimes fq fails to decode or you know there is valid data buried inside some binary or maybe
you know the format of some unknown value. Then you can decode manually.

<pre>
# try decode a `mp3_frame` that failed to decode
$ fq file.mp3 .unknown0 mp3_frame
# skip first 10 bytes then decode as `mp3_frame`
$ fq file.mp3 .unknown0._bytes[10:] mp3_frame
</pre>

### Use `.` as input and in a positional argument

The expression `.a | f(.b)` might not work as expected. `.` is `.a` when evaluating the arguments so
the positional argument will end up being `.a.b`. Instead do `. as $c | .a | f($c.b)`.

### Building array is slow

Try to use `map` or `foreach` to avoid rebuilding the whole array for each append.

### Use `print` and `println` to produce more friendly compact output

```
> [[0,"a"],[1,"b"]]
[
  [
    0,
    "a"
  ],
  [
    1,
    "b"
  ]
]
> [[0,"a"],[1,"b"]] | .[] | "\(.[0]): \(.[1])" | println
0: a
1: b
```

### `repl` argument using function or variable causes `variable not defined`

`true as $verbose | repl({verbose: $verbose})` will currently fail as `repl` is
implemented by rewriting the query to  `map(true as $verbose | .) | repl({verbose: $verbose})`.

### `error` produces no output

`null | error` behaves as `empty`.
