# eval is implemented as an internal function evaluting $e for input and
# returns an array with all generated values, we then each over the values
# to make it behave as a normal jq generator.
def eval($e): _eval($e)[];

def trim: capture("^\\s*(?<a>.*?)\\s*$"; "").a;

# does +1 and [:1] as " "*0 is null
def rpad($w;$s): . + ($s * (([0,$w-(.|length)] | max)+1))[1:];
