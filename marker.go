// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sm

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
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
	Label    string
}

// NewMarker creates a new Marker
func NewMarker(pos s2.LatLng, col color.Color, size float64) *Marker {
	m := new(Marker)
	m.Position = pos
	m.Color = col
	m.Size = size
	m.Label = ""
	return m
}

func parseSizeString(s string) (float64, error) {
	switch {
	case s == "mid":
		return 16.0, nil
	case s == "small":
		return 12.0, nil
	case s == "tiny":
		return 8.0, nil
	}

	if ss, err := strconv.ParseFloat(s, 64); err != nil && ss > 0 {
		return ss, nil
	}

	return 0.0, fmt.Errorf("Cannot parse size string: %s", s)
}

// ParseMarkerString parses a string and returns an array of markers
func ParseMarkerString(s string) ([]*Marker, error) {
	markers := make([]*Marker, 0, 0)

	var color color.Color = color.RGBA{0xff, 0, 0, 0xff}
	size := 16.0
	label := ""

	for _, ss := range strings.Split(s, "|") {
		if ok, suffix := hasPrefix(ss, "color:"); ok {
			var err error
			color, err = ParseColorString(suffix)
			if err != nil {
				return nil, err
			}
		} else if ok, suffix := hasPrefix(ss, "label:"); ok {
			label = suffix
		} else if ok, suffix := hasPrefix(ss, "size:"); ok {
			var err error
			size, err = parseSizeString(suffix)
			if err != nil {
				return nil, err
			}
		} else {
			lat, lng, err := coordsparser.Parse(ss)
			if err != nil {
				return nil, err
			}
			m := NewMarker(s2.LatLngFromDegrees(lat, lng), color, size)
			m.Label = label
			markers = append(markers, m)
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

	if m.Label != "" {
		if Luminance(m.Color) >= 0.5 {
			gc.SetColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
		} else {
			gc.SetColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
		}
		gc.DrawStringAnchored(m.Label, x, y-m.Size, 0.5, 0.5)
	}
}
