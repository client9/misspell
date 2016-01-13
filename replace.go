package misspell

import (
	"log"
	"strings"
	"text/scanner"
)

/**
 * Need to redo this so its more similar to how golang Flag works
 * there is a default global, but then you can make your own object if needed
 */

var replacer *strings.Replacer

func init() {
	replacer = strings.NewReplacer(dictWikipedia...)
	if replacer == nil {
		panic("unable to create strings.Replacer")
	}
}

// Diff is datastructure showing what changed in a single line
type Diff struct {
	Filename  string
	Line      int
	Column    int
	Original  string
	Corrected string
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func commonPrefixWordLength(a, b string) int {
	// re-order so len(a) <= len(b) always
	if len(a) > len(b) {
		b, a = a, b
	}
	lastWhite := 0
	for i := 0; i < len(a); i++ {
		if a[i] == ' ' || a[i] == '\t' {
			lastWhite = i
		}
		if a[i] != b[i] {
			if lastWhite == 0 {
				return 0
			}
			return min(lastWhite+1, len(a))
		}
	}
	return len(a)
}

// commonSuffixWordLength
func commonSuffixWordLength(a, b string) int {
	alen, blen := len(a), len(b)
	n := min(alen, blen)
	lastWhite := 0
	for i := 0; i < n; i++ {
		if a[alen-i-1] == ' ' || a[alen-i-1] == '\t' {
			lastWhite = i
		}
		if a[alen-i-1] != b[blen-i-1] {
			if lastWhite == 0 {
				return 0
			}
			return min(lastWhite+1, n)
		}
	}
	return n
}

// Return one words that are corrected in a single line
//
//
func corrected(instr, outstr string) (orig, corrected string, column int) {
	prefixLen := commonPrefixWordLength(instr, outstr)
	suffixLen := commonSuffixWordLength(instr, outstr)

	a := instr[prefixLen : len(instr)-suffixLen]
	b := outstr[prefixLen : len(outstr)-suffixLen]

	// Normal, we found the right snippet and it seems sane
	if len(a) < 30 && len(b) < 30 {
		return a, b, prefixLen
	}

	// some lines have no spaces and triggers a huge output
	// trim down
	var col int
	for i := 0; i < len(a); i++ {
		if i < len(b) && a[i] != b[i] {
			col = i - 10
			if col < 0 {
				col = 0
			}
			amax := i + 10
			if amax > len(a) {
				amax = len(a)
			}
			bmax := i + 10
			if bmax > len(b) {
				bmax = len(b)
			}
			return a[col:amax], b[col:bmax], col
		}
	}
	return a, b, col
}

// DiffLines produces a grep-like diff between two strings showing
// filename, linenum and change.  It is not meant to be a comprehensive diff.
func DiffLines(filename, input, output string) (string, []Diff) {
	var changes []Diff
	if output == input {
		return output, changes
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
		s1, s2, col := corrected(inlines[i], outlines[i])
		changes = append(changes, Diff{
			Filename:  filename,
			Line:      i + 1, // lines start at 1
			Column:    col,
			Original:  s1,
			Corrected: s2,
		})
	}
	return output, changes
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

// ReplaceGo is a specialized routine for correcting Golang source
// files.  Currently only checks comments, not identifiers for
// spelling.
//
// Other items:
//   - check strings, but need to ignore
//      * import "statements" blocks
//      * import ( "blocks" )
//   - skip first comment (line 0) if build comment
//
func ReplaceGo(input string, debug bool) string {
	var s scanner.Scanner
	s.Init(strings.NewReader(input))
	s.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanChars | scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanComments
	lastPos := 0
	output := ""
	for {

		switch s.Scan() {
		case scanner.Comment:
			origComment := s.TokenText()

			var newComment string
			if debug {
				newComment = ReplaceDebug(origComment)
			} else {
				newComment = Replace(origComment)
			}

			if origComment != newComment {
				// s.Pos().Offset is the end of the current token
				//  subtract len(origComment) to get the start of token
				offset := s.Pos().Offset
				output = output + input[lastPos:offset-len(origComment)] + newComment
				lastPos = offset
			}
		case scanner.EOF:
			// no changes, no copies
			if lastPos == 0 {
				return input
			}
			if lastPos >= len(input) {
				return output
			}

			return output + input[lastPos:]
		}
	}
}

func inArray(haystack []string, needle string) bool {
	for _, word := range haystack {
		if needle == word {
			return true
		}
	}
	return false
}

// Ignore removes a correction rule
//   WARNING: multiple calls to this will unset the previous calls.
//    thats not so good.
func Ignore(words []string) {
	newwords := make([]string, 0, len(dictWikipedia))
	for i := 0; i < len(dictWikipedia); i += 2 {
		if inArray(words, dictWikipedia[i]) {
			continue
		}
		newwords = append(newwords, dictWikipedia[i])
		newwords = append(newwords, dictWikipedia[i+1])
	}
	replacer = strings.NewReplacer(newwords...)
	if replacer == nil {
		panic("unable to create strings.Replacer")
	}
}
