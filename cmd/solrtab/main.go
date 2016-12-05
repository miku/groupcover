package main

import (
	"flag"
	"log"
	"net/url"
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("a full SOLR query URL is required")
	}

	u, err := url.Parse(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v", u.Query())
}
