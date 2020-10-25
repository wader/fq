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
func (r *Registry) Probe(rootFieldName string, bb *bitbuf.Buffer, forceFormats []*Format) (*Value, interface{}, []error) {
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

	startPos := bb.Pos

	var errs []error
	for _, f := range probeable {
		cbb := bb.Copy()

		// TODO: how to pass regsiters? do later? current field?

		d := (&D{Endian: BigEndian}).FieldStructBitBuf(rootFieldName, cbb)
		decodeErr, dv := d.SafeDecodeFn(f.DecodeFn)
		if decodeErr != nil {
			d.value.Error = decodeErr

			errs = append(errs, decodeErr)
			if !forceOne {
				continue
			}
		}

		// TODO: nicer
		d.value.Desc = f.Name
		d.value.Range = Range{Start: startPos, Stop: cbb.Pos}

		if d.value.Parent == nil {
			d.value.Sort()
		}

		// TODO: wrong keep track of largest?
		_ = cbb.TruncateRel(0)

		return d.value, dv, errs
	}

	return nil, nil, errs
}

func (r *Registry) FindFormat(name string) *Format {
	for _, f := range r.Formats {
		if f.Name == name {
			return f
		}
	}
	return nil
}
