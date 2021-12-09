include "internal";
include "options";
include "binary";

def _display_default_opts:
  options({depth: 1});

def display($opts):
  ( options($opts) as $opts
  | if _can_display then _display($opts)
    else
      ( if type == "string" and $opts.raw_string then print
        else _print_color_json($opts)
        end
      , ( $opts.join_string
        | if . then print else empty end
        )
      )
    end
  | error("unreachable")
  );
def display: display({});


def hexdump($opts): _hexdump(options({display_bytes: 0} + $opts));
def hexdump: hexdump({display_bytes: 0});
def hd($opts): hexdump($opts);
def hd: hexdump;

def intdiv(a; b): _intdiv(a; b);

def trim: capture("^\\s*(?<str>.*?)\\s*$"; "").str;

# does +1 and [:1] as " "*0 is null
def rpad($s; $w): . + ($s * ($w+1-length))[1:];

# like group but groups streaks based on condition
def streaks_by(f):
  ( . as $a
  | length as $l
  | if $l == 0 then []
    else
      ( [ foreach $a[] as $v (
            {cf: ($a[0] | f), index: 0, start: 0, extract: null};
            ( ($v | f) as $vf
            | (.index == 0 or (.cf == $vf)) as $equal
            | if $equal then
                ( .extract = null
                )
              else
                ( .cf = $vf
                | .extract = [.start, .index]
                | .start = .index
                )
              end
            | .index += 1
            );
            ( if .extract then .extract else empty end
            , if .index == $l then [.start, .index] else empty end
            )
          )
        ]
      | map($a[.[0]:.[1]])
      )
    end
  );
# [1, 2, 2, 3] => [[1], [2, 2], [3]]
def streaks: streaks_by(.);

# same as group_by but counts, array or pairs with [value, count]
def count_by(exp):
  group_by(exp) | map([(.[0] | exp), length]);
def count: count_by(.);

# array of result of applying f on all consecutive pairs
def delta_by(f):
  ( . as $a
  | if length < 1 then []
    else
      [ range(length-1) as $i
      | {a: $a[$i], b: $a[$i+1]}
      | f
      ]
    end
  );
# array of minus between all consecutive pairs
def delta: delta_by(.b - .a);

# split array or string into even chunks, except maybe the last
def chunk($size):
  if length == 0 then []
  else
    [ ( range(
          ( (length / $size)
          | ceil
          | if . == 0 then 1 end
          )
        ) as $i
      | .[$i * $size:($i + 1) * $size]
      )
    ]
  end;

# [{a: 123, ...}, ...]
# colmap maps something into [col, ...]
# render maps [{column: 0, string: "coltext", maxwidth: 12}, ..] into a row
def table(colmap; render):
  def _column_widths:
    [ . as $rs
    | range($rs[0] | length) as $i
    | [$rs[] | colmap | (.[$i] | length)]
    | max
    ];
  if length == 0 then ""
  else
    ( _column_widths as $cw
    | . as $rs
    | ( $rs[]
      | . as $r
      | [ range($r | length) as $i
        | ($r | colmap | {column: $i, string: .[$i], maxwidth: $cw[$i]})
        ]
      | render
      )
    )
  end;

def fromradix($base; $table):
  ( if type != "string" then error("cannot fromradix convert: \(.)") end
  | split("")
  | reverse
  | map($table[.])
  | if . == null then error("invalid char \(.)") end
  # state: [power, ans]
  | reduce .[] as $c ([1,0];
      ( (.[0] * $base) as $b
      | [$b, .[1] + (.[0] * $c)]
      )
    )
  | .[1]
  );
def fromradix($base):
  fromradix($base; {
    "0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
    "a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15, "g": 16,
    "h": 17, "i": 18, "j": 19, "k": 20, "l": 21, "m": 22, "n": 23,
    "o": 24, "p": 25, "q": 26, "r": 27, "s": 28, "t": 29, "u": 30,
    "v": 31, "w": 32, "x": 33, "y": 34, "z": 35,
    "A": 36, "B": 37, "C": 38, "D": 39, "E": 40, "F": 41, "G": 42,
    "H": 43, "I": 44, "J": 45, "K": 46, "L": 47, "M": 48, "N": 49,
    "O": 50, "P": 51, "Q": 52, "R": 53, "S": 54, "T": 55, "U": 56,
    "V": 57, "W": 58, "X": 59, "Y": 60, "Z": 61,
    "@": 62, "_": 63,
  });

def toradix($base; $table):
  ( if type != "number" then error("cannot toradix convert: \(.)") end
  | if . == 0 then "0"
    else
      ( [ recurse(if . > 0 then intdiv(.; $base) else empty end) | . % $base]
      | reverse
      | .[1:]
      | if $base <= ($table | length) then
          map($table[.]) | join("")
        else
          error("base too large")
        end
      )
    end
  );
