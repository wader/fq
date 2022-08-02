package interp

import (
	"encoding/json"

	"github.com/wader/gojq"
)

func init() {
	RegisterFunc0("_query_fromstring", (*Interp)._queryFromString)
	RegisterFunc0("_query_tostring", (*Interp)._queryToString)
}

func (i *Interp) _queryFromString(c string) any {
	q, err := gojq.Parse(c)
	if err != nil {
		p := queryErrorPosition(c, err)
		return compileError{
			err:  err,
			what: "parse",
			pos:  p,
		}
	}

	// TODO: use mapstruct?
	b, err := json.Marshal(q)
	if err != nil {
		return err
	}
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	return v
}

func (i *Interp) _queryToString(c any) any {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	var q gojq.Query
	if err := json.Unmarshal(b, &q); err != nil {
		return err
	}

	return q.String()
}
