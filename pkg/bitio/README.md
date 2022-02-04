The bitio package tries to mimic the standard library packages io and bytes as much as possible.

- bitio.Buffer same as bytes.Buffer
- bitio.IOBitReadSeeker is a bitio.ReaderAtSeeker that from a io.ReadSeeker
- bitio.IOBitWriter a bitio.BitWriter that write bytes to a io.Writer, use Flush() to write possible unaligned byte
- bitio.IOReader is a io.Reader that reads bytes from a bit reader, will zero pad on unaligned byte eof
- bitio.IOReadSeeker is a io.ReadSeeker that read/seek bytes in a bit stream, will zero pad on unaligned - bitio.NewBitReader same as bytes.NewReader
- bitio.LimitReader same as io.LimitReader
- bitio.MultiReader same as io.MultiReader
- bitio.SectionReader same as io.SectionReader
- bitio.Copy* same as io.Copy*
- bitio.ReadFull same as io.ReadFull

TODO:
- bitio.IOBitReader bitio.Reader that reads from a io.Reader
- bitio.IOBitWriteSeeker bitio.BitWriteSeeker that writes to a io.WriteSeeker
- bitio.CopyN
- speed up read by using a cache somehow ([]byte or just a uint64?)
