## Basic usage

fq tries to behave the same way as jq as much as possible, so you can do:

```sh
fq . file
fq < file
cat file | fq
fq . < file
fq . *.png *.mp3
fq '.frames[0]' *.mp3
```

### Common usages

```sh
# recursively display decode tree but truncate long arrays
fq d file
# same as
fq display file

# display all bytes for each value
fq 'd({display_bytes: 0})' file
# display 200 bytes for each value
fq 'd({display_bytes: 200})' file

# recursively display decode tree without truncating
fq da file

# recursively and verbosely display decode tree
fq dv file

# JSON repersenation for whole file
fq tovalue file

# recursively look for decode value roots for a format
fq '.. | select(format=="jpeg")' file
# can also use grep_by
fq 'grep_by(format=="jpeg")' file

# recursively look for first decode value root for a format
fq 'first(.. | select(format=="jpeg"))' file
fq 'first(grep_by(format=="jpeg"))' file

# recursively look for objects fullfilling condition
fq '.. | select(.type=="trak")?' file
fq 'grep_by(.type=="trak")' file

# grep whole tree
fq 'grep("^prefix")' file
fq 'grep(123)' file
fq 'grep_by(. >= 100 and . =< 100)' file

# decode file as mp4 and return a result even if there are some errors
fq -d mp4 file.mp4
# decode file as mp4 and also ignore validity assertions
fq -o force=true -d mp4 file.mp4
```

### Display output

`display` or `d` is the main function for displying values and is also the function that will be used if no other output function is explicitly used. If its input is a decode value it will output a dump and tree structure or otherwise it will output as JSON.

Below demonstrates some usages:

First and second example does the same thing, inputs `"hello"` to  `display`.

![fq demo](display_json.svg)

In the next few examples we select out the first "edit list" box in an mp4 file, it's a list of which part of media track to be included during playback, and displays it in various ways.

Default if not explicitly used `display` will only show the root level:

![fq demo](display_decode_value.svg)

First row shows ruler with byte offset into the line and jq path for the value.

The columns are:
- Start address for the line. For example we see that `type` starts at `0xd60`+`0x09`.
- Hex repersenation of input bits for value. Will show the whole byte even if the value only partially uses bits from it.
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

Bit range uses `bytes.bits` notation. For example `type` start at byte `0xd69` bit `0` (left out if zero) and ends at `0xd6c` bit `7` (inclusive) and have byte size of `4`.

There are also some other `display` aliases:
- `da` same as `display({array_truncate: 0})` which will not truncate long arrays.
- `dd` same as `display({array_truncate: 0, display_bytes: 0})` which will not truncate long ranges.
- `dv` same as `display({array_truncate: 0, verbose: true})`
- `ddv` same as `display({array_truncate: 0, display_bytes: 0 verbose: true})` which will not truncate long and also display verbosely.

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

Recursively look for first value that is a `jpeg` decode value root. Use `tobytes` to get bytes for value. Redirect bytes to a file.

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

Common beginner gotcha are:
- jq's use of `;` and `,`. jq uses `;` as argument separator
and `,` as output separator. To call a function `f` with two arguments use `f(1; 2)`. If you do `f(1, 2)` you pass a
single argument `1, 2` (a lambda expression that output `1` and then output `2`) to `f`.
- Expressions can return or "output" zero or more values. This is how loops, foreach etc is
achieved.
- Expressions have one implicit input and output value. This how pipelines like `1 | . * 2` work.


## Types specific to fq

fq has two additional types compared to jq, decode value and binary. In standard jq expressions they will in most case behave as some standard jq type.

### Decode value

This type is returned by decoders and it used to represent parts of the decoed input. It can act as all standard jq types, object, array, number, string etc.

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

The value of a decode value is the symbolic value if available and otherwise the actual value. To explicitly access the value use `tovalue`. In most expression this is not needed as it will be done automactically.

### Binary

Binaries are raw bits with a unit size, 1 (bits) or 8 (bytes), that can have a non-byte aligned size. Will act as byte padded strings in standard jq expressions.

Use `tobits` and `tobytes` to create them from a decode values, strings, numbers or binary arrays. `tobytes` will if needed zero pad most significant bits to be byte aligned.

