package ignore

import (
	"testing"
)

func TestParseMatchSingle(t *testing.T) {
	testCases := []struct {
		pattern  string
		filename string
		want     bool
	}{
		{pattern: "*.c", filename: "foo.c", want: true},
		{pattern: "*.c", filename: "foo/bar.c", want: true},
		{pattern: "Documentation/*.html", filename: "Documentation/git.html", want: true},
		{pattern: "Documentation/*.html", filename: "Documentation/ppc/ppc.html"},
		{pattern: "/*.c", filename: "cat-file.c", want: true},
		{pattern: "/*.c", filename: "mozilla-sha1/sha1.c"},
		{pattern: "foo", filename: "foo", want: true},
		{pattern: "**/foo", filename: "./foo", want: true}, // <--- leading './' required
		{pattern: "**/foo", filename: "junk/foo", want: true},
		{pattern: "**/foo/bar", filename: "./foo/bar", want: true}, // <--- leading './' required
		{pattern: "**/foo/bar", filename: "junk/foo/bar", want: true},
		{pattern: "abc/**", filename: "abc/foo", want: true},
		{pattern: "abc/**", filename: "abc/foo/bar", want: true},
		{pattern: "a/**/b", filename: "a/b", want: true},
		{pattern: "a/**/b", filename: "a/x/b", want: true},
		{pattern: "a/**/b", filename: "a/x/y/b", want: true},

		{pattern: "*_test*", filename: "foo_test.go", want: true},
		{pattern: "*_test*", filename: "junk/foo_test.go", want: true},
		{pattern: "junk\n!junk", filename: "foo"},
		{pattern: "junk\n!junk", filename: "junk"},

		{pattern: "*.html\n!foo.html", filename: "junk.html", want: true},
		{pattern: "*.html\n!foo.html", filename: "foo.html"},

		{pattern: "/*\n!/foo\n/foo/*\n!/foo/bar", filename: "crap", want: true},
		{pattern: "/*\n!/foo\n/foo/*\n!/foo/bar", filename: "foo/crap", want: true},
		{pattern: "/*\n!/foo\n/foo/*\n!/foo/bar", filename: "foo/bar"},
		{pattern: "/*\n!/foo\n/foo/*\n!/foo/bar", filename: "foo/bar/other"},
		{pattern: "/*\n!/foo\n/foo/*\n!/foo/bar", filename: "foo"},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.pattern, func(t *testing.T) {
			t.Parallel()

			matcher, err := Parse([]byte(test.pattern))
			if err != nil {
				t.Errorf("error: %s", err)
			}

			got := matcher.Match(test.filename)
			if test.want != got {
				t.Errorf("%q.Match(%q) = %v, got %v", test.pattern, test.filename, test.want, got)
			}
		})
	}
}
