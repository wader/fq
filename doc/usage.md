## Basic usage

fq tries to behave the same way as jq as much as possible, so you can do:
```sh
fq . file
fq < file
file | fq
fq . < file
fq . *.png *.mp3
fq '.frames[0]' file.mp3
```

Common usages:
```sh
# recursively display decode tree but truncate long arrays
fq d file
fq display file

# display all or more bytes for each value
fq 'd({display_bytes: 0})' file
fq 'd({display_bytes: 200})' file

# recursively display decode tree
fq f file
fq full file

# recursively verbosely display decode tree
fq v file
fq verbose file

# JSON for whole file
fq tovalue file

# recursively look for decode value roots for a format
fq '.. | select(format=="jpeg")' file

# recursively look for first decode value root for a format
fq 'first(.. | select(format=="jpeg"))' file

# recursively look for objects fullfilling condition
fq '.. | select(.type=="trak")?' file
```

## Interactive REPL

The interactive [REPL](https://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop)
has auto completion and nested REPL support:

```
# start REPL with null input
$ fq -i
null>
# same as
$ fq -ni
null>
# start REPL with one file as input
$ fq -i . doc/file.mp3
mp3>
```

In the REPL you will see a prompt indicating current input and you can type jq expression to evaluate.

```
$ fq -i . doc/file.mp3
# basic arithmetics and jq expressions
mp3> 1+1
2
mp3> 1, 2, 3 | . * 2
2
4
6
mp3> [1, 2, 3] | add
6
# "." is the identity function which just returns current input, the mp3 file.
mp3> .
# access the first frame in the mp3 file
mp3> .frames[0]
# start a new nested REPL with first frame as input
mp3> .frames[0] | repl
# prompt shows "path" to current input and that it's an mp3_frame.
# Ctrl-D to exit REPL or to shell if last REPL
> .frames[0] mp3_frame> ^D
# "jq" value of layer in first frame
mp3> .frames[0].header.layer | tovalue
3
mp3> .frames[0].header.layer * 2
6
# symbolic value, same as "jq" value
mp3> .frames[0].header.layer | tosym
3
# actual underlaying decoded value
mp3> .frames[0].header.layer | toactual
1
# description of value
mp3> .frames[0].header.layer | todescription
"MPEG Layer 3"
mp3> ^D
$
```

Use Ctrl-D to exit and Ctrl-C to interrupt current evaluation.

## Example usages

#### Second mp3 frame header as JSON
```sh
fq '.frames[1].header | tovalue' file.mp3
```

#### Byte start position for the first 10 mp3 frames in an array
```sh
fq '.frames[0:10] | map(tobytesrange.start)' file.mp3
```

#### Decode at range
```sh
# decode byte range 100 to end
fq -d raw 'tobytes[100:] | mp3_frame | d' file.mp3
# decode byte range 10 bytes into .somefield and preseve relative position in file
fq '.somefield | tobytesrange[10:] | mp3_frame | d' file.mp3
```

#### Show AVC SPS difference between two mp4 files

`-n` tells fq to not have an implicit `input`, `f` is function to select out some interesting value, call `diff` with two arguments,
decoded value for `a.mp4` and `b.mp4` filtered thru `f`.

```sh
fq -n 'def f: .. | select(format=="avc_sps"); diff(input|f; input|f)' a.mp4 b.mp4
```
#### Extract first JPEG found in file

Recursively look for first value that is a `jpeg` decode value root. Use `tobytes` to get bytes buffer for value. Redirect bytes to a file.

```sh
fq 'first(.. | select(format=="jpeg")) | tobytes' file > file.jpeg
```

#### Sample size histogram

Recursively look for a all sample size boxes "stsz" and use `?` to ignore errors when doing `.type` on arrays etc. Save reference to box, count unique values, save the max, output the path to the box and output a historgram scaled to 0-100.

```sh
fq '.. | select(.type=="stsz")? as $stsz | .entries | count | max_by(.[1])[1] as $m | ($stsz | topath | path_to_expr), (.[] | "\(.[0]): \((100*.[1]/$m)*"=") \(.[1])") | println' file.mp4
```

#### Find TCP streams that looks like HTTP GET requests in a PCAP file

Use `grep` to recursively find strings matching a regexp.

```sh
fq '.tcp_connections | grep("GET /.* HTTP/1.?")' file.pcap
```

###

#### Widest PNG in a directory
```sh
$ fq -rn '[inputs | [input_filename, first(.chunks[] | select(.type=="IHDR") | .width)]] | max_by(.[1]) | .[0]' *.png
```

#### What values include the byte at position 0x123
```sh
$ fq '.. | select(scalars and in_bytes_range(0x123))' file
```

## The jq language

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
fq - jq for binary formats
Tool, language and decoders for inspecting binary data.
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
--raw-file NAME PATH     Set variable $NAME to string content of file
--raw-input,-R           Read raw input strings (don't decode)
--raw-output,-r          Raw string output (without quotes)
--repl,-i                Interactive REPL
--slurp,-s               Read (slurp) all inputs into an array
--version,-v             Show version
</pre>

## Color and unicode output

fq by default tries to use colors if possible, this can be disabled with `-M`. You can also
enable useage of unicode characters for improved output by setting the environment
variable `CLIUNICODE`.

## Configuration

To add own functions you can use `init.fq` that will be read from
- `$HOME/Library/Application Support/fq/init.jq` on macOS
- `$HOME/.config/fq/init.jq` on Linux, BSD etc
- `%AppData%\fq\init.jq` on Windows

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
- Has bitwise operators, `band`, `bor`, `bxor`, `bsl`, `bsr`, `bnot`.
- Try include `include "file?";` that don't fail if file is missing.
- Some values can act as a object with keys even when it's an array, number etc.
- There can be keys hidden from `keys` and `[]`.
- Some values are readonly and can't be updated.

## Functions

- All standard library functions from jq
- Adds a few new general functions:
  - `print/0`, `println/0`, `printerr/0`, `printerrln` prints to stdout and stderr.
  - `streaks/0`, `streaks_by/1` like `group` but groups streaks based on condition.
  - `count/0`, `count_by/1` like `group` but counts groups lengths.
  - `debug/1` like `debug/0` but uses arg to produce debug message. `{a: 123} | debug({a}) | ...`.
  - `path_to_expr/0` from `["key", 1]` to `".key[1]"`.
  - `expr_to_path/0` from `".key[1]"` to `["key", 1]`.
  - `diff/2` produce diff object between two values.
  - `delta/0`, `delta_by/1`, array with difference between all consecutive pairs.
  - `chunk/1`, split array or string into even chunks
- Adds some decode value specific functions:
  - `root/0` tree root for value
  - `buffer_root/0` root value of buffer for value
  - `format_root/0` root value of format for value
  - `parent/0` parent value
  - `parents/0` output parents of value
  - `topath/0` path of value. Use `path_to_expr` to get a string representation.
  - `tovalue/0`, `tovalue/1` symbolic value if available otherwise actual value
  - `toactual/0` actual value (decoded etc)
  - `tosym/0` symbolic value (mapped etc)
  - `todescription/0` description of value
  - All regexp functions work with buffers as input and pattern argument with these differences
  from the string versions:
    - All offset and length will be in bytes.
    - For `capture` the `.string` value is a buffer.
    - If pattern is a buffer it will be matched literally and not as a regexp.
    - If pattern is a buffer or flags include "b" each input byte will be read as separate code points
  - `scan_toend/1`, `scan_toend/2` works the same as `scan` but output buffer are from start of match to
  end of buffer.
  instead of possibly multi-byte UTF-8 codepoints. This allows to match raw bytes. Ex: `match("\u00ff"; "b")`
  will match the byte `0xff` and not the UTF-8 encoded codepoint for 255, `match("[^\u00ff]"; "b")` will match
  all non-`0xff` bytes.
  - `grep` functions take 1 or 2 arguments. First is a scalar to match, where a string is
  treated as a regexp. A buffer scalar will be matches exact bytes. Second argument are regexp
  flags with addition that "b" will treat each byte in the input buffer as a code point, this
  makes it possible to match exact bytes.
    - `grep/1`, `grep/2` recursively match value and buffer
    - `vgrep/1`, `vgrep/2` recursively match value
    - `bgrep/1`, `bgrep/2` recursively match buffer
    - `fgrep/1`, `fgrep/2` recursively match field name
  - Buffers:
    - `tobits/0` - Transform input into a bits buffer not preserving source range, will start at zero.
    - `tobitsrange/0` - Transform input into a bits buffer preserving source range if possible.
    - `tobytes/0` - Transform input into a bytes buffer not preserving source range, will start at zero.
    - `tobytesrange/0` - Transform input into a byte buffer preserving source range if possible.
    - `buffer[start:end]`, `buffer[:end]`, `buffer[start:]` - Create a sub buffer from start to end in buffer units preserving source range.
- `open` open file for reading
- All decode function takes a optional option argument. The only option currently is `force` to ignore decoder asserts.
For example to decode as mp3 and ignore assets do `mp3({force: true})` or `decode("mp3"; {force: true})`, from command line
you currently have to do `fq -d raw 'mp3({force: true})' file`.
- `decode/0`, `decode/1`, `decode/2` decode format
- `probe/0`, `probe/1` probe and decode format
- `mp3/0`, `mp3/1`, ..., `<name>/0`, `<name>/1` same as `decode(<name>)/1`, `decode(<name>; <opts>)/2`  decode as format
- Display shows hexdump/ASCII/tree for decode values and JSON for other values.
  - `d/0`/`display/0` display value and truncate long arrays
  - `f/0`/`full/0` display value and don't truncate arrays
  - `v/0`/`verbose/0` display value verbosely and don't truncate array
- `p/0`/`preview/0` show preview of field tree
- `hd/0`/`hexdump/0` hexdump value
- `repl/0` nested REPL, must be last in a pipeline. `1 | repl`, can "slurp" multiple outputs `1, 2, 3 | repl`.

## Decoded values

When you decode something you will get a decode value. A decode values work like
normal jq values but has special abilities and is used to represent a tree structure of the decoded
binary data. Each value always has a name, type and a bit range.

A value has these special keys (TODO: remove, are internal)

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
$ fq -d mp3 '.unknown0 | mp3_frame' file.mp3
# skip first 10 bytes then decode as `mp3_frame`
$ fq -d raw 'tobytes[10:] | mp3_frame' file.mp3
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
