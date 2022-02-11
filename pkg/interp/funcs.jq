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

def formats:
  _registry.formats;

def intdiv(a; b): _intdiv(a; b);

# TODO: escape for safe key names
# path ["a", 1, "b"] -> "a[1].b"
def path_to_expr:
  ( if length == 0 or (.[0] | type) != "string" then
      [""] + .
    end
  | map(
      if type == "number" then "[", ., "]"
      else
        ( "."
        , # empty (special case for leading index or empty path) or key
          if . == "" or _is_ident then .
          else
            ( "\""
            , _escape_ident
            , "\""
            )
          end
        )
      end
    )
  | join("")
  );

# TODO: don't use eval? should support '.a.b[1]."c.c"' and escapes?
def expr_to_path:
  ( if type != "string" then error("require string argument") end
  | eval("null | path(\(.))")
  );

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

# helper to build path query/generate functions for tree structures with
# non-unique children, ex: mp4_path
def tree_path(children; name; $v):
  def _lookup:
    # add implicit zeros to get first value
    # ["a", "b", 1] => ["a", 0, "b", 1]
    def _normalize_path:
      ( . as $np
      | if $np | last | type == "string" then $np+[0] end
      # state is [path acc, possible pending zero index]
      | ( reduce .[] as $np ([[], []];
          if $np | type == "string" then
            [(.[0]+.[1]+[$np]), [0]]
          else
            [.[0]+[$np], []]
          end
        ))
      )[0];
    ( . as $c
    | $v
    | expr_to_path
    | _normalize_path
    | reduce .[] as $n ($c;
        if $n | type == "string" then
          children | map(select(name == $n))
        else
          .[$n]
        end
      )
    );
  def _path:
    [ . as $r
    | $v._path as $p
    | foreach range(($p | length)/2) as $i (null; null;
        ( ($r | getpath($p[0:($i+1)*2]) | name) as $name
        | [($r | getpath($p[0:($i+1)*2-1]))[] | name][0:$p[($i*2)+1]+1] as $before
        | [ $name
          , ($before | map(select(. == $name)) | length)-1
          ]
        )
      )
    | [ ".", .[0],
      (.[1] | if . == 0 then empty else "[", ., "]" end)
      ]
    ]
    | flatten
    | join("");
  if $v | type == "string" then _lookup
  else _path
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
