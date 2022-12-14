package vpx

import "github.com/wader/fq/pkg/scalar"

var vpxLevelNames = scalar.UintMapSymStr{
	10: "level_1",
	11: "level_1.1",
	20: "level_2",
	21: "level_2.1",
	30: "level_3",
	31: "level_3.1",
	40: "level_4",
	41: "level_4.1",
	50: "level_5",
	51: "level_5.1",
	52: "level_5.2",
	60: "level_6",
	61: "level_6.1",
	62: "level_6.2",
}

var vpxChromeSubsamplingNames = scalar.UintMapSymStr{
	0: "4:2:0 vertical",
	1: "4:2:0 colocated with luma (0,0)",
	2: "4:2:2",
	3: "4:4:4",
}
