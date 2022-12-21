# https://en.wikipedia.org/wiki/Privacy-Enhanced_Mail
def from_pem:
  ( tobytes
  | tostring
  | capture("-----BEGIN(.*?)-----(?<s>.*?)-----END(.*?)-----"; "mg").s
  | _from_base64({encoding: "std"})
  ) // error("no pem header or footer found");

def to_pem($label):
  ( tobytes
  | _to_base64({encoding: "std"})
  | ($label | if $label != "" then " " + $label end) as $label
  | [ "-----BEGIN\($label)-----"
    , .
    , "-----END\($label)-----"
    , ""
    ]
  | join("\n")
  );
def to_pem: to_pem("");