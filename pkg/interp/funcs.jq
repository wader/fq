include "internal";
include "options";
include "binary";
include "decode";

def intdiv(a; b): _intdiv(a; b);

# does +1 and [:1] as " "*0 is null
def rpad($s; $w): . + ($s * ($w+1-length))[1:];

# add missing group/0 function
# https://github.com/stedolan/jq/issues/2444
def group: group_by(.);

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

# TODO: rename keys and add more, ascii/utf8/utf16/codepoint name?, le/be, signed/unsigned?
# TODO: move?
def iprint:
  { bin: "0b\(to_radix(2))"
  , oct: "0o\(to_radix(8))"
  , dec: "\(.)"
  , hex: "0x\(to_radix(16))"
  , str: (try ([.] | implode) catch null)
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

def expr_to_path: _expr_to_path;
def path_to_expr: _path_to_expr;

def torepr:
  ( format as $f
  | if $f == null then error("value is not a format root") end
  | _format_func($f; "torepr")
  );
