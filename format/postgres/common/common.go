package common

import "github.com/wader/fq/pkg/scalar"

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
