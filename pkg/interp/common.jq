def default_options: _eval_state("default_options");
def default_options($opts): _eval_state("default_options"; $opts);

def push_options($opts): _eval_state("options_stack"; [$opts] + (_eval_state("options_stack") // []));
def pop_options: _eval_state("options_stack"; _eval_state("options_stack")[1:]);

# eval f and finally eval fin even on empty or error
def finally(f; fin):
	( try f // (fin | empty)
	  catch (fin as $_ | error(.))
  )
  | fin as $_
  | .;

def with_options($opts; f):
	push_options($opts) as $_ | finally(f; pop_options);

# TODO: escape for safe key names
# path ["a", 1, "b"] -> "a[1].b"
def path_to_expr:
	map(if type == "number" then "[", ., "]" else ".", . end) | join("");

# TODO: don't use eval? should support '.a.b[1]."c.c"' and escapes?
def expr_to_path:
	if . | type != "string" then error("require string argument") end
	| eval("null | path(\(.))");

def trim: capture("^\\s*(?<str>.*?)\\s*$"; "").str;

# does +1 and [:1] as " "*0 is null
def rpad($s; $w): . + ($s * ($w+1-length))[1:];

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
    if (. | length) == 0 then ""
    else
      _column_widths as $cw
      | . as $rs
      | ( ($rs[]
          | . as $r
          | [ range($r | length) as $i
            | ($r | colmap | {column: $i, string: .[$i], maxwidth: $cw[$i]})
            ]
          | render
          )
        )
      end;

def _parsed_args: _global_state("parsed_args");
def _parsed_args($v): _global_state("parsed_args"; $v);

# TODO: probe format opt
# TODO: isempty?
def input:
  ( _global_state("inputs")
  | if length == 0 then error("break") end
  | [.[0], .[1:]] as [$h, $t]
  | _global_state("input_filename"; $h)
  | _global_state("inputs"; $t)
  | open($h)
  | decode(_parsed_args.decode)
  );

def inputs:
  try repeat(input)
  catch if . == "break" then empty else error end;

def inputs($v):
  ( _global_state("input_filename"; $v[0])
  | _global_state("inputs"; $v)
  );

def input_filename: _global_state("input_filename");

