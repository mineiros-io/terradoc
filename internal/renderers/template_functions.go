package renderers

import (
	"regexp"
	"strings"
	"text/template"
)

var TemplatesFuncMap = template.FuncMap{
	"urlfragment": urlfragment,
	"indent":      indent,
	"repeat":      repeat,
	"multiply":    func(x, y int) int { return x * y },
	"getIndent":   GetIndent,
	"newline":     newLine,
}

var urlfragmentRegex *regexp.Regexp

func urlfragment(str string) string {
	val := urlfragmentRegex.ReplaceAllString(str, "")

	return strings.ReplaceAll(strings.ToLower(val), " ", "-")
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

func init() {
	urlfragmentRegex = regexp.MustCompile("[^a-zA-Z0-9 -]+")
}
