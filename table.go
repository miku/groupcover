package groupcover

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
)

// Preference is just a ordered group names.
type Preference []string

// PreferenceMap maps keys to a Preference.
type PreferenceMap map[string]*Preference

func (p *Preference) Position(choice string) int {
	for i, item := range *p {
		if item == choice {
			return i
		}
	}
	return 0
}

func (p *Preference) Preferred(choices ...string) string {
	var preferred string
	k := len(*p)
	for _, c := range choices {
		pos := p.Position(c)
		if pos <= k {
			preferred = c
			k = pos
		}
	}
	return preferred
}

// Entry represent the relevant attributes of a record.
type Entry struct {
	ID    string
	Group string
	Attr  string
	Keys  []string
}

// UnmarshalText unwraps a line into an Entry.
func (e *Entry) UnmarshalText(text []byte) error {
	parts := bytes.Split(bytes.TrimSpace(text), []byte("\t"))
	if len(parts) == 0 {
		return nil
	}
	if len(parts) != 4 {
		return fmt.Errorf("expected %v columns, got %v", 4, len(parts))
	}
	e.ID = string(parts[0])
	e.Group = string(parts[1])
	e.Attr = string(parts[2])

	var ks []string
	for _, b := range bytes.Split(parts[3], []byte(",")) {
		ks = append(ks, string(b))
	}
	e.Keys = ks
	return nil
}

// Table is a complete list of items to deduplicate. Won't fly for large sets.
type Table struct {
	Entries []Entry
}

// String represenation of a table.
func (t *Table) String() string {
	var buf bytes.Buffer
	for _, entry := range t.Entries {
		s := fmt.Sprintf("%s\t%s\t%s\t%s\n", entry.ID, entry.Group, entry.Attr, strings.Join(entry.Keys, ", "))
		io.WriteString(&buf, s)
	}
	return buf.String()
}

// TableFromReader reads from a reader and builds up a table.
func TableFromReader(r io.Reader) (*Table, error) {
	br := bufio.NewReader(r)
	var table Table
	for {
		text, err := br.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		var entry Entry
		if err := entry.UnmarshalText(text); err != nil {
			log.Fatal(err)
		}
		table.Entries = append(table.Entries, entry)
	}
	return &table, nil
}
