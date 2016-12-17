package main

import (
	"log"
	"os"

	"github.com/miku/groupcover"
)

func main() {
	// Preference per key.
	preferences := PreferenceMap{
		"K1": lexChoice,
		"K2": lexChoice,
		"K3": lexChoice,
	}
	// Use the third column as grouping criteria.
	thirdColumn := groupcover.Column(2)
	// A simple rewriter, that has considers per key preferences.
	rewriter := SimpleRewriter(preferences)

	if err := GroupLines(os.Stdin, os.Stdout, thirdColumn, rewriter); err != nil {
		log.Fatal(err)
	}
}
