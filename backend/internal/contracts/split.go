package contracts

import (
	"fmt"
	"strings"
)

type Clause struct {
	ID    string
	Index int
	Text  string
}

func Split(text string) []Clause {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	clauses := make([]Clause, 0)
	current := make([]string, 0)

	flush := func() {
		content := strings.TrimSpace(strings.Join(current, "\n"))
		if content != "" {
			index := len(clauses) + 1
			clauses = append(clauses, Clause{
				ID:    fmt.Sprintf("clause-%d", index),
				Index: index,
				Text:  content,
			})
		}
		current = current[:0]
	}

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			flush()
			continue
		}
		current = append(current, line)
	}
	flush()
	return clauses
}
