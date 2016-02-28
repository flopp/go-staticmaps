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
	MapObject
	Position s2.LatLng
	Color    color.Color
	Size     float64
}

// NewMarker creates a new Marker
func NewMarker(pos s2.LatLng, col color.Color, size float64) *Marker {
	m := new(Marker)
	m.Position = pos
	m.Color = col
	m.Size = size
	return m
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

// ParseMarkerString parses a string and returns an array of markers
func ParseMarkerString(s string) ([]*Marker, error) {
	markers := make([]*Marker, 0, 0)

	var color color.Color = color.RGBA{0xff, 0, 0, 0xff}
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
			markers = append(markers, NewMarker(s2.LatLngFromDegrees(lat, lng), color, size))
		}
	}
	return markers, nil
}

func (m *Marker) extraMarginPixels() float64 {
	return 1.0 + 1.5*m.Size
}

func (m *Marker) bounds() s2.Rect {
	r := s2.EmptyRect()
	r = r.AddPoint(m.Position)
	return r
}

func (m *Marker) draw(gc *gg.Context, trans *transformer) {
	gc.ClearPath()

	gc.SetLineJoin(gg.LineJoinRound)
	gc.SetLineWidth(1.0)

	radius := 0.5 * m.Size
	x, y := trans.ll2p(m.Position)
	gc.DrawArc(x, y-m.Size, radius, (90.0+60.0)*math.Pi/180.0, (360.0+90.0-60.0)*math.Pi/180.0)
	gc.LineTo(x, y)
	gc.ClosePath()
	gc.SetColor(m.Color)
	gc.FillPreserve()
	gc.SetRGB(0, 0, 0)
	gc.Stroke()
}
