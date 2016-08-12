package main

// reads in a reddit comment archive URL
//  and just extracts the body field
import (
	"compress/bzip2"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
)

// globals used for worker queues and global counts
var wg sync.WaitGroup

// Reddit is struct used to unmarshal the reddit comment
type Reddit struct {
	Body string `json:"body"`
}

// doit does the following
//   reads a URL
//   uncompresses it (bzip2)
//   json decodes
//   extracts comment body
//   writes to output file as mini-json
//
//
func doit(prefix, url string) {
	const maxLines = 1 << 31

	// generate outputfile name
	// blah/RC_2015-06.gz --> RC_2015-06-counts.csv
	base := path.Base(url)
	ext := filepath.Ext(base)
	base = base[0:len(base)-len(ext)] + "-body.json.gz"
	log.Printf("[%s] %s -> %s", prefix, url, base)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("[%s] url error: %s", prefix, err)
	}
	defer resp.Body.Close()
	file := bzip2.NewReader(resp.Body)

	// no need to buffer this since raw network and bzip2 will
	// naturally buffer the input
	jsonin := json.NewDecoder(file)

	// set up output file
	fo, err := os.Create(base)
	if err != nil {
		log.Fatalf("[%s] unable to write: %s", prefix, err)
	}
	// gzip output
	bufout := gzip.NewWriter(fo)
	// steam out json
	jsonout := json.NewEncoder(bufout)

	obj := Reddit{}
	lines := 0
	for jsonin.More() && lines < maxLines {
		lines++
		// decode an array value (Message)
		obj.Body = ""
		err := jsonin.Decode(&obj)
		if err != nil {
			log.Fatalf("[%s] unable to unmarshal object: %s", prefix, err)
		}
		if obj.Body == "[deleted]" || obj.Body == "" {
			continue
		}
		err = jsonout.Encode(&obj)
		if err != nil {
			log.Fatalf("[%s] unable to marshal object: %s", prefix, err)
		}
	}

	bufout.Close()
	fo.Close()
	log.Printf("[%s] done %d lines", prefix, lines)
}

func worker(id int, jobs <-chan string) {
	for j := range jobs {
		doit(fmt.Sprintf("%d:%s", id, j), j)
	}
	wg.Done()
}

type yearmonth struct {
	year  int
	month int
}

func yearmonthRange(start, end yearmonth) []yearmonth {
	out := []yearmonth{}
	m := start.month
	y := start.year
	for {
		out = append(out, yearmonth{y, m})
		m++
		if m == 13 {
			m = 1
			y++
		}
		if y > end.year {
			break
		}
		if y == end.year && m > end.month {
			break
		}
	}
	return out
}

func main() {
	dates := yearmonthRange(yearmonth{2010, 1}, yearmonth{2011, 12})
	for _, ym := range dates {
		arg := fmt.Sprintf("http://files.pushshift.io/reddit/comments/RC_%d-%02d.bz2", ym.year, ym.month)
		args = append(args, arg)
	}

	jobs := make(chan string, len(args))

	numCPU := runtime.NumCPU()
	for w := 1; w <= numCPU; w++ {
		wg.Add(1)
		go worker(w, jobs)
	}

	for _, arg := range args {
		log.Printf("[MASTER]: adding %s", arg)
		jobs <- arg
	}
	close(jobs)
	wg.Wait()
	log.Printf("[MASTER]: done")
}
