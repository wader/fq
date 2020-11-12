package rangegap_test

import (
	"fq/internal/rangegap"
	"log"
	"testing"
)

func Test(t *testing.T) {

	l := rangegap.Find(0, 10, [][2]int64{{1, 1}, {9, 10}})
	log.Printf("l: %#+v\n", l)

}
