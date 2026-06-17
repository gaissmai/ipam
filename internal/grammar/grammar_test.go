package grammar_test

import (
	"testing"

	g "github.com/gaissmai/ipam/internal/grammar"
)

func TestGrammar(t *testing.T) {
	t.Parallel()

	for name := range g.AllRules() {
		rx := g.MustRx(name)
		t.Log(name, rx)
	}
}

func TestGrammarWords(t *testing.T) {
	t.Parallel()
	tt := []struct {
		rule  string
		input string
		want  bool
	}{
		{
			rule:  "word",
			input: "with-hyphen",
			want:  true,
		},
		{
			rule:  "cswords",
			input: "std,open",
			want:  true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			got := g.MustRx(tc.rule).MatchString(tc.input)
			if got != tc.want {
				t.Fatalf("rule: %s, input: %s, got %v, want %v", tc.rule, tc.input, got, tc.want)
			}
		})
	}
}
