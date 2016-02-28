// Copyright 2016 Florian Pigorsch. All rights reserved.
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

// Path represents a path or area on the map
type Path struct {
	MapObject
	Positions []s2.LatLng
	Color     color.Color
	Weight    float64
}

// ParsePathString parses a string and returns a path
func ParsePathString(s string) (*Path, error) {
	path := new(Path)
	path.Color = color.RGBA{0xff, 0, 0, 0xff}
	path.Weight = 5.0

	for _, ss := range strings.Split(s, "|") {
		if ok, suffix := hasPrefix(ss, "color:"); ok {
			var err error
			path.Color, err = ParseColorString(suffix)
			if err != nil {
				return nil, err
			}
		} else if ok, suffix := hasPrefix(ss, "weight:"); ok {
			var err error
			path.Weight, err = strconv.ParseFloat(suffix, 64)
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

func (p *Path) draw(gc *gg.Context, trans *transformer) {
	if len(p.Positions) <= 1 {
		return
	}

	gc.ClearPath()
	gc.SetLineWidth(p.Weight)
	gc.SetLineCap(gg.LineCapRound)
	gc.SetLineJoin(gg.LineJoinRound)
	for _, ll := range p.Positions {
		gc.LineTo(trans.ll2p(ll))
	}
	gc.SetColor(p.Color)
	gc.Stroke()
}
