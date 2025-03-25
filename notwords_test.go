package misspell

import (
	"testing"
)

func TestNotWords(t *testing.T) {
	testCases := []struct {
		word string
		want string
	}{
		{word: " /foo/bar abc", want: "          abc"},
		{word: "X/foo/bar abc", want: "X/foo/bar abc"},
		{word: "[/foo/bar] abc", want: "[        ] abc"},
		{word: "/", want: "/"},
		{word: "x nickg@client9.xxx y", want: "x                   y"},
		{word: "x fqdn.example.org. y", want: "x                 . y"},
		{word: "x infinitie.net y", want: "x               y"},
		{word: "x infinitie.net ", want: "x               "},
		{word: "x infinitie.net", want: "x              "},
		{word: "x foo.example.com y", want: "x                 y"},
		{word: "x foo.example.com ", want: "x                 "},
		{word: "x foo.example.com", want: "x                "},
		{word: "foo.example.com y", want: "                y"},
		{word: "foo.example.com", want: "               "},
		{word: "(s.svc.GetObject(", want: "(s.svc.GetObject("},
		{word: "defer file.Close()", want: "defer file.Close()"},
		{word: "defer file.c()", want: "defer file.c()"},
		{word: "defer file.cl()", want: "defer        ()"},       // false negative
		{word: "defer file.close()", want: "defer           ()"}, // false negative
		{word: "\\nto", want: "  to"},
	}

	for _, test := range testCases {
		t.Run(test.word, func(t *testing.T) {
			t.Parallel()

			got := RemoveNotWords(test.word)
			if got != test.want {
				t.Errorf("want %q got %q", test.want, got)
			}
		})
	}
}
