package script_test

import (
	"log"
	"regexp"
	"testing"

	"github.com/wader/fq/internal/deepequal"
	"github.com/wader/fq/internal/script"
)

func TestSectionParser(t *testing.T) {
	actualSections := script.SectionParser(
		regexp.MustCompile(`^(?:(a:)|(b:))$`),
		`
a:
c
c
b:
a:
c
a:
`[1:])

	expectedSections := []script.Section{
		{LineNr: 1, Name: "a:", Value: "c\nc\n"},
		{LineNr: 4, Name: "b:", Value: ""},
		{LineNr: 5, Name: "a:", Value: "c\n"},
		{LineNr: 7, Name: "a:", Value: ""},
	}

	deepequal.Error(t, "sections", expectedSections, actualSections)
}

func TestUnescape(t *testing.T) {

	s := script.Unescape(`asd\n\r\t \0b11110000 asd \0xffcb sdfd `)
	log.Printf("s: %v\n", []byte(s))

}
