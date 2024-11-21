## Basic usage

fq tries to behave the same way as jq as much as possible, so you can do:

```sh
fq . file
fq < file
cat file | fq
fq . < file
fq . *.png *.mp3
fq '.frames[0]' *.mp3
fq '.frames[-1] | tobytes' file.mp3 > last_frame
```

### Common usages

```sh
# recursively display decode tree but truncate long arrays
fq d file
# same as
fq display file

# display all bytes for each value
fq dd file
# same as
fq 'd({display_bytes: 0})' file

# display 200 bytes for each value
fq 'd({display_bytes: 200})' file

# recursively display decode tree without truncating
fq da file
# same as
fq 'd({array_truncate: 0})' file

# display a specific decode tree one level
fq '.path[1].to.value' file
# display a specific decode tree all levels
fq '.path[1].to.value | d' file
fq '.path[1].to.value | dd' file
fq '.path[1].to.value | da' file

# recursively and verbosely display decode tree
fq dv file
# same as
fq 'd({verbose: true})' file

# JSON representation for whole file
fq tovalue file
# or use -V (--value-output) that does tovalue automatically
fq -V . file
# or -Vr if the value is a string and you want a "raw" string
fq -Vr .path.to.string file
# JSON but raw bit fields truncated
fq -o bits_format=truncate tovalue file
# JSON but raw bit fields as md5 hex string
fq -o bits_format=md5 tovalue file
# JSON but raw bit fields as byte arrays
fq -o bits_format=byte_array tovalue file
# look up a path
fq '.some[1].path' file
# look up a path and output JSON
fq -V '.some[1].path' file
# can be a query that outputs multiple values
# this outputs first and last value in .same array and .path, three values in total
fq -V '.some[0,-1], .path' file

# grep whole tree by value
fq 'grep("^prefix")' file
fq 'grep(123)' file
# grep whole tree by condition
fq 'grep_by(. >= 100 and . =< 100)' file

# recursively look for values fulfilling some condition
fq '.. | select(.type=="trak")?' file
fq 'grep_by(.type=="trak")' file
# grep_by(f) is alias for .. | select(f)?, that is: recurse, select and ignore errors

# recursively look for decode value roots for a format
fq '.. | select(format=="jpeg")' file
# can also use grep_by
fq 'grep_by(format=="jpeg")' file

# recursively look for first decode value root for a format
fq 'first(.. | select(format=="jpeg"))' file
fq 'first(grep_by(format=="jpeg"))' file

# decode file as mp4 and return a result even if there are some errors
fq -d mp4 file.mp4
# decode file as mp4 and also ignore validity assertions
fq -o force=true -d mp4 file.mp4
```

### CLI arguments

Most of jq's CLI arguments work with fq. But here are some additional ones specific to fq:

#### Decode format `--decode`, `-d NAME`

Force format to decode instead of probing.

`NAME` is a name of a format, ex `-d mp4`, see `-h formats` for list of formats.

#### Interactive REPL `--repl`,`-i`

Start interactive REPL.

Can be used with no input, one and multiple inputs, for example just `fq -i ` starts a REPL with `null` input, `fq -i 123` with the number 123 as input, `fq -i . a b` with two files as input. This also works with `--slurp`. In the REPL it is also possible to start a sub-REPLs by ending a query with `<query> | repl`, use ctrl-D to exit the sub-REPL. The sub-REPL will evaluate separately on each output from the query it was started. Use `[<query>] | repl` if you want to "slurp" into an array.

#### Set option `--options`,`-o KEY=VALUE|@PATH`

`KEY` is name of option

`VALUE` will be interpreted as a JSON value if possible otherwise a string, ex -o `name=abc` and `-o name='"abc"'` is the same.

`@PATH` will read string from file at `PATH`.

Specify a global option or a format option, ex: `-o decode_samples=false` would for some container decoders like `mp4` and `matroska` disable decoding of samples.

#### Value output `--value-output`, `-V`

Output JSON value instead of decode tree. Use `-Vr` if you want raw string (no quotes).

### Display output

`display` or `d` is the main function for displaying values and is also the function that will be used if no other output function is explicitly used. If its input is a decode value it will output a dump and tree structure or otherwise it will output as JSON.

Below demonstrates some usages:

First and second example does the same thing, inputs `"hello"` to  `display`.

![fq demo](display_json.svg)

In the next few examples we select out the first "edit list" box in an mp4 file, it's a list of which part of the media track to be included during playback, and displays it in various ways.

Default if not explicitly used `display` will only show the root level:

![fq demo](display_decode_value.svg)

First row shows a ruler with byte offset into the line and jq path for the value.

