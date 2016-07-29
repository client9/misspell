package misspell

import (
	"path/filepath"
	"strings"
)

// The number of possible binary formats is very large
// items that might be checked into a repo or be an
// artifact of a build.  Additions welcome.
//
// Golang's internal table is very small and can't be
// relied on.  Even then things like ".js" have a mime
// type of "application/javascipt" which isn't very helpful.
//
var binary = map[string]bool{
	".a":     true, // archive
	".bin":   true, // binary
	".bz2":   true, // compression
	".dll":   true, // shared library
	".exe":   true, // binary
	".gif":   true, // image
	".gz":    true, // compression
	".ico":   true, // image
	".jar":   true, // archive
	".jpeg":  true, // image
	".jpg":   true, // image
	".mp3":   true, // audio
	".mp4":   true, // video
	".mpeg":  true, // video
	".o":     true, // object file
	".pdf":   true, // pdf -- might be possible to use this later
	".png":   true, // image
	".so":    true, // shared library
	".tar":   true, // archive
	".tiff":  true, // image
	".woff":  true, // font
	".woff2": true, // font
	".xz":    true, // compression
	".z":     true, // compression
	".zip":   true, // archive
}

// IsBinaryFile returns true if the file is likely to be binary
//
// Better heuristics could be done here, in particular a binary
// file is unlikely to be UTF-8 encoded.  However this is cheap
// and will solve the immediate need of making sure common
// binary formats are not corrupted by mistake.
func IsBinaryFile(s string) bool {
	return binary[strings.ToLower(filepath.Ext(s))]
}
