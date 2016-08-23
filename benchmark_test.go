package misspell

import (
	"bytes"
	"io/ioutil"
	"testing"
)

var (
	sampleClean string
	sampleDirty string
	tmpCount    int
	tmp         string
	rep         *Replacer
)

func init() {

	buf := bytes.Buffer{}
	for i := 0; i < len(DictMain); i += 2 {
		buf.WriteString(DictMain[i+1] + "\n")
	}
	sampleClean = buf.String()
	sampleDirty = sampleClean + DictMain[0] + "\n"
	rep = New()

}

// BenchmarkCleanString takes a clean string (one with no errors)
func BenchmarkCleanString(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	var updated string
	var diffs []Diff
	var count int
	for n := 0; n < b.N; n++ {
		updated, diffs = rep.Replace(sampleClean)
		count += len(diffs)
	}

	// prevent compilier optimizations
	tmpCount = count
	tmp = updated
}

// BenchmarkCleanStream takes a clean reader (no misspells) and outputs to a buffer
func BenchmarkCleanStream(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	count := 0
	buf := bytes.NewBufferString(sampleClean)
	out := bytes.NewBuffer(make([]byte, 0, len(sampleClean)+100))
	for n := 0; n < b.N; n++ {
		buf.Reset()
		buf.WriteString(sampleClean)
		out.Reset()
		diffs := rep.ReplaceReader(buf, out)
		count += len(diffs)
	}
	tmpCount = count
}

// BenchmarkCleanStreamDiscard takes a clean reader and discards output
func BenchmarkCleanStreamDiscard(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	buf := bytes.NewBufferString(sampleClean)
	count := 0
	for n := 0; n < b.N; n++ {
		buf.Reset()
		buf.WriteString(sampleClean)
		diffs := rep.ReplaceReader(buf, ioutil.Discard)
		count += len(diffs)
	}

	tmpCount = count
}

// BenchmarkCleanString takes a clean string (one with no errors)
func BenchmarkDirtyString(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	var updated string
	var diffs []Diff
	var count int
	for n := 0; n < b.N; n++ {
		updated, diffs = rep.Replace(sampleDirty)
		count += len(diffs)
	}

	// prevent compilier optimizations
	tmpCount = count
	tmp = updated
}

// BenchmarkCleanStreamDiscard takes a clean reader and discards output
func BenchmarkDirtyStreamDiscard(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	buf := bytes.NewBufferString(sampleDirty)
	count := 0
	for n := 0; n < b.N; n++ {
		buf.Reset()
		buf.WriteString(sampleDirty)
		diffs := rep.ReplaceReader(buf, ioutil.Discard)
		count += len(diffs)
	}

	tmpCount = count
}
