package main

import "testing"

func TestCleanInput(t *testing.T) {
	tests := []struct {
		name  int
		input string
		want  []string
	}{
		{name: 1, input: "   hello   world  ", want: []string{"hello", "world"}},
		{name: 2, input: "  hello,universe  ", want: []string{"hello,universe"}},
		{name: 3, input: "  CHARMANDER  ", want: []string{"charmander"}},
		{name: 4, input: "\thello\nworld\t", want: []string{"hello", "world"}},
		{name: 5, input: " ", want: []string{}},
		{name: 6, input: "How many TESTS until we are done?  ", want: []string{"how", "many", "tests", "until", "we", "are", "done?"}},
	}

	for _, tc := range tests {
		got := cleanInput(tc.input)
		if len(got) != len(tc.want) {
			t.Errorf("Failed test %d. Wanted length %d, got length %d", tc.name, len(tc.want), len(got))
			continue
		}

		for i := range got {
			word := got[i]
			expectedWord := tc.want[i]
			if word != expectedWord {
				t.Errorf("Failed test: %d. word at index %d is not correct. wanted %v but got %v", tc.name, i, tc.want[i], got[i])
			}
		}
	}
}
