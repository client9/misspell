package main

// finds types.. O(n^2)
// looks for 0.01% (1 in 10,000) frequency or more
import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/xrash/smetrics"
)

type pair struct {
	word  string
	count int
}

type counts []pair

func (s counts) Len() int      { return len(s) }
func (s counts) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s counts) Less(i, j int) bool {
	if s[i].count < s[j].count {
		return true
	}
	if s[i].count > s[j].count {
		return false
	}
	return s[i].word < s[j].word
}

// LoadCSV loads in a CSV in form of WORD,COUNT
func loadCSV(fname string) (counts, error) {
	fi, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	fizip, err := gzip.NewReader(fi)
	if err != nil {
		return nil, err
	}
	defer fizip.Close()
	words := make(counts, 0, 700000)
	scanner := bufio.NewScanner(fizip)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ",", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("Got extra junk in line: %q", line)
		}
		c, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("Number conversion failed: %q", line)
		}
		words = append(words, pair{parts[0], c})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	sort.Sort(sort.Reverse(words))
	return words, nil
}

// LoadWordList loads in a list of known-good words
func LoadWordList(fname string) (map[string]bool, error) {
	fi, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	out := make(map[string]bool)
	intro := true
	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		line := scanner.Text()
		if intro {
			if line == "---" {
				intro = false
			}
			continue
		}
		out[strings.ToLower(line)] = true
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func main() {

	// load known-good words`
	truewords, err := LoadWordList("dict.txt")
	if err != nil {
		log.Fatalf("Unable to read wordlist")
	}
	log.Printf("Got %d real words from dictionary", len(truewords))

	// load frequency counts
	words, err := loadCSV("RC_2015-total.csv.gz")
	if err != nil {
		log.Fatalf("Unable to freq counts: %s", err)
	}

	// make total count
	total := 0
	for _, kv := range words {
		total += kv.count
	}
	log.Printf("Got %d uniques, %d total", len(words), total)

	for top := 0; top < len(words); top++ {
		a := words[top]

		// must have at least this many occurances to
		// have an entry

		// TODO: what is the break point?  A CDF/PDF
		// would be better instead of absolute
		// TODO: move comparison out of loop
		if a.count < 20000 {
			break
		}
		aword := a.word
		for bottom := top + 1; bottom < len(words); bottom++ {
			b := words[bottom]

			// misspelling must occur twice
			// TODO: again this is a fixed point in the list
			if b.count < 2 {
				break
			}
			ratio := 100.0 * float64(b.count) / float64(a.count)

			// great than one percent probably means
			// bword is just a different work that is similar
			if ratio > 1.0 {
				continue
			}
			// if less than 0.01% ignore, too rare
			if ratio < 0.01 {
				break
			}

			bword := b.word

			// handle "foobar"/"foobars"
			if strings.HasPrefix(aword, bword) || strings.HasPrefix(bword, aword) {
				continue
			}

			// handle "foobars", "foobarr"
			if aword[:len(aword)-1] == bword[:len(bword)-1] {
				continue
			}
			val := smetrics.JaroWinkler(aword, bword, 0.7, 4)
			if val > 0.96 {
				fmt.Printf("%s,%s,%d,%d,%d,%d,%f,%f\n",
					aword, bword,
					top, bottom,
					a.count, b.count, ratio,
					val)
			}
		}
	}
}
