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
	"net/http"
	"os"
	"strings"
)

var (
	server = flag.String("server", "", "SOLR server, hostport plus core, e.g. http://1.2.3.4:8081/solr/biblio")
	prefs  = flag.String("prefs", "85 55 89 60 50 105 34 101 53 49 28 48 121", "most preferred first")
)

// WithLabels for intermediate schema fragment.
type WithLabels struct {
	SourceID string   `json:"finc.source_id"`
	Labels   []string `json:"x.labels"`
	DOI      string   `json:"doi"`
}

// SelectResponse with reduced fields.
type SelectResponse struct {
	Response struct {
		Docs []struct {
			Institution []string `json:"institution"`
			SourceID    string   `json:"source_id"`
		} `json:"docs"`
		NumFound int64 `json:"numFound"`
		Start    int64 `json:"start"`
	} `json:"response"`
	ResponseHeader struct {
		Params struct {
			Q  string `json:"q"`
			Wt string `json:"wt"`
		} `json:"params"`
		QTime  int64
		Status int64 `json:"status"`
	} `json:"responseHeader"`
}

// fetchDocuments fetches documents for a given identifier.
func fetchDocuments(doi string) (*SelectResponse, error) {
	// The identifier has no single field, we assume the DOI is unique.
	// Prefixes might be a problem.
	link := fmt.Sprintf(`http://172.18.113.7:8085/solr/biblio/select?wt=json&q="%s"`, doi)
	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var sr SelectResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, err
	}
	return &sr, nil
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
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
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
		sr, err := fetchDocuments(doc.DOI)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, indexed := range sr.Response.Docs {
			fmt.Printf("[%s] %s, [%s] %s\n",
				doc.SourceID, doc.Labels, indexed.SourceID, indexed.Institution)
		}
		// fmt.Printf("%s => %d\n", doc.DOI, sr.Response.NumFound)
	}
}
