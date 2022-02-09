// Package bitio tries to mimic the standard library packages io and bytes but for bits.
//
// - bitio.Buffer same as bytes.Buffer
//
// - bitio.IOBitReadSeeker is a bitio.ReaderAtSeeker that reads from a io.ReadSeeker
//
// - bitio.IOBitWriter a bitio.BitWriter that writes to a io.Writer, use Flush() to write possible zero padded unaligned byte
//
// - bitio.IOReader is a io.Reader that reads bytes from a bitio.Reader, will zero pad unaligned byte at EOF
//
// - bitio.IOReadSeeker is a io.ReadSeeker that reads from a bitio.ReadSeeker, will zero pad unaligned byte at EOF
//
// - bitio.NewBitReader same as bytes.NewReader
//
// - bitio.LimitReader same as io.LimitReader
//
// - bitio.MultiReader same as io.MultiReader
//
// - bitio.SectionReader same as io.SectionReader
//
// - bitio.Copy* same as io.Copy*
//
// - bitio.ReadFull same as io.ReadFull
//
// TODO:
//
// - bitio.IOBitReader bitio.Reader that reads from a io.Reader
//
// - bitio.IOBitWriteSeeker bitio.BitWriteSeeker that writes to a io.WriteSeeker
//
// - bitio.CopyN
//
// - Speed up read by using a cache somehow ([]byte or just a uint64?)
package bitio
