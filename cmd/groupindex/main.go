// The groupindex tool can be applied to an intermediate schema or solr file,
// that is about to be indexed. It will query the index for potential
// duplicates and will adjust the document labels (x.labels, institution)
// accordingly. With groupindex it should be possible add documents to an
// index, without the need to run groupcover on the complete dataset.
package main

import "log"

func main() {
	log.Println("groupindex")
}
