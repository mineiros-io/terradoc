package renderers

import (
	"testing"

	"github.com/madlambda/spells/assert"
)

func TestURLFragment(t *testing.T) {
	input := "Backwards compatibility in `0.0.z` and `0.y.z` version"
	want := "backwards-compatibility-in-00z-and-0yz-version"

	assert.EqualStrings(t, want, urlfragment(input))
}

func TestNewLine(t *testing.T) {
	assert.EqualStrings(t, "\n\n", newLine())
}

func TestIndent(t *testing.T) {
	tests := []struct {
		desc        string
		input       string
		want        string
		indentLevel int
	}{
		{
			input:       "one line stuff",
			want:        "    one line stuff",
			indentLevel: 4,
		},
		{
			input: `multi

line


stuff`,
			want: `  multi

  line


  stuff`,
			indentLevel: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := indent(tt.indentLevel, tt.input)

			assert.EqualStrings(t, tt.want, got)
		})
	}
}

func TestRepeat(t *testing.T) {
	assert.EqualStrings(t, "$$$$$", repeat("$", 5))
	assert.EqualStrings(t, "##", repeat("#", 2))
}
