package event

import (
	"os"

	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/search"
	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/source"
	tea "github.com/charmbracelet/bubbletea"
	fzf "github.com/junegunn/fzf/src"
	"github.com/junegunn/fzf/src/util"
)

type Runner struct {
	itemCount *int
	data      *source.FormatLogData

	eventBox  *util.EventBox
	matcher   *fzf.Matcher
	chunkList *fzf.ChunkList

	query *[]rune
}

func NewRunner() *Runner {
	count := 0

	eventBox := util.NewEventBox()
	patternBuilder := search.NewPatternBuilder()
	matcher := search.NewMatcher(patternBuilder, eventBox)
	chunkList := fzf.NewChunkList(func(i *fzf.Item, b []byte) bool {
		*i = fzf.NewBaiscItem(string(b), count)
		count += 1
		return true
	})

	return &Runner{
		itemCount: &count,
		eventBox:  eventBox,
		matcher:   matcher,
		chunkList: chunkList,
		query:     &[]rune{},
	}
}

func (r *Runner) Go(p *tea.Program) {
	// async launch fzf matcher
	go r.matcher.Loop()

	// await events
	var snapshot []*fzf.Chunk
	ticks := 0
	for {
		ticks += 1
		r.eventBox.Wait(func(events *util.Events) {
			for evt, value := range *events {
				switch evt {
				case fzf.EvtQuit:
					p.Quit()
					os.Exit(value.(int))
				case fzf.EvtReadNew:
					path := value.(string)
					session, err := source.LoadFromFile(path)
					if err != nil {
						p.Send(ErrorMsg{Err: err})
						continue
					}
					r.data = source.NewFormatLogData(session)
					*r.itemCount = 0
					r.chunkList.Clear()
					for i := 0; i < r.data.Len(); i += 1 {
						r.chunkList.Push([]byte(*r.data.GetItem(i).Text))
					}
					p.Send(FormatLogDataMsg{Data: r.data})
				case fzf.EvtSearchNew:
					p.Send(LoadingMsg{IsLoading: true})
					snapshot, _ = r.chunkList.Snapshot()
					r.matcher.Reset(snapshot, *r.query, true, true, false, 0)
				case fzf.EvtSearchFin:
					switch val := value.(type) {
					case *fzf.Merger:
						matches := r.getMergerMatches(val)
						p.Send(SearchFinishedMsg{Matches: matches})
						p.Send(LoadingMsg{IsLoading: false})
					}
				}
			}
			events.Clear()
		})
	}
}

func (r Runner) getMergerMatches(m *fzf.Merger) *[]source.EntryMatch {
	var matches []source.EntryMatch = make([]source.EntryMatch, m.Length())
	for i := 0; i < m.Length()-1; i += 1 {
		result := m.Get(i)
		index := int(result.GetItem().Index())
		item := r.data.GetItem(index)
		matches[i] = source.EntryMatch{
			Index: index,
			Entry: item.EntryIndex,
			Line:  item.LineIndex,
		}
	}

	return &matches
}

func (r Runner) SetReadEvent(path string) {
	r.eventBox.Set(fzf.EvtReadNew, path)
}

func (r Runner) SetSearchEvent(query string) {
	*r.query = []rune(query)
	r.eventBox.Set(fzf.EvtSearchNew, query)
}
