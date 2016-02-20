// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package staticmaps

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"github.com/flopp/go-coordsparser"
	"github.com/golang/geo/s2"
	"github.com/llgcode/draw2d/draw2dimg"
)

type Marker struct {
	Position s2.LatLng
	Color    color.RGBA
	Size     float64
}

func parseSizeString(s string) (float64, error) {
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
			var err error
			color, err = ParseColorString(strings.TrimPrefix(ss, "color:"))
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(ss, "label:") {
			// TODO
		} else if strings.HasPrefix(ss, "size:") {
			size_, err := parseSizeString(strings.TrimPrefix(ss, "size:"))
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

func (m *Marker) draw(gc *draw2dimg.GraphicContext, trans *transformer) {
	gc.SetStrokeColor(color.RGBA{0, 0, 0, 0xff})
	gc.SetFillColor(m.Color)
	gc.SetLineWidth(1.0)
	radius := 0.5 * m.Size
	x, y := trans.ll2p(m.Position)
	gc.ArcTo(x, y-m.Size, radius, radius, 150.0*math.Pi/180.0, 240.0*math.Pi/180.0)
	gc.LineTo(x, y)
	gc.Close()
	gc.FillStroke()
}