The columns are:
- Start address for the line. For example we see that `type` starts at `0xd60` (row) + `0x09` (column).
- Hex representation of input bits for value. Will show the whole byte even if the value only partially uses bits from it.
- ASCII representation of input bits for value. Will show the whole byte even if the value only partially uses bits from it.
- Tree structure of decoded value, symbolic value and description.

Notation:
- `{}` value is an object that might have nested values.
- `[start:end]` value is an array with index starting at `start` and ending at `end` (exclusive).


With `display` or `d` it will recursively show the whole tree:

![fq demo](display_decode_value_d.svg)

Same but verbose `dv`:

![fq demo](display_decode_value_dv.svg)

In verbose mode bit ranges and array element names as shown.

Bit ranges uses `<start-byte>[.<bits>]-<end-byte>[.<bits>]` as notation where `.<bits>` is left out if byte aligned. For example `type` starts at byte `0xd69` bit `0` (`.0` is left out) and ends at `0xd6d` bit `0` (exclusive) and have a size of `4` bytes.

This verbosely displays the header of the second frame in an mp3 file which has a bunch of non-byte-aligned fields:

![fq demo](display_decode_value_dv2.svg)

Here the `sync` pattern starts at `0xb79` (bit `0`) and ends at `0xb7a.3` (exclusive) and has a size of `1` byte and `3` bits, `11` bits in total (`8+3`).

There are also some other `display` aliases:
- `da` is `display({array_truncate: 0})` don't truncate arrays.
- `dd` is `display({array_truncate: 0, string_truncate: 0, display_bytes: 0})` don't truncate arrays and raw bytes.
- `dv` is `display({array_truncate: 0, string_truncate: 0, verbose: true})` don't truncate arrays and display verbosely.
- `ddv` is `display({array_truncate: 0, string_truncate: 0, display_bytes: 0 verbose: true})` don't truncate arrays and raw bytes. and display verbosely.

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

# in the REPL you will see a prompt indicating current input and you can type jq expression to evaluate.

# start REPL with one file as input
$ fq -i . doc/file.mp3
mp3>

