package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

// Attr extract an attribute from a CSV record.
type AttrFunc func([]string) (string, error)

// RewriterFunc rewrites a list of lines.
type RewriterFunc func([][]string) ([][]string, error)

// ChoiceFunc presented with a list of choices, chooses one.
type ChoiceFunc func([]string) string

// PreferenceMap maps a string to a ChoiceFunc.
type PreferenceMap map[string]ChoiceFunc

// GroupLines takes a reader (over CSV) and a attribtue value extractor and
// groups lines that have the same attribute. As with uniq, the lines must be
// sorted by the attribute.
func GroupLines(r io.Reader, w io.Writer, attrFunc AttrFunc, rewriterFunc RewriterFunc) error {
	cr := csv.NewReader(r)
	cw := csv.NewWriter(w)

	var prev string
	var group [][]string

	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		attr, err := attrFunc(record)
		if err != nil {
			return err
		}
		if attr != prev {
			regroup, err := rewriterFunc(group)
			if err != nil {
				return err
			}
			if err := cw.WriteAll(regroup); err != nil {
				return err
			}
			group = nil
		}
		group = append(group, record)
		prev = attr
	}

	// final group
	regroup, err := rewriterFunc(group)
	if err != nil {
		return err
	}
	return cw.WriteAll(regroup)
}

func main() {
	// extract the third column
	thirdColumn := func(record []string) (string, error) {
		if len(record) != 4 {
			return "", fmt.Errorf("invalid column count: got %d, want 4", len(record))
		}
		return record[2], nil
	}

	// lexChoice chooses the key with the highest lexicographic order.
	lexChoice := func(s []string) string {
		if len(s) == 0 {
			return ""
		}
		sort.Strings(s)
		return s[len(s)-1]
	}

	preferences := PreferenceMap{
		"K1": lexChoice,
		"K2": lexChoice,
		"K3": lexChoice,
	}

	sampleRewriter := func(s [][]string) ([][]string, error) {
		if len(s) < 2 {
			return s, nil
		}

		// 1. For each ISIL (keys) get the available groups.
		keyGroups := make(map[string][]string)
		for _, record := range s {
			group := record[1]
			for _, key := range strings.Split(record[3], ",") {
				keyGroups[key] = append(keyGroups[key], group)
			}
		}

		// 2. For each ISIL (keys) get the preferred group.
		preferredGroup := make(map[string]string)
		for key, groups := range keyGroups {
			f, ok := preferences[key]
			if !ok {
				log.Fatalf("no preference defined for %s", key)
			}
			preferredGroup[key] = f(groups)
		}

		// 3. For each lines, check the group and list the ISIL (keys) for which this group is the preferred.
		for _, record := range s {
			var updated []string
			group := record[1]
			for _, key := range strings.Split(record[3], ",") {
				if preferredGroup[key] == group {
					updated = append(updated, key)
				}
			}
			record[3] = strings.Join(updated, ",")
		}

		return s, nil
	}

	if err := GroupLines(os.Stdin, os.Stdout, thirdColumn, sampleRewriter); err != nil {
		log.Fatal(err)
	}
}
