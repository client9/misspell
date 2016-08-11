package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/xrash/smetrics"
)

var badWord = map[string]bool{
	"shitter":       true,
	"shittier":      true,
	"lightsaber":    true,
	"einstein":      true,
	"samsung":       true,
	"googled":       true,
	"dragonborn":    true,
	"circlejerk":    true,
	"cinderhulk":    true,
	"shithead":      true,
	"sanchez":       true,
	"ronaldo":       true,
	"retards":       true,
	"michael":       true,
	"killionaires":  true,
	"villionaires":  true,
	"sverige":       true,
	"snapchat":      true,
	"siempre":       true,
	"shittiest":     true,
	"shitposting":   true,
	"warpig":        true,
	"shitlord":      true,
	"illidan":       true,
	"magneto":       true,
	"martinez":      true,
	"darkrai":       true,
	"fuckton":       true,
	"deandre":       true,
	"whatsapp":      true,
	"wrestlemania":  true,
	"youtuber":      true,
	"godzilla":      true,
	"charles":       true,
	"cortana":       true,
	"monsanto":      true,
	"horseshit":     true,
	"balotelli":     true,
	"aquaman":       true,
	"antonio":       true,
	"happend":       true,
	"arduino":       true,
	"acuerdo":       true,
	"ableton":       true,
	"marshawn":      true,
	"metroid":       true,
	"starcraft":     true,
	"sarkeesian":    true,
	"fatebringer":   true,
	"frickin":       true,
	"minecraft":     true,
	"shoulda":       true,
	"tristana":      true,
	"pcpartpicker":  true,
	"linkedin":      true,
	"myfitnesspal":  true,
	"goddamned":     true,
	"dipshit":       true,
	"katarina":      true,
	"programa":      true,
	"mayweather":    true,
	"superstars":    true, // FP
	"redditor":      true,
	"muhammad":      true,
	"redditors":     true,
	"nintendo":      true,
	"nosotros":      true,
	"goddamnit":     true,
	"tolkien":       true,
	"stephanie":     true,
	"warlock":       true,
	"circlejerking": true,
	"ciudadanos":    true, // Spanish
	"judgment":      true, // misspelling
	"importante":    true, // Spanish
	"iphones":       true, // brand
	"soundcloud":    true, // company
	"roberts":       true, // name
	"nicholas":      true,
	"doritos":       true, // brand
	"radiohead":     true, // name
	"bullshitting":  true,
	"katowice":      true,
	"chromebook":    true,
	"abrahamic":     true,
	"greninja":      true,
	"reddits":       true,
	"mitchell":      true,
	"templatejs":    true,
	"expandjs":      true, // javascript something?
	"caitlyn":       true,
	"sigelei":       true,
	"andreas":       true,
	"tarantino":     true,
	"amiibos":       true,
	"patreon":       true, // misspelling of patron
	"herrera":       true,
	"skeltal":       true, // misspelling of skeletal?
	"scientology":   true,
	"spongebob":     true,
	"swordbearer":   true,
	"coutinho":      true,
	"porsche":       true,
	"michaels":      true,
	"junglers":      true,
	"manziel":       true,
	"griffin":       true,
	"thunderlord":   true,
	"witcher":       true,
	"batshit":       true,
	"clusterfuck":   true,
	"atlantis":      true,
	"netflix":       true,
	"greatsword":    true,
	"karambit":      true,
	"sephora":       true,
	"freakin":       true,
	"jungler":       true,
	"hawkmoon":      true,
	"benjamin":      true,
	"westeros":      true,
	"hermione":      true,
	"runescape":     true,
	"bioshock":      true,
	"rosalina":      true,
	"superheroes":   true, // FP
	"chipotle":      true,
	"dragonball":    true,
	"gilbert":       true,
	"johnson":       true,
	"christopher":   true,
	"mcdonalds":     true,
	"santorin":      true,
	"pikachu":       true,
	"okcupid":       true,
	"shitpost":      true,
	"jessica":       true,
	"chromecast":    true,
	"francis":       true,
	"boohoooooo":    true,
	"william":       true,
	"charmander":    true,
	"douchebags":    true,
	"overheat":      true,
	"presser":       true,
	"trudeau":       true,
	"zelnite":       true,
	"cualquier":     true,
	"ventura":       true,
	"problema":      true,
	"gracias":       true,
	"outside":       true,
	"hammers":       true,
	"seinfeld":      true,
	"hecarim":       true,
	"barbara":       true,
	"mariota":       true,
	"hornets":       true, // FP
	"ferguson":      true,
	"machina":       true,
	"charlie":       true,
	"bernard":       true,
	"alienware":     true,
	"daenerys":      true,
	"hahahah":       true,
	"sejuani":       true,
	"leblanc":       true,
	"targaryen":     true,
	"charizard":     true,
	"bundesliga":    true,
	"mohammed":      true,
	"daniels":       true,
	"walmart":       true,
	"chomsky":       true,
	"netanyahu":     true,
	"voldemort":     true,
	"rodriguez":     true,
	"gjallarhorn":   true,
	"bayonetta":     true,
	"morgana":       true,
	"vladimir":      true,
	"motherfucker":  true,
	"motherfucking": true,
	"megaman":       true,
	"rngesus":       true,
	"collins":       true,
	"motherfuckers": true,
	"hahahahaha":    true,
	"spaces":        true,
	"murphy":        true,
	"phones":        true,
	"radios":        true, // FP
	"partido":       true,
	"kingston":      true,
	"kassadin":      true,
	"gregory":       true,
	"shitlords":     true,
	"fittit":        true,
	"lucario":       true,
	"floats":        true,
	"masses":        true,
	"grinds":        true,
	"christie":      true,
	"stephen":       true, // name
	"stomps":        true, // lots of FP
	"graphs":        true, // lots of FP
	"pacquiao":      true,
	"baretta":       true,
	"aguero":        true,
	"hahahaha":      true,
	"kendrick":      true, // name, multiple spellings
	"kalista":       true, // name
	"muchos":        true, // spanish?
	"aldrig":        true, // ??
	"mourinho":      true,
	"dennis":        true,
	"monedero":      true,
	"magicka":       true,
	"golems":        true,
	"terraria":      true,
	"closes":        true, // lots of FP
}

