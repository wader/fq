def from_radix($base; $table):
  ( if _is_string | not then error("cannot from_radix convert: \(.)") end
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
def from_radix($base):
  from_radix($base; {
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

def to_radix($base; $table):
  ( if type != "number" then error("cannot to_radix convert: \(.)") end
  | if . == 0 then "0"
    else
      ( [ recurse(if . > 0 then _intdiv(.; $base) else empty end) | . % $base]
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
def to_radix($base):
  to_radix($base; "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ@_");