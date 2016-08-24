package misspell

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"
)

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func inArray(haystack []string, needle string) bool {
	for _, word := range haystack {
		if needle == word {
			return true
		}
	}
	return false
}

var wordRegexp = regexp.MustCompile(`[a-zA-Z0-9']+`)

/*
line1 and line2 are different
extract words from each line1

replace word -> newword
if word == new-word
  continue
if new-word in list of replacements
  continue
new word not original, and not in list of replacements
  some substring got mixed up.  UNdo
*/
func recheckLine(s string, lineNum int, buf *bytes.Buffer, rep *strings.Replacer, corrected map[string]string) []Diff {
	// pre-allocate up to 4 corrections per line
	diffs := make([]Diff, 0, 4)

	first := 0
	redacted := RemoveNotWords(s)

	idx := wordRegexp.FindAllStringIndex(redacted, -1)
	for _, ab := range idx {
		word := s[ab[0]:ab[1]]
		newword := rep.Replace(word)
		if newword == word {
			// no replacement done
			continue
		}
		if corrected[strings.ToLower(word)] == strings.ToLower(newword) {
			// word got corrected into something we know
			buf.WriteString(s[first:ab[0]])
			buf.WriteString(newword)
			first = ab[1]
			diffs = append(diffs, Diff{
				FullLine:  s,
				Line:      lineNum,
				Original:  word,
				Corrected: newword,
				Column:    ab[0],
			})
			continue
		}
		// Word got corrected into something unknown. Ignore it
	}
	buf.WriteString(s[first:])
	return diffs
}

// Diff is datastructure showing what changed in a single line
type Diff struct {
	Filename  string
	FullLine  string
	Line      int
	Column    int
	Original  string
	Corrected string
}

// diffLines produces a grep-like diff between two strings showing
// filename, linenum and change.  It is not meant to be a comprehensive diff.
func diffLines(input, output string, r *strings.Replacer, c map[string]string) (string, []Diff) {
	changes := make([]Diff, 0, 16)
	buf := bytes.NewBuffer(make([]byte, 0, max(len(input), len(output))+100))

	// line by line to make nice output
	// This is horribly slow.
	outlines := strings.SplitAfter(output, "\n")
	inlines := strings.SplitAfter(input, "\n")
	for i := 0; i < len(inlines); i++ {
		if inlines[i] == outlines[i] {
			buf.WriteString(outlines[i])
			continue
		}
		linediffs := recheckLine(inlines[i], i+1, buf, r, c)
		changes = append(changes, linediffs...)
	}

	return buf.String(), changes
}

// Replacer is the main struct for spelling correction
type Replacer struct {
	Replacements []string
	Debug        bool
	engine       *strings.Replacer
	corrected    map[string]string
}

// New creates a new default Replacer using the main rule list
func New() *Replacer {
	r := Replacer{
		Replacements: DictMain,
	}
	r.Compile()
	return &r
}

// RemoveRule deletes existings rules.
// TODO: make inplace to save memory
func (r *Replacer) RemoveRule(ignore []string) {
	newwords := make([]string, 0, len(r.Replacements))
	for i := 0; i < len(r.Replacements); i += 2 {
		if inArray(ignore, r.Replacements[i]) {
			continue
		}
		newwords = append(newwords, r.Replacements[i:i+2]...)
	}
	r.engine = nil
	r.Replacements = newwords
}

// AddRuleList appends new rules.
// Input is in the same form as Strings.Replacer: [ old1, new1, old2, new2, ....]
// Note: does not check for duplictes
func (r *Replacer) AddRuleList(additions []string) {
	r.engine = nil
	r.Replacements = append(r.Replacements, additions...)
}

// Compile compiles the rules.  Required before using the Replace functions
func (r *Replacer) Compile() {

	r.corrected = make(map[string]string)
	for i := 0; i < len(r.Replacements); i += 2 {
		r.corrected[strings.ToLower(r.Replacements[i])] = strings.ToLower(r.Replacements[i+1])
	}
	r.engine = strings.NewReplacer(r.Replacements...)
}

// Replace makes spelling corrections to the input string
func (r *Replacer) Replace(input string) (string, []Diff) {
	news := r.engine.Replace(input)
	if input == news {
		return input, nil
	}

	// changes were made, diffLines rechecks and undoes bad corrections
	return diffLines(input, news, r.engine, r.corrected)
}

// ReplaceReader applies spelling corrections to a reader stream
func (r *Replacer) ReplaceReader(raw io.Reader, w io.Writer) []Diff {
	var (
		err     error
		orig    string
		changes = make([]Diff, 0, 16)
		lineNum int
	)
	buf := bytes.Buffer{}
	reader := bufio.NewReader(raw)
	writer := bufio.NewWriter(w)
	for {
		lineNum++
		orig, err = reader.ReadString('\n')
		if err != nil {
			break
		}
		// easily 5x faster than regexp+map
		if orig == r.engine.Replace(orig) {
			//w.Write([]byte(orig))
			writer.WriteString(orig)
			continue
		}

		// but it can be inaccurate, so we need to double check
		buf.Reset()
		linediffs := recheckLine(orig, lineNum, &buf, r.engine, r.corrected)
		changes = append(changes, linediffs...)
		//w.Write(buf.Bytes())
		writer.Write(buf.Bytes())
	}
	if err == io.EOF {
		return changes
	}
	return nil
}
