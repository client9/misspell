package main

import (
	"bufio"
	//"compress/gzip"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/xrash/smetrics"
)

var badWord = map[string]bool{
	"alexandria":    true,
	"alexander":     true,
	"donovan":       true,
	"kakashi":       true, // Japanese for "scarecrow"/
	"welbeck":       true,
	"remeber":       true, // Misspelling!!
	"mcmahon":       true,
	"ganondorf":     true,
	"jeffrey":       true,
	"forsett":       true,
	"shyvana":       true,
	"froakie":       true,
	"jenkins":       true,
	"walgreens":     true,
	"buzzfeed":      true,
	"kershaw":       true,
	"rebecca":       true,
	"completly":     true, // Misspelling!!
	"gandalf":       true,
	"apeshit":       true,
	"cazorla":       true,
	"maximus":       true,
	"reginald":      true,
	"andrews":       true,
	"abdullah":      true,
	"ableton":       true,
	"abrahamic":     true,
	"accutane":      true, //brand (drug)
	"acuerdo":       true,
	"aguero":        true,
	"aldrig":        true, // ??
	"alguien":       true, // unknown
	"alienware":     true,
	"alistar":       true, // name
	"ambrose":       true, // name
	"amendola":      true, // unknown
	"amiibos":       true,
	"amirite":       true, // slang
	"anderson":      true, // name
	"andreas":       true,
	"antonio":       true,
	"aquaman":       true,
	"arduino":       true,
	"atlantis":      true,
	"atletico":      true, // ?
	"baldwin":       true, // name
	"balotelli":     true,
	"baratheon":     true, // ??
	"barbara":       true,
	"baretta":       true,
	"batshit":       true,
	"bayonetta":     true,
	"beckham":       true,
	"beleive":       true, // FP
	"benjamin":      true,
	"bernard":       true,
	"bioshock":      true,
	"bitchin":       true,
	"bjergsen":      true,
	"bledsoe":       true, // name
	"blitzcrank":    true, // ??
	"boohoooooo":    true,
	"bortles":       true, // ??
	"bradley":       true, // name
	"bradshaw":      true,
	"bridgewater":   true, // brand
	"brienne":       true,
	"britney":       true,
	"bronies":       true,
	"bulbasaur":     true,
	"bullshitting":  true,
	"bundesliga":    true,
	"caitlyn":       true,
	"cartoony":      true, // FP
	"castlevania":   true,
	"catelyn":       true,
	"celestia":      true,
	"charizard":     true,
	"charles":       true,
	"charleston":    true, //place name with possible variations
	"charlie":       true,
	"charlotte":     true, // name
	"charmander":    true,
	"chipotle":      true,
	"chomsky":       true,
	"christie":      true,
	"christina":     true, // name
	"christine":     true, // name
	"christopher":   true,
	"chromebook":    true,
	"chromecast":    true,
	"chronos":       true,
	"cinderhulk":    true,
	"circlejerk":    true,
	"circlejerking": true,
	"circlejerks":   true,
	"ciudadanos":    true, // Spanish
	"closes":        true, // lots of FP
	"clusterfuck":   true,
	"collins":       true,
	"connery":       true,
	"constantine":   true,
	"copypasta":     true,
	"courtney":      true,
	"craigslist":    true,
	"crawford":      true,
	"chrysler":      true, // brand
	"cortana":       true,
	"coutinho":      true,
	"crossfit":      true, // brand
	"cualquier":     true,
	"dawkins":       true,
	"dratini":       true,
	"dempsey":       true,
	"desserts":      true, // FP
	"deserts":       true, // FP
	"daenerys":      true,
	"dailymotion":   true, // brand
	"daniels":       true,
	"darkrai":       true,
	"deandre":       true,
	"dennis":        true,
	"dipshit":       true,
	"doritos":       true, // brand
	"douchebags":    true,
	"dragonball":    true,
	"dragonborn":    true,
	"dragonite":     true,
	"einstein":      true,
	"evangelion":    true, // brand, name
	"expandjs":      true, // javascript something?
	"edwards":       true, // name
	"distros":       true, // too small,
	"deviantart":    true, // brand
	"dignitas":      true,
	"dogecoin":      true,
	"disneyland":    true, // brand
	"fatebringer":   true,
	"ferguson":      true,
	"fernando":      true, // name
	"fittit":        true,
	"fitzgerald":    true, //name
	"floats":        true,
	"francis":       true,
	"freakin":       true,
	"frickin":       true,
	"friendzone":    true, // causual word
	"fuckton":       true,
	"gabriel":       true, //name
	"galactica":     true,
	"gilbert":       true,
	"gjallarhorn":   true,
	"goddamned":     true,
	"goddamnit":     true,
	"godzilla":      true,
	"golems":        true,
	"googled":       true,
	"gracias":       true,
	"graphs":        true, // lots of FP
	"greatsword":    true,
	"gregory":       true,
	"greninja":      true,
	"griffin":       true,
	"grinds":        true,
	"hammers":       true,
	"happend":       true,
	"harbaugh":      true,
	"hawkmoon":      true,
	"hecarim":       true,
	"hermione":      true,
	"herrera":       true,
	"historia":      true, //??
	"hornets":       true, // FP
	"horseshit":     true,
	"illidan":       true,
	"importante":    true, // Spanish
	"iniesta":       true, // ??
	"iphones":       true, // brand
	"jessica":       true,
	"jigglypuff":    true,
	"johnson":       true,
	"judgment":      true, // misspelling
	"jungler":       true,
	"junglers":      true,
	"kalista":       true, // name
	"kanthal":       true,
	"karambit":      true,
	"kassadin":      true,
	"katarina":      true,
	"katowice":      true,
	"kendrick":      true, // name, multiple spellings
	"killionaires":  true,
	"kingston":      true,
	"krugman":       true, // name
	"lamborghini":   true, // name, brand
	"leblanc":       true,
	"lebowski":      true, // name
	"leonardo":      true, // name
	"lightsaber":    true,
	"linkedin":      true,
	"lucario":       true,
	"machina":       true,
	"magicka":       true,
	"magneto":       true,
	"manziel":       true,
	"mariota":       true,
	"marshawn":      true,
	"martinez":      true,
	"masses":        true,
	"mayweather":    true,
	"mcdonalds":     true,
	"megaman":       true,
	"metroid":       true,
	"michael":       true,
	"michaels":      true,
	"minecraft":     true,
	"mitchell":      true,
	"mohammed":      true,
	"monedero":      true,
	"monsanto":      true,
	"morgana":       true,
	"motherfucker":  true,
	"motherfuckers": true,
	"motherfucking": true,
	"mourinho":      true,
	"muchos":        true, // spanish?
	"muhammad":      true,
	"murican":       true, // slang
	"murphy":        true,
	"myfitnesspal":  true,
	"mythbusters":   true, // brand
	"natalie":       true, // name
	"netanyahu":     true,
	"netflix":       true,
	"neville":       true, // name
	"nicholas":      true,
	"nickelback":    true, // brand
	"nicolas":       true, // name
	"nintendo":      true,
	"nosotros":      true,
	"okcupid":       true,
	"outside":       true,
	"overheat":      true,
	"pacquiao":      true,
	"partido":       true,
	"patreon":       true, // misspelling of patron
	"pcpartpicker":  true,
	"phones":        true,
	"pikachu":       true,
	"planetside":    true, // brand
	"porsche":       true,
	"presser":       true,
	"problema":      true,
	"programa":      true,
	"radiohead":     true, // name
	"radios":        true, // FP
	"redditor":      true,
	"redditors":     true,
	"reddits":       true,
	"retards":       true,
	"rngesus":       true,
	"roberts":       true, // name
	"rodriguez":     true,
	"ronaldo":       true,
	"rosalina":      true,
	"runescape":     true,
	"samsung":       true,
	"sanchez":       true,
	"santorin":      true,
	"sarkeesian":    true,
	"scientology":   true,
	"seinfeld":      true,
	"sejuani":       true,
	"sephora":       true,
	"shithead":      true,
	"shitlord":      true,
	"shitlords":     true,
	"shitpost":      true,
	"shitposting":   true,
	"shitposts":     true,
	"shitter":       true,
	"shittier":      true,
	"shittiest":     true,
	"shoulda":       true,
	"siempre":       true,
	"sigelei":       true,
	"skeltal":       true, // misspelling of skeletal?
	"snapchat":      true,
	"somthing":      true, // FP
	"soundcloud":    true, // company
	"spaces":        true,
	"spongebob":     true,
	"starcraft":     true,
	"stephanie":     true,
	"stephen":       true, // name
	"stomps":        true, // lots of FP
	"superheroes":   true, // FP
	"superstars":    true, // FP
	"sverige":       true,
	"swordbearer":   true,
	"tarantino":     true,
	"targaryen":     true,
	"templatejs":    true,
	"terran":        true,
	"terrans":       true, // (in science fiction) an inhabitant of the planet Earth.
	"terraria":      true,
	"thunderlord":   true,
	"tolkien":       true,
	"tristana":      true,
	"trudeau":       true,
	"ventura":       true,
	"veronica":      true, // name
	"villionaires":  true,
	"vladimir":      true,
	"voldemort":     true,
	"vonnegut":      true, // name
	"douglas":       true, // name
	"walmart":       true,
	"warlock":       true,
	"dreamcast":     true,
	"dreamhack":     true,
	"dumbasses":     true,
	"dumbledore":    true,
	"fabregas":      true,
	"fucktard":      true,
	"fleshlight":    true,
	"futurama":      true,
	"gangsta":       true,
	"gambino":       true,
	"genesect":      true,
	"ghostbusters":  true,
	"grimoire":      true,
	"gobierno":      true,
	"gonzalez":      true,
	"heisenberg":    true,
	"hendricks":     true,
	"hitchens":      true,
	"hogwarts":      true,
	"honedge":       true,
	"hyundai":       true,
	"instagram":     true,
	"jirachi":       true,
	"warpig":        true,
	"westeros":      true,
	"kaepernick":    true,
	"kardashian":    true,
	"karthus":       true,
	"katrina":       true,
	"whatcha":       true, // slang
	"limbaugh":      true,
	"whatsapp":      true,
	"william":       true,
	"witcher":       true,
	"wrestlemania":  true,
	"youtuber":      true,
	"zelnite":       true,
	"lololol":       true,
	"macklemore":    true,
	"magikarp":      true,
	"malphite":      true,
	"materia":       true,
	"mcdonald":      true,
	"melissa":       true,
	"metacritic":    true,
	"micheal":       true,
	"microsoft":     true,
	"nietzsche":     true,
	"nordstrom":     true,
	"orianna":       true,
	"optimus":       true,
	"phillips":      true,
	"playstation":   true,
	"archers":       true,
	"rhaegar":       true,
	"reccomend":     true,
	"rediculous":    true,
	"reddiquette":   true,
	"richardson":    true,
	"richards":      true,
	"roosevelt":     true,
	"scarlet":       true,
	"shadowbanned":  true,
	"shadowrun":     true,
	"schneider":     true,
	"shitload":      true,
	"skrillex":      true,
	"sturridge":     true,
	"scarlett":      true,
	"sakurai":       true,
	"arsehole":      true,
	"sylvanas":      true,
	"thalmor":       true,
	"tmobile":       true,
	"triforce":      true,
	"tryndamere":    true,
	"ubisoft":       true,
	"victoria":      true,
	"virginia":      true,
	"warlocks":      true,
	"wannabe":       true,
	"westboro":      true,
	"wikipedia":     true,
	"winterfell":    true,
	"wolfenstein":   true,
	"xpecial":       true,
	"yogscast":      true,
}

