package crypto

import (
	"crypto/md5"
	//nolint: gosec
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"embed"
	"fmt"
	"hash"
	"io"

	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/interp"

	//nolint: staticcheck
	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/sha3"
)

//go:embed hash.jq
var hashFS embed.FS

func init() {
	interp.RegisterFunc1("_to_hash", toHash)
	interp.RegisterFS(hashFS)
}

func hashFn(s string) hash.Hash {
	switch s {
	case "md4":
		return md4.New()
	case "md5":
		return md5.New()
	case "sha1":
		return sha1.New()
	case "sha256":
		return sha256.New()
	case "sha512":
		return sha512.New()
	case "sha3_224":
		return sha3.New224()
	case "sha3_256":
		return sha3.New256()
	case "sha3_384":
		return sha3.New384()
	case "sha3_512":
		return sha3.New512()
	default:
		return nil
	}
}

type toHashOpts struct {
	Name string
}

func toHash(_ *interp.Interp, c any, opts toHashOpts) any {
	inBR, err := interp.ToBitReader(c)
	if err != nil {
		return err
	}

	h := hashFn(opts.Name)
	if h == nil {
		return fmt.Errorf("unknown hash function %s", opts.Name)
	}
	if _, err := io.Copy(h, bitio.NewIOReader(inBR)); err != nil {
		return err
	}

	outBR := bitio.NewBitReader(h.Sum(nil), -1)

	bb, err := interp.NewBinaryFromBitReader(outBR, 8, 0)
	if err != nil {
		return err
	}
	return bb
}
