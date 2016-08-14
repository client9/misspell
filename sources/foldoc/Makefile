
words-foldoc.txt: download.go
	go run download.go

test: words-foldoc.txt
	misspell words-foldoc.txt
