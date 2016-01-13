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

	"github.com/client9/misspell"
)

var (
	defaultWrite *template.Template
	defaultRead  *template.Template
	stdout       *log.Logger // see below in init()
)

const (
	defaultWriteTmpl = `{{ .Filename }}:{{ .Line }} corrected "{{ js .Original }}" to "{{ js .Corrected }}"`
	defaultReadTmpl  = `{{ .Filename }}:{{ .Line }} found "{{ js .Original }}" a misspelling of "{{ js .Corrected }}"`
)

func init() {
	defaultWrite = template.Must(template.New("defaultWrite").Parse(defaultWriteTmpl))
	defaultRead = template.Must(template.New("defaultRead").Parse(defaultReadTmpl))

	// we cant't just write to os.Stdout directly since we have multiple goroutine
	// all writing at the same time causing broken output.  Log is routine safe.
	// we see it so it doesn't use a prefix or include a time stamp.
	stdout = log.New(os.Stdout, "", 0)
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
			updated = misspell.ReplaceGo(orig, debug)
		} else if debug {
			updated = misspell.ReplaceDebug(orig)
		} else {
			updated = misspell.Replace(orig)
		}

		updated, changes := misspell.DiffLines(filename, orig, updated)
		if len(changes) == 0 {
			continue
		}
		for _, diff := range changes {
			// output can be done by doing multiple goroutines
			// and can clobber os.Stdout.
			//
			// the log package can be used simultaneously from multiple goroutines
			var output bytes.Buffer
			if writeit {
				defaultWrite.Execute(&output, diff)
			} else {
				defaultRead.Execute(&output, diff)
			}

			// goroutine-safe print to os.Stdout
			stdout.Println(output.String())
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
		err := misspell.Ignore(strings.Split(*ignores, ","))
		if err != nil {
			log.Fatalf("unable to process ignores: %s", err)
		}
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
