# Wikimedia Pages Dump Analysis

This source takes a wikimedia page dump file, extracts words, and 
performs misspelling analysis.

The code is nearly identical to Reddit code and needs to be consolidated.

Notes:

* No attempt at formally parsing the wikipedia dump file.   we grab words blindly.
* Wikipedia errors are completely different than Reddit
* Words are mostly well-formed, very little if any slang or vocalized words (yooooouuu)
* Capalitization is mostly good.
* However, there are many quotes from Old English, that uses an "f" instead of a "s"
* Brand names with odd spellings can bias things, so the extractor only takes
  lower case words.
* Sample size is much smaller than Reddit (maybe 40x).  The scoring mechanism only
  looks at 80th percentile, since after than the words become more and more unusal.
  Interestingly, for both wikipedia and reddit, it takes about 12,000 words to hit
  the 90th percentile.
