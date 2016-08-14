package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

// >computable</a>
var reTerm = regexp.MustCompile(`>[^<]+</a>`)
var reAZ = regexp.MustCompile(`[a-z]+`)

// foldoc has some intentional misspellings
var knownBad = map[string]bool{
	"algorithim":  true,
	"protocal":    true,
	"interupt":    true,
	"syncronous":  true,
	"asyncronous": true,
	"terrabyte":   true,
	"hexidecimal": true,
}

func main() {
	// unique word list
	uniques := make(map[string]bool)

	resp, err := http.Get("http://foldoc.org/contents/all.html")
	if err != nil {
		log.Fatalf("Unable to get foldoc: %s", err)
	}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		// for each line
		line := scanner.Text()
		terms := reTerm.FindAllString(line, -1)
		// for each term
		for _, t := range terms {
			//log.Printf("TERM: %s", t)
			words := reAZ.FindAllString(strings.ToLower(t), -1)
			// for each word in term
			for _, w := range words {
				if len(w) > 4 && !knownBad[w] {
					uniques[w] = true
				}
			}
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("reading foldoc : %s", err)
	}
	resp.Body.Close()
	log.Printf("Got %d uniques", len(uniques))
	keys := make([]string, 0, len(uniques))
	for k := range uniques {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fo, err := os.Create("words-foldoc.txt")
	if err != nil {
		log.Fatalf("Unable to create output: %s", err)
	}
	for _, k := range keys {
		fo.Write([]byte(k + "\n"))
	}
	fo.Close()
	log.Printf("Wrote %s", "words-foldoc.txt")
}
