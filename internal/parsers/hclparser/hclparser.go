package hclparser

import "github.com/hashicorp/hcl/v2"

func GetAttribute(attrs hcl.Attributes, name string) *HCLAttribute {
	attr, exists := attrs[name]
	if exists {
		return &HCLAttribute{attr}
	}

	return nil
}
