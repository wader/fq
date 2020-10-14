package decode

import (
	"fmt"
	"fq/pkg/bitbuf"
	"runtime"
)

type Registry struct {
	Formats []*Format
}

func NewRegistryWithFormats(formats []*Format) *Registry {
	r := &Registry{
		Formats: formats,
	}

	return r
}

type ProbeError struct {
	Format        *Format
	Err           error
	PanicHandeled bool
	PanicStack    string
}

func (pe *ProbeError) Error() string { return fmt.Sprintf("%s probe: %s", pe.Format.Name, pe.Err) }
func (pe *ProbeError) Unwrap() error { return pe.Err }

// Probe probes all probeable formats and turns first found Decoder and all other decoder errors
func (r *Registry) Probe(parent Decoder, rootFieldName string, parentRange Range, bb *bitbuf.Buffer, forceFormats []*Format) (*Field, int64, Decoder, []error) {
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
		rootField := &Field{
			Name: rootFieldName,
			Value: Value{
				V:      []*Field{},
				Range:  Range{}, // TODO:
				BitBuf: cbb,
				Desc:   f.Name,
			},
		}
		var common *Common
		d.Prepare(func(c *Common) {
			common = c
		})

		common.Parent = parent
		common.format = f
		common.registry = r
		common.bitBuf = cbb
		common.root = rootField
		common.current = rootField

		decodeErr := func() (err error) {
			defer func() {
				if recoverErr := recover(); recoverErr != nil {
					// https://github.com/golang/go/blob/master/src/net/http/server.go#L1770
					const size = 64 << 10
					buf := make([]byte, size)
					buf = buf[:runtime.Stack(buf, false)]

					pe := &ProbeError{
						Format:     f,
						PanicStack: string(buf),
					}
					switch panicErr := recoverErr.(type) {
					case BitBufError:
						pe.Err = panicErr
						pe.PanicHandeled = true
					case ValidateError:
						pe.Err = panicErr
						pe.PanicHandeled = true
					default:
						pe.Err = fmt.Errorf("%s", panicErr)
						pe.PanicHandeled = false
					}

					err = pe
				}
			}()

			d.Decode()

			return nil
		}()

		if decodeErr != nil {
			common.current.Error = decodeErr

			errs = append(errs, decodeErr)
			if !forceOne {
				continue
			}
		}

		// TODO: will resort
		rootField.Sort()
		// TODO: wrong keep track of largest?
		_ = cbb.TruncateRel(0)

		return rootField, cbb.Pos, d, errs
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
