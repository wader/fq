package wasm

import (
	"errors"

	"github.com/wader/fq/pkg/scalar"
)

var sectionIDToSym = &sectionIDToSymMapper{}

type sectionIDToSymMapper struct {
}

func (m *sectionIDToSymMapper) MapScalar(s scalar.S) (scalar.S, error) {
	v, ok := s.Actual.(uint64)
	if !ok {
		return s, errors.New("unexpected data type for section ID")
	}

	switch v {
	case sectionIDCustom:
		s.Sym = "custom section"
	case sectionIDType:
		s.Sym = "type section"
	case sectionIDImport:
		s.Sym = "import section"
	case sectionIDFunction:
		s.Sym = "function section"
	case sectionIDTable:
		s.Sym = "table section"
	case sectionIDMemory:
		s.Sym = "memory section"
	case sectionIDGlobal:
		s.Sym = "global section"
	case sectionIDExport:
		s.Sym = "export section"
	case sectionIDStart:
		s.Sym = "start section"
	case sectionIDElement:
		s.Sym = "element section"
	case sectionIDCode:
		s.Sym = "code section"
	case sectionIDData:
		s.Sym = "data section"
	case sectionIDDataCount:
		s.Sym = "data count section"
	default:
		s.Sym = "unknown section"
	}

	return s, nil
}

var valtypeToSymMapper = &valtypeToSym{}

type valtypeToSym struct {
}

// valtype ::= t:numtype => t
//          |  t:vectype => t
//          |  t:reftype => t
func (m *valtypeToSym) MapScalar(s scalar.S) (scalar.S, error) {
	v, ok := s.Actual.(uint64)
	if !ok {
		return s, errors.New("unexpected data type for valtype")
	}

	switch v {
	case 0x7f:
		s.Sym = "i32"
	case 0x7e:
		s.Sym = "i64"
	case 0x7d:
		s.Sym = "f32"
	case 0x7c:
		s.Sym = "f64"
	case 0x7b:
		s.Sym = "v128"
	case 0x70:
		s.Sym = "funcref"
	case 0x6f:
		s.Sym = "externref"
	default:
		s.Sym = "unknown valtype"
	}

	return s, nil
}

var importdescTagToSym = &importdescTagToSymMapper{}

type importdescTagToSymMapper struct {
}

func (m *importdescTagToSymMapper) MapScalar(s scalar.S) (scalar.S, error) {
	v, ok := s.Actual.(uint64)
	if !ok {
		return s, errors.New("unexpected data type for importdesc tag")
	}

	switch v {
	case 0x00:
		s.Sym = "func"
	case 0x01:
		s.Sym = "table"
	case 0x02:
		s.Sym = "mem"
	case 0x03:
		s.Sym = "global"
	default:
		s.Sym = "unknown importdesc"
	}

	return s, nil
}

var exportdescTagToSym = &exportdescTagToSymMapper{}

type exportdescTagToSymMapper struct {
}

func (m *exportdescTagToSymMapper) MapScalar(s scalar.S) (scalar.S, error) {
	v, ok := s.Actual.(uint64)
	if !ok {
		return s, errors.New("unexpected data type for exportdesc tag")
	}

	switch v {
	case 0x00:
		s.Sym = "funcidx"
	case 0x01:
		s.Sym = "tableidx"
	case 0x02:
		s.Sym = "memidx"
	case 0x03:
		s.Sym = "globalidx"
	default:
		s.Sym = "unknown exportdesc"
	}

	return s, nil
}

var reftypeTagToSym = &reftypeTagToSymMapper{}

type reftypeTagToSymMapper struct {
}

// reftype ::= 0x70 => funcref
//          |  0x6F => externref
func (m *reftypeTagToSymMapper) MapScalar(s scalar.S) (scalar.S, error) {
	v, ok := s.Actual.(uint64)
	if !ok {
		return s, errors.New("unexpected data type for reftype tag")
	}

	switch v {
	case 0x70:
		s.Sym = "funcref"
	case 0x6f:
		s.Sym = "externref"
	default:
		s.Sym = "unknown reftype"
	}

	return s, nil
}

var mutToSym = &mutToSymMapper{}

type mutToSymMapper struct {
}

func (m *mutToSymMapper) MapScalar(s scalar.S) (scalar.S, error) {
	v, ok := s.Actual.(uint64)
	if !ok {
		return s, errors.New("unexpected data type for mut")
	}

	switch v {
	case 0x00:
		s.Sym = "const"
	case 0x01:
		s.Sym = "var"
	default:
		s.Sym = "unknown mut"
	}

	return s, nil
}

var elemkindTagToSym = &elemkindTagToSymMapper{}

type elemkindTagToSymMapper struct {
}

func (m *elemkindTagToSymMapper) MapScalar(s scalar.S) (scalar.S, error) {
	v, ok := s.Actual.(uint64)
	if !ok {
		return s, errors.New("unexpected data type for elemkind tag")
	}

	switch v {
	case 0x00:
		s.Sym = "funcref"
	default:
		s.Sym = "unknown elemkind"
	}

	return s, nil
}
