# convert curl trace to {send: <binary>, recv: <binary>}
# Trace format:
# == Info:   Trying 0.0.0.0:8080...
# == Info: Connected to 0.0.0.0 (127.0.0.1) port 8080 (#0)
# => Send header, 156 bytes (0x9c)
# 0000: 50 4f 53 54 20 2f 6f 6b 20 48 54 54 50 2f 31 2e POST /ok HTTP/1.
# <= Recv header, 17 bytes (0x11)
# 0000: 48 54 54 50 2f 31 2e 31 20 32 30 30 20 4f 4b 0d HTTP/1.1 200 OK.
# 0010: 0a                                              .
def from_curl_trace:
  ( reduce split("\n")[] as $l (
      {state: "send",  send: [], recv: []};
      if $l | startswith("=>") then .state="send"
      elif $l | startswith("<=") then .state="recv"
      elif $l | test("^\\d") then .[.state] += [$l]
      end
    )
  | (.send, .recv) |=
      ( map(capture(": (?<hex>.{1,47})").hex
      | gsub(" "; ""))
      | add
      | hex
      )
  );
