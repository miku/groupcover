//  Copyright 2016 by Leipzig University Library, http://ub.uni-leipzig.de
//                    The Finc Authors, http://finc.info
//                    Martin Czygan, <martin.czygan@uni-leipzig.de>
//
// This file is part of some open source application.
//
// Some open source application is free software: you can redistribute
// it and/or modify it under the terms of the GNU General Public
// License as published by the Free Software Foundation, either
// version 3 of the License, or (at your option) any later version.
//
// Some open source application is distributed in the hope that it will
// be useful, but WITHOUT ANY WARRANTY; without even the implied warranty
// of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Foobar.  If not, see <http://www.gnu.org/licenses/>.
//
// @license GPL-3.0+ <http://spdx.org/licenses/GPL-3.0+>
//
package groupcover

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"reflect"
	"sort"
	"strings"
)

// ChoiceFunc presented with a list of choices, chooses one.
type ChoiceFunc func([]string) string

// PreferenceMap groups many choices by key.
type PreferenceMap map[string]ChoiceFunc

// AttrFunc extracts an attribute value from a CSV record. Example values
// could be a single column, part of a column or a value spanning multiple
// columns.
type AttrFunc func(record []string) (string, error)

// RewriterFunc rewrites a list of records.
type RewriterFunc func(records [][]string) ([][]string, error)

// LexChoice chooses the key with the highest lexicographic order.
func LexChoice(s []string) string {
	if len(s) == 0 {
		return ""
	}
	sort.Strings(s)
	return s[len(s)-1]
}

// Column returns an AttrFunc, that extracts the value at a given column,
// zero-based.
func Column(k int) AttrFunc {
	f := func(record []string) (string, error) {
		if k >= len(record) {
			return "", fmt.Errorf("invalid column: got %d, record has only %d", k, len(record))
		}
		return strings.TrimSpace(record[k]), nil
	}
	return f
}

// GroupRewrite reads CSV records from a given reader, extracts attribute
// values with attrFunc, groups subsequent records with the same attribute
// value and passes these groups to a rewriter. The altered records are
// written as CSV to the given writer.
func GroupRewrite(r io.Reader, w io.Writer, attrFunc AttrFunc, rewriterFunc RewriterFunc) error {
	cw := csv.NewWriter(w)
	cr := csv.NewReader(r)
	// If FieldsPerRecord is negative, no check is made and records may have a
	// variable number of fields.
	cr.FieldsPerRecord = -1

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
		value, err := attrFunc(record)
		if err != nil {
			return err
		}
		if value == "" {
			continue
		}
		if value != prev {
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
		prev = value
	}

	// Final group.
	regroup, err := rewriterFunc(group)
	if err != nil {
		return err
	}
	return cw.WriteAll(regroup)
}

// SimpleRewriter takes a preference map (which key is interested in which
// group) and returns a rewriter, which drops certain keys that are assigned
// to records from multiple groups with the same attribute value. This
// rewriter returns only differing records.
func SimpleRewriter(preferences PreferenceMap) RewriterFunc {
	f := func(records [][]string) ([][]string, error) {
		// A single entry does not need any deduplication.
		if len(records) < 2 {
			return nil, nil
		}

		// Only keep comparable records.
		var valid [][]string

		for _, record := range records {
			if len(record) < 4 {
				continue
			}
			valid = append(valid, record)
		}

		records = valid

		// For each key get the associated groups.
		groupsPerKey := make(map[string][]string)
		for _, record := range records {
			for _, key := range record[3:] {
				groupsPerKey[key] = append(groupsPerKey[key], record[1])
			}
		}

		// For each key determine the preferred group.
		preferred := make(map[string]string)
		for key, groups := range groupsPerKey {
			if _, ok := preferences[key]; !ok {
				preferences[key] = LexChoice
				log.Printf("no preference for %s, using lexicographic default", key)
			}
			f := preferences[key]
			preferred[key] = f(groups)
		}

		// Collect changed records here.
		var changedRecords [][]string

		// For each record, check the group and list the ISIL (keys) for which
		// this group is the preferred.
		for _, record := range records {
			var updated []string
			id, group, keys := record[0], record[1], record[3:]

			for _, key := range keys {
				if preferred[key] == group {
					updated = append(updated, key)
				}
			}

			// Keep only lines that changed.
			current := make([]string, len(keys))
			copy(current, keys)

			sort.Strings(current)
			sort.Strings(updated)

			if !reflect.DeepEqual(current, updated) {
				log.Printf("keys changed from %s to %s for %s", current, updated, id)
				// Assemble a new record.
				record := []string{record[0], record[1], record[2]}
				record = append(record, updated...)
				changedRecords = append(changedRecords, record)
			}
		}

		return changedRecords, nil
	}
	return f
}

// LastRow rewriter that only keeps the last row, similar to uniq(1).
//
// if err := GroupRewrite(os.Stdin, os.Stdout, Column(0), LastRow); err != nil {
//     log.Fatal(err)
// }
//
func LastRow(records [][]string) ([][]string, error) {
	if len(records) == 0 {
		return nil, nil
	}
	return [][]string{records[len(records) - 1]}, nil
}
