// lazyre lazily compiles a *regexp.Regexp in concurrency safe way
// Use &lazyre.RE{S: `...`} or call New
package lazyre

import (
	"regexp"
	"sync"
)

type RE struct {
	S string

	m  sync.RWMutex
	re *regexp.Regexp
}

// New creates a new *lazyRE
func New(s string) *RE {
	return &RE{S: s}
}

// Must compiles regexp, returned *regexp.Regexp can be stored away and reused
func (lr *RE) Must() *regexp.Regexp {
	lr.m.Lock()
	defer lr.m.Unlock()
	if lr.re == nil {
		lr.re = regexp.MustCompile(lr.S)
	}
	return lr.re
}
