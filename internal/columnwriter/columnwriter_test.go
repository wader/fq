package columnwriter_test

import (
	"fmt"
	"fq/internal/columnwriter"
	"os"
	"testing"
)

func TestColumnWriter(t *testing.T) {
	cw := columnwriter.New(os.Stdout, []int{3, 4})

	fmt.Fprintln(cw.Columns[0], "aaaaa")
	fmt.Fprintln(cw.Columns[0], "bb")
	fmt.Fprint(cw.Columns[0], "cc")

	fmt.Fprintln(cw.Columns[1], "11111")
	fmt.Fprintln(cw.Columns[1], "22")
	fmt.Fprintln(cw.Columns[1], "33")

	cw.Flush()

	fmt.Fprintln(cw.Columns[1], "aaaaa")
	fmt.Fprintln(cw.Columns[1], "bb")
	fmt.Fprint(cw.Columns[1], "cc")

	fmt.Fprintln(cw.Columns[0], "11111")
	fmt.Fprintln(cw.Columns[0], "22")
	fmt.Fprintln(cw.Columns[0], "33")

	cw.Flush()
}
