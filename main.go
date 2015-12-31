package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/client9/misspell/lib"
)

func worker(writeit bool, debug bool, files <-chan string, results chan<- int) {
	fails := 0
	for filename := range files {
		//log.Printf("Scanning %q", filename)
		raw, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("Unable to read %q: %s", filename, err)
			continue
		}
		orig := string(raw)
		updated := lib.Replace(orig, debug)

		count := lib.DiffLines(filename, orig, updated, os.Stdout)
		if count == 0 {
			continue
		}
		fails += count
		//log.Printf("Updating %q", filename)
		if writeit {
			ioutil.WriteFile(filename, []byte(updated), 0)
		}
	}
	results <- fails
}

func main() {
	workers := flag.Int("j", 0, "Number of workers, 0 = number of CPUs")
	writeit := flag.Bool("w", false, "Write correction to file (default is just to display)")
	debug := flag.Bool("debug", false, "Debug matching, very slow")
	flag.Parse()

	if *workers < 0 {
		log.Fatalf("-j must >= 0")
	}
	if *workers == 0 {
		*workers = runtime.NumCPU()
	}

	if *debug {
		*workers = 1
	}
	args := flag.Args()

	/*
		if len(args) == 0 {
			rawin, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				panic(err)
			}
			replacer.WriteString(os.Stdout, string(rawin))
			return
		}
	*/
	c := make(chan string, len(args))
	results := make(chan int, *workers)

	for i := 0; i < *workers; i++ {
		go worker(*writeit, *debug, c, results)
	}

	for _, filename := range args {
		c <- filename
	}
	close(c)

	count := 0
	for i := 0; i < *workers; i++ {
		changed := <-results
		count += changed
	}
	if count != 0 {
		// log.Printf("Got %d", count)
		// error
		os.Exit(1)
	}
}
