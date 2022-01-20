package decoders

import (
	"errors"
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/scalar"
	"time"
)

type Precision int

const (
	SECOND = iota
	MILLISECOND
	MICROSECOND
	NANOSECOND
)

func logicalMapperForSchema(schema schema.SimplifiedSchema) scalar.Mapper {
	switch schema.LogicalType {
	case "timestamp":
		return TimestampMapper{Precision: SECOND}
	case "timestamp-millis":
		return TimestampMapper{Precision: MILLISECOND}
	case "timestamp-micros":
		return TimestampMapper{Precision: MICROSECOND}
	case "timestamp-nanos":
		return TimestampMapper{Precision: NANOSECOND}
	case "time":
		return TimeMapper{Precision: SECOND}
	case "time-millis":
		return TimeMapper{Precision: MILLISECOND}
	case "time-micros":
		return TimeMapper{Precision: MICROSECOND}
	case "time-nanos":
		return TimeMapper{Precision: NANOSECOND}
	case "date":
		return DateMapper{}
	default:
		return nil
	}
}

type TimestampMapper struct {
	Precision Precision
}

func (t TimestampMapper) MapScalar(s scalar.S) (scalar.S, error) {
	v := s.ActualS()
	if t.Precision == SECOND {
		s.Sym = time.Unix(v, 0)
	} else if t.Precision == MILLISECOND {
		s.Sym = time.UnixMilli(v)
	} else if t.Precision == MICROSECOND {
		s.Sym = time.UnixMicro(v)
	} else if t.Precision == NANOSECOND {
		s.Sym = time.Unix(0, v)
	} else {
		return s, errors.New("unknown precision")
	}
	s.Sym = time.UnixMilli(v)
	return s, nil
}

type TimeMapper struct {
	Precision Precision
}

func (t TimeMapper) MapScalar(s scalar.S) (scalar.S, error) {
	v, ok := s.Actual.(int64)
	if !ok {
		return s, errors.New("not an int64")
	}

	if t.Precision == SECOND {
		s.Sym = time.Unix(v, 0).Format("15:04:05")
	} else if t.Precision == MILLISECOND {
		s.Sym = time.UnixMilli(v).Format("15:04:05.000")
	} else if t.Precision == MICROSECOND {
		s.Sym = time.UnixMicro(v).Format("15:04:05.000000")
	} else if t.Precision == NANOSECOND {
		s.Sym = time.Unix(0, v).Format("15:04:05.000000000")
	} else {
		return s, errors.New("unknown precision")
	}
	return s, nil
}

type DateMapper struct {
}

func (d DateMapper) MapScalar(s scalar.S) (scalar.S, error) {
	v, ok := s.Actual.(int64)
	if !ok {
		return s, errors.New("not an int64")
	}
	s.Sym = time.Unix(0, 0).AddDate(0, 0, int(v)).Format("2006-01-02")
	return s, nil
}

// Todo Decimal: https://github.com/linkedin/goavro/blob/master/logical_type.go
// Todo Duration
