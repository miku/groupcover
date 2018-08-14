package groupcover

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// choice helper for tests.
type choice struct {
	given  []string
	result string
}

// TestListChooser tests a selection from a list of elements given preferences
// (most preferred first).
func TestListChooser(t *testing.T) {
	var cases = []struct {
		preferences []string
		choices     []choice
	}{
		{
			preferences: []string{"A", "B", "C", "D"},
			choices: []choice{
				{given: []string{"A"}, result: "A"},
				{given: []string{"B", "C"}, result: "B"},
				{given: []string{"D", "C"}, result: "C"},
				{given: []string{"A", "B", "A"}, result: "A"},
				{given: []string{}, result: ""},
				{given: []string{"X"}, result: "X"},
			},
		},
		{
			preferences: []string{"A"},
			choices: []choice{
				{given: []string{"A"}, result: "A"},
				{given: []string{"B", "C"}, result: "C"},
				{given: []string{"D", "C"}, result: "D"},
				{given: []string{"C", "D"}, result: "D"},
				{given: []string{"A", "B", "A"}, result: "A"},
				{given: []string{}, result: ""},
				{given: []string{"X"}, result: "X"},
			},
		},
	}

	for _, c := range cases {
		choiceFunc := ListChooser(c.preferences)
		for _, ch := range c.choices {
			r := choiceFunc(ch.given)
			if r != ch.result {
				t.Errorf("with prefs %v, given %v, got %v, want %v",
					c.preferences, ch.given, r, ch.result)
			}
		}
	}
}

// TestLexChoice tests lexical choice selection.
func TestLexChoice(t *testing.T) {
	var cases = []struct {
		choices []choice
	}{
		{
			choices: []choice{
				{given: []string{"A"}, result: "A"},
				{given: []string{"B", "C"}, result: "C"},
				{given: []string{"D", "C"}, result: "D"},
				{given: []string{"A", "B", "A"}, result: "B"},
				{given: []string{}, result: ""},
				{given: []string{"X"}, result: "X"},
			},
		},
	}

	for _, c := range cases {
		for _, ch := range c.choices {
			r := LexChoice(ch.given)
			if r != ch.result {
				t.Errorf("with LexChoice, given %v, got %v, want %v", ch.given, r, ch.result)
			}
		}
	}
}

// TestPreferencesWithDefaults checks, whether withDefaults returns a given
// type (using reflection).
func TestPreferencesWithDefaults(t *testing.T) {
	// Test default default.
	var cases = []struct {
		about    string
		prefs    Preferences
		fragment string
	}{
		{
			about:    "test for a default",
			prefs:    Preferences{},
			fragment: "groupcover.LexChoice",
		},
		{
			about:    "test for custom default fallback",
			prefs:    Preferences{Default: ListChooser([]string{"X"})},
			fragment: "groupcover.ListChooser",
		},
	}

	for _, c := range cases {
		f := c.prefs.withDefaults("<MISS>")
		name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		if !strings.Contains(name, c.fragment) {
			t.Errorf("got %s, want something with %s", name, c.fragment)
		}
	}
}

func TestLowerCaseDeduplication(t *testing.T) {
	var cases = []struct {
		about    string
		reader   io.Reader
		attrFunc AttrFunc
		rewriter RewriterFunc
		want     string
		err      error
	}{
		{
			about:    "a basic deduplication example, two records, d0 is dropped from entry",
			reader:   strings.NewReader("i0,s0,v0,d0,d1\ni1,s1,v0,d0\n"),
			attrFunc: ColumnLower(3),
			rewriter: SimpleRewriter(Preferences{}),
			want:     "i0,s0,v0,d1\n",
			err:      nil,
		},
		{
			about:    "by default we are case-sensitive",
			reader:   strings.NewReader("i0,s0,v0,d0,d1\ni1,s1,V0,d0\n"),
			attrFunc: Column(0),
			rewriter: SimpleRewriter(Preferences{}),
			want:     "",
			err:      nil,
		},
		{
			about:    "case insensitive, so d0 is dropped from first entry",
			reader:   strings.NewReader("i0,s0,v0,d0,d1\ni1,s1,V0,d0\n"),
			attrFunc: ColumnLower(3),
			rewriter: SimpleRewriter(Preferences{}),
			want:     "i0,s0,v0,d1\n",
			err:      nil,
		},
		{
			about: "case from #12755",
			reader: strings.NewReader(`ai-60-MTAuMzQxNC9NRTEzLTAxLTAxMzQ,60,https://doi.org/10.3414/ME13-01-0134,DE-15,DE-14
ai-49-aHR0cDovL2R4LmRvaS5vcmcvMTAuMzQxNC9tZTEzLTAxLTAxMzQ,49,https://doi.org/10.3414/me13-01-0134,DE-15,DE-14`),
			attrFunc: ColumnLower(3),
			rewriter: SimpleRewriter(Preferences{}),
			want:     "ai-49-aHR0cDovL2R4LmRvaS5vcmcvMTAuMzQxNC9tZTEzLTAxLTAxMzQ,49,https://doi.org/10.3414/me13-01-0134",
			err:      nil,
		},
	}

	for _, c := range cases {
		var buf bytes.Buffer
		err := GroupRewrite(c.reader, &buf, c.attrFunc, c.rewriter)
		if err != c.err {
			t.Errorf("GroupRewrite (%s): got %v, want %v", c.about, err, c.err)
		}
		if buf.String() != c.want {
			t.Errorf("GroupRewrite (%s): got %v, want %v", c.about, buf.String(), c.want)
		}
	}
}

func BenchmarkSimpleRewriter(b *testing.B) {
	file, err := os.Open(`fixtures/input.10k`)
	if err != nil {
		b.Errorf(err.Error())
	}
	for i := 0; i < b.N; i++ {
		rewriter := SimpleRewriter(Preferences{Default: LexChoice})
		if err := GroupRewrite(file, ioutil.Discard, Column(2), rewriter); err != nil {
			b.Errorf(err.Error())
		}
	}
}

func BenchmarkDiscardRows(b *testing.B) {
	file, err := os.Open(`fixtures/input.10k`)
	if err != nil {
		b.Errorf(err.Error())
	}
	for i := 0; i < b.N; i++ {
		if err := GroupRewrite(file, ioutil.Discard, Column(2), DiscardRows); err != nil {
			b.Errorf(err.Error())
		}
	}
}
