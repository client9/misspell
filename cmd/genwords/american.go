package main

func dictAmerican() map[string]string {
	dict := parseWikipediaFormat(american)
	return dict
}

// http://www.oxforddictionaries.com/us/words/american-and-british-spelling-american
// http://www.oxforddictionaries.com/us/words/american-and-british-terms-american
var american = `
centre->center
fibre->fiber
litre->liter
metre->meter
theatre->theater
colour->color
flavour->flavor
humour->humor
labour->labor
neighbour->neighbor
apologise->apologize
organise->organize
recognise->recognize
analyse->analyze
breathalyse->breathalyze
paralyse->paralyze
travelled->traveled
travelling->traveling
traveller->traveler
fuelled->fueled
fuelling->fueling
leukaemia->leukemia
manoeuvre->maneuver
oestrogen->estrogen
paediatric->pediatric
defence->defense
license->license
offence->offense
pretence->pretense
analogue->analog
catalogue->catalog
catalogues->catalogs
catalogued->cataloged
uncatalogued->uncataloged
miscatalogued->miscataloged
dialogue->dialog
aeroplane->airplane
aluminium->aluminum
`
