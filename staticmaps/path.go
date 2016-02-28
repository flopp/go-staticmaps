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
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

// Path represents a path or area on the map
type Path struct {
	MapObject
	Positions []s2.LatLng
	Color     color.Color
	IsFilled  bool
	FillColor color.Color
	Weight    float64
}

// ParsePathString parses a string and returns a path
func ParsePathString(s string) (*Path, error) {
	path := new(Path)
	path.Color = color.RGBA{0xff, 0, 0, 0xff}
	path.IsFilled = false
	path.FillColor = color.Transparent
	path.Weight = 5.0

	for _, ss := range strings.Split(s, "|") {
		if strings.HasPrefix(ss, "color:") {
			var err error
			path.Color, err = ParseColorString(strings.TrimPrefix(ss, "color:"))
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(ss, "fillcolor:") {
			path.IsFilled = true
			var err error
			path.FillColor, err = ParseColorString(strings.TrimPrefix(ss, "fillcolor:"))
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(ss, "weight:") {
			var err error
			path.Weight, err = strconv.ParseFloat(strings.TrimPrefix(ss, "weight:"), 64)
			if err != nil {
				return nil, err
			}
		} else {
			lat, lng, err := coordsparser.Parse(ss)
			if err != nil {
				return nil, err
			}
			path.Positions = append(path.Positions, s2.LatLngFromDegrees(lat, lng))
		}

	}
	return path, nil
}

func (p *Path) extraMarginPixels() float64 {
	return 0.5 * p.Weight
}

func (p *Path) bounds() s2.Rect {
	r := s2.EmptyRect()
	for _, ll := range p.Positions {
		r = r.AddPoint(ll)
	}
	return r
}

func (p *Path) draw(dc *gg.Context, trans *transformer) {
	if len(p.Positions) <= 1 {
		return
	}

	dc.ClearPath()

	dc.SetLineWidth(p.Weight)
	dc.SetLineCap(gg.LineCapRound)
	dc.SetLineJoin(gg.LineJoinRound)

	for _, ll := range p.Positions {
		dc.LineTo(trans.ll2p(ll))
	}

	if p.IsFilled {
		dc.ClosePath()
		dc.SetColor(p.FillColor)
		dc.FillPreserve()
		dc.SetColor(p.Color)
		dc.Stroke()
	} else {
		dc.SetColor(p.Color)
		dc.Stroke()
	}
}
