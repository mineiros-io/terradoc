package schema

import "github.com/hashicorp/hcl/v2"

var RootSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "section",
			LabelNames: []string{},
		},
	},
}

var SectionSchema = &hcl.BodySchema{
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

var VariableSchema = &hcl.BodySchema{
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

var AttributeSchema = &hcl.BodySchema{
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
	},
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "attribute",
			LabelNames: []string{"name"},
		},
	},
}
