package interp

import (
	"encoding/json"

	"github.com/wader/gojq"
)

func init() {
	functionRegisterFns = append(functionRegisterFns, func(i *Interp) []Function {
		return []Function{
			{"_query_fromstring", 0, 0, i.queryFromString, nil},
			{"_query_tostring", 0, 0, i.queryToString, nil},
		}
	})
}

func (i *Interp) queryFromString(c interface{}, a []interface{}) interface{} {
	s, err := toString(c)
	if err != nil {
		return err
	}
	q, err := gojq.Parse(s)
	if err != nil {
		p := queryErrorPosition(s, err)
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
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	return v

}

func (i *Interp) queryToString(c interface{}, a []interface{}) interface{} {
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
