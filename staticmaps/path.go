// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package staticmaps

import (
	"image/color"
	"strconv"
	"strings"

	"github.com/flopp/go-coordsparser"
	"github.com/golang/geo/s2"
)

type Path struct {
	Positions []s2.LatLng
	Color     color.RGBA
	IsFilled  bool
	FillColor color.RGBA
	Weight    float64
}

func ParsePathString(s string) (Path, error) {
	path := Path{Positions: nil, Color: color.RGBA{0xff, 0, 0, 0xff}, IsFilled: false, FillColor: color.RGBA{}, Weight: 5.0}

	for _, ss := range strings.Split(s, "|") {
		if strings.HasPrefix(ss, "color:") {
			color, err := ParseColorString(strings.TrimPrefix(ss, "color:"))
			if err != nil {
				return Path{}, err
			}
			path.Color = *color
		} else if strings.HasPrefix(ss, "fillcolor:") {
			color, err := ParseColorString(strings.TrimPrefix(ss, "fillcolor:"))
			if err != nil {
				return Path{}, err
			}
			path.FillColor = *color
			path.IsFilled = true
		} else if strings.HasPrefix(ss, "weight:") {
			weight, err := strconv.ParseFloat(strings.TrimPrefix(ss, "weight:"), 64)
			if err != nil {
				return Path{}, err
			}
			path.Weight = weight
		} else {
			lat, lng, err := coordsparser.Parse(ss)
			if err != nil {
				return Path{}, err
			}
			path.Positions = append(path.Positions, s2.LatLngFromDegrees(lat, lng))
		}

	}
	return path, nil
}
