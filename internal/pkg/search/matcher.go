package search

import (
	fzf "github.com/junegunn/fzf/src"
	"github.com/junegunn/fzf/src/util"
)

const (
	reqRetry util.EventType = iota
	reqReset
)

// MatchRequest represents a search request
type MatchRequest struct {
	chunks   []*fzf.Chunk
	pattern  *fzf.Pattern
	final    bool
	sort     bool
	revision int
}

// Matcher is responsible for performing search
type Matcher struct {
	patternBuilder func([]rune) *fzf.Pattern
	sort           bool
	tac            bool
	eventBox       *util.EventBox
	reqBox         *util.EventBox
	partitions     int
	slab           []*util.Slab
	mergerCache    map[string]*fzf.Merger
	revision       int
}
