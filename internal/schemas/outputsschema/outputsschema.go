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
			{
				Name:     "type",
				Required: true,
			},
			{
				Name:     "description",
				Required: false,
			},
		},
	}
}
