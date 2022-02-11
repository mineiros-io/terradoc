package entities_test

import (
	"testing"

	"github.com/madlambda/spells/assert"
	"github.com/mineiros-io/terradoc/internal/entities"
	"github.com/mineiros-io/terradoc/internal/types"
)

func TestTypeAsString(t *testing.T) {
	tests := []struct {
		desc       string
		ty         entities.Type
		wantString string
	}{
		{
			ty: entities.Type{
				TFType: types.TerraformList,
				Nested: &entities.Type{
					TFType: types.TerraformNumber,
				},
			},
			wantString: "list(number)",
		},
		{
			ty: entities.Type{
				TFType: types.TerraformString,
			},
			wantString: "string",
		},
		{
			ty: entities.Type{
				TFType: types.TerraformNumber,
			},
			wantString: "number",
		},
		{
			ty: entities.Type{
				TFType: types.TerraformBool,
			},
			wantString: "bool",
		},
		{
			ty: entities.Type{
				TFType: types.TerraformObject,
				Label:  "foo",
			},
			wantString: "object(foo)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert.EqualStrings(t, tt.wantString, tt.ty.AsString())
		})
	}
}
