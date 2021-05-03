include "@builtin/common.jq";

def args_parse($args;$opts):
	def _parse($args;$flagmap;$r):
		def _parse_with_arg($new_args;$optname;$value;$opt):
			if $opt.object then
				( $value
				| capture("^(?<key>.*?)=(?<value>.*)$")
					// error("\($value): should be key=value")
				)
				as {$key, $value} |
                # TODO: validate option name key
				_parse($new_args;$flagmap;($r|.parsed.[$optname][$key] |= $value))
			elif $opt.array then
				_parse($new_args;$flagmap;($r|.parsed.[$optname] += [$value]))
			else
				_parse($new_args;$flagmap;($r|.parsed.[$optname] = $value))
			end;
		def _parse_without_arg($new_args;$optname):
			_parse($new_args;$flagmap;($r|.parsed.[$optname] = true));
		($args[0] | index("=")) as $assign_i
		| (if $assign_i then $args[0][0:$assign_i]
		   else $args[0] end
		  ) as $arg
		| if $arg == null then
			$r
		  else
			if $arg == "--" then
				$r | .rest += $args[1:]
			elif $arg | test("^--?[^-]") then
				$flagmap[$arg] as $optname
				| ($opts[$optname]? // null) as $opt
				| if $opt == null then
					if $arg | test("^-[^-]") then
						$arg[0:2] as $arg
						| $flagmap[$arg] as $optname
						| ($opts[$optname]? // null) as $opt
						| if $opt and $opt.bool then
							_parse_without_arg((["-"+$args[0][2:]]+$args[1:]);$optname)
						  else
							error("\($arg): needs an argument")
						  end
					else
						error("\($arg): no such argument")
					end
				  elif $opt.string or $opt.array or $opt.object then
					if $assign_i then
						_parse_with_arg($args[1:];$optname;$args[0][$assign_i+1:];$opt)
					elif ($args | length) < 2 then
						error("\($arg): needs an argument")
					else
						_parse_with_arg($args[2:];$optname;$args[1];$opt)
					end
				  else
					if $assign_i then error("\($arg): takes no argument")
					else _parse_without_arg($args[1:];$optname) end
				  end
			else
				_parse($args[1:];$flagmap;($r|.rest += [$args[0]]))
			end
		end;
	# build {"-s": "name", "--long": "name", ...}
	def _flagmap:
		( $opts
		| to_entries
		| map(
			({(.value.short): .key}? // {}) +
			({(.value.long): .key}? // {})
		  )
		| add
		);
	def _defaults:
		( $opts
		| to_entries
		| map(select(.value.default))
		| map({(.key): .value.default})
		| add
		);
	_parse($args;_flagmap;{parsed: _defaults, rest: []});

def args_help_text($opts):
	def _opthelp:
		( [ .long
		  , .short
		  ] | map(select(strings)) | join(",")
		) +
		if .value or .array or .object then "=ARG"
		else null end;
	def _maxoptlen:
		[$opts[] | (.|_opthelp|length)] | max;
	_maxoptlen as $l
	| $opts
	| to_entries[]
	| (.value | .help_default // .default) as $default
	| [ "\(.value|_opthelp|rpad(" ";$l))  \(.value.description)"
	  , if $default then
			if .value.object then
				[ "\n",
				  ([$default | to_entries[] | "\(" "*$l)    \(.key)=\(.value)"] | join("\n"))
				]
			else
				" (\($default))"
			end
		else
			""
	    end
	  ]
	| flatten
	| join("");