def toradix($base):
  toradix($base; "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ@_");

# TODO: rename keys and add more, ascii/utf8/utf16/codepoint name?, le/be, signed/unsigned?
def iprint:
  {
    bin: "0b\(toradix(2))",
    oct: "0o\(toradix(8))",
    dec: "\(.)",
    hex: "0x\(toradix(16))",
    str: (try ([.] | implode) catch null),
  };

# produce a/b pairs for diffing values
def diff($a; $b):
  ( ( $a | type) as $at
  | ( $b | type) as $bt
  | if $at != $bt then {a: $a, b: $b}
    elif ($at == "array" or $at == "object") then
      ( [ ((($a | keys) + ($b | keys)) | unique)[] as $k
        | {
          ($k | tostring):
            ( [($a | has($k)), ($b | has($k))]
            | if . == [true, true] then diff($a[$k]; $b[$k])
              elif . == [true, false] then {a: $a[$k]}
              elif . == [false, true] then {b: $b[$k]}
              else empty # TODO: can't happen? error?
              end
            )
          }
        ]
      | add
      | if . == null then empty end
      )
    else
      if $a == $b then empty else {a: $a, b: $b} end
    end
  );

# https://en.wikipedia.org/wiki/Privacy-Enhanced_Mail
# TODO: add test
def frompem:
  ( tobytes
  | tostring
  | capture("-----BEGIN(.*?)-----(?<s>.*?)-----END(.*?)-----"; "mg").s
  | base64
  ) // error("no pem header or footer found");

def topem($label):
  ( tobytes
  | base64
  | ($label | if $label != "" then " " + $label end) as $label
  | [ "-----BEGIN\($label)-----"
    , .
    , "-----END\($label)-----"
    , ""
    ]
  | join("\n")
  );
def topem: topem("");

def paste:
  if _is_completing | not then
    ( [ _repeat_break(
          try _stdin(64*1024)
          catch if . == "eof" then error("break") end
        )
      ]
    | join("")
    )
  end;

def tojq($style):
  def _is_ident: test("^[a-zA-Z_][a-zA-Z_0-9]*$");
  def _key: if _is_ident | not then tojson end;
  def _f($style):
    def _r($indent):
      ( type as $t
      | if $t == "null" then tojson
        elif $t == "string" then tojson
        elif $t == "number" then tojson
        elif $t == "boolean" then tojson
        elif $t == "array" then
          [ "[", $style.compound_newline
          , ( [ .[]
              | $indent, $style.indent
              , _r($indent+$style.indent), $style.array_sep
              ]
            | .[0:-1]
            )
          , $style.compound_newline
          , $indent, "]"
          ]
        elif $t == "object" then
          [ "{", $style.compound_newline
          , ( [ to_entries[]
              | $indent, $style.indent
              , (.key | _key), $style.key_sep
              , (.value | _r($indent+$style.indent)), $style.value_sep
              ]
            | .[0:-1]
            )
          , $style.compound_newline
          , $indent, "}"
          ]
        else error("unknown type \($t)")
        end
      );
    _r("");
  ( {
      compact: {
        indent: "",
        key_sep: ":",
        value_sep: ",",
        array_sep: ",",
        compound_newline: "",
      },
      fancy_compact: {
        indent: "",
        key_sep: ": ",
        value_sep: ", ",
        array_sep: ", ",
        compound_newline: "",
      },
      verbose: {
        indent: "  ",
        key_sep: ": ",
        value_sep: ",\n",
        array_sep: ",\n",
        compound_newline: "\n",
      }
    } as $styles
  | _f(
      ( $style // "compact"
      | if type == "string" then $styles[.]
        elif type == "object" then .
        else error("invalid style")
        end
      )
    )
  | flatten
  | join("")
  );
def tojq: tojq(null);

# very simple markdown to text converter
# assumes very basic markdown as input
def _markdown_to_text:
  ( .
  # ```
  # code
  # ```
  # -> code
  | gsub("\\n```\\n"; "\n"; "m")
  # #, ##, ###, ... -> #
  | gsub("(?<line>\\n)?#+(?<title>.*)\\n"; "\(.line // "")#\(.title)\n"; "m")
  # [title](url) -> title (url)
  | gsub("\\[(?<title>.*)\\]\\((?<url>.*)\\)"; "\(.title) (\(.url))")
  # `code` -> code
  | gsub("`(?<code>.*)`"; .code)
  );

def expr_to_path: _expr_to_path;
def path_to_expr: _path_to_expr;

def torepr:
  ( format as $f
  | if $f == null then error("value is not a format root") end
  | _format_func($f; "torepr")
  );
