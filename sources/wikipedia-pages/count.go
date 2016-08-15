package main

// this takes a list of gzipped Reddit comment-body files
// and returns a frequency count of words
// as a gzipped CSV file
import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

// freqCount is mapping of string->count
type freqCount map[string]int

// make a new counter with some minor preallocation
//  each month has about 2.2M uniques
func newFreqCount() freqCount {
	return make(freqCount, 3000000)
}

// isAZ returns true if input is only consists of letters a-z (case SENSITIVE)
func isAZ(s string) bool {
	// NO: for _, ch := range s {
	//   This iterates over UTF-8 characters
	slen := len(s)
	for i := 0; i < slen; i++ {
		ch := s[i]
		if ch < 'a' || ch > 'z' {
			return false
		}
	}
	return true
}

// text2words takes raw text and extracts English-list words
//   split on whitespace
//   remove starting or ending  `'".,:;-()`
//   lower case
//   check if all a-z
//
// this doesn't catch everything, but that's ok for this purpose.
//
func text2words(raw string) []string {
	parts := strings.Fields(raw)
	out := make([]string, 0, len(parts))
	for _, word := range parts {
		word = strings.Trim(word, `[]'".,:;-()`)
		//	word = strings.ToLower(word)
		if isAZ(word) {
			out = append(out, word)
		}
	}
	return out
}

// doit does the following
//   opens a gizpped up reddit file
//   extracts words
//   updates counts
//
// Every 1M lines, it updates the global freq count
//
func main() {
	const maxLines = 1 << 32
	counts := newFreqCount()
	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 0)
	lines := 0
	total := 0
	for scanner.Scan() {
		lines++
		if lines > maxLines {
			break
		}
		if lines&0x000FFFF == 0 {
			log.Printf("lines=%d uniques=%d words=%d", lines, len(counts), total)
		}
		// turn into words
		line := scanner.Text()
		words := text2words(line)

		// take only ones that are reasonable
		// and add to counts
		for _, word := range words {
			if len(word) > 4 && len(word) < 21 {
				counts[word]++
				total++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Invalid input: %s", err)
	}

	log.Printf("Done lines=%d uniques=%d words=%d", lines, len(counts), total)
	fo, err := os.Create("wikipedia-counts.csv.gz")
	if err != nil {
		log.Fatalf("OH NO, unable to write: %s", err)
	}
	fout := gzip.NewWriter(fo)

	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fout.Write([]byte(fmt.Sprintf("%s,%d\n", k, counts[k])))
	}
	fout.Close()
	fo.Close()
}