There is also `tobitsrange` and `tobytesrange` which does the same thing but will preserve it's source range when displayed.

- `"string" | tobytes` produces a binary with UTF8 codepoint bytes.
- `1234 | tobits` produces a binary with the unsigned big-endian integer 1234 with enough bits to represent the number. Use `tobytes` to get the same but with enough bytes to represent the number. This is different to how numbers works inside binary arrays where they are limited to 0-255.
- `["abc", 123, ...]  | tobytes` produce a binary from a binary array. See [binary array](#binary-array) below.
- `.[index]` access bit or byte at index `index`. Index is in units.
  - `[0x12, 0x34, 0x56] | tobytes[1]` is `0x35`
  - `[0x12, 0x34, 0x56] | tobits[3]` is `1`
- `.[start:]`, `.[start:end]` or `.[:end]` is normal jq slice syntax and will slice the binary from `start` to `end`. `start` and `end` is in units.
  - `[0x12, 0x34, 0x56] | tobytes[1:2]` will be a binary with the byte `0x34`
  - `[0x12, 0x34, 0x56] | tobits[4:12]` will be a binary with the byte `0x23`
  - `[0x12, 0x34, 0x56] | tobits[4:20]` will be a binary with the byte `0x23`, `0x45`
  - `[0x12, 0x34, 0x56] | tobits[4:20] | tobytes[1:]` will be a binary with the byte `0x45`,

Both `.[index]` and `.[start:end]` support negative indices to index from end.

TODO: tobytesrange, padding

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

The difference between `tobits` and `tobytes` is

TODO: padding and alignment

## Functions

- All standard library functions from jq
- Adds a few new general functions:
  - `print`, `println`, `printerr`, `printerrln` prints to stdout and stderr.
  - `group` group values, same as `group_by(.)`.
  - `streaks`, `streaks_by(f)` like `group` but groups streaks based on condition.
  - `count`, `count_by(f)` like `group` but counts groups lengths.
  - `debug(f)` like `debug` but uses arg to produce debug message. `{a: 123} | debug({a}) | ...`.
  - `path_to_expr` from `["key", 1]` to `".key[1]"`.
  - `expr_to_path` from `".key[1]"` to `["key", 1]`.
  - `diff($a; $b)` produce diff object between two values.
  - `delta`, `delta_by(f)`, array with difference between all consecutive pairs.
  - `chunk(f)`, split array or string into even chunks
- Bitwise functions `band`, `bor`, `bxor`, `bsl`, `bsr` and `bnot`. Works the same as jq math functions,
unary uses input and if more than one argument all as arguments ignoring the input. Ex: `1 | bnot` `bsl(1; 3)`
- Adds some decode value specific functions:
  - `root` tree root for value
  - `buffer_root` root value of buffer for value
  - `format_root` root value of format for value
  - `parent` parent value
  - `parents` output parents of value
  - `topath` path of value. Use `path_to_expr` to get a string representation.
  - `tovalue`, `tovalue($opts)` symbolic value if available otherwise actual value
  - `toactual` actual value (decoded etc)
  - `tosym` symbolic value (mapped etc)
  - `todescription` description of value
  - `torepr` convert decode value into what it reptresents. For example convert msgpack decode value
  into a value representing its JSON representation.
  - All regexp functions work with binary as input and pattern argument with these differences
  compared to when using string input:
    - All offset and length will be in bytes.
    - For `capture` the `.string` value is a binary.
    - If pattern is a binary it will be matched literally and not as a regexp.
    - If pattern is a binary or flags include "b" each input byte will be read as separate code points
  - String function are not overloaded to support binary for now as some of them are bahaviours that might be confusing.
  - `explode` is overloaded to work with binary. Will explode into array of the unit of the binary.
  end of binary.
  instead of possibly multi-byte UTF-8 codepoints. This allows to match raw bytes. Ex: `match("\u00ff"; "b")`
  will match the byte `0xff` and not the UTF-8 encoded codepoint for 255, `match("[^\u00ff]"; "b")` will match
  all non-`0xff` bytes.
  - `grep` functions take 1 or 2 arguments. First is a scalar to match, where a string is
  treated as a regexp. A binary will be matches exact bytes. Second argument are regexp
  flags with addition that "b" will treat each byte in the input binary as a code point, this
  makes it possible to match exact bytes.
    - `grep($v)`, `grep($v; $flags)` recursively match value and binary
    - `vgrep($v)`, `vgrep($v; $flags)` recursively match value
    - `bgrep($v)`, `bgrep($v; $flags)` recursively match binary
    - `fgrep($v)`, `fgrep($v; $flags)` recursively match field name
  - `grep_by(f)` recursively match using a filter. Ex: `grep_by(. > 180 and . < 200)`, `first(grep_by(format == "id3v2"))`.
  - Binary:
    - `tobits` - Transform input to binary with bit as unit, does not preserving source range, will start at zero.
    - `tobitsrange` - Transform input to binary with bit as unit, preserves source range if possible.
    - `tobytes` - Transform input to binary with byte as unit, does not preserving source range, will start at zero.
    - `tobytesrange` - Transform input binary with byte as unit, preserves source range if possible.
    - `.[start:end]`, `.[:end]`, `.[start:]` - Slice binary from start to end preserving source range.
- `open` open file for reading
- All decode function takes a optional option argument. The only option currently is `force` to ignore decoder asserts.
For example to decode as mp3 and ignore assets do `mp3({force: true})` or `decode("mp3"; {force: true})`, from command line
you currently have to do `fq -d raw 'mp3({force: true})' file`.
- `decode`, `decode("<format>")`, `decode("<format>"; $opts)` decode format
- `probe`, `probe($opts)` probe and decode format
- `mp3`, `mp3($opts)`, ..., `<format>`, `<format>($opts)` same as `decode("<format>")`, `decode("<format>"; $opts)`  decode as format
- Display shows hexdump/ASCII/tree for decode values and jq value for other types.
  - `d`/`d($opts)` display value and truncate long arrays and binaries
  - `da`/`da($opts)` display value and don't truncate arrays
  - `dd`/`dd($opts)` display value and don't truncate arrays or binaries
  - `dv`/`dv($opts)` verbosely display value and don't truncate arrays but truncate binaries
  - `ddv`/`ddv($opts)` verbosely display value and don't truncate arrays or binaries
- `p`/`preview` show preview of field tree
- `hd`/`hexdump` hexdump value
- `repl`/`repl($opts)` nested REPL, must be last in a pipeline. `1 | repl`, can "slurp" outputs. Ex: `1, 2, 3 | repl`, `[1,2,3] | repl({compact: true})`.
- `slurp("<name>")` slurp outputs and save them to `$name`, must be last in pipeline. Will be available as global array `$name`. Ex `1,2,3 | slurp("a")`, `$a[]` same as `spew("a")`.
- `spew`/`spew("<name>")` output previously slurped values. `spew` outputs all slurps as an object, `spew("<name>")` outouts one slurp. Ex: `spew("a")`.
- `paste` read string from stdin until ^D. Useful for pasting text.
    - Ex: `paste | frompem | asn1_ber | repl` read from stdin then decode and start a new sub-REPL with result.

### Encodings, serializations and hashes

In an addition to binary formats fq also support reading to and from encodings and serialization formats.

At the moment fq does not have any dedicated argument for serialization formats but raw string input `-R` slurp `-s` and raw string output `-r` can make things easier. The combination `-Rs` will read all inputs into one string (same as jq).

Note that `from*` functions output jq values and `to*` takes jq values as input so in some cases not all information will properly preserved, for example for XML element and attribute order might change and text and comment nodes might move and will be merged. [yq](https://github.com/mikefarah/yq) might be a better tool if that is needed.

Some example usages:

```sh
# read yml (format is probed, use -d yaml to force) and do some query
$ fq '...' file.yml

# convert YAML to JSON
# note -r for raw string output, without a JSON string with JSON would outputted
$ fq -r 'tojson({indent:2})' file.yml

# add token to URL
$ echo -n "https://host.org" | fq -Rsr 'fromurl | .user.username="token" | tourl'
https://token@host.org

# top 3 hosts in src or href attributes:
# -d to decode as html, can't be probed as html5 parsers always produce some parse tree
# [...] to start collect values into an array
# .. | ."-src"?, ."-href"? | values, recurse and try (?) to get src and href attributes and filter out nulls
# fromurl.host | values, parse as url and filter out those without a host
# count to count unique values, returns [[key, count], ...]
# reverse sort by count and pick first 3
# map [key, count] tuples into {key: key, values: count}
# from_entries, convert into object
$ curl -s https://www.discogs.com/ | fq -d html '[.. | ."-src"?, ."-href"? | values | fromurl.host | values] | count | sort_by(-.[1])[0:3] | map({key: .[0], value: .[1]}) | from_entries'
{
  "blog.discogs.com": 9,
  "st.discogs.com": 10,
  "www.discogs.com": 14
}

# shows how serialization functions can be used on any string, how to transform values and output som other format
# read and decode zip file and start an interactive REPL
$ fq  -i . <(curl -sL https://github.com/stefangabos/world_countries/archive/master.zip)
# select from interesting xml file
zip> .local_files[] | select(.file_name == "world_countries-master/data/countries/en/world.xml").uncompressed | repl
# convert xml into jq value
> .local_files[95].uncompressed string> fromxml | repl
# sort countries by and select the first one
>> object> .countries.country | sort_by(."-name") | first | repl
# see what current input is
>>> object> .
{
  "-alpha2": "af",
  "-alpha3": "afg",
  "-id": "4",
  "-name": "Afghanistan"
}
# remove "-" prefix from keys and convert to YAML and print it
>>> object> with_entries(.key |= .[1:]) | toyaml | print
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

- `fromxml`/`fromxml($opts)` Parse XML into jq value.<br>
  `{seq: true}` preserve element ordering if more than one sibling.<br>
  `{array: true}` use nested `[name, attributes, children]` arrays to represent elements. Attributes will be `null` if none and children will be `[]` if none, this is to make it easier to work it. `toxml` does not require this.<br>
- `fromhtml`/`fromhtml($opts)` Parse HTML into jq value.<br>
  Similar to `fromxml` but parses html5 in non-script mode. Will always have a `html` root with `head` and `body` elements.<br>
  `{array: true}` use nested arrays to represent elements.<br>
  `{seq: true}` preserve element ordering if more than one sibling.<br>
- `toxml`/`toxml($opts})` Serialize jq value into XML.<br>
  `{indent: number}` indent child elements.<br>
  Assumes object representation if input is an object, and nested arrays if input is an array.<br>
  Will automatically add a root `doc` element if jq value has more then one root element.<br>
  If a `#seq` is found on at least one element all siblings will be sort by sequence number. Attributes are always sorted.<br>

  XML elements can be represented as jq value in two ways, as objects (inspired by [mxj](https://github.com/clbanning/mxj) and [xml.com's Converting Between XML and JSON
](https://www.xml.com/pub/a/2006/05/31/converting-between-xml-and-json.html)) or nested arrays. Both representations are lossy and might lose ordering of elements, text nodes and comments. In object representation `fromxml`, `fromhtml` and `toxml` support `{seq:true}` option to parse/serialize `{"#seq"=<number>}` attributes to preserve element sibling ordering.

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
  - Attributes as dash prefixed `-<key>` keys.
  - Text nodes as `#text`.
  - Comment nodes as `#comment` keys.
  - For explicit sibling ordering `#seq` keys with a number, can be negative, assumed zero if missing.
  - Child element with only text as `<name>` key with text as value.
  - Child element with more than just text as `<name>` key with value an object.
  - Multiple child element sibling with same name as `name` key with value as array with strings and objects.
  ```jq
  > $xml | fromxml
  {
    "doc": {
      "child": [
        {
          "-attr": "1"
        },
        {
          "#text": "text",
          "-attr": "2"
        }
      ],
      "other": "text"
    }
  }
  ```

  With nested array representation, an array with these values `["<name>", {attributes...}, [children...]]`
  - Index zero is element name.
  - Optional first object attributes (including `#text` and `#comment` keys).
  - Optional first array are child elements.
  #
  ```jq
  > $xml | fromxml({array: true})
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
  > $xml | fromxml({seq:true})
  {
    "doc": {
      "child": [
        {
          "#seq": 0,
          "-attr": "1"
        },
        {
          "#seq": 1,
          "#text": "text",
          "-attr": "2"
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
  > $xml | fromxml.doc | del(.child) | .new = "abc" | {root: .} | toxml({indent: 2}) | println
  <root>
    <new>abc</new>
    <other>text</other>
  </root>
  ```

JSON and jq-flavoured JSON
- `fromjson` Parse JSON into jq value.
- `tojson`/`tojson($opt)`  Serialize jq value into JSON.<br>
  `{indent: number}` indent array/object values.<br>
- `fromjq` Parse jq-flavoured JSON into jq value.
- `tojq`/`tojq($opt)`  Serialize jq value into jq-flavoured JSON<br>
  `{indent: number}` indent array/object values.<br>
  jq-flavoured JSON has optional key quotes, `#` comments and can have trailing comma in arrays.
- `fromjsonl` Parse JSON lines into jq array.
- `tojsonl` Serialize jq array into JSONL.

YAML
- `fromyaml` Parse YAML into jq value.
- `toyaml`  Serialize jq value into YAML.

TOML
- `fromtoml` Parse TOML into jq value.
- `totoml`  Serialize jq value into TOML.

CSV
- `fromcsv`/`fromcvs($opts)` Parse CSV into jq value.<br>
  `{comma: string}` field separator, default ",".<br>
  `{comment: string}` comment line character, default "#".<br>
  To work with tab separated values you can use `fromcvs({comma: "\t"})` or `fq -d csv -o 'comma="\t"'`
- `tocsv`/`tocsv($opts)` Serialize jq value into CSV.<br>
  `{comma: string}` field separator, default ",".<br>

XML encoding
- `fromxmlentities` Decode XML entities.
- `toxmlentities` Encode XML entities.

URL parts and XML encodings
- `fromurlpath` Decode URL path component.
- `tourlpath` Encode URL path component. Whitespace as %20.
- `fromurlencode` Decode URL query encoding.
- `tourlencode` Encode URL to query encoding. Whitespace as "+".
- `fromurlquery` Decode URL query into object. For duplicates keys value will be an array.
- `tourlquery` Encode objet into query string.
- `fromurl` Decode URL into object.
  ```jq
  > "schema://user:pass@host/path?key=value#fragment" | fromurl
  {
    "fragment": "fragement",
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
- `tourl` Encode object into URL string.
- `fromhex` Decode hexstring to binary.

Binary encodings like hex and base64
- `tohex` Encode binay into hexstring.
- `frombase64`/`frombase64($opts)` Decode base64 encodings into binary.<br>
  `{encoding:string}` encoding variant: `std` (default), `url`, `rawstd` or `rawurl`
- `tobase64`/`tobase64($opts)` Encode binary into base64 encodings.<br>
  `{encoding:string}` encoding variant: `std` (default), `url`, `rawstd` or `rawurl`

Hash functions
- `tomd4` Hash binary using md4.
- `tomd5` Hash binary using md5.
- `tosha1` Hash binary using sha1.
- `tosha256` Hash binary using sha256.
- `tosha512` Hash binary using sha512.
- `tosha3_224` Hash binary using sha3 224.
- `tosha3_256` Hash binary using sha3 256.
- `tosha3_384` Hash binary using sha3 384.
- `tosha3_512` Hash binary using sha3 512.

Text encodings
- `toiso8859_1` Decode binary as ISO8859-1 into string.
- `fromiso8859_1` Encode string as ISO8859-1 into binary.
- `toutf8` Encode string as UTF8 into binary.
- `fromutf8` Decode binary as UTF8 into string.
- `toutf16` Encode string as UTF16 into binary.
- `fromutf16` Decode binary as UTF16 into string.
- `toutf16le` Encode string as UTF16 little-endian into binary.
- `fromutf16le` Decode binary as UTF16 little-endian into string.
- `toutf16be` Encode string as UTF16 big-endian into binary.
- `fromutf16be` Decode binary as UTF16 big-endian into string.

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
- Try include `include "file?";` that don't fail if file is missing.
- Some values can act as a object with keys even when it's an array, number etc.
- There can be keys hidden from `keys` and `[]`.
- Some values are readonly and can't be updated.

## Decoded values

When you decode something you will get a decode value. A decode value work like
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
