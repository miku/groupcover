package groupcover

import (
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
