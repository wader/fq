def dv($p):
    . as $c | [$p, $c] | debug | $c;

def trim: capture("^\\s*(?<a>.*?)\\s*$"; "").a;

# does +1 and [:1] as " "*0 is null
def rpad($w;$s): . + ($s * (([0,$w-(.|length)] | max)+1))[1:];
