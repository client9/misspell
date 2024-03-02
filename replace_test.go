package misspell

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
		r := New()
		r.RemoveRule(strings.Split(tt.ignore, ","))
		r.Compile()
		got, _ := r.Replace(tt.text)
		if got != tt.text {
			t.Errorf("%d: Replace files want %q got %q", line, tt.text, got)
		}
	}
}

func TestReplaceLocale(t *testing.T) {
	cases := []struct {
		orig string
		want string
	}{
		{orig: "The colours are pretty", want: "The colors are pretty"},
		{orig: "summaries", want: "summaries"},
	}

	r := New()
	r.AddRuleList(DictAmerican)
	r.Compile()

	for _, test := range cases {
		test := test
		t.Run(test.orig, func(t *testing.T) {
			t.Parallel()

			got, _ := r.Replace(test.orig)
			if got != test.want {
				t.Errorf("ReplaceLocale want %q got %q", test.orig, got)
			}
		})
	}
}

func TestReplace(t *testing.T) {
	cases := []struct {
		orig string
		want string
	}{
		{orig: "I live in Amercia", want: "I live in America"},
		{orig: "grill brocoli now", want: "grill broccoli now"},
		{orig: "There is a zeebra", want: "There is a zebra"},
		{orig: "foo other bar", want: "foo other bar"},
		{orig: "ten fiels", want: "ten fields"},
		{orig: "Closeing Time", want: "Closing Time"},
		{orig: "closeing Time", want: "closing Time"},
		{orig: " TOOD: foobar", want: " TODO: foobar"},
		{orig: " preceed ", want: " precede "},
		{orig: "preceeding", want: "preceding"},
		{orig: "functionallity", want: "functionality"},
	}

	r := New()

	for _, test := range cases {
		test := test
		t.Run(test.orig, func(t *testing.T) {
			t.Parallel()

			got, _ := r.Replace(test.orig)
			if got != test.want {
				t.Errorf("Replace files want %q got %q", test.orig, got)
			}
		})
	}
}

func TestCheckReplace(t *testing.T) {
	t.Run("Nothing at all", func(t *testing.T) {
		t.Parallel()

		r := Replacer{
			engine:    NewStringReplacer("foo", "foobar"),
			corrected: map[string]string{"foo": "foobar"},
		}

		s := "nothing at all"
		news, diffs := r.Replace(s)
		if s != news || len(diffs) != 0 {
			t.Errorf("Basic recheck failed: %q vs %q", s, news)
		}
	})

	t.Run("Single, correct,.Correctedacements", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			orig     string
			expected string
		}{
			{
				orig:     "foo",
				expected: "foobar",
			},
			{
				orig:     "foo junk",
				expected: "foobar junk",
			},
			{
				orig:     "junk foo",
				expected: "junk foobar",
			},
			{
				orig:     "junk foo junk",
				expected: "junk foobar junk",
			},
		}

		for _, test := range testCases {
			orig := test.orig
			expected := test.expected
			t.Run(orig, func(t *testing.T) {
				t.Parallel()

				r := Replacer{
					engine:    NewStringReplacer("foo", "foobar", "runing", "running"),
					corrected: map[string]string{"foo": "foobar", "runing": "running"},
				}

				news, diffs := r.Replace(orig)
				if news != expected || len(diffs) != 1 || diffs[0].Original != "foo" && diffs[0].Corrected != expected && diffs[0].Column != 0 {
					t.Errorf("basic recheck failed %q vs %q", expected, news)
				}
			})
		}
	})

	t.Run("Incorrect.Correctedacements", func(t *testing.T) {
		t.Parallel()

		r := Replacer{
			engine:    NewStringReplacer("foo", "foobar"),
			corrected: map[string]string{"foo": "foobar"},
		}

		s := "food pruning"
		news, _ := r.Replace(s)
		if news != s {
			t.Errorf("incorrect.Correctedacement failed: %q vs %q", s, news)
		}
	})
}
