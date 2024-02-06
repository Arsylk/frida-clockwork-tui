package search

import (
	fzf "github.com/junegunn/fzf/src"
	"github.com/junegunn/fzf/src/algo"
)

type PatternBuilder = func([]rune) *fzf.Pattern

func NewPatternBuilder() PatternBuilder {
	return func(runes []rune) *fzf.Pattern {
		return fzf.BuildPattern(
			false,
			algo.FuzzyMatchV2,
			true,
			fzf.CaseIgnore,
			true,
			true,
			false,
			true,
			make([]fzf.Range, 0),
			fzf.Delimiter{},
			runes,
		)
	}
}
