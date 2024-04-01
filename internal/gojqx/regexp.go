package gojqx

import (
	"fmt"
	"regexp"
	"strings"
)

// from gojq, see https://github.com/itchyny/gojq/blob/main/LICENSE
func CompileRegexp(re, allowedFlags, flags string) (*regexp.Regexp, error) {
	if strings.IndexFunc(flags, func(r rune) bool {
		return !strings.ContainsAny(string([]rune{r}), allowedFlags)
	}) >= 0 {
		return nil, fmt.Errorf("unsupported regular expression flag: %q", flags)
	}
	re = strings.ReplaceAll(re, "(?<", "(?P<")
	if strings.ContainsRune(flags, 'i') {
		re = "(?i)" + re
	}
	if strings.ContainsRune(flags, 'm') {
		re = "(?s)" + re
	}
	r, err := regexp.Compile(re)
	if err != nil {
		return nil, fmt.Errorf("invalid regular expression %q: %w", re, err)
	}
	return r, nil
}
