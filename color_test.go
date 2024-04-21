package sm

import (
	"image/color"
	"testing"
)

type string_color_err struct {
	input          string
	expected_color color.Color
	expected_error bool
}

func TestParseColor(t *testing.T) {
	for _, test := range []string_color_err{
		{"WHITE", color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"white", color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, false},
		{"yellow", color.RGBA{0xFF, 0xFF, 0x00, 0xFF}, false},
		{"transparent", color.RGBA{0x00, 0x00, 0x00, 0x00}, false},
		{"#FF00FF42", color.RGBA{0xFF, 0x00, 0xFF, 0x42}, false},
		{"#ff00ff42", color.RGBA{0xFF, 0x00, 0xFF, 0x42}, false},
		{"#ff00ff", color.RGBA{0xFF, 0x00, 0xFF, 0xFF}, false},
		{"#f0f", color.RGBA{0xFF, 0x00, 0xFF, 0xFF}, false},
		{"FF00FF42", color.RGBA{0xFF, 0x00, 0xFF, 0x42}, false},
		{"ff00ff42", color.RGBA{0xFF, 0x00, 0xFF, 0x42}, false},
		{"ff00ff", color.RGBA{0xFF, 0x00, 0xFF, 0xFF}, false},
		{"f0f", color.RGBA{0xFF, 0x00, 0xFF, 0xFF}, false},
		{"bad-name", color.RGBA{0x00, 0x00, 0x00, 0x00}, true},
		{"#FF00F", color.RGBA{0x00, 0x00, 0x00, 0x00}, true},
		{"#GGGGGG", color.RGBA{0x00, 0x00, 0x00, 0x00}, true},
		{"", color.RGBA{0x00, 0x00, 0x00, 0x00}, true},
	} {
		c, err := ParseColorString(test.input)
		if test.expected_error {
			if err == nil {
				t.Errorf("error expected when parsing '%s'", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("unexpected error when parsing '%s': %v", test.input, err)
			}
			if c != test.expected_color {
				t.Errorf("unexpected color when parsing '%s': %v expected: %v", test.input, c, test.expected_color)
			}
		}
	}
}
