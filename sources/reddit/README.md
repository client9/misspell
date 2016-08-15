# Reddit comment analysis


Notes:

* bzip2 is horrible.  Uncompression is painfully slow.  Files are downloaded
  and recompressed using gzip.  All files takes about 150GB
* Uses poor mans Map-Reduce.  Each file is converted to a word/count file
  and at the end they are merged.
* the scoring is really slow and O(n^2).  This could be done in parallel
* Tried to eliminate the Huge numbers of brand names, tv show characters, fantasy stuff.
* Some misspelling are so common they make the 90% percentile!