$ fq -i . doc/file.mp3
# basic arithmetic and jq expressions
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
# actual underlying decoded value
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
# decode byte range 100 to end as mp3_frame
fq -d bytes '.[100:] | mp3_frame | d' file.mp3
# decode byte range 10 bytes from .somefield and preserve relative position in file
fq '.somefield | tobytesrange[10:] | mp3_frame | d' file.mp3
```

#### Show AVC SPS difference between two mp4 files

`-n` tells fq to not have an implicit `input`, `f` is a function to select out some interesting value, call `diff` with two arguments,
decoded value for `a.mp4` and `b.mp4` filtered thru `f`.

```sh
fq -n 'def f: .. | select(format=="avc_sps"); diff(input|f; input|f)' a.mp4 b.mp4
```

#### Extract first JPEG found in file

Recursively look for the first value that is a `jpeg` decode value root. Use `tobytes` to get bytes for value. Redirect bytes to a file.

```sh
fq 'first(.. | select(format=="jpeg")) | tobytes' file > file.jpeg
```

#### Sample size histogram

Recursively look for a all sample size boxes "stsz" and use `?` to ignore errors when doing `.type` on arrays etc. Save reference to box, count unique values, save the max, output the path to the box and output a histogram scaled to 0-100.

```sh
fq '.. | select(.type=="stsz")? as $stsz | .entries | count | max_by(.[1])[1] as $m | ($stsz | topath | path_to_expr), (.[] | "\(.[0]): \((100*.[1]/$m)*"=") \(.[1])") | println' file.mp4
```

#### Find TCP streams that looks like HTTP GET requests in a PCAP file

Use `grep` to recursively find strings matching a regexp.

```sh
fq '.tcp_connections | grep("GET /.* HTTP/1.?")' file.pcap
```

#### Use representation of a format

Some formats like `msgpack`, `bson` etc are used to represent some data structure. In those cases the `torepr`
function can be used to get the representation.

```sh
# whole represented value
fq -d msgpack torepr file.msgpack
# value of the key "field" from the represented value
fq -d msgpack `torepr.field` file.msgpack
# query or transform represented value
fq -d msgpack 'torepr | ...' file.msgpack
```

#### Widest PNG in a directory
```sh
$ fq -rn '[inputs | [input_filename, first(.chunks[] | select(.type=="IHDR") | .width)]] | max_by(.[1]) | .[0]' *.png
```

#### What values include the byte at position 0x123
```sh
$ fq '.. | select(scalars and in_bytes_range(0x123))' file
```

## Support formats

See [formats](formats.md)

## The jq language

fq is based on the [jq language](https://stedolan.github.io/jq/) and for basic usage its syntax
is similar to how object and array access looks in JavaScript or JSON path, `.food[10]` etc. but
it can do much more and is a very expressive language.

To get the most out of fq it's recommended to learn more about jq, here are some good starting points:

- [jq manual](https://stedolan.github.io/jq/manual/)
- [Peter Koppstein's A Stream oriented Introduction to jq](https://github.com/pkoppstein/jq/wiki/A-Stream-oriented-Introduction-to-jq)
- [jq wiki: Language Description](https://github.com/stedolan/jq/wiki/jq-Language-Description)
- [jq wiki: page Cookbook](https://github.com/stedolan/jq/wiki/Cookbook)
- [jq wiki: Pitfalls](https://github.com/stedolan/jq/wiki/How-to:-Avoid-Pitfalls)
- [FAQ](https://github.com/stedolan/jq/wiki/FAQ)

For a more convenient jq experience these might be interesting:

- [jq Dash docset](https://github.com/wader/jq-dash-docset)
- [vscode-jq](https://github.com/wader/vscode-jq)
- [jq-lsp](https://github.com/wader/jq-lsp)

Common beginner gotchas are:
- jq's use of `;` and `,`. jq uses `;` as argument separator
and `,` as output separator. To call a function `f` with two arguments use `f(1; 2)`. If you do `f(1, 2)` you pass a
single argument `1, 2` (a lambda expression that outputs `1` and then outputs `2`) to `f`.
- Expressions can return or "output" zero or more values. This is how loops, foreach etc is
achieved.
- Expressions have one implicit input and output value. This how pipelines like `1 | . * 2` work.

## Differences to jq

- All [gojq's differences to jq](https://github.com/itchyny/gojq#difference-to-jq).
Notably it adds support for arbitrary-precision integers.
- Supports raw strings using back-ticks. <code>\`\\(1234)\u1234\`</code> results in the string `"\\(1234)\\u1234"`.
- Supports hexadecimal `0xab`, octal `0o77` and binary `0b101` integer literals. Also support grouping using underscore `0xaa_bb_cc_dd`.
- Try include using a ending question mark `include "file?";` that doesn't fail if file is missing or has errors.
- Some values can act as an object with keys even when it's an array, number etc.
- There can be keys hidden from `keys` and `[]`.
- Some values are readonly and can't be updated or will convert to JSON on update.

### Types specific to fq

fq has two additional types compared to jq, decode value and binary. In standard jq expressions they will in most cases behave as some standard jq type.

### Decode value

This type is returned by decoders and it is used to represent parts of the decoded input. It can act as all standard jq types, object, array, number, string etc.

Each decode value has these properties:
- A bit range in the input
  - Can be accessed as a binary using `tobits`/`tobytes`. Use the `start` and `size` keys to access position and size.
  - `.name` as bytes `.name | tobytes`
  - Bit 4-8 of `.name` as bits `.name | tobits[4:8]`

Each non-compound decode value has these properties:
- An actual value:
  - This is the decoded representation of the bits, a number, string, bool etc.
  - Can be accessed using `toactual`.
- An optional symbolic value:
  - Is usually a mapping of the actual to symbolic value, ex: map number to a string value.
  - Can be accessed using `tosym`.
- An optional description:
  - Can be accessed using `todescription`
- `parent` is the parent decode value
- `parents` is the all parent decode values
- `topath` is the jq path for the decode value
- `torepr` convert decode value to its representation if possible

The value of a decode value is the symbolic value if available and otherwise the actual value. To explicitly access the value use `tovalue`. In most expressions this is not needed as it will be done automatically.

### Binary

Binaries are raw bits with a unit size, 1 (bits) or 8 (bytes), that can have a non-byte aligned size. Will act as byte padded strings in standard jq expressions.

Use `tobits` and `tobytes` to create them from decode values, strings, numbers or binary arrays. `tobytes` will, if needed zero pad most significant bits to be byte aligned.

There is also `tobitsrange` and `tobytesrange` which does the same thing but will preserve its source range when displayed.

- `"string" | tobytes` produces a binary with UTF8 codepoint bytes.
- `1234 | tobits` produces a binary with the unsigned big-endian integer 1234 with enough bits to represent the number. Use `tobytes` to get the same but with enough bytes to represent the number. This is different to how numbers work inside binary arrays where they are limited to 0-255.
- `["abc", 123, ...]  | tobytes` produce a binary from a binary array. See [binary array](#binary-array) below.
- `.[index]` access bit or byte at index `index`. Index is in units.
  - `[0x12, 0x34, 0x56] | tobytes[1]` is `0x35`
  - `[0x12, 0x34, 0x56] | tobits[3]` is `1`
- `.[start:]`, `.[start:end]` or `.[:end]` is normal jq slice syntax and will slice the binary from `start` to `end`. `start` and `end` is in units.
  - `[0x12, 0x34, 0x56] | tobytes[1:2]` will be a binary with the byte `0x34`
  - `[0x12, 0x34, 0x56] | tobits[4:12]` will be a binary with the byte `0x23`
  - `[0x12, 0x34, 0x56] | tobits[4:20]` will be a binary with the byte `0x23`, `0x45`
  - `[0x12, 0x34, 0x56] | tobits[4:20] | tobytes[1:]` will be a binary with the byte `0x45`,
  - Both `.[index]` and `.[start:end]` support negative indices to index from end.
- `explode` output an array with all byte or bits as integers.

#### Binary array

Is an array of numbers, strings, binaries or other nested binary arrays. When used as input to `tobits`/`tobytes` the following rules are used:
- Number is a byte with value be 0-255
- String it's UTF8 codepoint bytes
- Binary as is
- Binary array used recursively

Binary arrays are similar to and inspired by [Erlang iolist](https://www.erlang.org/doc/man/erlang.html#type-iolist).

Some examples:

`[0, 123, 255] | tobytes` will be binary with 3 bytes 0, 123 and 255

`[0, [123, 255]] | tobytes` same as above

`[0, 1, 1, 0, 0, 1, 1, 0 | tobits]`  will be binary with 1 byte, 0x66 an "f"

`[(.a | tobytes[-10:]), 255, (.b | tobits[:10])] | tobytes` the concatenation of the last 10 bytes of `.a`, a byte with value 255 and the first 10 bits of `.b`.

## Functions

fq has all the same standard library jq functions and in addition some new ones.

### Additional generic functions

#### `grep_by(f)`
Recursively select using a filter. Ex: `grep_by(. > 180 and . < 200)`, `first(grep_by(format == "id3v2"))`.

#### `group`
Group values, same as `group_by(.)`.

#### `streaks`, `streaks_by(f)`
Like `group` but groups streaks based on condition.

#### `count`, `count_by(f)`
Like `group` but counts groups lengths based on condition.

#### `delta`, `delta_by(f)`
Array with difference between consecutive. `delta` is same as `delta_by(.b - .a)`.

#### `chunk($size)`
Split array or string into `$size` length chunks. Last chunk might be shorter.

#### `path_to_expr`
Converts a path value `["key", 1]` to a string `".key[1]"`.

#### `expr_to_path`
Converts from a string `".key[1]"` to path value `["key", 1]`.

#### `diff($a; $b)`
Produce a diff between `$a` and `$b`. Differences are represented as a object `{a: <value from a>, b: <value from b>}`.

#### `band`, `bor`, `bxor`, `bsl`, `bsr`, `bnot`.
Bitwise functions. Works the same as jq math functions. Functions with no arguments like `1 | bnot` uses only input, functions with more than one argument ignores input, `bsl(1; 3)`.

#### `repl`/`repl($opts)`
Nested REPL. Must be last in a pipeline. `1 | repl`, can "slurp" outputs. Ex: `1, 2, 3 | repl`, `[1,2,3] | repl({compact: true})`.

#### `slurp("<name>")`
Slurp outputs and save them to `$name`. Must be last in the pipeline. Will be available as a global array `$name`. Ex `1,2,3 | slurp("a")`, `$a[]` same as `spew("a")`.

#### `spew`/`spew("<name>")`
Output previously slurped values.

#### `spew`
Outputs all slurps as an object. `spew("<name>")` outputs one slurp. Ex: `spew("a")`.

#### `paste`
Read string from stdin until ^D. Useful for pasting text. Ex: `paste | from_pem | asn1_ber | repl` read from stdin then decode and start a new sub-REPL with result.

### Format decode functions

Format decode functions are available in two forms, just `mp3` or `mp3($opts)` that returns a decode value even on error and `from_mp3` or `from_mp3($opts)` which throws error on decode error.

The the only general format option currently is `force` to ignore decoder asserts.
For example to decode as mp3 and ignore assets do `mp3({force: true})` or `decode("mp3"; {force: true})`. From command line you can either do `fq -d mp3 -o force=true . file.mp3` or `fq -d bytes 'mp3({force: true})' file.mp3`.

Some formats has own options that can be specificed as part of `$opts` or as `-o name=value`. Too see options for a format do `fq -h mp3` or `help(mp3)` in a REPL. From command line you can either do `fq -d mp3 -o max_sync_seek=100 . file.mp3` or `fq -d bytes 'mp3({max_sync_seek: 100})' file.mp3`.

#### `decode`, `decode("<format>")`, `decode("<format>"; $opts)`
Decode format.

#### `probe`, `probe($opts)`
Probe and decode format.

#### `<format>`, `<format>($opts)`
Same as `decode("<format>")` and `decode("<format>"; $opts)`. Decode as format and return decode value even on decode error.

#### `from_<format>`, `from_<format>($opts)`
Same as `decode("<format>")` and `decode("<format>"; $opts)` decode as format but throw error on decode error.

Note that jq sometimes uses the notation `name/0`, `name/1` etc in error messages and documentation which means `<function-name>/<arity>`. Same function names with different arity are treated as separate functions, but are usually related in some way in practice.

#### `print`, `println`, `printerr`, `printerrln`
Print string or if not a string compact JSON value to stdout or stderr.

#### `root`
Root decode value for decode value.

#### `buffer_root`
Root decode value of buffer for decode value.

#### `format_root`
Root decode value of format for decode value.

#### `parent`
Parent decode value for decode value.

#### `parents`
Outputs all parent decode values from decode value.

#### `topath`
Path for decode value. Use `path_to_expr` to get a string representation.

#### `tovalue`, `tovalue($opts)`
Symbolic, if available, or actual value for decode value.

#### `toactual`, `toactual($opts)`
Actual value for decode value.

#### `tosym`, `tosym($opts)`
Symbolic value for decode value.

#### `todescription`
Description for decode value.

#### `torepr`
Converts decode value into what it represents. For example converts msgpack decode value into a value representing its JSON representation.

### Display functions

Display shows hexdump, ASCII and tree column dump for decode values and jq value for other types.

#### `d`/`d($opts)`
display value and truncate long arrays and binaries.

#### `da`/`da($opts)`
Display value and don't truncate arrays.

#### `dd`/`dd($opts)`
Display value and don't truncate arrays or binaries.

#### `dv`/`dv($opts)`
Verbosely display value and don't truncate arrays but truncate binaries.

#### `ddv`/`ddv($opts)`
Verbosely display value and don't truncate arrays or binaries.

#### `hd`/`hexdump`
Hexdump value.

### Binary values

Binary values represents raw bits or bytes. When used in standard jq expressions they will behave as strings (UTF-8) with some exceptions listed below.

- All regexp functions work with binary as input and pattern argument with these differences
compared to when using string input:
  - All offset and length will be in bytes.
  - For `capture` the `.string` value is a binary.
  - If pattern is a binary it will be matched literally and not as a regexp.
  - If pattern is a binary or flags include "b" each input byte will be read as separate code points
- `explode` is overloaded to work with binary. Will explode into array of the unit of the binary.
- `.[start:end]`, `.[:end]`, `.[start:]` - Slice binary from start to end preserve source range.

#### `grep($v)`, `grep($v; $flags)`, `vgrep($v)`, `vgrep($v; $flags)`, `bgrep($v)`, `bgrep($v; $flags)`
Recursively match `$v`.

`$v` is a scalar to match, where a string is treated as a regexp. A binary will match exact bytes.
`$flags` argument are regexp flags with additional flag "b" that will treat each byte in the input binary
as a code point. This makes it possible to match exact bytes.

#### `fgrep($v)`, `fgrep($v; $flags)`
Recursively match field name in for decode value.

#### `tobits`
Transform input to binary with bit as unit and don't preserve source range.

#### `tobitsrange`
Transform input to binary with bit as unit and preserve source range.

#### `tobytes`
Transform input to binary with byte as unit and don't preserve source range.

#### `tobytesrange`
Transform input to binary with byte as unit and preserve source range.

#### `open`
Open file for reading.

### Naming inconsistencies

jq's naming conversion is a bit inconsistent, some standard library functions are named `tojson` while others `from_entries`. fq follows this tradition a bit by but tries to use snake_case unless there is a good reason.

Here are all the non-snake_case functions added by. Most of them deal with decode and binary values which are new "primitive" types:
- `toactual`
- `tobits`
- `tobitsrange`
- `tobytes`
- `tobytesrange`
- `todescription`
- `topath`
- `torepr`
- `tosym`
- `tovalue`

### Encodings, serializations and hashes

In addition to binary formats fq also support reading to and from encodings and serialization formats.

At the moment fq does not have any dedicated argument for serialization formats but raw string input `-R` slurp `-s` and raw string output `-r` can make things easier. The combination `-Rs` will read all inputs into one string (same as jq).

Note that `from*` functions output jq values and `to*` takes jq values as input so in some cases not all information will be properly preserved. For example, for the element and attribute order might change and text and comment nodes might move or be merged. [yq](https://github.com/mikefarah/yq) might be a better tool if that is needed.

Some example usages:

```sh
# read yml (format is probed, use -d yaml to force) and do some query
$ fq '...' file.yml

