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

	pm := make(groupcover.PreferenceMap)
	pm["K1"] = &groupcover.Preference{"G1", "G2"}
	pm["K2"] = &groupcover.Preference{"G2", "G1"}
	pm["K3"] = &groupcover.Preference{}

	fmt.Println(table)
	cleaner := groupcover.SampleCleaner{Preferences: pm}
	entries := cleaner.Clean(table.Entries)
	fmt.Println(&groupcover.Table{Entries: entries})
}
