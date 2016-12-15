package groupcover

import (
	"fmt"

	"github.com/miku/groupcover/container"
)

type Cleaner interface {
	Clean([]Entry) []Entry
}

// SampleCleaner naive approach.
type SampleCleaner struct {
	Preferences PreferenceMap
}

// Clean the entries.
func (c *SampleCleaner) Clean(entries []Entry) []Entry {
	// first collect groups per key
	groups := make(map[string]*container.StringSet)
	for _, e := range entries {
		for _, key := range e.Keys {
			// TODO(miku): make string set null type usable
			if groups[key] == nil {
				groups[key] = container.NewStringSet()
			}
			groups[key].Add(e.Group)
		}
	}

	// group available choices by key
	for k, v := range groups {
		keys := v.Values()
		// preferred value from available values for key k
		preferred := c.Preferences[k].Preferred(keys...)
		fmt.Println(k, v, preferred)
	}
	return entries
}
