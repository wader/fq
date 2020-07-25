package columnwriter_test

import (
	"fmt"
	"fq/internal/columnwriter"
	"os"
	"testing"
)

func TestColumnWriter(t *testing.T) {
	cw := columnwriter.New(os.Stdout, []int{3, 4})

	fmt.Fprintln(cw, "aaaaa")
	fmt.Fprintln(cw, "bbb")
	fmt.Fprint(cw, "cc")

	cw.Next()

	fmt.Fprintln(cw, "11111")
	fmt.Fprintln(cw, "22")
	fmt.Fprintln(cw, "33")

	cw.Next()

}
