package lib

import (
	"testing"
)

func TestReplaceIgnore(t *testing.T) {
	cases := []struct {
		ignore string
		text   string
	}{
		{"knwo", "https://github.com/Unknwon"},
		{"gae", "github.com/hnakamur/gaesessions"},
	}

	for line, tt := range cases {
		Ignore(tt.ignore)
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
