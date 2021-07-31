# TODO: introspect and show doc, reflection somehow?
def help:
  ( builtins[]
  , "^C interrupt"
  , "^D exit REPL"
  ) | println;

# TODO: escape for safe key names
# path ["a", 1, "b"] -> "a[1].b"
def path_to_expr:
  map(if type == "number" then "[", ., "]" else ".", . end) | join("");

# TODO: don't use eval? should support '.a.b[1]."c.c"' and escapes?
def expr_to_path:
  ( if type != "string" then error("require string argument") end
  | eval("null | path(\(.))")
  );

def trim: capture("^\\s*(?<str>.*?)\\s*$"; "").str;

# does +1 and [:1] as " "*0 is null
def rpad($s; $w): . + ($s * ($w+1-length))[1:];

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
          children | map(select(name==$n))
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

# convert number to array of bytes
def number_to_bytes($bits):
  def _number_to_bytes($d):
    if . > 0 then
      . % $d, (. div $d | _number_to_bytes($d))
    else
      empty
    end;
  if . == 0 then [0]
  else [_number_to_bytes(1 bsl $bits)] | reverse
  end;
def number_to_bytes:
  number_to_bytes(8);

def from_radix($base; $table):
  ( split("")
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

def to_radix($base; $table):
  if . == 0 then "0"
  else
    ( [ recurse(if . > 0 then . div $base else empty end) | . % $base]
    | reverse
    | .[1:]
    | if $base <= ($table | length) then
        map($table[.]) | join("")
      else
        error("base too large")
      end
    )
  end;

def radix($base; $to_table; $from_table):
  if . | type == "number" then to_radix($base; $to_table)
  elif . | type == "string" then from_radix($base; $from_table)
  else error("needs to be number of string")
  end;

def radix2: radix(2; "01"; {"0": 0, "1": 1});
def radix8: radix(8; "01234567"; {"0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7});
def radix16:radix(16; "0123456789abcdef"; {
    "0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
    "a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15
  });
def radix62: radix(62; "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"; {
    "0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
    "A": 10, "B": 11, "C": 12, "D": 13, "E": 14, "F": 15, "G": 16,
    "H": 17, "I": 18, "J": 19, "K": 20, "L": 21, "M": 22, "N": 23,
    "O": 24, "P": 25, "Q": 26, "R": 27, "S": 28, "T": 29, "U": 30,
    "V": 31, "W": 32, "X": 33, "Y": 34, "Z": 35,
    "a": 36, "b": 37, "c": 38, "d": 39, "e": 40, "f": 41, "g": 42,
    "h": 43, "i": 44, "j": 45, "k": 46, "l": 47, "m": 48, "n": 49,
    "o": 50, "p": 51, "q": 52, "r": 53, "s": 54, "t": 55, "u": 56,
    "v": 57, "w": 58, "x": 59, "y": 60, "z": 61
  });
def radix62: radix(62; "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"; {
    "A": 0, "B": 1, "C": 2, "D": 3, "E": 4, "F": 5, "G": 6,
    "H": 7, "I": 8, "J": 9, "K": 10, "L": 11, "M": 12, "N": 13,
    "O": 14, "P": 15, "Q": 16, "R": 17, "S": 18, "T": 19, "U": 20,
    "V": 21, "W": 22, "X": 23, "Y": 24, "Z": 25,
    "a": 26, "b": 27, "c": 28, "d": 29, "e": 30, "f": 31, "g": 32,
    "h": 33, "i": 34, "j": 35, "k": 36, "l": 37, "m": 38, "n": 39,
    "o": 40, "p": 41, "q": 42, "r": 43, "s": 44, "t": 45, "u": 46,
    "v": 47, "w": 48, "x": 49, "y": 50, "z": 51,
    "0": 52, "1": 53, "2": 54, "3": 55, "4": 56, "5": 57, "6": 58, "7": 59, "8": 60, "9": 61
  });
def radix64: radix(64; "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"; {
    "A": 0, "B": 1, "C": 2, "D": 3, "E": 4, "F": 5, "G": 6,
    "H": 7, "I": 8, "J": 9, "K": 10, "L": 11, "M": 12, "N": 13,
    "O": 14, "P": 15, "Q": 16, "R": 17, "S": 18, "T": 19, "U": 20,
    "V": 21, "W": 22, "X": 23, "Y": 24, "Z": 25,
    "a": 26, "b": 27, "c": 28, "d": 29, "e": 30, "f": 31, "g": 32,
    "h": 33, "i": 34, "j": 35, "k": 36, "l": 37, "m": 38, "n": 39,
    "o": 40, "p": 41, "q": 42, "r": 43, "s": 44, "t": 45, "u": 46,
    "v": 47, "w": 48, "x": 49, "y": 50, "z": 51,
    "0": 52, "1": 53, "2": 54, "3": 55, "4": 56, "5": 57, "6": 58, "7": 59, "8": 60, "9": 61,
    "+": 62, "/": 63
  });

# like iprint
def i:
  {
    bin: "0b\(radix2)",
    oct: "0o\(radix8)",
    dec: "\(.)",
    hex: "0x\(radix16)",
    str: (try ([.] | implode) catch null),
  };

# produce a/b pairs for diffing values
def diff($a; $b):
  ( ( $a | type) as $at
  | ($b | type) as $bt
  | if $at != $bt then {a: $a, b: $b}
    elif ($at == "array" or $at == "object") then
      ( [ ((($a | keys) + ($b | keys)) | unique)[] as $k
        | {
          ($k | tostring): (
            [($a | has($k)), ($b | has($k))]
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

def in_bits_range($p):
  select(scalars and ._start? and ._start <= $p and $p < ._stop);
def in_bytes_range($p):
  select(scalars and ._start? and ._start/8 <= $p and $p < ._stop/8);

