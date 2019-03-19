// The groupindex tool can be applied to an intermediate schema file. It will
// query the index for potential duplicates and will print out changes, just as
// groupcover. With groupindex it should be possible add documents to an index,
// without the need to run groupcover on the complete dataset.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	server = flag.String("server", "", "SOLR server, hostport plus core, e.g. http://1.2.3.4:8081/solr/biblio")
)

// WithLabels for intermediate schema fragment.
type WithLabels struct {
	Labels []string `json:"x.labels"`
	DOI    string   `json:"doi"`
}

func main() {
	flag.Parse()

	// Input documents.
	var br = bufio.NewReader(os.Stdin)

	if flag.NArg() > 0 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		br = bufio.NewReader(f)
	}

	// For each document in input, extract the DOI, ask index for records with
	// the same doi.
	for {
		b, err := br.ReadBytes('\n')
		if err != nil {
			log.Fatal(err)
		}
		if err == io.EOF {
			break
		}
		if len(b) == 0 {
			continue
		}
		var doc WithLabels
		if err := json.Unmarshal(b, &doc); err != nil {
			log.Fatal(err)
		}
		if strings.TrimSpace(doc.DOI) == "" {
			continue
		}
		fmt.Println(doc.DOI)
	}
}
