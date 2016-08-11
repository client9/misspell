package main

// merges multiple frequency count files into one
// this is could use a merge sort and be smart

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// freqCount is mapping of string->count
type freqCount map[string]int

// make a new counter with some minor preallocation
//  each month has about 2.2M uniques
func newFreqCount() freqCount {
	return make(freqCount, 3000000)
}

// LoadCSV loads in a CSV in form of WORD,COUNT
func loadCSV(counts freqCount, fname string) error {
	fi, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fi.Close()

	fizip, err := gzip.NewReader(fi)
	if err != nil {
		return err
	}
	defer fizip.Close()

	scanner := bufio.NewScanner(fizip)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ",", 2)
		if len(parts) != 2 {
			return fmt.Errorf("Got extra junk in line: %q", line)
		}
		c, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("Number conversion failed: %q", line)
		}

		counts[parts[0]] += c
	}
	return scanner.Err()
}

func main() {
	outfile := flag.String("o", "", "output file name")
	mincount := flag.Int("mincount", 0, "only output if freqcount >=, 0 = all")
	minlen := flag.Int("minlen", 0, "only output if word is >=, 0 = all")
	flag.Parse()
	if *outfile == "" {
		log.Fatalf("Must specificy outfile")
	}
	args := flag.Args()
	counts := newFreqCount()
	for _, arg := range args {
		log.Printf("Loading %s", arg)
		err := loadCSV(counts, arg)
		if err != nil {
			log.Fatalf("%s: %s", arg, err)
		}
	}
	fo, err := os.Create(*outfile)
	if err != nil {
		log.Fatalf("OH NO, unable to write: %s", err)
	}
	fout := gzip.NewWriter(fo)

	keys := make([]string, 0, len(counts))
	total := 0
	for k, v := range counts {
		total += v
		if v >= *mincount && len(k) >= *minlen {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		fout.Write([]byte(fmt.Sprintf("%s,%d\n", k, counts[k])))
	}

	fout.Close()
	fo.Close()
	log.Printf("DONE: wrote %s got %d unique words from %d", *outfile, len(keys), total)
}
