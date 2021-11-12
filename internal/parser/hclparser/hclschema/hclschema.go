package hclschema

import "github.com/hashicorp/hcl/v2"

func RootSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "section",
				LabelNames: []string{},
			},
		},
	}
}

func SectionSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "title",
				Required: true,
			},
			{
				Name:     "description",
				Required: false,
			},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "section",
				LabelNames: []string{},
			},
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
				Name:     "readme_type",
				Required: false,
			},
			{
				Name:     "description",
				Required: false,
			},
			{
				Name:     "default",
				Required: false,
			},
			{
				Name:     "required",
				Required: false,
			},
			{
				Name:     "forces_recreation",
				Required: false,
			},
			{
				Name:     "readme_example",
				Required: false,
			},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "attribute",
				LabelNames: []string{"name"},
			},
		},
	}
}

func AttributeSchema() *hcl.BodySchema {
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
				Name:     "required",
				Required: false,
			},
			{
				Name:     "forces_recreation",
				Required: false,
			},
			{
				Name:     "readme_type",
				Required: false,
			},
			{
				Name:     "default",
				Required: false,
			},
			{
				Name:     "readme_example",
				Required: false,
			},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "attribute",
				LabelNames: []string{"name"},
			},
		},
	}
}
