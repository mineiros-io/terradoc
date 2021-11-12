package renderers

import (
	"strings"
	"text/template"
)

var TemplatesFuncMap = template.FuncMap{
	"indent":    indent,
	"repeat":    repeat,
	"multiply":  func(x, y int) int { return x * y },
	"getIndent": GetIndent,
	"newline":   newLine,
}

func newLine() string {
	return "\n\n"
}

func GetIndent(level int) int {
	return level*2 + 2
}
func indent(level int, v string) string {
	if strings.TrimSpace(v) == "" {
		return ""
	}

	indent := strings.Repeat(" ", level)

	lines := strings.SplitAfter(v, "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}

	return strings.Join(append([]string{""}, lines...), indent)
}

func repeat(str string, n int) string {
	return strings.Repeat(str, n)
}
