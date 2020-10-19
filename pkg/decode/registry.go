package decode

import (
	"fmt"
	"fq/pkg/bitbuf"
)

type ProbeError struct {
	Format        *Format
	Err           error
	PanicHandeled bool
	PanicStack    string
}

func (pe *ProbeError) Error() string { return fmt.Sprintf("%s probe: %s", pe.Format.Name, pe.Err) }
func (pe *ProbeError) Unwrap() error { return pe.Err }

type Registry struct {
	Formats []*Format
}

func NewRegistryWithFormats(formats []*Format) *Registry {
	r := &Registry{
		Formats: formats,
	}

	return r
}

// Probe probes all probeable formats and turns first found Decoder and all other decoder errors
func (r *Registry) Probe(parent Decoder, rootFieldName string, parentRange Range, bb *bitbuf.Buffer, forceFormats []*Format) (*Value, int64, Decoder, []error) {
	var probeable []*Format
	var forceOne = len(forceFormats) == 1
	if forceFormats != nil {
		probeable = forceFormats
	} else {
		for _, f := range r.Formats {
			if f.SkipProbe {
				continue
			}
			probeable = append(probeable, f)
		}
	}

	// TODO: order..

	var errs []error
	for _, f := range probeable {
		cbb := bb.Copy()

		// TODO: how to pass regsiters? do later? current field?
		d := f.New()
		rootValue := &Value{
			V:      Struct{}, // TODO: array or struct?
			Range:  Range{},  // TODO:
			BitBuf: cbb,
			Name:   rootFieldName,
			Desc:   f.Name,
		}
		common := d.GetCommon()
		common.registry = r
		common.bitBuf = cbb
		common.current = rootValue

		decodeErr := d.GetCommon().SafeDecodeFn(d.Decode)
		// TODO: wrap in ProbeError?

		if decodeErr != nil {
			common.current.Error = decodeErr

			errs = append(errs, decodeErr)
			if !forceOne {
				continue
			}
		}

		// TODO: will resort
		rootValue.Sort()
		// TODO: wrong keep track of largest?
		_ = cbb.TruncateRel(0)

		return rootValue, cbb.Pos, d, errs
	}

	return nil, 0, nil, errs
}

func (r *Registry) FindFormat(name string) *Format {
	for _, f := range r.Formats {
		if f.Name == name {
			return f
		}
	}
	return nil
}
