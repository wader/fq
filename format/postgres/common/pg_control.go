package common

import (
	"fmt"
	"github.com/wader/fq/pkg/scalar"
	"time"
)

//typedef enum DBState
//{
//	DB_STARTUP = 0,
//	DB_SHUTDOWNED,
//	DB_SHUTDOWNED_IN_RECOVERY,
//	DB_SHUTDOWNING,
//	DB_IN_CRASH_RECOVERY,
//	DB_IN_ARCHIVE_RECOVERY,
//	DB_IN_PRODUCTION
//} DBState;
var DBState = scalar.UToScalar{
	0: {Sym: "DB_STARTUP"},
	1: {Sym: "DB_SHUTDOWNED"},
	2: {Sym: "DB_SHUTDOWNED_IN_RECOVERY"},
	3: {Sym: "DB_SHUTDOWNING"},
	4: {Sym: "DB_IN_CRASH_RECOVERY"},
	5: {Sym: "DB_IN_ARCHIVE_RECOVERY"},
	6: {Sym: "DB_IN_PRODUCTION"},
}

//typedef enum WalLevel
//{
//	WAL_LEVEL_MINIMAL = 0,
//	WAL_LEVEL_REPLICA,
//	WAL_LEVEL_LOGICAL
//} WalLevel;
var WalLevel = scalar.SToScalar{
	0: {Sym: "WAL_LEVEL_MINIMAL"},
	1: {Sym: "WAL_LEVEL_REPLICA"},
	2: {Sym: "WAL_LEVEL_LOGICAL"},
}

type icuVersionMapper struct{}

func (m icuVersionMapper) MapScalar(s scalar.S) (scalar.S, error) {
	a := s.ActualU()
	major := a & 0xff
	minor := (a >> 8) & 0xff
	v1 := (a >> 16) & 0xff
	v2 := (a >> 24) & 0xff
	s.Sym = fmt.Sprintf("%d.%d.%d.%d", major, minor, v1, v2)
	return s, nil
}

var IcuVersionMapper = icuVersionMapper{}

type xLogRecPtrMapper struct{}

func (m xLogRecPtrMapper) MapScalar(s scalar.S) (scalar.S, error) {
	v := s.ActualU()
	s.Sym = fmt.Sprintf("%X/%X", v>>32, uint32(v))
	return s, nil
}

var XLogRecPtrMapper = xLogRecPtrMapper{}
var LocPtrMapper = xLogRecPtrMapper{}

type timeMapper struct{}

func (m timeMapper) MapScalar(s scalar.S) (scalar.S, error) {
	ut := s.ActualS()
	t := time.Unix(ut, 0)
	s.Sym = t.UTC().Format(time.RFC1123)
	return s, nil
}

var TimeMapper = timeMapper{}

// typedef enum
//{
//	PG_UNKNOWN					= 0xFFFF,
//	PG_ORIGINAL					= 0,
//	PGPRO_STANDARD				= ('P'<<8|'P'),
//	PGPRO_ENTERPRISE			= ('P'<<8|'E'),
//} PgEdition;
const (
	PG_UNKNOWN       = 0xFFFF
	PG_ORIGINAL      = 0
	PGPRO_STANDARD   = (uint32('P') << 8) | uint32('P')
	PGPRO_ENTERPRISE = (uint32('P') << 8) | uint32('E')

	PG_UNKNOWN_STR       = "(unknown edition)"
	PG_ORIGINAL_STR      = "PostgreSQL"
	PGPRO_STANDARD_STR   = "Postgres Pro Standard"
	PGPRO_ENTERPRISE_STR = "Postgres Pro Enterprise"
)

type versionMapper struct{}

func (m versionMapper) MapScalar(s scalar.S) (scalar.S, error) {
	v := s.ActualU()
	v1 := uint32(v >> 16)
	v2 := uint32(v & 0xffff)
	switch v1 {
	case PG_UNKNOWN:
		s.Sym = fmt.Sprintf("%s %d", PG_UNKNOWN_STR, v2)
	case PG_ORIGINAL:
		s.Sym = fmt.Sprintf("%s %d", PG_ORIGINAL_STR, v2)
	case PGPRO_STANDARD:
		s.Sym = fmt.Sprintf("%s %d", PGPRO_STANDARD_STR, v2)
	case PGPRO_ENTERPRISE:
		s.Sym = fmt.Sprintf("%s %d", PGPRO_ENTERPRISE_STR, v2)
	}
	return s, nil
}

var VersionMapper = versionMapper{}

type hexMapper struct{}

func (m hexMapper) MapScalar(s scalar.S) (scalar.S, error) {
	v := s.ActualU()
	s.Sym = fmt.Sprintf("%X", v)
	return s, nil
}

var HexMapper = hexMapper{}
