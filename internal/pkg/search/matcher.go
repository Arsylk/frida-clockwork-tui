package search

import (
	fzf "github.com/junegunn/fzf/src"
	"github.com/junegunn/fzf/src/util"
)

func NewMatcher(patternBuilder PatternBuilder, eventBox *util.EventBox) *fzf.Matcher {
	return fzf.NewMatcher(patternBuilder, false, false, eventBox, 0)
}
