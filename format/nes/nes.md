### Limitations

- `prg_rom`, `chr_rom` and `trainer` fields may contain data that is just random
  junk from the memory chips, since they are of a fixed size.
- The `nes_toasm` function outputs ALL opcodes, including the unofficial ones,
  which means that none of the regular assemblers can recompile it.
- The `nes_tokitty` function works on tiles in `chr_rom` but only outputs a Kitty
  graphics compatible string. You need to manually `printf` that string to get
  Kitty (or another compatible terminal) to output the graphics.

### Decompile PRG ROM
```
$ fq -r '.prg_rom[] | nes_toasm' file.nes
```

### Print out first CHR ROM tile in Kitty (or Konsole, wayst, WezTerm) at size 5
```
$ printf $(fq -r -d nes '.chr_rom[0] | nes_tokitty(5)' file.nes)
```

### Print out all CHR ROM tiles in Kitty (with Bash) at size 5
```
$ for line in $(fq -r '.chr_rom[] | nes_tokitty(5)' file.nes);do printf "%b%s" "$line";done
```

### Authors
- Mikael Lofjärd mikael.lofjard@gmail.com, original author

### References
- https://www.nesdev.org/wiki/INES
- https://www.nesdev.org/wiki/NES_2.0
- https://www.nesdev.org/wiki/CPU
- https://bugzmanov.github.io/nes_ebook/chapter_6_3.html
