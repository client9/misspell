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
		dict["Dont "] = "Don't "
	*/
	return dict
}

// arent
var additions = `
Accomodate->Accommodate
COMMERICIAL->COMMERCIAL
Constructur->Constructor
Depdending->Depending
Disclamer->Disclaimer
Dont->Don't
Elimintates->Eliminates
Embeded->Embedded
Externalise->Externalize
Fowrards->Forwards
Futher->Further
IMPOSTER->IMPOSTOR
Implementor->Implementer
Instalation->Installation
LICENCE->LICENSE
Licence->License
Numerious->Numerous
Perfomance->Performance
Recieve->Receive
Runing->Running
Specifcation->Specification
Wheter->Whether
acknowledgement->acknowledgment
aforementioend->aforementioned
aggregate->aggregate
america->America
annonymouse->anonymous
apologise->apologize
approstraphe->apostrophe
apporach->approach
aribtrary->arbitrary
artefact->artifact
asychronous->asynchronous
australia->Australia
avaiable->available
cahched->cached
calback->callback
canadian->Canadian
capitalised->capitalized
careflly->carefully
chinese->Chinese
commmand->command
compatibilty->compatibility
comptability->compatibility
conatins->contains
conditon->condition
configuraiton->configuration
conflic ->conflict
consitency->consistency
contructed->constructed
contructor->constructor
convertor->converter
couldn ->couldn't
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
organisation->organization
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
sanitise->sanitize
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
