Correct commonly misspelled English words... quickly.

## FAQ

### What problem does this solve?

This attempts to correct commonly misspelled English words in source
code and documentation.  It is not a complete spell-checking program
or grammar checker.

### What are other misspelling correctors and what's wrong with them?

I took a look at the following:

* https://github.com/vlajos/misspell_fixer
* https://github.com/lyda/misspell-check
* https://github.com/lucasdemarchi

They are work but had problems that prevented me from using them at scale:

* slow, all of the above check one misspelling at a time (i.e. linear) using regexps
* not MIT/Apache2 licensed (or equivalent)
* have dependencies that don't work for me (python3, bash, linux sed, etc)

That said, they might be perfect for you and many have more features
that this project!

### How much faster is this project?

Easily 100x to 100x faster.  You should be able to check and correct
1000 files in under 250ms.

### What license is this?

MIT

### What are the depedencies?

You need golang 1.5 to compile this, but the resulting binary has no
dependencies.  If people want precompiled binaries for various
platforms, let me know.

### Where do the word lists come from?

It's currently pulled from
[Wikipedia](https://en.wikipedia.org/wiki/Wikipedia:Lists_of_common_misspellings/For_machines)
and then edited to remove false positives.

### Why is this so fast?

This uses the mighty power of golang's
[strings.Replacer](https://golang.org/pkg/strings/#Replacer) which is
a implementation or variation of the
[Aho–Corasick algorithm](https://en.wikipedia.org/wiki/Aho–Corasick_algorithm).
This makes multiple substring matches *simultanouesly*

In addition this uses multiple CPU cores to works on multiple files

### What problems does it have?

Unlike the other projects, this doesn't know what a "word" is.  There
may be more false positives and false negatives due to this.  Please
file bugs and we'll fix them!

Since it operates in parallel to make corrections, it can be non
obvious to determine exactly what word was corrected.

### It's making mistakes.  How can I debug?

Run using `-debug -dry` flag on the file you want.  It should then
print what word it is trying to correct.  Then file a bug describing
the problem.  Thanks!

### Does this work with gometalinter?

[gometalinter](https://github.com/alecthomas/gometalinter) runs multiple golang linters.

And yes, it works with it.

After `go get -u github.com/client9/misspell`

You need to add it, then enable it like so:

```bash
gometalinter --disable-all \
   --linter='misspell:misspell -dry ./*.go:PATH:LINE:MESSAGE' --enable=misspell \
   ./...
```
