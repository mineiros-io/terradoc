{{- range .AllVariables}}{{template "variable" .}}{{end}}

{{define "variable"}}
variable "{{.Name}}" {
  type = {{.Type}}
  {{if .Description}}description = "{{.Description}}"{{end}}
  {{if .Default}}default = {{.Default}}{{end}}
}

{{end}}
