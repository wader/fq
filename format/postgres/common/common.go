package common

import (
	"fmt"
	"github.com/wader/fq/pkg/scalar"
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
