package entities

import "github.com/mineiros-io/terradoc/internal/types"

// Type represents a variable or attribute type with its readme and Terraform type data
type Type struct {
	// TFType is the specific Terraform type definition for this type
	TFType types.TerraformType `json:"type"`
	// TFTypeLabel is an optional label for the TerraformType
	TFTypeLabel string `json:"type_label"`
	// TFType is an optional Terraform type definition for the nested type
	NestedTFType types.TerraformType `json:"nested_type"`
	// TFTypeLabel is an optional label for the nested TerraformType
	NestedTFTypeLabel string `json:"nested_type_label"`
}

func (t Type) HasNestedType() bool {
	return t.NestedTFType != types.TerraformEmptyType
}
