package modeline

import (
	"github.com/kick-project/kick/internal/resources/file"
	"github.com/kick-project/kick/internal/resources/modeline/parser"
)

// Parse parses mode line from a file
func Parse(path string, input any, lines int) (*ModeLine, error) {
	empty := true
	ml := ModeLine{
		options_uniq: map[string]any{},
		labels_uniq:  map[string]any{},
	}
	src, err := file.ReadLines(path, input, lines)
	if err != nil {
		return nil, err
	}
	items := parser.Parse(path, string(src), lines)
	for _, i := range items {
		if i.Type == parser.OPTION {
			if _, ok := ml.options_uniq[i.Value]; !ok {
				ml.options = append(ml.options, i.Value)
				ml.options_uniq[i.Value] = struct{}{}
				empty = false
			}
		}
		if i.Type == parser.LABEL {
			if _, ok := ml.labels_uniq[i.Value]; !ok {
				ml.labels = append(ml.labels, i.Value)
				ml.labels_uniq[i.Value] = struct{}{}
				empty = false
			}
		}
	}
	if empty {
		return nil, nil
	}
	return &ml, nil
}

type ModeLine struct {
	options      []string
	options_uniq map[string]any
	labels       []string
	labels_uniq  map[string]any
}

func (m ModeLine) GetOptions() []string {
	return m.options
}

func (m ModeLine) Option(option string) bool {
	_, ok := m.options_uniq[option]
	return ok
}

func (m ModeLine) GetLabel() []string {
	return m.labels
}

func (m ModeLine) Label(label string) bool {
	_, ok := m.labels_uniq[label]
	return ok
}
