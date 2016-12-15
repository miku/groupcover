package main

import (
	"fmt"
	"log"
	"os"

	"github.com/miku/groupcover"
)

func main() {
	table, err := groupcover.TableFromReader(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(table)
}
