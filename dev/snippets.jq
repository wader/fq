def protobuf_to_value:
  .fields | map({(.name | tostring): (.enum // .value)}) | add;

# hack to parse just a box
# <binary> | mp4_box
def mp4_box:
  [0, 0, 0, 16, "ftyp", "isom", 0, 0 , 2 , 0, .] | mp4.boxes;

def flac_dump:
  [ "fLaC"
  , first(.. | select(._format == "flac_metadatablock"))
  , (.. | select(._format == "flac_frame"))
  ] | bits;

def urldecode:
  gsub(
    "%(?<c>..)";
    ( .c
    | ascii_downcase
    | explode
    | map(.-48 | if .>=49 then .-39 else . end)
    | [.[0]*16+.[1]]
    | implode
    )
  );

def radix62sp: radix(62; "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"; {
    "0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
    "a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15, "g": 16,
    "h": 17, "i": 18, "j": 19, "k": 20, "l": 21, "m": 22, "n": 23,
    "o": 24, "p": 25, "q": 26, "r": 27, "s": 28, "t": 29, "u": 30,
    "v": 31, "w": 32, "x": 33, "y": 34, "z": 35,
    "A": 36, "B": 37, "C": 38, "D": 39, "E": 40, "F": 41, "G": 42,
    "H": 43, "I": 44, "J": 45, "K": 46, "L": 47, "M": 48, "N": 49,
    "O": 50, "P": 51, "Q": 52, "R": 53, "S": 54, "T": 55, "U": 56,
    "V": 57, "W": 58, "X": 59, "Y": 60, "Z": 61
  });

# "01:09:55.76" -> 4195.76 ->
# 4195.76 -> "01:09:55.76"
def duration:
  def lpad($s; $w): ($s * ($w+1-length))[1:] + .;
  def _string:
    ( split(":")
    | map(tonumber)
    | reverse
    | [foreach .[] as $n (0; . + 1; pow(60; .-1) * $n)]
    # sum smallest to largest seem to improve precision
    | sort
    | reverse
    | add
    );
  def _number:
    if . == 0 then 0
    else
      [ ( [ (recurse(if . > 0 then . div 60 else empty end) | . % 60)]
        | reverse
        | .[1:]
        | map(tostring | lpad("0"; 2))
        | join(":")
        )
        # ugly but float is not accurate enough
      , ( tostring
        | split(".")[1] // empty
        | ".", .
        )
      ] | join("")
    end;
  if . | type == "string" then _string
  elif . | type == "number" then _number
  else error("expected string or number") end;

# def recurse_foreach(init; update; extract):
#     def _recurse_foreach($state; $c):
#         (. as $c | ["_recurse_foreach $state", $state] | debug | $c) |
#         (. as $b | ["_recurse_foreach $c", $c] | debug | $b) |
#         foreach $c[]? as $e (
#             $state
#             ;
#             (. as $c | ["update $state", $state] | debug | $c) |
#             (. as $b | ["update $e", $e] | debug | $b) |
#             null | _recurse_foreach({state: $state, e: $e} | update; $e)
#             ;
#             (. as $c | ["extract $state", $state] | debug | $c) |
#             (. as $b | ["extract $e", $e] | debug | $b) |
#             {state: $state, e: .} |
#             extract
#         )
#         // $c;
#     _recurse_foreach(init; .);

# # TODO: split? can't really switch on type
# def grep(f):
# 	if f | type == "string" then
# 		.. | select((._name | contains(f)) or (contains(f)? // false))
# 	elif f | type == "number" then
# 		.. | select(. == f)
# 	else
# 		.. | debug | select(f)?
# 	end;
