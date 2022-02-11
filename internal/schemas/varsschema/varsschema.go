package varsschema

import "github.com/hashicorp/hcl/v2"

func RootSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
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
	}
}
