package hclparser

import (
	"fmt"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/madlambda/spells/assert"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/types"
)

func TestProcessType(t *testing.T) {
	var testCases = []struct {
		expr string
		want entities.Type
	}{
		{
			expr: "number",
			want: entities.Type{
				TFType: types.TerraformNumber,
			},
		},
		{
			expr: "string",
			want: entities.Type{
				TFType: types.TerraformString,
			},
		},
		{
			expr: "bool",
			want: entities.Type{
				TFType: types.TerraformBool,
			},
		},
		{
			expr: "object(object_label)",
			want: entities.Type{
				TFType: types.TerraformObject,
				Label:  "object_label",
			},
		},
		{
			expr: "resource(google_resource_blah)",
			want: entities.Type{
				TFType: types.TerraformResource,
				Label:  "google_resource_blah",
			},
		},
		{
			expr: "module(a_badass_mineiros_module)",
			want: entities.Type{
				TFType: types.TerraformModule,
				Label:  "a_badass_mineiros_module",
			},
		},
		{
			expr: "list(number)",
			want: entities.Type{
				TFType: types.TerraformList,
				Nested: &entities.Type{
					TFType: types.TerraformNumber,
				},
			},
		},
		{
			expr: "list(string)",
			want: entities.Type{
				TFType: types.TerraformList,
				Nested: &entities.Type{
					TFType: types.TerraformString,
				},
			},
		},
		// {
		// 	// TODO: test is breaking when trying to evaluate this expr
		// 	expr: "list(my_object)",
		// 	want: entities.Type{
		// 		TFType: types.TerraformList,
		// 		Nested: &entities.Type{
		// 			TFType: types.TerraformObject,
		// 			Label:  "my_object",
		// 		},
		// 	},
		// },
		{
			expr: "list(object(my_object))",
			want: entities.Type{
				TFType: types.TerraformList,
				Nested: &entities.Type{
					TFType: types.TerraformObject,
					Label:  "my_object",
				},
			},
		},
		{
			expr: "list(module(another_module))",
			want: entities.Type{
				TFType: types.TerraformList,
				Nested: &entities.Type{
					TFType: types.TerraformModule,
					Label:  "another_module",
				},
			},
		},
		{
			expr: "list(resource(some_aws_resource))",
			want: entities.Type{
				TFType: types.TerraformList,
				Nested: &entities.Type{
					TFType: types.TerraformResource,
					Label:  "some_aws_resource",
				},
			},
		},
		{
			expr: "list(list(list(object(foo))))",
			want: entities.Type{
				TFType: types.TerraformList,
				Nested: &entities.Type{
					TFType: types.TerraformList,
					Nested: &entities.Type{
						TFType: types.TerraformList,
						Nested: &entities.Type{
							TFType: types.TerraformObject,
							Label:  "foo",
						},
					},
				},
			},
		},
		{
			expr: "set(number)",
			want: entities.Type{
				TFType: types.TerraformSet,
				Nested: &entities.Type{
					TFType: types.TerraformNumber,
				},
			},
		},
		{
			expr: "set(string)",
			want: entities.Type{
				TFType: types.TerraformSet,
				Nested: &entities.Type{
					TFType: types.TerraformString,
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.expr, func(t *testing.T) {
			hclExpr, parseDiags := hclsyntax.ParseExpression([]byte(tt.expr), "", hcl.Pos{Line: 1, Column: 1, Byte: 0})
			if parseDiags.HasErrors() {
				t.Fatalf("Error parsing expression: %v", parseDiags.Errs())
			}

			got, err := processType(hclExpr, types.TerraformEmptyType)
			if err != nil {
				t.Logf("[%s] ERROR: %v", tt.expr, err)
			}
			assert.NoError(t, err, fmt.Sprintf("[ERR] evaluating expr %s", tt.expr))

			assertEqualTypes(t, tt.want, got)
		})
	}
}

func assertEqualTypes(t *testing.T, want, got entities.Type) {
	t.Helper()

	assert.EqualStrings(t, want.TFType.String(), got.TFType.String())
	assert.EqualStrings(t, want.Label, got.Label)

	if want.Nested != nil && got.Nested != nil {
		assertEqualTypes(t, *want.Nested, *got.Nested)
	}
}
