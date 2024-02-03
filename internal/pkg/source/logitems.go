package source

type LogItem struct {
	text  *string
	index int
}

func (i LogItem) GetText() string {
	return *i.text
}
func (i LogItem) GetIndex() int {
	return i.index
}

func (data *ParsedLogData) GetItems() *LogItems {
	length := data.GetEntryIndex()
	items := make(LogItems, length)
	for i := 0; i < length; i += 1 {
		entry := data.GetEntry(i)
		text := RenderParsedLogEntry(entry, data)
		items[i] = LogItem{
			text:  &text,
			index: i,
		}
	}
	return &items
}

type LogItems []LogItem

func (items LogItems) ItemString(i int) string {
	return *items[i].text
}
func (items LogItems) Len() int {
	return len(items)
}

func (items LogItems) Get(index int) *LogItem {
	return &items[index]
}
