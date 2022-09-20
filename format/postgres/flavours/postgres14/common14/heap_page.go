package common14

// HeapPage used in tables, indexes...

// type = struct ItemIdData {
/*    0: 0   |     4 */ // unsigned int lp_off: 15
/*    1: 7   |     4 */ // unsigned int lp_flags: 2
/*    2: 1   |     4 */ // unsigned int lp_len: 15
//
/* total size (bytes):    4 */

type ItemIdData struct {
	Off   uint32 // unsigned int lp_off: 15
	Flags uint32 // unsigned int lp_flags: 2
	Len   uint32 // unsigned int lp_len: 15
}
