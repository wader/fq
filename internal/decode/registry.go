package decode

import (
	"fmt"
	"fq/internal/bitbuf"
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
func (r *Registry) Probe(parent Decoder, bb *bitbuf.Buffer, forceFormats []*Format) (Decoder, []error) {
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
		// TODO: how to pass regsiters? do later? current field?
		d := f.New()
		rootField := &Field{Name: f.Name}
		d.Prepare(Common{
			Parent:   parent,
			Format:   f,
			Registry: r,
			BitBuf:   bb.Copy(),

			Root:    rootField,
			Current: rootField,
		})
		err := func() (err error) {
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

		if err != nil {
			errs = append(errs, err)
			if !forceOne {
				continue
			}
		}

		return d, errs
	}

	return nil, errs
}

func (r *Registry) FindFormat(name string) *Format {
	for _, f := range r.Formats {
		if f.Name == name {
			return f
		}
	}
	return nil
}
