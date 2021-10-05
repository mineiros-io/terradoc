package entities

import (
	"bytes"
	"encoding/json"
	"strings"
)

type Variable struct {
	Name             string       `json:"name"`
	Type             string       `json:"type"`
	Description      string       `json:"description,omitempty"`
	ReadmeType       string       `json:"readme_type,omitempty"`
	Default          string       `json:"default,omitempty"`
	Required         bool         `json:"required,omitempty"`
	ForcesRecreation bool         `json:"forces_recreation,omitempty"`
	ReadmeExample    string       `json:"readme_example,omitempty"`
	Attributes       []*Attribute `json:"attributes,omitempty"`
}

func (v *Variable) readmeExampleJSON() []byte {
	a, _ := json.MarshalIndent(v.ReadmeExample, "", " ")
	b := bytes.NewBuffer(a)
	err := json.Indent(b, a, "-----", "  ")

	if err != nil {
		panic(err)
	}

	return b.Bytes()
}

func (v *Variable) HasAttributes() bool {
	// TODO: check type?
	return len(v.Attributes) > 0
}

func (v *Variable) IsCollection() bool {
	// TODO: refactor
	return strings.Contains(v.ReadmeType, "list") ||
		strings.Contains(v.ReadmeType, "map") ||
		strings.Contains(v.ReadmeType, "set")
}

func (v *Variable) IsComplex() bool {
	return v.Type == "any" ||
		(v.ReadmeType != "" && !strings.Contains(v.ReadmeType, "string"))
}
