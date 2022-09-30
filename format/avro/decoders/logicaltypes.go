package decoders

import (
	"errors"
	"time"

	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/scalar"
)

type Precision int

const (
	SECOND = iota
	MILLISECOND
	MICROSECOND
	NANOSECOND
)

func logicalTimeMapperForSchema(schema schema.SimplifiedSchema) scalar.SintMapper {
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

func (t TimestampMapper) MapSint(s scalar.Sint) (scalar.Sint, error) {
	v := s.Actual
	var ts time.Time
	if t.Precision == SECOND {
		ts = time.Unix(v, 0)
	} else if t.Precision == MILLISECOND {
		ts = time.UnixMilli(v)
	} else if t.Precision == MICROSECOND {
		ts = time.UnixMicro(v)
	} else if t.Precision == NANOSECOND {
		ts = time.Unix(0, v)
	} else {
		return s, errors.New("unknown precision")
	}
	s.Sym = ts.UTC().Format(time.RFC3339Nano)
	return s, nil
}

type TimeMapper struct {
	Precision Precision
}

func (t TimeMapper) MapSint(s scalar.Sint) (scalar.Sint, error) {
	v := s.Actual

	if t.Precision == SECOND {
		s.Sym = time.Unix(v, 0).UTC().Format("15:04:05")
	} else if t.Precision == MILLISECOND {
		s.Sym = time.UnixMilli(v).UTC().Format("15:04:05.000")
	} else if t.Precision == MICROSECOND {
		s.Sym = time.UnixMicro(v).UTC().Format("15:04:05.000000")
	} else if t.Precision == NANOSECOND {
		s.Sym = time.Unix(0, v).UTC().Format("15:04:05.000000000")
	} else {
		return s, errors.New("unknown precision")
	}
	return s, nil
}

type DateMapper struct {
}

func (d DateMapper) MapSint(s scalar.Sint) (scalar.Sint, error) {
	v := s.Actual
	s.Sym = time.Unix(0, 0).AddDate(0, 0, int(v)).UTC().Format("2006-01-02")
	return s, nil
}

// Todo Decimal: https://github.com/linkedin/goavro/blob/master/logical_type.go
// Todo Duration
