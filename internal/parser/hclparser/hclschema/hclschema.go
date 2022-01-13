package hclschema

import "github.com/hashicorp/hcl/v2"

func RootSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "header",
				LabelNames: []string{},
			},
			{
				Type:       "section",
				LabelNames: []string{},
			},
			{
				Type:       "references",
				LabelNames: []string{},
			},
		},
	}
}

func HeaderSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "image",
				Required: false,
			},
			{
				Name:     "url",
				Required: false,
			},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "badge",
				LabelNames: []string{"name"},
			},
		},
	}
}

func BadgeSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "image",
				Required: true,
			},
			{
				Name:     "url",
				Required: true,
			},
			{
				Name:     "text",
				Required: true,
			},
		},
	}
}

func ReferencesSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "ref",
				LabelNames: []string{"name"},
			},
		},
	}
}

func RefSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "value",
				Required: true,
			},
		},
	}
}

func SectionSchema() *hcl.BodySchema {
	return &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{
				Name:     "title",
				Required: false,
			},
			{
				Name:     "content",
				Required: false,
			},
			{
				Name:     "toc",
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
			{
				Type:       "output",
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
			{
				Name:     "readme_type",
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
				Name:     "default",
				Required: false,
			},
			{
				Name:     "readme_example",
				Required: false,
			},
			{
				Name:     "readme_type",
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
