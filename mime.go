package misspell

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
	".class": true, // Java class file
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
	".pyc":   true, // Python bytecode
	".pyo":   true, // Python bytecode
	".so":    true, // shared library
	".swp":   true, // vim swap file
	".tar":   true, // archive
	".tiff":  true, // image
	".woff":  true, // font
	".woff2": true, // font
	".xz":    true, // compression
	".z":     true, // compression
	".zip":   true, // archive
}

// isBinaryFilename returns true if the file is likely to be binary
//
// Better heuristics could be done here, in particular a binary
// file is unlikely to be UTF-8 encoded.  However this is cheap
// and will solve the immediate need of making sure common
// binary formats are not corrupted by mistake.
func isBinaryFilename(s string) bool {
	return binary[strings.ToLower(filepath.Ext(s))]
}

var scm = map[string]bool{
	".bzr": true,
	".git": true,
	".hg":  true,
	".svn": true,
	"CVS":  true,
}

// isSCMPath returns true if the path is likely part of a (private) SCM
//  directory.  E.g.  ./git/something  = true
func isSCMPath(s string) bool {
	parts := strings.Split(s, string(filepath.Separator))
	for _, dir := range parts {
		if scm[dir] {
			return true
		}
	}
	return false
}

func isTextFile(raw []byte) bool {
	// allow any text/ type with utf-8 encoding
	// DetectContentType sometimes returns charset=utf-16 for XML stuff
	//  in which case ignore.
	mime := http.DetectContentType(raw)
	return strings.HasPrefix(mime, "text/") && strings.HasSuffix(mime, "charset=utf-8")
}

// ReadTextFile returns the contents of a file, first testing if it is a text file
//  returns ("", nil) if not a text file
//  returns ("", error) if error
//  returns (string, nil) if text
//
// unfortunately, in worse case, this does
//   1 stat
//   1 open,read,close of 512 bytes
//   1 more stat,open, read everything, close (via ioutil.ReadAll)
//  This could be kinder to the filesystem.
//
// This uses some heuristics of the file's extenion (e.g. .zip, .txt) and
// uses a sniffer to determine if the file is text or not.
// Using file extensions isn't great, but probably
// good enough for real-world use.
// Golang's built in sniffer is problematic for differnet reasons.  It's
// optimized for HTML, and is very limited in detection.  It would be good
// to explicitly add some tests for ELF/DWARF formats to make sure we never
// corrupt binary files.
func ReadTextFile(filename string) (string, error) {
	if isBinaryFilename(filename) {
		return "", nil
	}

	if isSCMPath(filename) {
		return "", nil
	}

	fstat, err := os.Stat(filename)

	if err != nil {
		return "", fmt.Errorf("Unable to stat %q: %s", filename, err)
	}

	// directory: nothing to do.
	if fstat.IsDir() {
		return "", nil
	}

	// avoid reading in multi-gig files
	// if input is large, read the first 512 bytes to sniff type
	// if not-text, then exit
	isText := false
	if fstat.Size() > 50000 {
		fin, err := os.Open(filename)
		if err != nil {
			return "", fmt.Errorf("Unable to open large file %q: %s", filename, err)
		}
		defer fin.Close()
		buf := make([]byte, 512)
		_, err = io.ReadFull(fin, buf)
		if err != nil {
			return "", fmt.Errorf("Unable to read 512 bytes from %q: %s", filename, err)
		}
		if !isTextFile(buf) {
			return "", nil
		}

		// set so we don't double check this file
		isText = true
	}

	// read in whole file
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("Unable to read all %q: %s", filename, err)
	}

	if !isText && !isTextFile(raw) {
		return "", nil
	}
	return string(raw), nil
}
