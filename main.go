package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"text/template"

	"github.com/client9/misspell/lib"
)

var defaultWrite *template.Template
var defaultRead *template.Template

func init() {
	defaultWrite = template.Must(template.New("defaultWrite").Parse(`{{ .Filename }}:{{ .Line }} corrected "{{ js .Original }}" to "{{ js .Corrected }}"`))
	defaultRead = template.Must(template.New("defaultRead").Parse(`{{ .Filename }}:{{ .Line }} found "{{ js .Original }}" a misspelling of "{{ js .Corrected }}"`))

}

func worker(writeit bool, debug bool, mode string, files <-chan string, results chan<- int) {
	fails := 0
	for filename := range files {
		isGolang := strings.HasSuffix(filename, ".go")

		// ignore directories
		if f, err := os.Stat(filename); f.IsDir() && err == nil {
			continue
		}

		raw, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("Unable to read %q: %s", filename, err)
			continue
		}
		orig := string(raw)
		var updated string

		// GROSS
		if mode == "go" || (mode == "auto" && isGolang) {
			updated = lib.ReplaceGo(orig, debug)
		} else if debug {
			updated = lib.ReplaceDebug(orig)
		} else {
			updated = lib.Replace(orig)
		}

		changes := lib.DiffLines(filename, orig, updated)
		if len(changes) == 0 {
			continue
		}
		for _, diff := range changes {
			if writeit {
				defaultWrite.Execute(os.Stdout, diff)
			} else {
				// the log package can be used simultaneously from multiple goroutines
				var output bytes.Buffer
				defaultRead.Execute(&output, diff)
				log.Println(output.String())
			}
		}

		if writeit {
			ioutil.WriteFile(filename, []byte(updated), 0)
		}
	}
	results <- fails
}

func main() {
	workers := flag.Int("j", 0, "Number of workers, 0 = number of CPUs")
	writeit := flag.Bool("w", false, "Overwrite file with corrections (default is just to display)")
	format := flag.String("f", "", "Use Golang template for log message")
	ignores := flag.String("i", "", "Ignore the following corrections, comma separated")
	mode := flag.String("source", "auto", "Source mode: auto=guess, go=golang source, text=plain or markdown-like text")
	debug := flag.Bool("debug", false, "Debug matching, very slow")
	flag.Parse()

	switch *mode {
	case "auto":
	case "go":
	case "text":
	default:
		log.Fatalf("Mode must be one of auto=guess, go=golang source, text=plain or markdown-like text")
	}
	if len(*format) > 0 {
		t, err := template.New("custom").Parse(*format)
		if err != nil {
			log.Fatalf("Unable to compile log format: %s", err)
		}
		defaultWrite = t
		defaultRead = t
	}
	if len(*ignores) > 0 {
		lib.Ignore(strings.Split(*ignores, ","))
	}
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
		go worker(*writeit, *debug, *mode, c, results)
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
