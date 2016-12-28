package groupcover

import (
	"io/ioutil"
	"os"
	"testing"
)

type choice struct {
	given  []string
	result string
}

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
				{given: []string{"B", "C"}, result: "B"},
				{given: []string{"D", "C"}, result: "D"},
				{given: []string{"C", "D"}, result: "C"},
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
				t.Errorf("with prefs %v, given %v, got %v, want %v", c.preferences, ch.given, r, ch.result)
			}
		}
	}
}

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

func BenchmarkSimpleRewriter(b *testing.B) {
	file, err := os.Open(`fixtures/input.10k`)
	if err != nil {
		b.Errorf(err.Error())
	}
	for i := 0; i < b.N; i++ {
		if err := GroupRewrite(file, ioutil.Discard, Column(2), SimpleRewriter(Preferences{Default: LexChoice})); err != nil {
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
