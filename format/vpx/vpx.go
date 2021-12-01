package vpx

import "github.com/wader/fq/pkg/scalar"

var vpxLevelNames = scalar.UToSymStr{
	10: "Level 1",
	11: "Level 1.1",
	20: "Level 2",
	21: "Level 2.1",
	30: "Level 3",
	31: "Level 3.1",
	40: "Level 4",
	41: "Level 4.1",
	50: "Level 5",
	51: "Level 5.1",
	52: "Level 5.2",
	60: "Level 6",
	61: "Level 6.1",
	62: "Level 6.2",
}

var vpxChromeSubsamplingNames = scalar.UToSymStr{
	0: "4:2:0 vertical",
	1: "4:2:0 colocated with luma (0,0)",
	2: "4:2:2",
	3: "4:4:4",
}
