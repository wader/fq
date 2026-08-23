# from_base($base; $table) - Convert string to number in base and custom digits table
#
# Examples:
# "baab" | from_base(2; {"a": 0, "b": 1}) -> 9
#
# Tests:
# from_base(2; {"a": 0, "b": 1})
# "baab"
# 9
def from_base($base; $table):
  ( if type == "string" | not then error("cannot convert: \(.)") end
  | split("")
  | reverse
  | map(
      ( . as $c
      | $table[$c]
      | if . == null or . >= $base then
          error("invalid character '\($c)' for base \($base) number")
        end
      )
    )
  # state: [power, ans]
  | reduce .[] as $c ([1,0];
      ( (.[0] * $base) as $b
      | [$b, .[1] + (.[0] * $c)]
      )
    )
  | .[1]
  );
# from_base($base) - Convert string to number in base.
#
# Examples:
# "ff" | from_base(16) -> 255
#
# Tests:
# map(from_base(16))
# ["1","1","1","10000","20","10","11111111","377","ff","FF","100000000","400","100"]
# [1,1,1,65536,32,16,286331153,887,255,255,4294967296,1024,256]
def from_base($base):
  from_base(
    $base;
    if $base == 2 then {"0": 0, "1": 1}
    elif $base == 8 then
      {"0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7}
    elif $base == 10 then
      {"0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9}
    elif $base == 16 then
      { "0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9
      , "a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15
      , "A": 10, "B": 11, "C": 12, "D": 13, "E": 14, "F": 15
      }
    else
      { "0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9
      , "a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15, "g": 16
      , "h": 17, "i": 18, "j": 19, "k": 20, "l": 21, "m": 22, "n": 23
      , "o": 24, "p": 25, "q": 26, "r": 27, "s": 28, "t": 29, "u": 30
      , "v": 31, "w": 32, "x": 33, "y": 34, "z": 35
      , "A": 36, "B": 37, "C": 38, "D": 39, "E": 40, "F": 41, "G": 42
      , "H": 43, "I": 44, "J": 45, "K": 46, "L": 47, "M": 48, "N": 49
      , "O": 50, "P": 51, "Q": 52, "R": 53, "S": 54, "T": 55, "U": 56
      , "V": 57, "W": 58, "X": 59, "Y": 60, "Z": 61
      , "@": 62, "_": 63
      }
    end
  );
# from_base - Convert string to number and infer base.
#
# Examples:
# "0xff" | from_base -> 255
#
# Tests:
# map(from_base)
# ["0b10","0o10","10", "0x10", "0xff", "0xFF"]
# [2,8,10,16,255,255]
#
# try from_base catch .
# "0xGF"
# "invalid character 'G' for base 16 number"
def from_base:
  ( ( if startswith("0b") or startswith("0B") then [2, .[2:]]
      elif startswith("0o") or startswith("0O") then [8, .[2:]]
      elif startswith("0x") or startswith("0X") then [16, .[2:]]
      else [10, .]
      end
    ) as [$base, $s]
  | $s
  | from_base($base)
  );

# to_base($base; prefix; $table) - Convert number to string in base using custom prefix and digits table.
#
# Examples:
# 9 | to_base(2; "2#"; "ab") -> "2#baab"
#
# Tests:
# to_base(2; "2#"; "ab")
# 9
# "2#baab"
def to_base($base; prefix; $table):
  # integer division
  # inspired by https://github.com/itchyny/gojq/issues/63#issuecomment-765066351
  def _intdiv($a; $b):
    # TODO: figure out a saner way to force int
    def _to_int: (. % (. + 1));
    ( ($a | _to_int) as $a
    | ($b | _to_int) as $b
    | ($a - ($a % $b)) / $b
    );
  ( if type != "number" then error("cannot convert: \(.)") end
  | ( ($base | prefix)
    + if . == 0 then "0"
      else
        ( [recurse(if . > 0 then _intdiv(.; $base) else empty end) | . % $base]
        | reverse
        | .[1:]
        | if $base <= ($table | length) then
            map($table[.:.+1]) | join("")
          else
            error("base too large")
          end
        )
      end
    )
  );
# to_base($base; prefix) - Convert number to string in base using custom prefix.
#
# Examples:
# 255 | to_base(16; "") -> "ff"
#
# Tests:
# to_base(16; "")
# 255
# "ff"
def to_base($base; prefix):
  to_base($base; prefix; "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ@_");
# to_base($base) - Convert number to string in base.
#
# Examples:
# 255 | to_base(16; "") -> "0xff"
#
# Tests:
# map(to_base(2,8,16))
# [1,16,255,256]
# ["0b1","0o1","0x1","0b10000","0o20","0x10","0b11111111","0o377","0xff","0b100000000","0o400","0x100"]
def to_base($base):
  to_base(
    $base;
    if . == 2 then "0b"
    elif. == 8 then "0o"
    elif . == 16 then "0x"
    else ""
    end
  );
