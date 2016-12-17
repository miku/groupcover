package main

import (
	"log"
	"os"

	"github.com/miku/groupcover"
)

func main() {
	// Preference per key.
	// preferences := groupcover.PreferenceMap{
	// 	"K1": groupcover.LexChoice,
	// 	"K2": groupcover.LexChoice,
	// 	"K3": groupcover.LexChoice,
	// }

	preferences := groupcover.PreferenceMap{
		"DE-105":    groupcover.LexChoice,
		"DE-14":     groupcover.LexChoice,
		"DE-15":     groupcover.LexChoice,
		"DE-15-FID": groupcover.LexChoice,
		"DE-1972":   groupcover.LexChoice,
		"DE-540":    groupcover.LexChoice,
		"DE-Bn3":    groupcover.LexChoice,
		"DE-Brt1":   groupcover.LexChoice,
		"DE-Ch1":    groupcover.LexChoice,
		"DE-D13":    groupcover.LexChoice,
		"DE-D161":   groupcover.LexChoice,
		"DE-Gla1":   groupcover.LexChoice,
		"DE-Kn38":   groupcover.LexChoice,
		"DE-L152":   groupcover.LexChoice,
		"DE-L242":   groupcover.LexChoice,
		"DE-Zi4":    groupcover.LexChoice,
		"DE-Zwi2":   groupcover.LexChoice,
	}

	// Use the third column as grouping criteria.
	criteria := groupcover.Column(2)
	// A simple rewriter, that considers per-key preferences.
	rewriter := groupcover.SimpleRewriter(preferences)

	// Read from stdin, write to stdout, use third column as grouping criteria
	// and rewriter as rewriter.
	if err := groupcover.GroupRewrite(os.Stdin, os.Stdout, criteria, rewriter); err != nil {
		log.Fatal(err)
	}
}
