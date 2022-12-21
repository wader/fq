# read by jq-lsp to add additional builtins
def _can_display: empty;
def _decode($format; $opts): empty;
def _display($opts): empty;
def _eval($expr; $opts): empty;
def _extkeys: empty;
def _exttype: empty;
def _format_func($format; $func): empty;
def _global_state: empty;
def _global_state($v): empty;
def _hexdump($opts): empty;
def _is_completing: empty;
def _match_binary($regexp; $flags): empty;
def _print_color_json($opts): empty;
def _query_fromstring: empty;
def _query_tostring: empty;
def _readline: empty;
def _readline($opts): empty;
def _registry: empty;
def _stdio_info($name): empty;
def _stdio_read($name; $l): empty;
def _stdio_write($name): empty;
def _tobits($opts): empty;
def _tovalue($opts): empty;
def open: empty;
def scope: empty;

# used by help.jq
def to_jq: empty;
# used by funcs.jq iprint
def to_radix($base): empty;