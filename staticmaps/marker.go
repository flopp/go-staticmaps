// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package staticmaps

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/flopp/go-coordsparser"
	"github.com/golang/geo/s2"
)

type Marker struct {
	Position s2.LatLng
	Color    color.RGBA
	Size     float64
}

func ParseSizeString(s string) (float64, error) {
	if s == "mid" {
		return 16.0, nil
	} else if s == "small" {
		return 12.0, nil
	} else if s == "tiny" {
		return 8.0, nil
	}

	return 0.0, fmt.Errorf("Cannot parse size string: %s", s)
}

func ParseMarkerString(s string) ([]Marker, error) {
	markers := make([]Marker, 0, 0)

	color := color.RGBA{0xff, 0, 0, 0xff}
	size := 16.0

	for _, ss := range strings.Split(s, "|") {
		if strings.HasPrefix(ss, "color:") {
			color_, err := ParseColorString(strings.TrimPrefix(ss, "color:"))
			if err != nil {
				return nil, err
			}
			color = *color_
		} else if strings.HasPrefix(ss, "label:") {
			// TODO
		} else if strings.HasPrefix(ss, "size:") {
			size_, err := ParseSizeString(strings.TrimPrefix(ss, "size:"))
			if err != nil {
				return nil, err
			}
			size = size_
		} else {
			lat, lng, err := coordsparser.Parse(ss)
			if err != nil {
				return nil, err
			}
			marker := Marker{s2.LatLngFromDegrees(lat, lng), color, size}
			markers = append(markers, marker)
		}

	}
	return markers, nil
}
