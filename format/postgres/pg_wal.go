package postgres

// TO DO
// not ready yet

//func init() {
//	interp.RegisterFormat(decode.Format{
//		Name:        format.PG_WAL,
//		Description: "PostgreSQL write-ahead log file",
//		DecodeFn:    decodePgwal,
//		DecodeInArg: format.PostgresIn{
//			Flavour: "default",
//		},
//	})
//}
//
//// https://pgpedia.info/x/XLOG_PAGE_MAGIC.html
//const (
//	XLOG_PAGE_MAGIC_15 = uint16(0xD10F)
//	XLOG_PAGE_MAGIC_14 = uint16(0xD10D)
//	XLOG_PAGE_MAGIC_13 = uint16(0xD106)
//	XLOG_PAGE_MAGIC_12 = uint16(0xD101)
//	XLOG_PAGE_MAGIC_11 = uint16(0xD098)
//	XLOG_PAGE_MAGIC_10 = uint16(0xD097)
//	XLOG_PAGE_MAGIC_96 = uint16(0xD093)
//)
//
//func decodePgwal(d *decode.D, in any) any {
//	d.Endian = decode.LittleEndian
//
//	flavour := in.(format.PostgresIn).Flavour
//	switch flavour {
//	//case PG_FLAVOUR_POSTGRES11:
//	//	return postgres11.DecodePgControl(d, in)
//	case PG_FLAVOUR_POSTGRES14, PG_FLAVOUR_POSTGRES:
//		return postgres14.DecodePgwal(d, in)
//	//case PG_FLAVOUR_PGPROEE14:
//	//	return pgproee14.DecodePgControl(d, in)
//	default:
//		break
//	}
//
//	return probePgwal(d, in)
//}
//
//func probePgwal(d *decode.D, in any) any {
//	// read version
//	xlp_magic := uint16(d.U16())
//
//	// restore position
//	d.SeekAbs(0)
//
//	switch xlp_magic {
//	case XLOG_PAGE_MAGIC_14:
//		return postgres14.DecodePgwal(d, in)
//	default:
//		d.Fatalf("unsupported xlp_magic = %X\n", xlp_magic)
//	}
//	return nil
//}
