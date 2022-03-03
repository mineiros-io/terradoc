package varsschema

import "github.com/hashicorp/hcl/v2"

func RootSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
			},
			{
				Type:       "module",
				LabelNames: []string{"name"},
			},
			{
				Type:       "locals",
				LabelNames: []string{},
			},
			{
				Type:       "resource",
				LabelNames: []string{"source", "name"},
			},
			{
				Type:       "output",
				LabelNames: []string{"name"},
			},
			{
				Type:       "terraform",
				LabelNames: []string{},
			},
		},
	}
}

func VariableSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "type",
				Required: true,
			},
			{
				Name:     "description",
				Required: false,
			},
			{
				Name:     "default",
				Required: false,
			},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "validation",
				LabelNames: []string{},
			},
		},
	}
}
