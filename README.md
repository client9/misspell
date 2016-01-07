Correct commonly misspelled English words... quickly.

## FAQ

### What problem does this solve?

This attempts to correct commonly misspelled English words in source
code and documentation.  It is not a complete spell-checking program
nor a grammar checker.

### What are other misspelling correctors and what's wrong with them?

Some other misspelling correctors:

* https://github.com/vlajos/misspell_fixer
* https://github.com/lyda/misspell-check
* https://github.com/lucasdemarchi

They all work but had problems that prevented me from using them at scale:

* slow, all of the above check one misspelling at a time (i.e. linear) using regexps
* not MIT/Apache2 licensed (or equivalent)
* have dependencies that don't work for me (python3, bash, linux sed, etc)

That said, they might be perfect for you and many have more features
that this project!

### How much faster is this project?

Easily 100x to 1000x faster.  You should be able to check and correct
1000 files in under 250ms.

### What license is this?

[MIT](https://github.com/client9/misspell/blob/master/LICENSE)

### What are the dependencies?

You need [golang 1.5](https://golang.org/) to compile this, but the resulting binary has no
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
This makes multiple substring matches *simultaneously*

In addition this uses multiple CPU cores to works on multiple files.

### What problems does it have?

Unlike the other projects, this doesn't know what a "word" is.  There
may be more false positives and false negatives due to this.  On the
other hand, it sometimes catches things others don't.

Either way, please file bugs and we'll fix them!

Since it operates in parallel to make corrections, it can be
non-obvious to determine exactly what word was corrected.

### It's making mistakes.  How can I debug?

Run using `-debug` flag on the file you want.  It should then
print what word it is trying to correct.  Then [file a bug](https://github.com/client9/misspell/issues) describing the
problem.  Thanks!

### Why is making mistakes or missing items in golang files?

The matching function is *case-sensitive*, so variable names that are
multiple worlds either in all-upper or all-lower case sometimes can
cause false positives.  For instance a variable named `bodyreader`
could trigger a false positive since `yrea` is in the middle that
could be corrected to `year`.  Other problems happen if the variable
name uses a English contraction that should use an apostrophe.  The
best way of fixing this is to use the [Effective Go naming
conventions](https://golang.org/doc/effective_go.html#mixed-caps) and
use camelCase for variable names.  You can check your code using
[golint](https://github.com/golang/lint)

### Does this work with gometalinter?

[gometalinter](https://github.com/alecthomas/gometalinter) runs
multiple golang linters, and it works well with `misspell` too.

After `go get -u github.com/client9/misspell` you need to add it, then
enable it, like so:

```bash
gometalinter --disable-all \
   --linter='misspell:misspell ./*.go:PATH:LINE:MESSAGE' --enable=misspell \
   ./...
```
