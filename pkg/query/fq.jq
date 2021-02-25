# TODO;
# modules?

include "@builtin/common.jq";
include "@builtin/opts.jq";
include "@builtin/funcs.jq";

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
				bool: true
			},
			"noinput": {
				short: "-n",
				long: "--noinput",
				description: "No input",
				bool: true
			},
			"decode": {
				short: "-d",
				long: "--decode",
				description: "Decoder",
				default: "probe",
				string: true
			},
			"repl": {
				short: "-i",
				long: "--repl",
				description: "Interactive REPL",
				bool: true
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
				description: "Show version (\($VERSION))",
				bool: true
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
	# TODO: pass repl some other way
	options_expr($parsed.options + {repl: ($parsed.repl|tojson)}) |
	set_eval_options |
	if $parsed.version then
		$VERSION | print
	elif $parsed.help then
		("Usage: \($args[0]) [OPTIONS] [FILE] [EXPR]",
			opts_help_text(_opts),
		 	formats_help_text
		) | print
	else
		(if $parsed.file then open($parsed.file) | string
		 else
			(if $parsed.noinput then $rest[0] else $rest[1] end) // "."
		 end
		) as $expr |
		if $parsed.noinput then
			null
		else
			(if $rest[0] then $rest[0] else "-" end) as $filename |
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
