package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/client9/misspell"
)

var (
	defaultWrite *template.Template
	defaultRead  *template.Template

	stdout *log.Logger // see below in init()
)

const (
	defaultWriteTmpl = `{{ .Filename }}:{{ .Line }}:{{ .Column }}:corrected "{{ js .Original }}" to "{{ js .Corrected }}"`
	defaultReadTmpl  = `{{ .Filename }}:{{ .Line }}:{{ .Column }}:found "{{ js .Original }}" a misspelling of "{{ js .Corrected }}"`
	csvTmpl          = `{{ printf "%q" .Filename }},{{ .Line }},{{ .Column }},{{ .Original }},{{ .Corrected }}`
	csvHeader        = `file,line,column,typo,corrected`
	sqliteTmpl       = `INSERT INTO misspell VALUES({{ printf "%q" .Filename }},{{ .Line }},{{ .Column }},{{ printf "%q" .Original }},{{ printf "%q" .Corrected }});`
	sqliteHeader     = `PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE misspell(
	"file" TEXT, "line" INTEGER, "column" INTEGER, "typo" TEXT, "corrected" TEXT
);`
	sqliteFooter = "COMMIT;"
)

func init() {
	defaultWrite = template.Must(template.New("defaultWrite").Parse(defaultWriteTmpl))
	defaultRead = template.Must(template.New("defaultRead").Parse(defaultReadTmpl))

	// we cant't just write to os.Stdout directly since we have multiple goroutine
	// all writing at the same time causing broken output.  Log is routine safe.
	// we see it so it doesn't use a prefix or include a time stamp.
	stdout = log.New(os.Stdout, "", 0)
}

func worker(writeit bool, r *misspell.Replacer, mode string, files <-chan string, results chan<- int) {
	count := 0
	for filename := range files {
		orig, err := misspell.ReadTextFile(filename)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(orig) == 0 {
			continue
		}

		updated, changes := r.Replace(orig)

		if len(changes) == 0 {
			continue
		}
		count += len(changes)
		for _, diff := range changes {
			// add in filename
			diff.Filename = filename

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
	results <- count
}

func main() {
	quiet := flag.Bool("q", false, "quiet")
	workers := flag.Int("j", 0, "Number of workers, 0 = number of CPUs")
	writeit := flag.Bool("w", false, "Overwrite file with corrections (default is just to display)")
	format := flag.String("f", "", "'csv', 'sqlite3' or custom Golang template for output")
	ignores := flag.String("i", "", "ignore the following corrections, comma separated")
	locale := flag.String("locale", "", "Correct spellings using locale perferances for US or UK.  Default is to use a neutral variety of English.  Setting locale to US will correct the British spelling of 'colour' to 'color'")
	mode := flag.String("source", "auto", "Source mode: auto=guess, go=golang source, text=plain or markdown-like text")
	debug := flag.Bool("debug", false, "Debug matching, very slow")
	exitError := flag.Bool("error", false, "Exit with 2 if misspelling found")

	flag.Parse()

	r := misspell.Replacer{
		Replacements: misspell.DictMain,
		Debug:        *debug,
	}

	//
	// Figure out regional variations
	//
	switch strings.ToUpper(*locale) {
	case "":
		// nothing
	case "US":
		r.AddRuleList(misspell.DictAmerican)
	case "UK", "GB":
		r.AddRuleList(misspell.DictBritish)
	case "NZ", "AU", "CA":
		log.Fatalf("Help wanted.  https://github.com/client9/misspell/issues/6")
	default:
		log.Fatalf("Unknow locale: %q", *locale)
	}

	//
	// Stuff to ignore
	//
	if len(*ignores) > 0 {
		r.RemoveRule(strings.Split(*ignores, ","))
	}

	//
	// Source input mode
	//
	switch *mode {
	case "auto":
	case "go":
	case "text":
	default:
		log.Fatalf("Mode must be one of auto=guess, go=golang source, text=plain or markdown-like text")
	}

	//
	// Custom output
	//
	switch *format {
	case "csv":
		tmpl := template.Must(template.New("csv").Parse(csvTmpl))
		defaultWrite = tmpl
		defaultRead = tmpl
		stdout.Println(csvHeader)
	case "sqlite", "sqlite3":
		tmpl := template.Must(template.New("sqlite3").Parse(sqliteTmpl))
		defaultWrite = tmpl
		defaultRead = tmpl
		stdout.Println(sqliteHeader)
	default:
		if len(*format) > 0 {
			t, err := template.New("custom").Parse(*format)
			if err != nil {
				log.Fatalf("Unable to compile log format: %s", err)
			}
			defaultWrite = t
			defaultRead = t
		}
	}

	//
	// Number of Workers / CPU to use
	//
	if *workers < 0 {
		log.Fatalf("-j must >= 0")
	}
	if *workers == 0 {
		*workers = runtime.NumCPU()
	}
	if *debug {
		*workers = 1
	}

	//
	// Done with Flags.
	//  Compile the Replacer and process files
	//
	r.Compile()

	args := flag.Args()

	// unix style pipes: different output module
	// stdout: output of corrected text
	// stderr: output of log lines
	if len(args) == 0 {
		raw, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("Unable to read stdin")
		}
		orig := string(raw)
		updated, changes := r.Replace(orig)
		if !*quiet {
			for _, diff := range changes {
				diff.Filename = "stdin"
				var output bytes.Buffer
				if *writeit {
					defaultWrite.Execute(&output, diff)
				} else {
					defaultRead.Execute(&output, diff)
				}
				stdout.Println(output.String())
			}
		}
		if *writeit {
			stdout.Println(updated)
		}
		switch *format {
		case "sqlite", "sqlite3":
			stdout.Println(sqliteFooter)
		}
		return
	}

	c := make(chan string, 64)
	results := make(chan int, *workers)

	for i := 0; i < *workers; i++ {
		go worker(*writeit, &r, *mode, c, results)
	}

	for _, filename := range args {
		filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				c <- path
			}
			return nil
		})
	}
	close(c)

	count := 0
	for i := 0; i < *workers; i++ {
		changed := <-results
		count += changed
	}

	switch *format {
	case "sqlite", "sqlite3":
		stdout.Println(sqliteFooter)
	}

	if count != 0 && *exitError {
		// log.Printf("Got %d", count)
		// error
		os.Exit(2)
	}
}
