# Decoders take byte strings as input and output so-called I/O objects.
#
# An I/O object is either a byte string or an object {i, o}, where
# i is the consumed input and
# o is the produced output.
#
# Additionally, an I/O object {i, o} may contain a key "errs" that
# stores an array of objects {path, desc}, where
# path is an array to the source of the error, and
# desc is a description of the error (a string)
#
# A decoder that fails should add an error to errs (with empty path),
# then throw {i, errs} as error.
# This allows errors to be detected easily.

def isobject: type == "object";
def has_errs: isobject and has("errs");
def catch_errs(f): try f catch (if has_errs then . else error end);

def assert_len($i; $w):
  if $i | length != $w then
    .errs += [{desc: "expected \($w) bytes, found \($i | length)"}] | error
  end;

def update_output($o; $k; entry):
  ($o | has_errs) as $has_errs |
  if $has_errs then .errs += ($o.errs | (.[].path += [$k])) end |
  .v += entry |
  if $has_errs then error end;

def update_output2($o; entry):
  ($o | has_errs) as $has_errs |
  if $has_errs then .errs += ($o.errs | (.[].path += [])) end |
  .v += entry |
  if $has_errs then error end;

# Take an I/O object, consume exactly $w bytes from its remaining input and
# add a field $k to its output with the consumed bytes fed to f.
#
# If f yields an error, this adds the error to the I/O object
# with the appropriate path, and rethrows the I/O object as error.
def take($k; $w; f):
  .i[:$w] as $i | debug({take: {$k, $w, l: length}}) | assert_len($i; $w) | .i |= .[$w:] |
  catch_errs($i | f + {o: $i.start, l: $i.stop-$i.start}) as $v |
  update_output($v; $k; {($k): $v});

def take2($w; f):
  .i[:$w] as $i | debug({take2: {dot: .,$w, l: ($i | length)}})  | .i |= .[$w:] |
  catch_errs($i | f + {o: $i.start, l: $i.stop-$i.start}) as $v |
  update_output2($v; [$v]);

# Translate from i denoting remaining input to i denoting consumed input.
def set_consumed($input): .i |= .[:byteoffset($input)];

def takes(f):
  def rec:
    . as $input |
    if length == 0 then empty # we're done
    else f | set_consumed($input), (.i | rec)
    end;
  rec;

def rtrim_nul: (index(0 | tobytes) // length) as $i | .[:$i];

def u24be($k): take($k; 3; {i: ., v: (.[0]*65366 + .[1]*256 + .[2])});
def u32be($k): take($k; 4; {i: ., v: (.[0]*16777216 + .[1]*65366 + .[2]*256 + .[3])});
def str($k; $w): take($k; $w; {i: ., v: rtrim_nul});
def str($w): take2($w; {i: ., v: rtrim_nul});
def raw($k; $w): take($k; $w; {i: ., v: null});

def many($k; $w; f): take($k; $w; takes(f));

def decode_box:
  . as $input |
  {i: ., l:0, o:0} |
  # we use i as remaining input, which we later use to infer consumed input
  u32be("size") |
  str("type"; 4) |
  (.v.size.v-8) as $size |
  .v.type.v as $type |
  if $type == "ftyp" then
    str("major_brand"; 4) |
    u32be("minor_version") |
    many("brands"; $size-8; takes({i: ., l:0, o:0} | str(4)))
  elif $type == "moov" then many("boxes"; $size; decode_box)
  else raw("data"; $size)
  end
  # raw("data"; .v.size.s) | # use sym (decimal number)
;

def decode_mp4:
  tobytes |
  decode_box
;
