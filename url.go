package misspell

import (
	"github.com/mvdan/xurls"
)

// StripURL attemps to replace URLs with blank spaces, e.g.
//  "xxx http://foo.com/ yyy -> "xxx          yyyy"
func StripURL(s string) string {
	return xurls.Strict.ReplaceAllStringFunc(s, replaceWithBlanks)
}
