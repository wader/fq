# eval f and finally eval fin even on empty or error
def finally(f; fin):
	( try f // (fin | empty)
	  catch (fin as $_ | error(.))
	)
	| fin as $_
	| .;

def _default_options: _eval_state("default_options");
def _default_options($opts): _eval_state("default_options"; $opts);

def _push_options($opts): _eval_state("options_stack"; [$opts] + (_eval_state("options_stack") // []));
def _pop_options: _eval_state("options_stack"; _eval_state("options_stack")[1:]);

def _with_options($opts; f):
	_push_options($opts) as $_ | finally(f; _pop_options);

def _parsed_args: _global_state("parsed_args");
def _parsed_args($v): _global_state("parsed_args"; $v);

def input:
	( _global_state("inputs")
	| if length == 0 then error("break") end
	| [.[0], .[1:]] as [$h, $t]
	| _global_state("input_filename"; $h)
	| _global_state("inputs"; $t)
	| $h
	| open
	| decode(_parsed_args.decode_format)
	);

def inputs:
    def _inputs:
        try input, _inputs
		catch
			if . == "break" then empty
			else
				( (. as $err | ("error: ", ($err | tostring), "\n") | stderr)
				, _inputs
				)
			end;
	_inputs;

def inputs($v):
  ( _global_state("input_filename"; $v[0])
  | _global_state("inputs"; $v)
  );

def input_filename: _global_state("input_filename");

