# eval is implemented as an internal function evaluting $e for input and
# returns an array with all generated values, we then each over the values
# to make it behave as a normal jq generator.
def eval($e): _eval($e)[];

def print: _print[];

def display: _display[];
def display($opts): _display($opts)[];
def d: _display[];
def d($opts): _display($opts)[];

def verbose: _verbose[];
def verbose($opts): _verbose($opts)[];
def v: _verbose[];
def v($opts): _verbose($opts)[];

def preview: _preview[];
def preview($opts): _preview($opts)[];
def p: _preview[];
def p($opts): _preview($opts)[];

def hexdump: _hexdump[];
def hexdump($opts): _hexdump($opts)[];
def hd: _hexdump[];
def hd($opts): _hexdump($opts)[];
def h: _hexdump[];
def h($opts): _hexdump($opts)[];

def trim: capture("^\\s*(?<a>.*?)\\s*$"; "").a;

# does +1 and [:1] as " "*0 is null
def rpad($s;$w): . + ($s * ($w+1-length))[1:];

def maybe_each: if (. | type) == "array" then .[] end;

# [{a: 123, ...}, ...]
# colmap maps something into [col, ...]
# render maps [{string: "coltext", maxwidth: 12}, ..] into a row string
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
