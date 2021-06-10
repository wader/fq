def default_options: _state("default_options");
def default_options($opts): _state("default_options"; $opts);

def push_options($opts): _state("options_stack"; [$opts] + (_state("options_stack") // []));
def pop_options: _state("options_stack"; _state("options_stack")[1:]);

# eval f and finally eval fin even on empty or error
def finally(f; fin):
	( try f // (fin | empty)
	  catch (fin as $_ | error(.))
  )
  | fin as $_
  | .;

def with_options($opts; f):
	push_options($opts) as $_ | finally(f; pop_options);

def trim: capture("^\\s*(?<str>.*?)\\s*$"; "").str;

# does +1 and [:1] as " "*0 is null
def rpad($s;$w): . + ($s * ($w+1-length))[1:];

# [{a: 123, ...}, ...]
# colmap maps something into [col, ...]
# render maps [{column: 0, string: "coltext", maxwidth: 12}, ..] into a row
def table(colmap;render):
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
