package lib

import (
	"strings"
	"testing"
)

func TestReplaceIgnore(t *testing.T) {
	cases := []struct {
		ignore string
		text   string
	}{
		{"knwo,gae", "https://github.com/Unknwon, github.com/hnakamur/gaesessions"},
	}

	for line, tt := range cases {
		Ignore(strings.Split(tt.ignore, ","))
		got := ReplaceDebug(tt.text)
		if got != tt.text {
			t.Errorf("%d: Replace files want %q got %q", line, tt.text, got)
		}
	}
}

func TestReplace(t *testing.T) {
	cases := []struct {
		orig string
		want string
	}{
		{"I live in Amercia", "I live in America"},
		{"There is a zeebra", "There is a zebra"},
		{"foo other bar", "foo other bar"},
		{"ten fiels", "ten fields"},
		{"Closeing Time", "Closing Time"},
		{"closeing Time", "closing Time"},
	}
	for line, tt := range cases {
		got := Replace(tt.orig)
		if got != tt.want {
			t.Errorf("%d: Replace files want %q got %q", line, tt.orig, got)
		}
	}
}

func TestFalsePositives(t *testing.T) {
	cases := []string{
		"disguise",
		"begging",
		"cmo",
		"cmos",
		"borked",
		"hadn't",
		"Iceweasel",
		"summarised",
		"autorenew",
		"travelling",
		"republished",
		"fallthru",
		"pruning",
		"deb.VersionDontCare",
		"authtag",
		"intrepid",
		"usefully",
		"there",
		"definite",
		"earliest",
		"Japanese",
		"international",
		"excellent",
		"gracefully",
		"carefully",
		"class",
		"include",
		"process",
		"address",
		"attempt",
		"large",
		"although",
		"specific",
		"taste",
		"against",
		"successfully",
		"unsuccessfully",
		"occurred",
		"agree",
		"controlled",
		"publisher",
		"strategy",
		"geoposition",
		"paginated",
		"happened",
		"relative",
		"computing",
		"language",
		"manual",
		"token",
		"into",
		"nothing",
		"datatool",
		"propose",
		"learnt",
		"tolerant",
		"whitehat",
		"monotonic",
		"comprised",
		"indemnity",
		"flattened",
		"interrupted",
		"inotify",
		"occasional",
		"forging",
		"ampersand",
		"decomposition",
		"commit",
		"programmer", // "grammer"
		//		"requestsinserted",
		"seeked",      // technical word
		"bodyreader",  // variable name
		"cantPrepare", // variable name
		"dontPrepare", // variable name
	}
	for casenum, tt := range cases {
		got := ReplaceDebug(tt)
		if got != tt {
			t.Errorf("%d: %q got converted to %q", casenum, tt, got)
		}
	}
}

func TestReplaceGo(t *testing.T) {
	cases := []struct {
		orig string
		want string
	}{
		{
			orig: `
// I am a zeebra
var foo 10
`,
			want: `
// I am a zebra
var foo 10
`,
		},
		{
			orig: `
var foo 10
// I am a zeebra`,
			want: `
var foo 10
// I am a zebra`,
		},
		{
			orig: `
// I am a zeebra
var foo int
/* multiline
 * zeebra
 */
`,
			want: `
// I am a zebra
var foo int
/* multiline
 * zebra
 */
`,
		},
	}

	for casenum, tt := range cases {
		got := ReplaceGo(tt.orig, true)
		if got != tt.want {
			t.Errorf("%d: %q got converted to %q", casenum, tt, got)
		}
	}
}
