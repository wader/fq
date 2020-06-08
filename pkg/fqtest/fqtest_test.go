package fqtest_test

import (
	"fq/internal/deepequal"
	"fq/pkg/fqtest"
	"log"
	"regexp"
	"testing"
)

func TestSectionParser(t *testing.T) {
	actualSections := fqtest.SectionParser(
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

	expectedSections := []fqtest.Section{
		{LineNr: 1, Name: "a:", Value: "c\nc\n"},
		{LineNr: 4, Name: "b:", Value: ""},
		{LineNr: 5, Name: "a:", Value: "c\n"},
		{LineNr: 7, Name: "a:", Value: ""},
	}

	deepequal.Error(t, "sections", expectedSections, actualSections)
}

func TestUnescape(t *testing.T) {

	s := fqtest.Unescape(`asd \b123213 asd \xffcb sdfd `)
	log.Printf("s: %#+v\n", s)

}
