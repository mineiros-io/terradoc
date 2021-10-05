package renderers

import (
	"bytes"
	"strings"
	"text/template"
)

// Hey mr Reviewer, this looks pretty bad

type templateFunc func(name string, data interface{}) (string, error)

func TemplateFuncMap(t *template.Template) template.FuncMap {
	funcMap := make(template.FuncMap)

	funcMap["indent"] = indentFN
	funcMap["include"] = buildIncludeFN(t)
	funcMap["repeat"] = repeatFN
	funcMap["multiply"] = func(x, y int) int { return x * y } // TODO

	return funcMap
}

func buildIncludeFN(t *template.Template) templateFunc {
	return func(name string, data interface{}) (string, error) {
		buf := bytes.NewBuffer(nil)
		if err := t.ExecuteTemplate(buf, name, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
}

func indentFN(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)

	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}

func repeatFN(str string, n int) string {
	return strings.Repeat(str, n)
}
