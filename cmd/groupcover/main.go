package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Attr extract an attribute from a CSV record.
type AttrFunc func([]string) (string, error)

// RewriterFunc rewrites a list of lines.
type RewriterFunc func([][]string) ([][]string, error)

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

	// deduplicate keys
	sampleRewriter := func(s [][]string) ([][]string, error) {
		if len(s) < 2 {
			return s, nil
		}
		// TODO(miku):
		// 1. For each ISIL (keys) get the available groups.
		// 2. For each ISIL (keys) get the preferred group.
		// 3. For each lines, check the group and list the ISIL (keys) for which this group is the preferred.
		for _, line := range s {
			keys := strings.Split(line[3], ",")
			log.Println(keys)
		}
		return s, nil
	}

	if err := GroupLines(os.Stdin, os.Stdout, thirdColumn, sampleRewriter); err != nil {
		log.Fatal(err)
	}
}
