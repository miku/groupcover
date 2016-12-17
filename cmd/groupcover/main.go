package main

import (
	"log"
	"os"

	"github.com/miku/groupcover"
)

func main() {
	// Preference per key.
	preferences := groupcover.PreferenceMap{
		"K1": groupcover.LexChoice,
		"K2": groupcover.LexChoice,
		"K3": groupcover.LexChoice,
	}
	// Use the third column as grouping criteria.
	thirdColumn := groupcover.Column(2)
	// A simple rewriter, that considers per-key preferences.
	rewriter := groupcover.SimpleRewriter(preferences)

	// Read from stdin, write to stdout, use third column as grouping criteria
	// and rewriter as rewriter.
	if err := groupcover.GroupRewrite(os.Stdin, os.Stdout, thirdColumn, rewriter); err != nil {
		log.Fatal(err)
	}
}
