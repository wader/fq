include "internal";
include "options";
include "decode";

def _display_default_opts:
  options({depth: 1});

def _todisplay:
  ( format as $f
  # TODO: not sure about the error check here
  | if $f == null or ._error != null then error("value is not a format root or has errors") end
  | _format_func($f; "_todisplay")
  );

def display($opts):
  ( . as $c
  | options($opts) as $opts
  | try _todisplay catch $c
  | if ($opts.value_output | not) and _can_display then _display($opts)
    else
      ( if _is_string and $opts.raw_string then print
        else _print_color_json($opts)
        end
      , ( $opts.join_string
        | if . then print else empty end
        )
      )
    end
  | error("unreachable")
  );
def display: display({});

def d($opts): display($opts);
def d: display({});
def da($opts): display({array_truncate: 0} + $opts);
def da: da({});
def dd($opts): display({array_truncate: 0, display_bytes: 0} + $opts);
def dd: dd({});
def dv($opts): display({array_truncate: 0, verbose: true} + $opts);
def dv: dv({});
def ddv($opts): display({array_truncate: 0, display_bytes: 0, verbose: true} + $opts);
def ddv: ddv({});

def hexdump($opts): _hexdump(options({display_bytes: 0} + $opts));
def hexdump: hexdump({display_bytes: 0});
def hd($opts): hexdump($opts);
def hd: hexdump;
