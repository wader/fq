package colorjson

func newColor(c string) []byte {
	return []byte("\x1b[" + c + "m")
}

var (
	resetColor     = newColor("0")    // Reset
	nullColor      = newColor("90")   // Bright black
	falseColor     = newColor("33")   // Yellow
	trueColor      = newColor("33")   // Yellow
	numberColor    = newColor("36")   // Cyan
	stringColor    = newColor("32")   // Green
	objectKeyColor = newColor("34;1") // Bold Blue
	arrayColor     = []byte(nil)      // No color
	objectColor    = []byte(nil)      // No color
)
