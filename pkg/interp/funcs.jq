
# convert number to array of bytes
def number_to_bytes($bits):
	def _number_to_bytes($d):
		if . > 0 then
			. % $d, (. div $d | _number_to_bytes($d))
		else
			empty
		end;
	if . == 0 then [0]
	else [_number_to_bytes(1 bsl $bits)] | reverse end;
def number_to_bytes:
	number_to_bytes(8);


def from_radix($base; $table):
	split("")
	| reverse
	| map($table[.])
	| if . == null then error("invalid char \(.)") end
	| reduce .[] as $c
		# state: [power, ans]
		([1,0]; (.[0] * $base) as $b | [$b, .[1] + (.[0] * $c)])
	| .[1];

def to_radix($base; $table):
	if . == 0 then "0"
	else
		[ recurse(if . > 0 then . div $base else empty end) | . % $base]
		| reverse
		| .[1:]
		| if $base <= ($table | length) then
		 	map($table[.]) | join("")
		  else
		 	error("base too large")
		  end
	end;

def radix2:
	if (. | type) == "number" then to_radix(2;"01")
	else from_radix(2;{"0": 0, "1": 1}) end;

def radix8:
	if (. | type) == "number" then to_radix(8;"01234567")
	else from_radix(8;{"0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7}) end;

def radix16:
	if (. | type) == "number" then to_radix(16;"0123456789abcdef")
	else from_radix(16;{
		"0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15
	})
	end;

def radix62:
	if (. | type) == "number" then to_radix(62;"0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	else from_radix(62;{
		"0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"A": 10, "B": 11, "C": 12, "D": 13, "E": 14, "F": 15, "G": 16,
		"H": 17, "I": 18, "J": 19, "K": 20, "L": 21, "M": 22, "N": 23,
		"O": 24, "P": 25, "Q": 26, "R": 27, "S": 28, "T": 29, "U": 30,
		"V": 31, "W": 32, "X": 33, "Y": 34, "Z": 35,
		"a": 36, "b": 37, "c": 38, "d": 39, "e": 40, "f": 41, "g": 42,
		"h": 43, "i": 44, "j": 45, "k": 46, "l": 47, "m": 48, "n": 49,
		"o": 50, "p": 51, "q": 52, "r": 53, "s": 54, "t": 55, "u": 56,
		"v": 57, "w": 58, "x": 59, "y": 60, "z": 61
	})
	end;

def radix62sp:
	if (. | type) == "number" then to_radix(62;"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	else from_radix(62;{
		"0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15, "g": 16,
		"h": 17, "i": 18, "j": 19, "k": 20, "l": 21, "m": 22, "n": 23,
		"o": 24, "p": 25, "q": 26, "r": 27, "s": 28, "t": 29, "u": 30,
		"v": 31, "w": 32, "x": 33, "y": 34, "z": 35,
		"A": 36, "B": 37, "C": 38, "D": 39, "E": 40, "F": 41, "G": 42,
		"H": 43, "I": 44, "J": 45, "K": 46, "L": 47, "M": 48, "N": 49,
		"O": 50, "P": 51, "Q": 52, "R": 53, "S": 54, "T": 55, "U": 56,
		"V": 57, "W": 58, "X": 59, "Y": 60, "Z": 61
	})
	end;

def radix62:
	if (. | type) == "number" then to_radix(62;"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	else from_radix(62;{
		"A": 0, "B": 1, "C": 2, "D": 3, "E": 4, "F": 5, "G": 6,
		"H": 7, "I": 8, "J": 9, "K": 10, "L": 11, "M": 12, "N": 13,
		"O": 14, "P": 15, "Q": 16, "R": 17, "S": 18, "T": 19, "U": 20,
		"V": 21, "W": 22, "X": 23, "Y": 24, "Z": 25,
		"a": 26, "b": 27, "c": 28, "d": 29, "e": 30, "f": 31, "g": 32,
		"h": 33, "i": 34, "j": 35, "k": 36, "l": 37, "m": 38, "n": 39,
		"o": 40, "p": 41, "q": 42, "r": 43, "s": 44, "t": 45, "u": 46,
		"v": 47, "w": 48, "x": 49, "y": 50, "z": 51,
		"0": 52, "1": 53, "2": 54, "3": 55, "4": 56, "5": 57, "6": 58, "7": 59, "8": 60, "9": 61
	})
	end;

def radix64:
	if (. | type) == "number" then to_radix(64;"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	else from_radix(64;{
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
	})
	end;

# like iprint
def i:
	{
		bin: "0b\(radix2)",
		oct: "0o\(radix8)",
		dec: "\(.)",
		hex: "0x\(radix16)",
		str: (try ([.] | implode) catch null),
	};

def diff($a; $b):
    ($a | type) as $at
    | ($b | type) as $bt
    | if $at != $bt then {a: $a, b: $b}
      elif ($at == "array" or $at == "object" or $at == "struct") then
        [ ((($a | keys) + ($b | keys)) | unique)[] as $k
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
      else
        if $a == $b then empty else {a: $a, b: $b} end
      end;

def _formats_dot:
	def _record($title;$fields):
		[  "<"
		, "<table bgcolor=\"paleturquoise\" border=\"0\" cellspacing=\"0\">"
		, "<tr><td port=\"\($title)\">\($title)</td></tr>"
		, [$fields|flatten|map("<tr><td align=\"left\" bgcolor=\"lightgrey\" port=\"\(.)\">\(.)</td></tr>")]
		, "</table>"
		, ">"
		] | flatten | join("");
	"# ... | dot -Tsvg -o formats.svg"
	, "digraph formats {"
	, "  concentrate=True"
	, "  rankdir=TB"
	, "  graph ["
	, "  ]"
	, "  node [shape=\"none\"style=\"\"]"
	, "  edge [arrowsize=\"0.7\"]"
	, (.[]
	  | . as $f
	  | .dependencies|flatten?|.[]
	  | "  \"\($f.name)\":\(.) -> \(.)"
	  )
	, (.[]
	  | .name as $name
	  | .groups[]?
	  | "  \(.) -> \"\($name)\":\($name)"
	  )
	, (to_entries[]
	  | "  \(.key) [color=\"paleturquoise\", label=\(_record(.key;(.value.dependencies//[])))]")
	, ([.[].groups[]?] | unique[]
	  | "  \(.) [shape=\"record\",style=\"rounded,filled\",color=\"palegreen\"]"
	  )
	, "}";

def in_bits_range($p):
	select(scalars and ._start? and ._start <= $p and $p < ._stop);
def in_bytes_range($p):
	select(scalars and ._start? and ._start/8 <= $p and $p < ._stop/8);

# TODO: split? can't really switch on type
def grep(f):
	if f | type == "string" then
		.. | select((._name | contains(f)) or (._value | contains(f)? // false))
	elif f | type == "number" then
		.. | select(._value == f)
	else
		.. | debug | select(f)?
	end;
