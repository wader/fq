# https://en.wikipedia.org/wiki/Privacy-Enhanced_Mail
def frompem:
  ( tobytes
  | tostring
  | capture("-----BEGIN(.*?)-----(?<s>.*?)-----END(.*?)-----"; "mg").s
  | _frombase64({encoding: "std"})
  ) // error("no pem header or footer found");

def topem($label):
  ( tobytes
  | _tobase64({encoding: "std"})
  | ($label | if $label != "" then " " + $label end) as $label
  | [ "-----BEGIN\($label)-----"
    , .
    , "-----END\($label)-----"
    , ""
    ]
  | join("\n")
  );
def topem: topem("");