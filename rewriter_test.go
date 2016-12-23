package groupcover

import "testing"

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
				{given: []string{}, result: "A"},
				{given: []string{"X"}, result: "A"},
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
