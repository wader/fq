package ffv1

// ffmpeg -y -f lavfi -i testsrc -level 3 -t 10ms -c:v ffv1 ffv1.mkv

import (
	"log"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/ffv1/rangecoder"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

func init() {
	interp.RegisterFormat(
		format.FFV1_Frame,
		&decode.Format{
			Description: "FFV1 Video Coding Format frame",
			DecodeFn:    ffv1FrameDecode,
		})
}

type stateTransitions [256]byte

var defaultStateTransitionTable = stateTransitions{
	0, 0, 0, 0, 0, 0, 0, 0, 20, 21, 22, 23, 24, 25, 26, 27,
	28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 37, 38, 39, 40, 41, 42,
	43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 56, 57,
	58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73,
	74, 75, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88,
	89, 90, 91, 92, 93, 94, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103,
	104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 114, 115, 116, 117, 118,
	119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 133,
	134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149,
	150, 151, 152, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164,
	165, 166, 167, 168, 169, 170, 171, 171, 172, 173, 174, 175, 176, 177, 178, 179,
	180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 190, 191, 192, 194, 194,
	195, 196, 197, 198, 199, 200, 201, 202, 202, 204, 205, 206, 207, 208, 209, 209,
	210, 211, 212, 213, 215, 215, 216, 217, 218, 219, 220, 220, 222, 223, 224, 225,
	226, 227, 227, 229, 229, 230, 231, 232, 234, 234, 235, 236, 237, 238, 239, 240,
	241, 242, 243, 244, 245, 246, 247, 248, 248, 0, 0, 0, 0, 0, 0, 0,
}

// var alternativeStateTransitionTable = stateTransitions{
// 	0, 10, 10, 10, 10, 16, 16, 16, 28, 16, 16, 29, 42, 49, 20, 49,
// 	59, 25, 26, 26, 27, 31, 33, 33, 33, 34, 34, 37, 67, 38, 39, 39,
// 	40, 40, 41, 79, 43, 44, 45, 45, 48, 48, 64, 50, 51, 52, 88, 52,
// 	53, 74, 55, 57, 58, 58, 74, 60, 101, 61, 62, 84, 66, 66, 68, 69,
// 	87, 82, 71, 97, 73, 73, 82, 75, 111, 77, 94, 78, 87, 81, 83, 97,
// 	85, 83, 94, 86, 99, 89, 90, 99, 111, 92, 93, 134, 95, 98, 105, 98,
// 	105, 110, 102, 108, 102, 118, 103, 106, 106, 113, 109, 112, 114, 112, 116, 125,
// 	115, 116, 117, 117, 126, 119, 125, 121, 121, 123, 145, 124, 126, 131, 127, 129,
// 	165, 130, 132, 138, 133, 135, 145, 136, 137, 139, 146, 141, 143, 142, 144, 148,
// 	147, 155, 151, 149, 151, 150, 152, 157, 153, 154, 156, 168, 158, 162, 161, 160,
// 	172, 163, 169, 164, 166, 184, 167, 170, 177, 174, 171, 173, 182, 176, 180, 178,
// 	175, 189, 179, 181, 186, 183, 192, 185, 200, 187, 191, 188, 190, 197, 193, 196,
// 	197, 194, 195, 196, 198, 202, 199, 201, 210, 203, 207, 204, 205, 206, 208, 214,
// 	209, 211, 221, 212, 213, 215, 224, 216, 217, 218, 219, 220, 222, 228, 223, 225,
// 	226, 224, 227, 229, 240, 230, 231, 232, 233, 234, 235, 236, 238, 239, 237, 242,
// 	241, 243, 242, 244, 245, 246, 247, 248, 249, 250, 251, 252, 252, 253, 254, 255,
// }

type rangeCoderCtx struct {
	state     byte
	_range    int
	end       bool
	low       int
	zeroState stateTransitions
	oneState  stateTransitions
	// TODO: bitio.Reader somehow?
	readBits func(bits int) (int, error)
	bitsLeft func() (int64, error)
}

func newRangeCoderCtx(s stateTransitions, readBits func(bits int) (int, error), bitsLeft func() (int64, error)) (*rangeCoderCtx, error) {
	r := &rangeCoderCtx{
		readBits: readBits,
		bitsLeft: bitsLeft,

		state:  0,
		_range: 0xff00,
		end:    false,
	}

	r.oneState = s
	r.zeroState[0] = 0
	for i := range r.zeroState {
		r.zeroState[i] = -r.oneState[256-i]
	}

	low, err := readBits(16)
	if err != nil {
		return nil, err
	}
	r.low = low
	if low >= r._range {
		low = r._range
		r.end = true
	}

	return r, nil
}

func (r *rangeCoderCtx) reFill() error {
	if r._range >= 256 {
		return nil
	}
	r._range *= 256
	r.low *= 256
	if !r.end {
		n, err := r.readBits(8)
		if err != nil {
			return err
		}
		r.low += n
		if n, err := r.bitsLeft(); err != nil {
			return err
		} else if n == 0 {
			r.end = true
		}
	}
	return nil
}

func (r *rangeCoderCtx) getRac() (int, error) {
	rangeOff := (r._range * int(r.state)) / 256
	r._range -= rangeOff

	if r.low < r._range {
		// zero_state_i = 256 - one_state_(256-i)
		//
		//   Figure 23: Description of the coding of the state transition
		// 		   table for a "get_rac" readout value of 0.
		r.state = r.zeroState[r.state]
		if err := r.reFill(); err != nil {
			return 0, err
		}
		return 0, nil
	}

	// 	one_state_i =
	// 	default_state_transition_i + state_transition_delta_i
	//
	//   Figure 22: Description of the coding of the state transition
	// 		   table for a "get_rac" readout value of 1.
	r.low -= r._range
	r.state = r.oneState[r.state]
	r._range = rangeOff
	if err := r.reFill(); err != nil {
		return 0, err
	}
	return 1, nil
}

func (r *rangeCoderCtx) getSymbolU() (uint64, error) {
	n, err := r.getRac()
	if err != nil {
		return 0, err
	}
	if n == 1 {
		return 0, nil
	}

	// e := 0
	// for { //1..10
	// 	if _, err := r.getRac( /*1 + mathex.Min(e, 9)*/ ); err != nil {
	// 		return 0, err
	// 	}
	// 	e++
	// }

	// a := 1
	// for i := e - 1; i >= 0; i-- {
	// 	a = a*2 + r.getRac(22+mathex.Min(i, 9)) // 22..31
	// }

	// if !is_signed {
	// 	return a
	// }

	// if get_rac(c, state+11+min(e, 10)) { //11..21
	// 	return -a
	// } else {
	// 	return a
	// }

	return 1, nil
}

func fieldUR(d *decode.D, rc *rangecoder.Coder, state []uint8, name string) {
	d.TryFieldUintFn(name, func(d *decode.D) (uint64, error) {
		u, err := rc.UR(state)
		return uint64(u), err
	})
}

func ffv1FrameDecode(d *decode.D) any {
	const contextSize = 32

	state := make([]uint8, contextSize)
	for i := 0; i < contextSize; i++ {
		state[i] = 128
	}

	low := uint16(d.FieldU16("low"))

	rc, err := rangecoder.NewCoder(low, func() (byte, error) { n, err := d.TryU8(); log.Printf("n: %#+v\n", n); return byte(n), err }, int(d.BitsLeft())/8)
	if err != nil {
		d.IOPanic(err, "rangecoder.NewCoder")
	}
	fieldUR(d, rc, state, "version")
	fieldUR(d, rc, state, "micro_version")
	fieldUR(d, rc, state, "coder_type")
	// d.Field
	// fieldUR(d, rc, state, "micro_version")

	return nil
}
