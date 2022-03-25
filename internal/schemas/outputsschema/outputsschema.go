package outputsschema

import "github.com/hashicorp/hcl/v2"

func RootSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "output",
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

func OutputSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "value",
				Required: true,
			},
			{
				Name:     "description",
				Required: false,
			},
		},
	}
}
