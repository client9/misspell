package main

func dictAdditions() map[string]string {
	dict := parseWikipediaFormat(additions)
	/*
		// Additions... need to move these somewhere else
		addOrPanic(dict, "simultanoues", "simultaneous")
		addOrPanic(dict, "configuraiton", "configuration")
		addOrPanic(dict, "Didnt", "Didn't")
		addOrPanic(dict, "i'm", "I'm")

		// Additions
		// note some variable names are "cantFoo"
		dict[" cant "] = " can't "
		dict[" dont "] = " don't "
	*/
	dict["Dont "] = "Don't "
	return dict
}

// arent
var additions = `
retunred->returned
authenticor->authenticator
availabale->available
positve->positive
satifies->satisfies
capialized->capitalized
versoin->version
obvioulsy->obviously
fundemental->fundamental
crytopgraphic->cryptographic
appication->application
accending->ascending
consisent->consistent
percision->precision
determinsitic->deterministic
elasped->elapsed
udpated->updated
undescore->underscore
represenation->representation
registery->registry
redundent->redundant
puncutation->punctuation
genrates->generates
finallizes->finalizes
expoch->epoch
equivalant->equivalent
determinsitic->deterministic
normallized->normalized
elasped->elapsed
machiens->machines
demonstates->demonstrates
collumn->column
verical->vertical
refernece->reference
opartor->operator
elimiate->eliminate
coalese->coalesce
extenion->extension
affliated->affiliated
hesistate->hesitate
arrary->array
hunman->human
currate->curate
retuns->returns
interfce->interface
alrorythm->algorithm
credentaisl->credentials
closeing->closing
Constructur->Constructor
Depdending->Depending
Disclamer->Disclaimer
Elimintates->Eliminates
Externalise->Externalize
Fowrards->Forwards
IMPOSTER->IMPOSTOR
Implementor->Implementer
Instalation->Installation
Numerious->Numerous
Runing->Running
Specifcation->Specification
Wheter->Whether
acknowledgement->acknowledgment
aforementioend->aforementioned
annonymouse->anonymous
apologise->apologize
approstraphe->apostrophe
apporach->approach
aribtrary->arbitrary
artefact->artifact
asychronous->asynchronous
avaiable->available
cahched->cached
calback->callback
careflly->carefully
commmand->command
compatibilty->compatibility
comptability->compatibility
conatins->contains
conditon->condition
configuraiton->configuration
consitency->consistency
contructed->constructed
contructor->constructor
convertor->converter
customises->customizes
december->December
declareation->declaration
decomposeion->decomposition
deliviered->delivered
depedencies->dependencies
depedency->dependency
deperecation->deprecation
descendent->descendant
descriminant->discriminant
diffucult->difficult
documenation->documentation
dyamically->dynamically
embeded->embedded
everwhere->everywhere
exising->existing
explicitely->explicitly
explicity->explicitly
expliots->exploits
exprimental->experimental
extactly->exactly
functionlity->functionality
functtion->function
homogenous->homogeneous
idiosynchracies->idiosyncrasies
immidiate->immediate
implemention->implementation
implentation->implementation
implicitely->implicitly
implimenation->implementation
incldue->include
incorect->incorrect
incorectly->incorrectly
infeasible->infeasible
inferrence->inference
initialise->initialize
maximise->maximize
maximising->maximizing
milisecond->millisecond
mimimum->minimum
minimised->minimized
minimium->minimum
misinterpretting->misinterpreting
mississippi->Mississippi
momment->moment
muliple->multiple
mulitple->multiple
nubmers->numbers
officiallly->officially
otherhand->other hand
optimisation->optimization
optimising->optimizing
optinally->optimally
ouput->output
outputed->outputted
pacakge->package
packge->package
paramter->parameter
paramters->parameters
paricular->particular
parition->partition
performaces->performances
permisson->permission
precedeed->preceded
precendence->precedence
programattically->programmatically
programmar->programmer
programms->programs
properites->properties
propeties->properties
protototype->prototype
publsih->publish
queueing->queuing
quuery->query
recognise->recognize
recognising->recognizing
requried->required
retrived->retrieved
ridiculus->ridiculous
sceptical->skeptical
seperator->separator
similarlly->similarly
simplfy->simplify
singals->signals
spanish->Spanish
specifcally->specifically
specifed->specified
specifiy->specify
straitforward->straightforward
subsequant->subsequent
successfuly->successfully
supportied->supported
supression->suppression
synchornously->synchronously
syncronously->synchronously
tutorual->tutorial
unintuive->unintuitive
writting->writing
`
