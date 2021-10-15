package gojqextra

import (
	"fmt"
	"math/big"

	"github.com/wader/gojq"
)

// many of these based on errors from gojq
// TODO: refactor to use errors from gojq?
// TODO: preview from gojq?

type NonUpdatableTypeError struct {
	Typ string
	Key string
}

func (err NonUpdatableTypeError) Error() string {
	return fmt.Sprintf("update key %v cannot be applied to: %s", err.Key, err.Typ)
}

type FuncTypeError struct {
	Name string
	V    interface{}
}

func (err FuncTypeError) Error() string { return err.Name + " cannot be applied to: " + typeof(err.V) }

type FuncTypeNameError struct {
	Name string
	Typ  string
}

func (err FuncTypeNameError) Error() string {
	return err.Name + " cannot be applied to: " + err.Typ
}

type ExpectedObjectError struct {
	Typ string
}

func (err ExpectedObjectError) Error() string {
	return "expected an object but got: " + err.Typ
}

type ExpectedArrayError struct {
	Typ string
}

func (err ExpectedArrayError) Error() string {
	return "expected an array but got: " + err.Typ
}

type ExpectedObjectWithKeyError struct {
	Typ string
	Key string
}

func (err ExpectedObjectWithKeyError) Error() string {
	return fmt.Sprintf("expected an object with key %q but got: %s", err.Key, err.Typ)
}

type ExpectedArrayWithIndexError struct {
	Typ   string
	Index int
}

func (err ExpectedArrayWithIndexError) Error() string {
	return fmt.Sprintf("expected an array with index %d but got: %s", err.Index, err.Typ)
}

type IteratorError struct {
	Typ string
}

func (err IteratorError) Error() string {
	return "cannot iterate over: " + err.Typ
}

type HasKeyTypeError struct {
	L, R string
}

func (err HasKeyTypeError) Error() string {
	return "cannot check whether " + err.L + " has a key: " + err.R
}

type ArrayIndexTooLargeError struct {
	V interface{}
}

func (err *ArrayIndexTooLargeError) Error() string {
	return fmt.Sprintf("array index too large: %v", err.V)
}

func typeof(v interface{}) string {
	switch v := v.(type) {
	case nil:
		return "null"
	case bool:
		return "boolean"
	case int, float64, *big.Int:
		return "number"
	case string:
		return "string"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	case gojq.JQValue:
		return fmt.Sprintf("JQValue(%s)", v.JQValueType())
	default:
		panic(fmt.Sprintf("invalid value: %v", v))
	}
}
