include "@builtin/common.jq";
include "@builtin/opts.jq";
include "@builtin/funcs.jq";

# TODO: completionMode
def complete($e):
	($e | complete_query) as {$type, $query, $prefix}
	| {
		prefix: $prefix,
		names: (
			if $type == "function" or $type == "variable" then
				[.[] | eval($query) | scope] | add
			elif $type == "index" then
				[.[] | eval($query) | keys?, _value_keys?] | add
			else
				[]
			end
			| map(select(strings and startswith($prefix)))
			| unique
			| sort
		)
	};

def obj_to_csv_kv: [to_entries[] | [.key, .value] | join("=")] | join(",");

def build_default_options:
	{
		depth:        0,
		verbose:      false,
		color:        (tty.is_terminal and env.CLICOLOR!=null),
		unicode:      (tty.is_terminal and env.CLIUNICODE!=null),
		raw:          (tty.is_terminal | not),
		linebytes:    (if tty.is_terminal then [((tty.size[0] div 8) div 2) * 2, 4] | max else 16 end),
		displaybytes: (if tty.is_terminal then [((tty.size[0] div 8) div 2) * 2, 4] | max else 16 end),
		addrbase:     16,
		sizebase:     10,
		colors: ({
			null: "brightblack",
			false: "yellow",
			true: "yellow",
			number: "cyan",
			string: "green",
			objectkey: "brightblue",
			array: "white",
			object: "white",
			index: "white",
			value: "white",
			error: "brightred",
			frame: "yellow"
		} | obj_to_csv_kv),
		bytecolors: "0-255=brightwhite,0=brightblack,32-126:9-13=white",
	};

def parse_options:
	{
			depth:        (try (.depth | fromjson) catch null),
			verbose:      (try (.verbose | fromjson) catch null),
			color:        (try (.color | fromjson) catch null),
			unicode:      (try (.unicode | fromjson) catch null),
			raw:          (try (.raw | fromjson) catch null),
			linebytes:    (try (.linebytes | fromjson) catch null),
			displaybytes: (try (.displaybytes | fromjson) catch null),
			addrbase:     (try (.addrbase | fromjson) catch null),
			sizebase:     (try (.sizebase | fromjson) catch null),
			colors:       .colors,
			bytecolors:   .bytecolors,
	}
	| with_entries(select(.value));

def prompt:
	def _type_name_error:
		. as $c
		| try
			(. | display_name) +
				if ._error then "!" else "" end
		catch ($c | type);
	def _path_prefix:
		(._path? // ".")
		| if . == "." then "" else . + " " end;
	( if (. | length) == 1 then
		.[0] | _path_prefix + _type_name_error
	  else
		[ "["
		, ((.[0] | _type_name_error)
		, if (. | length) > 1 then ",..." else "" end)
		, "]"
		, "[\(length)]"
		] | join("")
	  end
	) + "> ";

def eval_f($e;f):
	default_options(build_default_options) as $_
	| try eval($e) | f
	  catch (. as $err | ("error: " + $err) | print);

def default_display: display({depth: 1});

def eval_print($e):
	eval_f($e;default_display);

# def read: #:: [a]| => string
# read with no prompt or completion
# def read: #:: [a]|(string) => string
# read with prompt and no completion
# def read: #:: [a]|(string;string) => string
# First argument is prompt to use.
# Second argument is name of completion function [a](string) => [string],
# it will be called with same input as read and a string argument being the
# current line from start to current cursor position. Should return possible completions.
def repl($opts):
	def _as_array: if (. | type) != "array" then [.] end;
	def _read_expr:
		read(prompt;"complete")
		| trim
		| if . == "" then "." end;
	def _repl:
		. as $c
		| try
			_read_expr as $e
			| (.[] | eval_print($e) | empty),
			_repl
		  catch
			if . == "interrupt" then $c | _repl
			elif . == "eof" then empty
			else error(.) end;
	with_options($opts; _as_array | _repl);

def repl: repl({});

def main:
	def _formats_list:
		[
			["Name:", "Description:"],
			( formats
			  | to_entries[]
			  | [(.key+" "), .value.description]
			)
		]
		| table(
			.;
			[.[] as $rc | $rc.string | rpad(" ";$rc.maxwidth)] | join("")
		);
	def _opts($version):
		{
			"version": {
				short: "-v",
				long: "--version",
				description: "Show version (\($version))",
				bool: true
			},
			"help": {
				short: "-h",
				long: "--help",
				description: "Show help",
				bool: true
			},
			"formats": {
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
			"rawstring": {
				short: "-r",
				description: "Raw strings",
				bool: true
			},
			"options": {
				short: "-o",
				long: "--option",
				description: "Set option, eg: color=true",
				object: true,
				default: {},
				help_default: build_default_options
			},
		};
	.version as $version
	| .args[0] as $arg0
	| opts_parse(.args[1:];_opts($version)) as {$parsed, $rest}
	# TODO: hack, pass opts some other way
	| default_options(build_default_options) as $_
	| push_options(
		($parsed.options | parse_options)
		+ {
			repl: $parsed.repl,
			rawstring: $parsed.rawstring,
		}
	)
	| if $parsed.version then
		$version | print
	  elif $parsed.formats then
		_formats_list | print
	  elif $parsed.help then
		"Usage: \($arg0) [OPTIONS] [FILE] [EXPR]...",
		opts_help_text(_opts($version))
		| print
	  else
		null
		# figure out filename and expressions
		| ( if $parsed.noinput then [null, $rest]
		    elif $rest[0] then [$rest[0], $rest[1:]]
		    else ["-", $rest]
		    end
		  ) as [$filename, $exprs]
		| if $filename then
			( open($filename)
			| decode($parsed.decode)
			)
		  end
		| if $parsed.file then
			( (open($parsed.file) | string) as $file_expr
			| eval_f($file_expr;.)
			)
		  end
		# this evaluates and combines all expression in order
		| (reduce $exprs[] as $expr ([.];[.[] | eval_f($expr;.)]))[]
		| if $parsed.repl then repl
		  else default_display end
	  end;
