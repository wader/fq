### basic_prog1.tap

The `basic_prog1.tap` test file was created directory from the FUSE emulator.

Inside the emulated ZX Spectrum a BASIC program was created:

```
10 PRINT "fq is the best!"
20 GOTO 10
```

and saved to tape:

```
SAVE "fqTestProg", LINE 10
```

Then from FUSE select the menu item `Media > Tape > Save As..`.

Any BASIC, machine code, screen image, or other data, can be saved directly
using the `SAVE` command. Further instructions can be found here:
https://worldofspectrum.org/ZXBasicManual/zxmanchap20.html
