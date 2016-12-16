package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
)

// sort -t ',' -k 3 fixtures/sample.tsv
func main() {
	r := csv.NewReader(bufio.NewReader(os.Stdin))
	r.Comma = ','

	var batch []string

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if len(batch) == 0 {
			batch = append(batch, record[2])
			continue
		}
		if record[2] != batch[len(batch)-1] {
			log.Println("processing batch of %d items", len(batch))
			log.Println(batch)
			batch = nil
		}
		batch = append(batch, record[2])
	}

	log.Println("rest", batch)
}
