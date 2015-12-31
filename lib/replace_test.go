package lib

import (
	"testing"
)

func TestReplace(t *testing.T) {
	cases := []struct{
		orig string
		want string
	}{
		{"I live in Amercia", "I live in America"},
		{"There is a zeebra", "There is a zebra"},
		{"foo other bar", "foo other bar"},
	}
	for line, tt := range cases {
		got := Replace(tt.orig, false)
		if got != tt.want {
			t.Errorf("%d: Replace files want %q got %q", line, tt.orig, got)
		}
	}
}
