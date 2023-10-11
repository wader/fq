def _esc: "\u001b";
def _ansi_codes:
  { bold: {set: "1", reset: "0"}
  , faint: {set: "2", reset: "0"}
  , black: {set: "30", reset: "39"}
  , red: {set: "31", reset: "39"}
  , green: {set: "32", reset: "39"}
  , yellow: {set: "33", reset: "39"}
  , blue: {set: "34", reset: "39"}
  , magenta: {set: "35", reset: "39"}
  , cyan: {set: "36", reset: "39"}
  , white: {set: "37", reset: "39"}
  , brightblack: {set: "90", reset: "39"}
  , brightred: {set: "91", reset: "39"}
  , brightgreen: {set: "92", reset: "39"}
  , brightyellow: {set: "93", reset: "39"}
  , brightblue: {set: "94", reset: "39"}
  , brightmagenta: {set: "95", reset: "39"}
  , brightcyan: {set: "96", reset: "39"}
  , brightwhite: {set: "97", reset: "39"}
  , default: {set: "39", reset: "39"}
  , bgblack: {set: "40", reset: "49"}
  , bgred: {set: "41", reset: "49"}
  , bggreen: {set: "42", reset: "49"}
  , bgyellow: {set: "43", reset: "49"}
  , bgblue: {set: "44", reset: "49"}
  , bgmagenta: {set: "45", reset: "49"}
  , bgcyan: {set: "46", reset: "49"}
  , bgwhite: {set: "47", reset: "49"}
  , bgbrightblack: {set: "100", reset: "49"}
  , bgbrightred: {set: "101", reset: "49"}
  , bgbrightgreen: {set: "102", reset: "49"}
  , bgbrightyellow: {set: "103", reset: "49"}
  , bgbrightblue: {set: "104", reset: "49"}
  , bgbrightmagenta: {set: "105", reset: "49"}
  , bgbrightcyan: {set: "106", reset: "49"}
  , bgbrightwhite: {set: "107", reset: "49"}
  , bold: {set: "1", reset: "22"}
  , italic: {set: "3", reset: "23"}
  , underline: {set: "4", reset: "24"}
  , inverse: {set: "7", reset: "27"}
  };

def _ansi_if($opts; $name):
  if $opts.color then
    ( ( $opts.colors[$name]
      | split("+")
      | reduce map(_ansi_codes[.])[] as $c (
          {set: [], reset: []};
          ( .set += [$c.set]
          | .reset += [$c.reset]
          )
        )
      ) as {$set, $reset}
    | "\(_esc)[\($set| join(";"))m\(.)\(_esc)[\($reset | join(";"))m"
    )
  end;

def _ansi:
  { clear_line: "\(_esc)[2K"
  };
