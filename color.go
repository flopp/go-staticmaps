// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sm

import (
	"image/color"

	"github.com/mazznoer/csscolorparser"
)

// ParseColorString parses hex color strings (i.e. `#RRGGBB`, `RRGGBBAA`, `#RRGGBBAA`), and named colors (e.g. 'black', 'blue', ...)
func ParseColorString(s string) (color.Color, error) {
	col, err := csscolorparser.Parse(s)
	if err != nil {
		return nil, err
	}

	r, g, b, a := col.RGBA255()
	return color.RGBA{r, g, b, a}, nil
}

// Luminance computes the luminance (~ brightness) of the given color. Range: 0.0 for black to 1.0 for white.
func Luminance(col color.Color) float64 {
	r, g, b, _ := col.RGBA()
	return (float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114) / float64(0xffff)
}
