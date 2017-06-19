package ignore

import (
	"bytes"
	"fmt"

	"github.com/gobwas/glob"
)

// Matcher defines an interface for filematchers
//
type Matcher interface {
	Match(string) bool
	True() bool
	MarshalText() ([]byte, error)
}

type MultiMatch struct {
	matchers []Matcher
}

func NewMultiMatch(matchers []Matcher) *MultiMatch {
	return &MultiMatch{matchers: matchers}
}

func (mm *MultiMatch) Match(arg string) bool {
	// Normal: OR
	// false, false -> false
	// false, true  -> true
	// true, false -> true
	// true, true -> true

	// Invert:
	// false, false -> false
	// false, true -> false
	// true, false -> true
	// true, true -> false
	use := false
	for _, m := range mm.matchers {
		if m.Match(arg) {
			use = m.True()
		}
	}
	return use

}

func (mm *MultiMatch) True() bool { return true }
func (mm *MultiMatch) MarshalText() ([]byte, error) {
	return []byte("multi"), nil
}

// GlobMatch handle glob matching
type GlobMatch struct {
	orig    string
	matcher glob.Glob
	normal  bool
}

func NewGlobMatch(arg []byte) (*GlobMatch, error) {
	truth := true
	if len(arg) > 0 && arg[0] == '!' {
		truth = false
		arg = arg[1:]
	}
	if bytes.IndexByte(arg, '/') == -1 {
		return NewBaseGlobMatch(string(arg), truth)
	}
	return NewPathGlobMatch(string(arg), truth)
}

// NewPathGlobMatch compiles a new matcher.
// Arg true should be set to false if the output is inverted.
func NewBaseGlobMatch(arg string, truth bool) (*GlobMatch, error) {
	g, err := glob.Compile(arg)
	if err != nil {
		return nil, err
	}
	return &GlobMatch{orig: arg, matcher: g, normal: truth}, nil
}

// NewPathGlobMatch compiles a new matcher.
// Arg true should be set to false if the output is inverted.
func NewPathGlobMatch(arg string, truth bool) (*GlobMatch, error) {
	// if starts with "/" then glob only applies to top level
	if len(arg) > 0 && arg[0] == '/' {
		arg = arg[1:]
	}

	// create path-aware glob
	g, err := glob.Compile(arg, '/')
	if err != nil {
		return nil, err
	}
	return &GlobMatch{orig: arg, matcher: g, normal: truth}, nil
}

// True returns true if this should be evaluated normally ("true is true")
//  and false if the result should be inverted ("false is true")
//
func (g *GlobMatch) True() bool { return g.normal }

// MarshalJSON is really a debug function
func (g *GlobMatch) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s: %v %s\"", "GlobMatch", g.normal, g.orig)), nil
}

// Match satisfies the Matcher interface
func (g *GlobMatch) Match(file string) bool {
	return g.matcher.Match(file)
}
