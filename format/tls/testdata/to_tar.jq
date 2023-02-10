def to_tar(g):
  def lpad($l; $n): [[range($l-length) | $n], .];
  def rpad($l; $n): [., [range($l-length) | $n]];
  def header($filename; $b):
    def checksum: [.[range(.size)]] | add;
    def h:
      [ ($filename | rpad(100; 0)) # name
      , ("000644 " | rpad(8; 0)) # mode
      , ("000000 " | rpad(8; 0)) # uid
      , ("000000 " | rpad(8; 0)) # gid
      , [($b.size | to_radix(8) | [lpad(11; "0")]), " "]  # size
      , [(0| to_radix(8) | lpad(11; "0")), " "] # mtime
      , "        " # chksum (blank spaces when adding checksum)
      , ("0") # typeflag
      , ("" | rpad(100; 0)) # linkname
      , ["ustar", 0] # magic
      , ("00") # version
      , ("user" | rpad(32; 0)) # uname
      , ("group" | rpad(32; 0)) # gname
      , ("000000 " | rpad(8; 0)) # devmajor
      , ("000000 " | rpad(8; 0)) # devminor
      , ("" | rpad(155; 0)) # prefix
      ] | tobytes;
    ( h as $h
    | [ $h[0:148]
      , [(($h | checksum) | to_radix(8) | lpad(6; "0")), 0, " "]
      , $h[148+8:]
      ]
    | tobytes
    );
  [ ( # per file
      ( g as {$filename, $data}
      | ($data | tobytes) as $b
      | ($filename | rpad(100; 0)) # name
      | header($filename; $b) as $header
      | $header
      , ("" | lpad((512 - ($header.size % 512)) % 512; 0))
      , $b
      , ("" | lpad((512 - ($b.size % 512)) % 512; 0))
      )
      # end_marker
    , [range(1024) | 0]
    )
  ] | tobytes;