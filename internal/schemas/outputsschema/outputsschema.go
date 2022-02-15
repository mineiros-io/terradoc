package outputsschema

import "github.com/hashicorp/hcl/v2"

func RootSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "output",
				LabelNames: []string{"name"},
			},
		},
	}
}

func OutputSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			// TODO: leaving `value` here even though we don't use it to not break parsing
			{
				Name:     "value",
				Required: false,
			},
			{
				Name:     "description",
				Required: false,
			},
		},
	}
}
