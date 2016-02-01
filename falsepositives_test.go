package misspell

import (
	"testing"
)

func TestFalsePositives(t *testing.T) {
	cases := []string{
		"setuptools", // python package
		"committing",
		"guises",
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
