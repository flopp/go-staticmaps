// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sm

import (
	"image/color"
	"strconv"
	"strings"

	"github.com/flopp/go-coordsparser"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

// Area represents a area or area on the map
type Area struct {
	MapObject
	Positions []s2.LatLng
	Color     color.Color
	Fill      color.Color
	Weight    float64
}

// NewArea creates a new Area
func NewArea(positions []s2.LatLng, col color.Color, fill color.Color, weight float64) *Area {
	a := new(Area)
	a.Positions = positions
	a.Color = col
	a.Fill = fill
	a.Weight = weight

	return a
}

// ParseAreaString parses a string and returns an area
func ParseAreaString(s string) (*Area, error) {
	area := new(Area)
	area.Color = color.RGBA{0xff, 0, 0, 0xff}
	area.Fill = color.Transparent
	area.Weight = 5.0

	for _, ss := range strings.Split(s, "|") {
		if ok, suffix := hasPrefix(ss, "color:"); ok {
			var err error
			area.Color, err = ParseColorString(suffix)
			if err != nil {
				return nil, err
			}
		} else if ok, suffix := hasPrefix(ss, "fill:"); ok {
			var err error
			area.Fill, err = ParseColorString(suffix)
			if err != nil {
				return nil, err
			}
		} else if ok, suffix := hasPrefix(ss, "weight:"); ok {
			var err error
			area.Weight, err = strconv.ParseFloat(suffix, 64)
			if err != nil {
				return nil, err
			}
		} else {
			lat, lng, err := coordsparser.Parse(ss)
			if err != nil {
				return nil, err
			}
			area.Positions = append(area.Positions, s2.LatLngFromDegrees(lat, lng))
		}
	}
	return area, nil
}

// ExtraMarginPixels returns the left, top, right, bottom pixel margin of the Area object, which is exactly the line width.
func (p *Area) ExtraMarginPixels() (float64, float64, float64, float64) {
	return p.Weight, p.Weight, p.Weight, p.Weight
}

// Bounds returns the geographical boundary rect (excluding the actual pixel dimensions).
func (p *Area) Bounds() s2.Rect {
	r := s2.EmptyRect()
	for _, ll := range p.Positions {
		r = r.AddPoint(ll)
	}
	return r
}

// Draw draws the object in the given graphical context.
func (p *Area) Draw(gc *gg.Context, trans *Transformer) {
	if len(p.Positions) <= 1 {
		return
	}

	gc.ClearPath()
	gc.SetLineWidth(p.Weight)
	gc.SetLineCap(gg.LineCapRound)
	gc.SetLineJoin(gg.LineJoinRound)
	for _, ll := range p.Positions {
		gc.LineTo(trans.LatLngToXY(ll))
	}
	gc.ClosePath()
	gc.SetColor(p.Fill)
	gc.FillPreserve()
	gc.SetColor(p.Color)
	gc.Stroke()
}