var badTypo = map[string]bool{
	"carnagie":       true, //name
	"granda":         true,
	"sourceid":       true,
	"messageid":      true,
	"falsey":         true, // technical word for "false value type"
	"progresse":      true, // progesses, progressives
	"amature":        true, // gets corrected to "armature" or maybe "a mature"
	"bogons":         true, // technical word
	"bogon":          true, // technical word
	"expandos":       true, // "technical" word, for something that expands (see jQuery)
	"accessors":      true, // technical word
	"accessor":       true, // technical word
	"sithlord":       true,
	"administratie":  true, // Dutch spelling
	"killionaires":   true,
	"villionaires":   true,
	"zillionaires":   true,
	"codifications":  true,
	"administracion": true,
	"potentiella":    true,
	"pistos":         true,
	"corpe":          true,
	"cleaner":        true, // real word
	"thirty":         true, // real word
	"reactjs":        true, // name of product
	"chroo":          true, // not clear what this is a typo of
	"cleane":         true, // not clear if "cleaner" or "cleanser", causes FP
	"matchers":       true, // real word
	"mongos":         true, // mongo database server
	"mongod":         true, // mongo database daemon
	"expresssion":    true, // real word
	"parens":         true, // common for parenthesis
	"thru":           true, // informal, style
	"warpig":         true,
	"governmnet":     true, // mapping this way for some reason, governmnet->governments
	"intereating":    true,
	"interdating":    true,
}

// LoadWordList loads in a list of known-good words
func LoadWordList(fname string) (map[string]bool, error) {
	fi, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	out := make(map[string]bool)
	intro := true
	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		line := scanner.Text()
		if intro {
			if line == "---" {
				intro = false
			}
			continue
		}
		out[strings.ToLower(line)] = true
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// LoadCSV loads a file in csv format of "word, typo, ..."
func LoadCSV(fname string, knownGood map[string]bool) (map[string][]string, error) {

	// map of word to mutiple typos
	dict := make(map[string][]string, 100000)

	// map of typo to word
	typos := make(map[string]string, 100000)

	fi, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	fizip, err := gzip.NewReader(fi)
	if err != nil {
		return nil, err
	}
	defer fizip.Close()
	scanner := bufio.NewScanner(fizip)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			return nil, fmt.Errorf("Got bad line: %q", line)
		}
		word, typo := parts[0], parts[1]
		// ignore words 5 chars or less ... too noisy
		if len(word) < 7 {
			continue
		}
		if badWord[word] {
			continue
		}
		if badTypo[typo] {
			continue
		}
		if knownGood[typo] {
			continue
		}
		// likely some javascript thing, e.g. reactjs
		if strings.HasSuffix(typo, "js") {
			continue
		}

		// for some reason "ism", "ist", "ish" make bad matches
		// libertaranism->libertarians
		if strings.HasSuffix(typo, "ism") && !strings.HasSuffix(word, "ism") {
			continue
		}
		if strings.HasSuffix(typo, "ish") && !strings.HasSuffix(word, "ish") {
			continue
		}
		if strings.HasSuffix(typo, "ist") && !strings.HasSuffix(word, "ist") {
			continue
		}

		// does one typo map to two correct words?
		if otherWord, ok := typos[typo]; ok {
			val1 := smetrics.JaroWinkler(word, typo, 0.7, 4)
			val2 := smetrics.JaroWinkler(otherWord, typo, 0.7, 4)
			if val1 == val2 {
				log.Printf("Typo %q has value %f for both %q and %q", typo, val1, otherWord, word)
				// remove so we are consistent (order is random)
				// and to reduce false positives.
				list := []string{}
				for _, val := range dict[otherWord] {
					if val != typo {
						list = append(list, val)
					}
				}
				dict[otherWord] = list
				continue
			}
			if val1 < val2 {
				// log.Printf("Typo %q keeping %q over %q", typo, otherWord, word)
				continue
			}
			//log.Printf("Typo %q picking %q over %q", typo, word, otherWord)
		}

		typos[typo] = word
		dict[word] = append(dict[word], typo)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// remove items with only 1 or 2 corrections
	for k, v := range dict {
		if len(v) < 3 {
			delete(dict, k)
		}
	}
	return dict, nil
}

func main() {
	knownGood, err := LoadWordList("dict.txt")
	if err != nil {
		log.Fatalf("Unable to load word list: %s", err)
	}
	log.Printf("Loaded %d known-good words", len(knownGood))

	dict, err := LoadCSV("RC_2015-total.csv.gz", knownGood)
	if err != nil {
		log.Fatalf("Unable to load csv: %s", err)
	}
	fmt.Printf(`package main

func dictReddit() map[string]string {
        dict := parseWikipediaFormat(additionsReddit)
        dict = expandCase(dict)
        return dict
}
`)

	lines := make([]string, 0, len(dict)*4)
	fmt.Printf("var additionsReddit = `\n")
	for goodWord, typoList := range dict {
		for _, typo := range typoList {
			lines = append(lines, fmt.Sprintf("%s->%s", typo, goodWord))
		}
	}
	sort.Strings(lines)
	for _, line := range lines {
		fmt.Println(line)
	}
	fmt.Printf("`\n")
}
