package source

import (
	"bufio"
	"fmt"
	"os"
)

type ParsedLogData struct {
	classes    map[JvmId]*JvmClass
	methods    map[JvmId]*JvmMethod
	entries    *[]LogEntry
	entryIndex int
}

func NewParsedLogData() *ParsedLogData {
	return &ParsedLogData{
		classes: make(map[JvmId]*JvmClass),
		methods: make(map[JvmId]*JvmMethod),
		entries: &[]LogEntry{},
	}
}

func LoadFromFile(path string) (*ParsedLogData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("os: %w", err)
	}
	defer file.Close()

	data := NewParsedLogData()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		entry, err := data.JsonToParsedLogEntry(line)
		if err == nil {
			data.AddNewEntry(entry)
		}
	}

	return data, err
}

func (data *ParsedLogData) AddNewEntry(entry LogEntry) {
	switch me := entry.(type) {
	case JvmClass:
		data.classes[me.id] = &me
	case JvmMethod:
		data.methods[me.id] = &me
	}

	*data.entries = append(*data.entries, entry)
	data.entryIndex += 1
}

func (data *ParsedLogData) GetMethod(id JvmId) *JvmMethod {
	if data == nil {
		return nil
	}
	if len(data.methods) == 0 {
		return nil
	}
	return data.methods[id]
}

func (data *ParsedLogData) GetArgType(id JvmId, index int) *string {
	if method := data.GetMethod(id); method != nil {
		return &method.a[index]
	}
	return nil
}

func (data *ParsedLogData) GetReturnType(id JvmId) *string {
	if method := data.GetMethod(id); method != nil {
		return &method.r
	}
	return nil
}

func (data *ParsedLogData) GetEntryIndex() int {
	return data.entryIndex - 1
}

func (data *ParsedLogData) GetEntry(index int) LogEntry {
	return (*data.entries)[index]
}

func (data *ParsedLogData) RenderEntry(e LogEntry) *[]string {
	var value *[]string
	if entry, ok := e.(LogDataRenderer); ok {
		value = entry.Render(data)
	} else {
		value = &[]string{fmt.Sprintf("%s %T", *e.GetTag(), e)}
	}

	return value
}
