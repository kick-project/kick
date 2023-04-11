package modeline

import (
	"github.com/kick-project/kick/internal/resources/file"
	"github.com/kick-project/kick/internal/resources/modeline/parser"
)

// Parse parses mode line from a file
func Parse(path string, input any, lines int) (*ModeLine, error) {
	empty := true
	ml := ModeLine{
		options: map[string]any{},
		labels:  map[string]any{},
	}
	src, err := file.ReadFile(path, input)
	if err != nil {
		return nil, err
	}
	items := parser.Parse(path, string(src), lines)
	for _, i := range items {
		if i.Type == parser.OPTION {
			ml.options[i.Value] = struct{}{}
			empty = false
		}
		if i.Type == parser.LABEL {
			ml.labels[i.Value] = struct{}{}
			empty = false
		}
	}
	if empty {
		return nil, nil
	}
	return &ml, nil
}

type ModeLine struct {
	options map[string]any
	labels  map[string]any
}

func (m ModeLine) Option(option string) bool {
	_, ok := m.options[option]
	return ok
}

func (m ModeLine) Label(label string) bool {
	_, ok := m.labels[label]
	return ok
}
