{{define "header" -}}
  [<img src="{{.Image}}" width="400"/>]({{.URL}})

{{range .Badges -}}
  {{template "badge" .}}
{{end}}
{{end -}}

{{define "badge" -}} [![{{.Text}}]({{.Image}})]({{.URL}}) {{- end}}
