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
	indent := strings.Repeat(" ", level)

	trimmedString := strings.TrimSpace(v)
	if trimmedString == "" {
		return ""
	}

	lines := strings.SplitAfter(trimmedString, "\n")

	var result string

	for _, line := range lines {
		if strings.Trim(line, " ") == "\n" {
			result += "\n"
		} else {
			result += indent + line
		}
	}

	return result
}

func repeat(str string, n int) string {
	return strings.Repeat(str, n)
}

func init() {
	urlfragmentRegex = regexp.MustCompile("[^a-zA-Z0-9 -]+")
}
