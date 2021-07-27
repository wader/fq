include "@builtin/common";
include "@builtin/args";
include "@builtin/funcs";
include "@builtin/format";

# optional user init
include "@config/init?";

# def read: #:: [a]| => string
# read with no prompt or completion

# def read($promp): #:: [a]|(string) => string
# read with prompt and no completion

# def read($promp; $completion): #:: [a]|(string;string) => string
# First argument is prompt to use.
# Second argument is name of completion function [a](string) => [string],
# it will be called with same input as read and a string argument being the
# current line from start to current cursor position. Should return possible completions.

# TODO: completionMode
def complete($e):
	( $e | complete_query) as {$type, $query, $prefix}
	| {
		prefix: $prefix,
		names: (
			if $type == "function" or $type == "variable" then
				[.[] | eval($query) | scope] | add
			elif $type == "index" then
				[.[] | eval($query) | keys?, extkeys?] | add
			else
				[]
			end
			| map(select(strings and startswith($prefix)))
			| unique
			| sort
		)
	};

def obj_to_csv_kv:
	[to_entries[] | [.key, .value] | join("=")] | join(",");

def color_themes:
	{
		default: {
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
				dumpheader: "yellow+underline",
				dumpaddr: "yellow"
			} | obj_to_csv_kv),
			bytecolors: "0-0xff=brightwhite,0=brightblack,32-126:9-13=white",
		},
		# TODO: more configurable? colors=neon?
		neon: {
			colors: ({
				null: "brightblack",
				false: "brightyellow",
				true: "brightyellow",
				number: "brightcyan",
				string: "brightgreen",
				objectkey: "brightblue",
				array: "brightwhite",
				object: "brightwhite",
				index: "brightwhite",
				value: "brightwhite",
				error: "brightred",
				dumpheader: "brightyellow+underline",
				dumpaddr: "brightyellow"
			} | obj_to_csv_kv),
			bytecolors: "0-0xff=brightwhite,0=brightblack,32-126:9-13=brightgreen",
		}
	};

def build_default_options:
	{
		depth:        0,
		verbose:      false,
		color:        (tty.is_terminal and env.CLICOLOR!=null),
		unicode:      (tty.is_terminal and env.CLIUNICODE!=null),
		raw:          (tty.is_terminal | not),
		linebytes:    (if tty.is_terminal then [((tty.width div 8) div 2) * 2, 4] | max else 16 end),
		displaybytes: (if tty.is_terminal then [((tty.width div 8) div 2) * 2, 4] | max else 16 end),
		addrbase:     16,
		sizebase:     10,
		colors:       color_themes.default.colors,
		bytecolors:   color_themes.default.bytecolors,
	};

def parse_options:
	{
			depth:        (.depth | if . then eval(.) else null end),
			verbose:      (.verbose | if . then eval(.) else null end),
			color:        (.color | if . then eval(.) else null end),
			unicode:      (.unicode | if . then eval(.) else null end),
			raw:          (.raw | if . then eval(.) else null end),
			linebytes:    (.linebytes | if . then eval(.) else null end),
			displaybytes: (.displaybytes | if . then eval(.) else null end),
			addrbase:     (.addrbase | if . then eval(.) else null end),
			sizebase:     (.sizebase | if . then eval(.) else null end),
			colors:       .colors,
			bytecolors:   .bytecolors,
	}
	| with_entries(select(.value != null));

def prompt:
	def _type_name_error:
		. as $c
		| try
			(. | display_name) +
				if ._error then "!" else "" end
		catch ($c | type);
	def _path_prefix:
		(._path? // [])
		| if . == [] then "" else path_to_expr + " " end;
	( options.repllevel
	  | if . > 1 then ((.-1) * ">") + " "
	    else "" end
	)
	+ (
	  if (. | length) == 1 then
		.[0] | _path_prefix + _type_name_error
	  else
		[ "["
		, ((.[0] | _type_name_error)
		, if (. | length) > 1 then ",..." else "" end)
		, "]"
		, "[\(length)]"
		] | join("")
	  end
	) +
	"> ";

def eval_debug:
	(["DEBUG", .] | tojson, "\n") | stderr;

def eval_f($e; f):
	default_options(build_default_options) as $_
	| try eval($e; "eval_debug") | f
	  catch (. as $err | ("error: " + ($err | tostring)) | println);

def default_display: display({depth: 1});

def eval_print($e):
	eval_f($e; default_display);

# run read-eval-print-loop
def repl($opts): #:: a|(Opts) => @
	def _as_array: if (. | type) != "array" then [.] end;
	def _read_expr: read(prompt; "complete") | trim;
	def _repl:
		. as $c
		| try
			_read_expr as $e
			| if $e != "" then (.[] | eval_print($e)) else empty end,
			_repl
		  catch
			if . == "interrupt" then $c | _repl
			elif . == "eof" then empty
			else error(.) end;
	with_options($opts | .repllevel = options.repllevel+1; _as_array | _repl);
# same as repl({})
def repl: repl({}); #:: a| => @

# TODO: introspect and show doc, reflection somehow?
def help:
    ( builtins[]
	, "^C interrupt"
	, "^D exit REPL"
    )
	| println;

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
			[.[] as $rc | $rc.string | rpad(" "; $rc.maxwidth)] | join("")
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
			"nullinput": {
				short: "-n",
				description: "Null input",
				bool: true
			},
			"decode": {
				short: "-d",
				long: "--decode",
				description: "Decode format",
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
			"compact": {
				short: "-c",
				long: "--compact",
				description: "Compact output",
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
	def _usage($arg0; $version):
		"Usage: \($arg0) [OPTIONS] [EXPR] [FILE...]",
		args_help_text(_opts($version));
	.version as $version
	| .args[0] as $arg0
	| args_parse(.args[1:]; _opts($version)) as {$parsed, $rest}
	| _args($parsed) as $_
	# TODO: hack, pass opts some other way
	| default_options(build_default_options) as $_
	| push_options(
		($parsed.options | parse_options)
		+ {
			repl: $parsed.repl,
			rawstring: ($parsed.rawstring == true),
			compact: ($parsed.compact == true),
			repllevel: 0,
		}
	)
	| if $parsed.version then
		$version | println
	  elif $parsed.formats then
		_formats_list | println
	  elif $parsed.help then
		_usage($arg0; $version) | println
	  else
		try


    # expr file...
	# -n expr
	# -f file expr
	# -nf scriptfile file


	# fq
	# fq . test.mp3
	# fq -i
	# fq -i . test.mp3
	# fq -n 2+2
	# fq -n -i

		  # figure out expression and filenames

		  {
			expr: $rest[0],
			filenames: $rest[1:],
		  }
		  # make -ni and -i without args act the same
		  | if $parsed.nullinput or ($parsed.repl and ($rest | length) == 0) then
		     .expr = $rest[0]
			| .filenames = $rest[1:]
		    end
		  | if $parsed.file then
			 .expr = (open($parsed.file) | string)
			 | .filenames = $rest
		    end
		  | if .filenames == [] then
			  .filenames = ["-"]
		    end
		  | inputs(.filenames) as $_ # store inputs
		  | .expr as $expr
		  | if $parsed.nullinput then null
		    else inputs end # will iterate inputs
		  | eval_f($expr; .)
		  | if $parsed.repl then repl
		    else default_display end
		catch tostring | halt_error(1)
	  end;
