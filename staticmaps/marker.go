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
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

// Marker represents a marker on the map
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
			var err error
			size, err = parseSizeString(strings.TrimPrefix(ss, "size:"))
			if err != nil {
				return nil, err
			}
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

func (m *Marker) draw(dc *gg.Context, trans *transformer) {
	dc.ClearPath()

	dc.SetLineJoin(gg.LineJoinRound)
	dc.SetLineWidth(1.0)

	radius := 0.5 * m.Size
	x, y := trans.ll2p(m.Position)
	dc.DrawArc(x, y-m.Size, radius, (90.0+60.0)*math.Pi/180.0, (360.0+90.0-60.0)*math.Pi/180.0)
	dc.LineTo(x, y)
	dc.ClosePath()
	dc.SetColor(m.Color)
	dc.FillPreserve()
	dc.SetRGB(0, 0, 0)
	dc.Stroke()
}
