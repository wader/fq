def _macho__help:
  { notes: "Supports decoding vanilla and FAT Mach-O binaries.",
    examples: [
      {comment: "Select 64bit load segments", shell: "fq '.load_commands[] | select(.cmd==\"segment_64\")' file"}
    ],
    links: [
      {url: "https://github.com/aidansteele/osx-abi-macho-file-format-reference"}
    ]
  };
