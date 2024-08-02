### basic_prog1.tzx

The `basic_prog1.tzx` test file was created directory from the FUSE emulator.

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


#### Archive Info

The FUSE emulator is not able to add the tape metadata. As this tape block is
very simple, it was added manually using a Hex editor.