# convert YAML to JSON
# note -r for raw string output, without a JSON string with JSON would outputted
$ fq -r 'tojson({indent:2})' file.yml

# add token to URL
$ echo -n "https://host.org" | fq -Rsr 'from_url | .user.username="token" | tourl'
https://token@host.org

# top 3 hosts in src or href attributes:
# -d to decode as html, can't be probed as html5 parsers always produce some parse tree
# [...] to start collect values into an array
# .. | ."@src"?, ."@href"? | values, recurse and try (?) to get src and href attributes and filter out nulls
# from_url.host | values, parse as url and filter out those without a host
# count to count unique values, returns [[key, count], ...]
# reverse sort by count and pick first 3
# map [key, count] tuples into {key: key, values: count}
# from_entries, convert into object
$ curl -s https://www.discogs.com/ | fq -d html '[.. | ."@src"?, ."@href"? | values | from_url.host | values] | count | sort_by(-.[1])[0:3] | map({key: .[0], value: .[1]}) | from_entries'
{
  "blog.discogs.com": 9,
  "st.discogs.com": 10,
  "www.discogs.com": 14
}

# shows how serialization functions can be used on any string, how to transform values and output some other format
# read and decode zip file and start an interactive REPL
$ fq  -i . <(curl -sL https://github.com/stefangabos/world_countries/archive/master.zip)
# select from interesting xml file
zip> .local_files[] | select(.file_name == "world_countries-master/data/countries/en/world.xml").uncompressed | repl
# convert xml into jq value
> .local_files[95].uncompressed string> from_xml | repl
# sort countries by and select the first one
>> object> .countries.country | sort_by(."@name") | first | repl
# see what current input is
>>> object> .
{
  "@alpha2": "af",
  "@alpha3": "afg",
  "@id": "4",
  "@name": "Afghanistan"
}
# remove "@" prefix from keys and convert to YAML and print it
>>> object> with_entries(.key |= .[1:]) | to_yaml | print
alpha2: af
alpha3: afg
id: "4"
name: Afghanistan
# exit all REPLs back to shell
>>> object> ^D
>> object> ^D
> .local_files[95].uncompressed string> ^D
zip> ^D
```

- `from_xml`/`from_xml($opts)` Parse XML into jq value. `$opts` are:
  - `{seq: true}` preserve element ordering if more than one sibling.<br>
  - `{array: true}` use nested `[name, attributes, children]` arrays to represent elements. Attributes will be `null` if none and children will be `[]` if none, this is to make it easier to work with as the array as 3 values. `to_xml` does not require this.<br>
- `from_html`/`from_html($opts)` Parse HTML into jq value.<br>
  Similar to `from_xml` but parses html5 in non-script mode. Will always have a `html` root with `head` and `body` elements.<br>
  `$opts` are:
  - `{array: true}` use nested arrays to represent elements.<br>
  - `{seq: true}` preserve element ordering if more than one sibling.<br>
- `to_xml`/`to_xml($opts})` Serialize jq value into XML.<br>
  Assumes object representation if input is an object, and nested arrays if input is an array.<br>
  Will automatically add a root `doc` element if jq value has more then one root element.<br>
  If a `#seq` is found on at least one element all siblings will be sort by sequence number. Attributes are always sorted.<br>
  `$opts` are:
  - `{indent: number}` indent child elements.<br>

  XML elements can be represented as jq value in two ways, as objects (inspired by [mxj](https://github.com/clbanning/mxj) and [xml.com's Converting Between XML and JSON
](https://www.xml.com/pub/a/2006/05/31/converting-between-xml-and-json.html)) or nested arrays. Both representations are lossy and might lose ordering of elements, text nodes and comments. In object representation `from_xml`, `from_html` and `to_xml` support `{seq:true}` option to parse/serialize `{"#seq"=<number>}` attributes to preserve element sibling ordering.

  The object version is denser and convenient to query, the nested arrays version is probably easier to use when generating XML.

  Let's assume `$xml` is this XML document as a string:
  ```xml
  <doc>
    <child attr="1"></child>
    <child attr="2">text</child>
    <other>text</other>
  </doc>
  ```

  With object representation an element is represented as:
  - Attributes as `@` prefixed `@<key>` keys.
  - Text nodes as `#text`.
  - Comment nodes as `#comment` keys.
  - For explicit sibling ordering `#seq` keys with a number, can be negative, assumed zero if missing.
  - Child element with only text as `<name>` key with text as value.
  - Child element with more than just text as `<name>` key with value an object.
  - Multiple child element siblings with same name as `name` key with value as array with strings and objects.
  ```jq
  > $xml | from_xml
  {
    "doc": {
      "child": [
        {
          "@attr": "1"
        },
        {
          "#text": "text",
          "@attr": "2"
        }
      ],
      "other": "text"
    }
  }
  ```

  With nested array representation, an array with these values `["<name>", {attributes...}, [children...]]`
  - Index 0 is an element name.
  - Index 1 object attributes (including `#text` and `#comment` keys).
  - Index 2 array of child elements.
  #
  ```jq
  > $xml | from_xml({array: true})
  [
    "doc",
    [
      [
        "child",
        {
          "attr": "1"
        }
      ],
      [
        "child",
        {
          "#text": "text",
          "attr": "2"
        }
      ],
      [
        "other",
        {
          "#text": "text"
        }
      ]
    ]
  ]
  ```
  Parse and include `#seq` attributes if needed:
  ```jq
  > $xml | from_xml({seq:true})
  {
    "doc": {
      "child": [
        {
          "#seq": 0,
          "@attr": "1"
        },
        {
          "#seq": 1,
          "#text": "text",
          "@attr": "2"
        }
      ],
      "other": {
        "#seq": 2,
        "#text": "text"
      }
    }
  }
  ````
  Select values in `<doc>`, remove `<child>`, add a `<new>` element, serialize to xml with 2 space indent and print the string
  ```jq
  > $xml | from_xml.doc | del(.child) | .new = "abc" | {root: .} | to_xml({indent: 2}) | println
  <root>
    <new>abc</new>
    <other>text</other>
  </root>
  ```

JSON and jq-flavoured JSON
- `fromjson` Parse JSON into jq value.
- `tojson`/`tojson($opts)` Serialize jq value into JSON. `$opts` are:
  - `{indent: number}` Indent depth.
- `from_jq` Parse jq-flavoured JSON into jq value.
- `to_jq`/`to_jq($opts)` Serialize jq value into jq-flavoured JSON. jq-flavoured JSON has optional key quotes, `#` comments and can have trailing comma in arrays. `$opts` are:
    - `{indent: number}` Indent depth.
- `from_jsonl` Parse JSON lines into jq array.
- `to_jsonl` Serialize jq array into JSONL.

Note that `fromjson` and `tojson` use different naming conventions as they originate from jq's standard library.

YAML
- `from_yaml` Parse YAML into jq value.
- `to_yaml`/`to_yaml($opts)` Serialize jq value into YAML. `$opts` are:
  - `{indent: number}` Indent depth.

TOML
- `from_toml` Parse TOML into jq value.
- `to_toml`/`to_toml($opts)` Serialize jq value into TOML. `$opts` are:
  - `{indent: number}` Indent depth.

CSV
- `from_csv`/`from_cvs($opts)` Parse CSV into jq value.<br>
  To work with tab separated values you can use `fromcvs({comma: "\t"})` or `fq -d csv -o 'comma="\t"'`<br>
  `$opts` are:
  - `{comma: string}` field separator, default ",".<br>
  - `{comment: string}` comment line character, default "#".<br>
- `to_csv`/`to_csv($opts)` Serialize jq value into CSV. `$opts` are:
  - `{comma: string}` field separator, default ",".<br>

XML encoding
- `from_xmlentities` Decode XML entities.
- `to_xmlentities` Encode XML entities.

URL parts and XML encodings
- `from_urlpath` Decode URL path component.
- `to_urlpath` Encode URL path component. Whitespace as %20.
- `from_urlencode` Decode URL query encoding.
- `to_urlencode` Encode URL to query encoding. Whitespace as "+".
- `from_urlquery` Decode URL query into object. For duplicates keys value will be an array.
- `to_urlquery` Encode object into query string.
- `from_url` Decode URL into object.
  ```jq
  > "schema://user:pass@host/path?key=value#fragment" | from_url
  {
    "fragment": "fragment",
    "host": "host",
    "path": "/path",
    "query": {
      "key": "value"
    },
    "rawquery": "key=value",
    "scheme": "schema",
    "user": {
      "password": "pass",
      "username": "user"
    }
  }
  ```
- `to_url` Encode object into URL string.

Binary encodings like hex and base64
- `from_hex` Decode hex string to binary.
- `to_hex` Encode binary into hex string.
- `from_base64`/`from_base64($opts)` Decode base64 encodings into binary. `$opts` are:
  - `{encoding:string}` encoding variant: `std` (default), `url`, `rawstd` or `rawurl`
- `to_base64`/`to_base64($opts)` Encode binary into base64 encodings. `$opts` are:
  - `{encoding:string}` encoding variant: `std` (default), `url`, `rawstd` or `rawurl`

Hash functions
- `to_md4` Hash binary using md4.
- `to_md5` Hash binary using md5.
- `to_sha1` Hash binary using sha1.
- `to_sha256` Hash binary using sha256.
- `to_sha512` Hash binary using sha512.
- `to_sha3_224` Hash binary using sha3 224.
- `to_sha3_256` Hash binary using sha3 256.
- `to_sha3_384` Hash binary using sha3 384.
- `to_sha3_512` Hash binary using sha3 512.

Text encodings
- `to_iso8859_1` Decode binary as ISO8859-1 into string.
- `from_iso8859_1` Encode string as ISO8859-1 into binary.
- `to_utf8` Encode string as UTF8 into binary.
- `from_utf8` Decode binary as UTF8 into string.
- `to_utf16` Encode string as UTF16 into binary.
- `from_utf16` Decode binary as UTF16 into string.
- `to_utf16le` Encode string as UTF16 little-endian into binary.
- `from_utf16le` Decode binary as UTF16 little-endian into string.
- `to_utf16be` Encode string as UTF16 big-endian into binary.
- `from_utf16be` Decode binary as UTF16 big-endian into string.

## Options

fq has some general options in addition to decode and decoders specific options. They all use the same `-o <name>=<value>` argument.

`<value>` is fuzzily parsed based on the type of the option. Ex: a string can be specified as `-o name=string` or `-o name="string"`.

### `-o bits_format=<string>`

How to represent raw binary as JSON.

- `-o bits_format=base64` Base64 string.
- `-o bits_format=byte_array` Array of bytes (zero bit padded if size is not byte aligned).
- `-o bits_format=hex` Hex string.
- `-o bits_format=md5` MD5 hex string (zero bit padded).
- `-o bits_format=snippet` Truncated Base64 string prefixed with bit length.
- `-o bits_format=string` String with raw bytes (zero bit padded if size is not byte aligned). The string is binary safe internally in fq but bytes not representable as UTF-8 will be lost if turn into JSON (default).
- `-p bits_format=truncate` Truncated string.

```sh
$ fq -V -o bits_format=base64 . file
```
In query
```jq
tovalue({bits_format: "md5"})
```

### `-o skip_gaps=<boolean>`

Skip gaps fields (`gap0` etc) when using `tovalue` or `-V`. Note that this might affect array indexes if one more more gaps fields are skipped in an array.

```sh
$ fq -V -o skip_gaps=true . file
```
In query
```jq
tovalue({skip_gaps: true})
```

### `-o array_truncate=<number>`

By default truncate long array when displaying decode value tree. Use `dd` or `d({array_truncate: 0})` to not truncate.

### `-o string_truncate=<number>`

By default truncate long strings when displaying decode value tree. Use `dd` or `d({string_truncate: 0})` to not truncate.


## Color and unicode output

fq by default tries to use colors if possible, this can be disabled with `-M`. You can also
enable usage of unicode characters for improved output by setting the environment
variable `CLIUNICODE`.

## Configuration

To add own functions you can use `init.fq` that will be read from
- `$HOME/Library/Application Support/fq` or `$HOME/.config/fq` on macOS
- `$HOME/.config/fq` on Linux, BSD etc
- `%AppData%` on Windows

## Use as script interpreter

fq can be used as a script interpreter:

`mp3_duration.jq`:
```jq
#!/usr/bin/env fq -d mp3 -rf
[.frames[].header | .sample_count / .sample_rate] | add
```

## Decoded values

When decoding something, using `decode` or `mp3` etc, you a decode value is returned. They behave like
normal jq values but has special abilities and is used to represent the decoded structure. Each value
always has a name, type and a bit range.

A value has these special keys (TODO: remove, are internal)

- `_actual` decoded (not symbol mapped value)
- `_bits` bits in range as a binary
- `_buffer_root` first decode value for current buffer
- `_bytes` bits in range as binary using byte units
- `_description` description of value (optional)
- `_error` error message (optional)
- `_format` name of decoded format (optional, only format root)
- `_format_root` first decode value for current format
- `_gap` is a bit range gap (was not decoded)
- `_index` index in parent array (only for values in arrays)
- `_len` bit range length (TODO: rename)
- `_name` name of value
- `_out` decoded out value
- `_parent` parent decode value
- `_path` jq path to decode value
- `_root` root decode value
- `_start` bit range start
- `_stop` bit range stop
- `_sym` symbolic value (optional)

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
you know the format of some gap field. Then you can decode manually.

<pre>
# try decode a `mp3_frame` that failed to decode
$ fq -d mp3 '.gap0 | mp3_frame' file.mp3
# skip first 10 bytes then decode as `mp3_frame`
$ fq -d bytes '.[10:] | mp3_frame' file.mp3
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