var badTypo = map[string]bool{
	"superpowder":      true,
	"proletara":        true,
	"pediction":        true,
	"motorola":         true,
	"moranian":         true,
	"misqualified":     true,
	"mingleplayer":     true,
	"melbournite":      true,
	"mathological":     true,
	"carnagie":         true, //name
	"granda":           true,
	"sourceid":         true,
	"messageid":        true,
	"falsey":           true, // technical word for "false value type"
	"progresse":        true, // progesses, progressives
	"amature":          true, // gets corrected to "armature" or maybe "a mature"
	"bogons":           true, // technical word
	"bogon":            true, // technical word
	"generalizaciones": true, //spanish
	"expandos":         true, // "technical" word, for something that expands (see jQuery)
	"accessors":        true, // technical word
	"accessor":         true, // technical word
	"sithlord":         true,
	"administratie":    true, // Dutch spelling
	"killionaires":     true,
	"villionaires":     true,
	"zillionaires":     true,
	"codifications":    true,
	"administracion":   true,
	"potentiella":      true,
	"pistos":           true,
	"corpe":            true,
	"cleaner":          true, // real word
	"thirty":           true, // real word
	"reactjs":          true, // name of product
	"chroo":            true, // not clear what this is a typo of
	"cleane":           true, // not clear if "cleaner" or "cleanser", causes FP
	"matchers":         true, // real word
	"mongos":           true, // mongo database server
	"mongod":           true, // mongo database daemon
	"expresssion":      true, // real word
	"parens":           true, // common for parenthesis
	"thru":             true, // informal, style
	"warpig":           true,
	"governmnet":       true, // mapping this way for some reason, governmnet->governments
	"intereating":      true,
	"interdating":      true,
	"dogspeed":         true,
	"laventine":        true,
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
func LoadCSV(fname string, knownGood map[string]bool, minTypo int) (map[string][]string, error) {

	// map of word to mutiple typos
	dict := make(map[string][]string, 100000)

	// map of typo to word
	typos := make(map[string]string, 100000)

	fi, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	/*
		fizip, err := gzip.NewReader(fi)
		if err != nil {
			return nil, err
		}
		defer fizip.Close()
	*/
	scanner := bufio.NewScanner(fi)
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

		// not sure how 'hahahahaha' keeps sneaking in
		if strings.Contains(word, "haha") {
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

	// remove items with too few variations
	for k, v := range dict {
		if len(v) < minTypo {
			delete(dict, k)
		}
	}
	return dict, nil
}

func main() {
	dictfile := flag.String("d", "dict.txt", "aspell wordlist")
	//outfile := flag.String("o", "RC-score.csv", "outfile")
	infile := flag.String("i", "RC-score.csv", "infile")
	minTypo := flag.Int("mintypo", 2, "require at least this many typos to include")
	flag.Parse()
	knownGood, err := LoadWordList(*dictfile)
	if err != nil {
		log.Fatalf("Unable to load word list: %s", err)
	}
	log.Printf("Loaded %d known-good words", len(knownGood))

	dict, err := LoadCSV(*infile, knownGood, *minTypo)
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
