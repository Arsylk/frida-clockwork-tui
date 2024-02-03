package event

import (
	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/source"
)

type (
	ErrorMsg           struct{ Err error }
	OnLoadLogData      struct{ Data *source.ParsedLogData }
	LoadedEntriesMsg   struct{ Entries *source.Entries }
	PreparedContentMsg struct {
		Content *string
		All     int
		Found   int
	}
)
