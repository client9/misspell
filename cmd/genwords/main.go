package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/client9/gospell"
)

func addOrPanic(dict map[string]string, key, value string) {
	if _, ok := dict[key]; ok {
		log.Printf("Already have %q", key)
	}
	dict[key] = value
}

func mergeDict(a, b map[string]string) {
	for k, v := range b {
		style := gospell.CaseStyle(k)
		kcase := gospell.CaseVariations(k, style)
		vcase := gospell.CaseVariations(v, style)
		for i := 0; i < len(kcase); i++ {
			addOrPanic(a, kcase[i], vcase[i])
		}
	}
}

func parseWikipediaFormat(text string) map[string]string {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	dict := make(map[string]string, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "->")
		if len(parts) != 2 {
			log.Fatalf(fmt.Sprintf("failed parsing %q", line))
		}
		spellings := strings.Split(parts[1], ",")
		dict[parts[0]] = strings.TrimSpace(spellings[0])
	}
	return dict
}

func main() {
	dict := make(map[string]string)
	mergeDict(dict, dictWikipedia())
	mergeDict(dict, dictAdditions())
	words := make([]string, 0, len(dict))
	for k := range dict {
		words = append(words, k)
	}
	sort.Strings(words)

	fmt.Printf("package misspell\n\n")
	fmt.Printf("var dictWikipedia = []string{\n")
	for _, word := range words {
		fmt.Printf("\t%q, %q,\n", word, dict[word])
	}
	fmt.Printf("}\n")
}
