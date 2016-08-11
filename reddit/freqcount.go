package main

// this takes a list of gzipped Reddit comment-body files
// and returns a frequency count of words
// as a gzipped CSV file
import (
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
)

// freqCount is mapping of string->count
type freqCount map[string]int

// make a new counter with some minor preallocation
//  each month has about 2.2M uniques
func newFreqCount() freqCount {
	return make(freqCount, 3000000)
}

// globals used for worker queues and global counts
var wg sync.WaitGroup

// Reddit is struct used to unmarshal the reddit comment
type Reddit struct {
	Body string `json:"body"`
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
		word = strings.Trim(word, `'".,:;-()`)
		word = strings.ToLower(word)
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
func doit(prefix, filename string) {
	const maxLines = 1 << 31
	const suffix = "-body.json.gz"
	// generate outputfile name
	// blah/RC_2015-06-body.json.gz --> RC_2015-06-counts.csv.gz
	base := filepath.Base(filename)
	base = base[0:len(base)-len(suffix)] + "-counts.csv.gz"

	log.Printf("[%s] starting. Output file %s", prefix, base)
	counts := newFreqCount()
	fi, err := os.Open(filename)
	if err != nil {
		log.Fatalf("[%s]: %s", prefix, err)
	}

	file, err := gzip.NewReader(fi)
	if err != nil {
		log.Fatalf("[%s]: gzip error %s", prefix, err)
	}

	// no need to buffer this since raw network and bzip2 will
	// naturally buffer the input
	jsonin := json.NewDecoder(file)
	obj := Reddit{}
	lines := 0
	total := 0
	for jsonin.More() && lines < maxLines {
		lines++
		// decode an array value (Message)
		err := jsonin.Decode(&obj)
		if err != nil {
			log.Fatalf("[%s] unable to unmarshal object: %s", prefix, err)
		}

		// turn into words
		words := text2words(obj.Body)

		// take only ones that are reasonable
		// and add to counts
		for _, word := range words {
			if len(word) > 4 && len(word) < 21 {
				counts[word]++
				total++
			}
		}

	}
	fi.Close()

	fo, err := os.Create(base)
	if err != nil {
		log.Fatalf("OH NO, unable to write: %s", err)
	}
	fout := gzip.NewWriter(fo)

	keys := make([]string, 0, len(counts))
	for k, _ := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fout.Write([]byte(fmt.Sprintf("%s,%d\n", k, counts[k])))
	}
	fout.Close()
	fo.Close()
	log.Printf("[%s] DONE: wrote %s got %d unique words from %d", prefix, base, len(counts), total)
}

func worker(id int, jobs <-chan string) {
	for j := range jobs {
		doit(fmt.Sprintf("%d:%s", id, j), j)
	}
	wg.Done()
}

func main() {
	numWorkers := runtime.NumCPU()

	flag.Parse()
	args := flag.Args()

	jobs := make(chan string, len(args))

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs)
	}

	for _, arg := range args {
		log.Printf("[MASTER}: adding %s", arg)
		jobs <- arg
	}
	close(jobs)
	wg.Wait()
	log.Printf("[MASTER]: done")
}
