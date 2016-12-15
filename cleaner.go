package groupcover

type Cleaner interface {
	Clean([]Entry) []Entry
}

type SimpleCleaner struct{}

func (c *SimpleCleaner) Clean(entries []Entry) []Entry {
    
}
