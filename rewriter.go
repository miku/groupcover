package groupcover

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strings"
)

// ChoiceFunc presented with a list of choices, chooses one.
type ChoiceFunc func([]string) string

// PreferenceMap maps a string to a ChoiceFunc.
type PreferenceMap map[string]ChoiceFunc

// Attr extracts an attribute from a CSV record.
type AttrFunc func([]string) (string, error)

// RewriterFunc rewrites a list of lines.
type RewriterFunc func([][]string) ([][]string, error)

// LexChoice chooses the key with the highest lexicographic order. These
// preference may come from external sources.
func LexChoice(s []string) string {
	if len(s) == 0 {
		return ""
	}
	sort.Strings(s)
	return s[len(s)-1]
}

// Column returns an AttrFunc, that extracts the given column value.
func Column(k int) AttrFunc {
	f := func(record []string) (string, error) {
		if k >= len(record) {
			return "", fmt.Errorf("invalid column: got %d, record has only %d", k, len(record))
		}
		return record[k], nil
	}
	return f
}

// GroupRewrite reads CSV records from reader, extracts an attribute with
// attrFunc, groups subsequent lines with the same attribute value and passes
// these groups to rewriterFunc. The rewritten lines are written as CSV to the
// given writer.
func GroupRewrite(r io.Reader, w io.Writer, attrFunc AttrFunc, rewriterFunc RewriterFunc) error {
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

// SimpleRewriter attempts key deduplication. SimpleRewriter takes a
// preference map and returns a rewriter.
func SimpleRewriter(preferences PreferenceMap) RewriterFunc {
	f := func(s [][]string) ([][]string, error) {
		// A single entry does not need any deduplication.
		if len(s) < 2 {
			return s, nil
		}

		// Basic sanity check.
		for _, record := range s {
			if len(record) != 4 {
				return nil, fmt.Errorf("expected 4 columns: ID, GROUP, ATTRIBUTE, KEYS")
			}
		}

		// For each key get the associated groups.
		groupsPerKey := make(map[string][]string)
		for _, record := range s {
			for _, key := range strings.Split(record[3], ",") {
				groupsPerKey[key] = append(groupsPerKey[key], record[1])
			}
		}

		// For each key determine the preferred group.
		preferred := make(map[string]string)
		for key, groups := range groupsPerKey {
			f, ok := preferences[key]
			if !ok {
				return nil, fmt.Errorf("no preference entry for %s", key)
			}
			preferred[key] = f(groups)
		}

		// For each record, check the group and list the ISIL (keys) for which
		// this group is the preferred.
		for _, record := range s {
			var updated []string
			group := record[1]
			for _, key := range strings.Split(record[3], ",") {
				if preferred[key] == group {
					updated = append(updated, key)
				}
			}
			record[3] = strings.Join(updated, ",")
		}

		return s, nil
	}
	return f
}
