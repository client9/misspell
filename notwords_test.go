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
		{word: "x infinitie.net y", want: "x               y"},
		{word: "(s.svc.GetObject(", want: "(               ("},
		{word: "\\nto", want: "  to"},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.word, func(t *testing.T) {
			t.Parallel()

			got := RemoveNotWords(test.word)
			if got != test.want {
				t.Errorf("want %q got %q", test.want, got)
			}
		})
	}
}
