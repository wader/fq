include "@builtin/common.jq";
include "@builtin/opts.jq";
include "@builtin/funcs.jq";

# TODO: completionMode
def complete($e):
	($e | complete_query) as {$type, $query, $prefix} |
	{
		prefix: $prefix,
		names: (
			(if $type == "function" or $type == "variable" then
				[.[] | eval($query) | scope] | add
			elif $type == "index" then
				[.[] | eval($query) | keys?, _value_keys?] | add
			else
				[]
			end) | map(select(strings and select(startswith($prefix)))) | unique | sort
		)
	};

def set_eval_options: options(options_expr | with_entries(.value |= eval(.)));

def prompt:
	def _type_name:
		. as $c | try (. | display_name) catch ($c | type);
	def _path_prefix:
		(._path? // ".") | if . == "." then "" else .+" " end;
	(if (. | length) == 1 then
		.[0] | (_path_prefix + _type_name)
	else
		"[" +
		((.[0] | _type_name) +
		if (. | length) > 1 then ",..." else "" end) +
		"]" + "[\(length)]"
	end
	) + "> ";


def eval_f($e;f):
	set_eval_options as $_ |
	try eval($e) | f
	catch (. as $err | ("error: " + $err) | print);

def eval_print($e):
	def _display:
		. as $c |
		try $c | display({maxdepth: 1})
		catch (
			if $c | type == "string" then $c
			elif $c | type == "number" then $c
			else $c | tojson end
		);
	eval_f($e;_display | print);


# def read: #:: [a]| => string
# read with no prompt or completion
# def read: #:: [a]|(string) => string
# read with prompt and no completion
# def read: #:: [a]|(string;string) => string
# First argument is prompt to use.
# Second argument is name of completion function [a](string) => [string],
# it will be called with same input as read and a string argument being the
# current line from start to current cursor position. Should return possible completions.
def repl:
	def _read_expr: read(prompt;"complete") | trim | if . == "" then "." end;
	def _as_array: if (. | type) != "array" then [.] end;
	def _repl:
		try _read_expr as $e |
		(.[] | eval_print($e) | empty),
		_repl;
    _as_array | _repl;

def main:
	def _formats_list:
		((formats | keys | map(length) | max)+2) as $m | [
			"\("Name:" | rpad($m;" "))Description:",
			(
				formats | to_entries[] |
				"\(.key|rpad($m;" "))\(.value.description)"
			)
		] | join("\n");
	def _opts($version):
		{
			"help": {
				short: "-h",
				long: "--help",
				description: "Show help",
				bool: true
			},
			"formats": {
				short: "-h",
				long: "--formats",
				description: "Show formats",
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
				string: true
			},
			"version": {
				short: "-v",
				long: "--version",
				description: "Show version (\($version))",
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
	.version as $version |
	.args[0] as $arg0 |
	opts_parse(.args[1:];_opts($version)) as {$parsed, $rest} |
	# TODO: pass repl some other way
	options_expr($parsed.options + {repl: ($parsed.repl|tojson)}) |
	set_eval_options |
	if $parsed.version then
		$version | print
	elif $parsed.formats then
		_formats_list | print
	elif $parsed.help then
		("Usage: \($arg0) [OPTIONS] [FILE] [EXPR]",
			opts_help_text(_opts($version))
		) | print
	else
		(if $parsed.file then open($parsed.file) | string
		 else (if $parsed.noinput then $rest[0] else $rest[1] end) // "." end
		) as $expr |
		if $parsed.noinput then
			null
		else
			(if $rest[0] then $rest[0] else "-" end) as $filename |
			open($filename) |
			decode($parsed.decode)
		end |
		if $parsed.repl then
			eval_f($expr;repl)
		else
			eval_print($expr)
		end
	end;
