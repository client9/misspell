package misspell

import (
	"reflect"
	"testing"
)

func TestCaseStyle(t *testing.T) {
	testCases := []struct {
		word string
		want WordCase
	}{
		{word: "lower", want: CaseLower},
		{word: "what's", want: CaseLower},
		{word: "UPPER", want: CaseUpper},
		{word: "Title", want: CaseTitle},
		{word: "CamelCase", want: CaseUnknown},
		{word: "camelCase", want: CaseUnknown},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.word, func(t *testing.T) {
			t.Parallel()

			got := CaseStyle(test.word)
			if test.want != got {
				t.Errorf("want %v got %v", test.want, got)
			}
		})
	}
}

func TestCaseVariations(t *testing.T) {
	testCases := []struct {
		word string
		want []string
	}{
		{word: "that's", want: []string{"that's", "That's", "THAT'S"}},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.word, func(t *testing.T) {
			t.Parallel()

			got := CaseVariations(test.word, CaseStyle(test.word))
			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("want %v got %v", test.want, got)
			}
		})
	}
}
