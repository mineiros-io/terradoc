package entities

import "github.com/mineiros-io/terradoc/internal/types"

// Type represents a variable or attribute type with its readme and Terraform type data
type Type struct {
	// Name is the name of the variable or attribute which contains this type definition
	Name string `json:"name"`
	// TerraformType is the specific Terraform type definition for this type
	TerraformType TerraformType `json:"terraform_type"`
	// ReadmeType is an optional type to be rendered on the README document. Can be any string value as opposed to the `Type` field
	ReadmeType string `json:"readme_type,omitempty"`
}

type TerraformType struct {
	Type       types.TerraformType // Type is the TerraformType
	NestedType types.TerraformType // NestedType is the nested TerraformType
}

func (t TerraformType) HasNestedType() bool {
	return t.NestedType != types.TerraformEmptyType
}
