# TODO;
# modules?
# done interrupt
# nicer options, some kind of "eval state object"?

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


def from_base($base;$table):
	split("")
	| reverse
	| map($table[.])
	| if . == null then error("invalid char \(.)") end
	| reduce .[] as $c
		# state: [power, ans]
		([1,0]; (.[0] * $base) as $b | [$b, .[1] + (.[0] * $c)])
	| .[1];

def to_base($base;$table):
	def stream:
		recurse(if . > 0 then . div $base else empty end) | . % $base;
	if . == 0 then
		"0"
	else
		[stream] |
		reverse  |
		.[1:] |
		if $base <= ($table | length) then
			map($table[.]) | join("")
		else
			error("base too large")
		end
	end;

def base2:
	if (. | type) == "number" then to_base(2;"01")
	else from_base(2;{"0": 0, "1": 1}) end;

def base8:
	if (. | type) == "number" then to_base(8;"01234567")
	else from_base(8;{"0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7}) end;

def base16:
	if (. | type) == "number" then to_base(16;"0123456789abcdef")
	else from_base(16;{
		"0": 0, "1": 1, "2": 2, "3": 3,"4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"a": 10, "b": 11, "c": 12, "d": 13, "e": 14, "f": 15
	})
	end;

def base62:
	if (. | type) == "number" then to_base(62;"0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	else from_base(62;{
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

def base62sp:
	if (. | type) == "number" then to_base(62;"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	else from_base(62;{
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

def base62:
	if (. | type) == "number" then to_base(62;"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	else from_base(62;{
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

def base64:
	if (. | type) == "number" then to_base(64;"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	else from_base(64;{
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

# from https://rosettacode.org/wiki/Non-decimal_radices/Convert#jq
# unknown author
# Convert the input integer to a string in the specified base (2 to 36 inclusive)
def _convert(base):
	def stream:
		recurse(if . > 0 then . div base else empty end) | . % base;
	if . == 0 then
		"0"
	else
		[stream] |
		reverse  |
		.[1:] |
		if base <  10 then
			map(tostring) | join("")
		elif base <= 36 then
			map(if . < 10 then 48 + . else . + 87 end) | implode
		else
			error("base too large")
		end
	end;

# input string is converted from "base" to an integer, within limits
# of the underlying arithmetic operations, and without error-checking:
def _to_i(base):
	explode
	| reverse
	| map(if . > 96  then . - 87 else . - 48 end)  # "a" ~ 97 => 10 ~ 87
	| reduce .[] as $c
		# state: [power, ans]
		([1,0]; (.[0] * base) as $b | [$b, .[1] + (.[0] * $c)])
	| .[1];

# like iprint
def i:
	{
		bin: "0b\(base2)",
		oct: "0o\(base8)",
		dec: "\(.)",
		hex: "0x\(base16)",
		str: ([.] | implode),
	};

def _formats_dot:
	"# ... | dot -Tsvg -o formats.svg",
	"digraph formats {",
	"  node [shape=\"box\",style=\"rounded,filled\"]",
	"  edge [arrowsize=\"0.7\"]",
	(.[] | "  \(.name) -> {\(.dependencies | flatten? | join(" "))}"),
	(.[] | .name as $name | .groups[]? | "  \(.) -> \($name)"),
	(keys[] | "  \(.) [color=\"paleturquoise\"]"),
	([.[].groups[]?] | unique[] | "  \(.) [color=\"palegreen\"]"),
	"}";

def field_inrange($p): ._type == "field" and ._range.start <= $p and $p < ._range.stop;


def dv($p):
    . as $c | [$p, $c] | debug | $c;

def trim: capture("^\\s*(?<a>.*?)\\s*$"; "").a;

# does +1 and [:1] as " "*0 is null
def rpad($w;$s): . + ($s * (([0,$w-(.|length)] | max)+1))[1:];



# TODO: completionMode
def complete($e):
	($e | complete_query) as {$type, $query, $prefix} |
	{
		prefix: $prefix,
		names: (
			if $type == "function" or $type == "variable" then
				[.[] | eval($query) | scope[] | select(startswith($prefix))]
			elif $type == "index" then
				[
					[.[] | eval($query) | keys?, _value_keys?] |
					add | unique | sort | .[] | strings | select(startswith($prefix))
				]
			else
				[]
			end
		)
	};

def set_eval_options: options(options_expr | with_entries(.value |= eval(.)));

def prompt:
	def _display_name:
		. as $c | try (. | display_name) catch ($c | type);
	((.[0] | _display_name) +
	if (. | length) > 1 then ",[\((. | length) - 1)]..." else "" end) + "> ";

def eval_print($e):
	set_eval_options as $_ |
	try eval($e) as $v |
		try ($v | display({maxdepth: 1}))
		catch ($v | tojson | print)
	catch (. as $err | ("ERR: " + $err) | print);


# def readline: #:: [a]|(string;string) => string
# First argument is name of completion function [a](string) => [string],
# it will be called with same input as readline and a string argument being the
# current line from start to current cursor position. Should return possible completions.
# Second argument is name of prompt function [a] => string, it will be called with
# same input as readline and should return a string.

def repl:
	def _readline_expr: readline("complete";"prompt") | trim | if . == "" then "." end;
	def _as_array: if (. | type) != "array" then [.] end;
	def _repl:
		try _readline_expr as $e |
		(.[] | eval_print($e) | empty),
		_repl;
    _as_array | _repl;

# TODO: validate option name? has key
# TODO: multi short -in
# TODO: parse -args until -- or end, collect unknown to rest
def opts_parse($args;$opts):
	def _parse($args;$flagmap;$parsed):
		def _parse_with_arg($argskip;$optname;$value;$opt):
			if $opt.object then
				($value | capture("^(?<key>.*?)=(?<value>.*)$") // error("\($value): should be key=value"))
				as {$key, $value} |
				_parse($args[$argskip:];$flagmap;($parsed|.[$optname][$key] |= $value))
			elif $opt.array then
				_parse($args[$argskip:];$flagmap;($parsed|.[$optname] += [$value]))
			else
				_parse($args[$argskip:];$flagmap;($parsed|.[$optname] = $value))
			end;
		def _parse_without_arg($optname):
			_parse($args[1:];$flagmap;($parsed|.[$optname] = true));
		($args[0] | index("=")) as $assigni |
		(
			if $assigni then $args[0][0:$assigni]
			else $args[0] end
		) as $arg |
		if $arg == null then
			{parsed: $parsed, rest: []}
		else
			if $arg == "--" then
				{parsed: $parsed, rest: $args[1:]}
			elif $arg | test("^--?.+") then
				$flagmap[$arg] as $optname |
				($opts[$optname]? // null) as $opt |
				if $opt == null then
					error("\($arg): no such argument")
				elif $opt.value or $opt.array or $opt.object then
					if $assigni then
						_parse_with_arg(1;$optname;$args[0][$assigni+1:];$opt)
					elif ($args | length) < 2 then
						error("\($arg): needs an argument")
					else
						_parse_with_arg(2;$optname;$args[1];$opt)
					end
				else
					if $assigni then error("\($arg): takes no argument")
					else _parse_without_arg($optname) end
				end
			else
				{parsed: $parsed, rest: $args}
			end
		end;
	def _flagmap:
		($opts | to_entries | map({(.value.short): .key, (.value.long): .key}) | add);
	def _defaults:
		($opts | to_entries | map(select(.value.default)) | map({(.key): .value.default}) | add);
	_parse($args;_flagmap;_defaults);

def opts_help_text($opts):
	def _opthelp:
		[
			"\(.long),\(.short)",
			if .value or .array or .object then "=ARG,\(.short) ARG" else "" end
		] | join("");
	def _maxoptlen:
		[$opts[] | (.|_opthelp|length)] | max;
	_maxoptlen as $l |
	[
		$opts | to_entries[] | [
		"\(.value|_opthelp|rpad($l;" "))  \(.value.description)",
		if .value.default then
			if .value.object then
				[
					"\n",
					if .value.default_eval then
						[.value.default | to_entries[] | "\(" "*$l)    \(.key)=\(eval(.value))\n"]
					else
						[.value.default | to_entries[] | "\(" "*$l)    \(.key)=\(.value)\n"]
					end
				]
			else
				" (\(.value.default))\n"
			end
		else
			"\n"
		end
		]
	] |
	flatten | join("");

def formats_help_text:
	((formats | keys | map(length) | max)+2) as $m | [
		"\("Name:" | rpad($m;" "))Description:", "\n",
		(
			formats | to_entries[] |
			"\(.key|rpad($m;" "))\(.value.description)", "\n"
		)
	] | join("");

def main($args):
	def _opts:
		{
			"help": {
				short: "-h",
				long: "--help",
				description: "Show help",
			},
			"noinput": {
				short: "-n",
				long: "--noinput",
				description: "No input",
			},
			"decode": {
				short: "-d",
				long: "--decode",
				description: "Decoder",
				default: "probe",
				value: true
			},
			"repl": {
				short: "-i",
				long: "--repl",
				description: "Interactive REPL",
			},
			"file": {
				short: "-f",
				long: "--file",
				description: "Read script from file",
				value: true
			},
			"version": {
				short: "-v",
				long: "--version",
				description: "Show version (\($VERSION))"
			},
			"options": {
				short: "-o",
				long: "--option",
				description: "Set option, eg: color=true",
				object: true,
				default_eval: true,
				default: {
					maxdepth:     "0",
					verbose:      "false",
					color:        "tty.is_terminal and env.CLICOLOR!=null",
					unicode:      "tty.is_terminal and env.CLIUNICODE!=null",
					raw:          "tty.is_terminal | not",
					linebytes:    "if tty.is_terminal then [((tty.size[0] div 8) div 2) * 2, 4] | max else 16 end",
					displaybytes: "if tty.is_terminal then [((tty.size[0] div 8) div 2) * 2, 4] | max else 16 end",
					addrbase:     "16",
					sizebase:     "10",
				}
			},
		};
	opts_parse($args[1:];_opts) as {$parsed, $rest} |
	options_expr($parsed.options) |
	set_eval_options |
	if $parsed.version then
		$VERSION | print
	elif $parsed.help then
		("Usage: \($args[0]) [OPTIONS] [FILE] [EXPR]",
			opts_help_text(_opts),
		 	formats_help_text
		) | print
	else
		(if $rest[0] then $rest[0] else "-" end) as $filename |
		(
			if $parsed.file then
				open($parsed.file) | string
			elif $rest[1] then $rest[1]
			else "." end
		) as $expr |
		if $parsed.noinput then
			null
		else
			open($filename) |
			decode($parsed.decode)
		end |
		if $parsed.repl then
			eval($expr) |
			repl
		else
			eval_print($expr) 
		end
	end;

