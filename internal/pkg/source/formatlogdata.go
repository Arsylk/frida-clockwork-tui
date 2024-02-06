package source

type FormatLogItem struct {
	Text       *string
	EntryIndex int
	LineIndex  int
}

func (i FormatLogItem) GetText() *string {
	return i.Text
}
func (i FormatLogItem) GetEntryIndex() int {
	return i.EntryIndex
}
func (i FormatLogItem) GetLineIndex() int {
	return i.LineIndex
}

type FormatLogData struct {
	*ParsedLogData
	items    *[]FormatLogItem
	itemsLen int
}

func NewFormatLogData(data *ParsedLogData) *FormatLogData {
	length := data.GetEntryIndex()
	items := new([]FormatLogItem)
	count := 0
	for i := 0; i < length; i += 1 {
		entry := data.GetEntry(i)
		lines := data.RenderEntry(entry)
		for j := 0; j < len(*lines); j += 1 {
			line := (*lines)[j]
			*items = append(*items, FormatLogItem{
				Text:       &line,
				EntryIndex: i,
				LineIndex:  j,
			})
			count += 1
		}
	}

	return &FormatLogData{
		ParsedLogData: data,
		items:         items,
		itemsLen:      count,
	}
}

func (data *FormatLogData) GetItem(index int) *FormatLogItem {
	return &(*data.items)[index]
}

func (data *FormatLogData) GetEntryIndices(index int) []int {
	match := (*data.items)[index]
	entryIndex := match.EntryIndex

	var start int
	for start = index; start >= 0 && entryIndex == (*data.items)[start].EntryIndex; start += -1 {
	}
	start += 1
	var end int
	for end = index; end < data.itemsLen && entryIndex == (*data.items)[end].EntryIndex; end += 1 {
	}

	indices := make([]int, end-start)
	for i := start; i < end; i += 1 {
		indices[i-start] = i
	}

	return indices
}

// for current FZF implementation
func (data FormatLogData) ItemString(i int) string {
	item := *data.GetItem(i)
	return *data.GetItem(i - item.LineIndex).Text
}

func (data FormatLogData) Len() int {
	return data.itemsLen
}

func (data FormatLogData) Get(index int) *FormatLogItem {
	return data.GetItem(index)
}
