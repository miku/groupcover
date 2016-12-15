package main

import (
	"fmt"
	"log"
	"os"

	"github.com/miku/groupcover"
)

func main() {
	table, err := groupcover.FromReader(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(table)
}
