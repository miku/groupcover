package groupcover

import "github.com/miku/groupcover/container"

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

	// preferred group for key
	preferredGroup := make(map[string]string)

	// group available choices by key
	for k, v := range groups {
		keys := v.Values()
		// preferred group for key k
		preferredGroup[k] = c.Preferences[k].Preferred(keys...)
	}

	// update entries

	var updatedEntries []Entry

	for _, e := range entries {
		var updatedKeys []string
		for _, k := range e.Keys {
			if e.Group == preferredGroup[k] {
				updatedKeys = append(updatedKeys, k)
			}
		}
		entry := Entry{
			ID:    e.ID,
			Group: e.Group,
			Attr:  e.Attr,
			Keys:  updatedKeys,
		}
		updatedEntries = append(updatedEntries, entry)
	}

	return updatedEntries
}
