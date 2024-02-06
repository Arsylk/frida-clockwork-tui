package event

import (
	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/source"
)

type (
	ErrorMsg          struct{ Err error }
	OnLoadLogData     struct{ Data *source.ParsedLogData }
	LoadedEntriesMsg  struct{ Entries *source.Entries }
	LoadingMsg        struct{ IsLoading bool }
	FormatLogDataMsg  struct{ Data *source.FormatLogData }
	SearchFinishedMsg struct{ Matches *[]source.EntryMatch }
)
