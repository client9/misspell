package lib

import (
	"fmt"
	"io"
	"log"
	"strings"
)

var replacer *strings.Replacer

func init() {
	replacer = strings.NewReplacer(dictWikipedia...)
	if replacer == nil {
		panic("unable to create strings.Replacer")
	}
}

// Return one word that was corrected in a line
//
//  NOTE: there may be multiple words corrected in a single line but
//  this is not meant to be a complete diff
func corrected(instr, outstr string) (orig, corrected string) {
	inparts := strings.Fields(instr)
	outparts := strings.Fields(outstr)
	for i := 0; i < len(inparts); i++ {
		if i < len(outparts) && inparts[i] != outparts[i] {
			return inparts[i], outparts[i]
		}
	}
	return "", ""
}

// DiffLines produces a grep-like diff between two strings showing
// filename, linenum and change.  It is not meant to be a comprehensive diff.
func DiffLines(filename, input, output string, w io.Writer) int {
	if output == input {
		return 0
	}
	count := 0
	// line by line to make nice output
	outlines := strings.Split(output, "\n")
	inlines := strings.Split(input, "\n")
	for i := 0; i < len(inlines); i++ {
		if inlines[i] == outlines[i] {
			continue
		}
		count++
		s1, s2 := corrected(inlines[i], outlines[i])
		io.WriteString(w, fmt.Sprintf("%s:%d: corrected %q -> %q\n",
			filename, i, s1, s2))
	}
	return count
}

// ReplaceDebug logs exactly what was matched and replaced for using
// in debugging
func ReplaceDebug(input string) string {
	for i := 0; i < len(dictWikipedia); i += 2 {
		idx := strings.Index(input, dictWikipedia[i])
		if idx != -1 {
			left := idx - 10
			if left < 0 {
				left = 0
			}
			right := idx + len(dictWikipedia[i]) + 10
			if right > len(input) {
				right = len(input)
			}
			snippet := input[left:right]
			log.Printf("Found %q in %q  (%q)", dictWikipedia[i], snippet, dictWikipedia[i+1])
		}
	}
	return Replace(input)
}

// Replace takes input string and does spelling corrections on
// commonly misspelled words.
func Replace(input string) string {
	// ok doesn't do much
	return replacer.Replace(input)
}