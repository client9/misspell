CONTAINER=nickg/misspell

all: install lint test

install:
	cp -f precommit.sh .git/hooks/pre-commit
	go install ./cmd/misspell

lint: 
	gometalinter \
		 --vendor \
		 --deadline=60s \
	         --disable-all \
		 --enable=vet \
		 --enable=golint \
		 --enable=gofmt \
		 --enable=goimports \
		 --enable=gosimple \
		 --enable=staticcheck \
		 --enable=ineffassign \
		 ./...

test: 
	go test .

bench:
	go test -bench '.*'

# the grep in line 2 is to remove misspellings in the spelling dictionary
# that trigger false positives!!
falsepositives: /scowl-wl
	cat /scowl-wl/words-US-60.txt | \
		grep -i -v -E "payed|Tyre|Euclidian|nonoccurence|dependancy|reenforced|accidently|surprize|dependance|idealogy|binominal|causalities|conquerer|withing|casette|analyse|analogue|dialogue|paralyse|catalogue|archaeolog|clarinettist|catalyses|cancell|chisell|ageing|cataloguing" | \
		misspell -locale=US -debug -error
	cat /scowl-wl/words-US-60.txt | tr '[:lower:]' '[:upper:]' | \
		grep -i -v -E "payed|Tyre|Euclidian|nonoccurence|dependancy|reenforced|accidently|surprize|dependance|idealogy|binominal|causalities|conquerer|withing|casette|analyse|analogue|dialogue|paralyse|catalogue|archaeolog|clarinettist|catalyses|cancell|chisell|ageing|cataloguing" | \
		 misspell -locale=US -debug -error
	cat /scowl-wl/words-GB-ise-60.txt | \
		grep -v -E "payed|nonoccurence|withing" | \
		misspell -locale=UK -debug -error
	cat /scowl-wl/words-GB-ise-60.txt | tr '[:lower:]' '[:upper:]' | \
		grep -i -v -E "payed|nonoccurence|withing" | \
		misspell -debug -error
#	cat /scowl-wl/words-GB-ize-60.txt | \
#		grep -v -E "withing" | \
#		misspell -debug -error
#	cat /scowl-wl/words-CA-60.txt | \
#		grep -v -E "withing" | \
#		misspell -debug -error

clean:
	rm -f *~
	go clean ./...
	git gc

ci-native: install lint test falsepositives

# when development is already in a docker contianer
ci:
	docker run --rm \
		--volumes-from=workspace \
		-w /go/src/github.com/client9/misspell \
		${CONTAINER} \
		make ci-native

# for on-mac and travis builds
ci-travis:
	docker --version
	docker run --rm \
		-v $(PWD):/go/src/github.com/client9/misspell \
		-w /go/src/github.com/client9/misspell \
		${CONTAINER} \
		make ci-native 

docker-build:
	docker build -t ${CONTAINER} .

console:
	docker run --rm -it \
		--volumes-from=workspace \
		-w /go/src/github.com/client9/misspell \
		${CONTAINER} sh

precommit:
	docker run --rm \
		$(shell dmnt .) \
		-w /go/src/github.com/client9/misspell \
		${CONTAINER} \
		make

.PHONY: ci ci-travis ci-native
