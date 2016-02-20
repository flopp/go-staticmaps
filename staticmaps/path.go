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
	"github.com/llgcode/draw2d/draw2dimg"
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
			var err error
			path.Color, err = ParseColorString(strings.TrimPrefix(ss, "color:"))
			if err != nil {
				return Path{}, err
			}
		} else if strings.HasPrefix(ss, "fillcolor:") {
			path.IsFilled = true
			var err error
			path.FillColor, err = ParseColorString(strings.TrimPrefix(ss, "fillcolor:"))
			if err != nil {
				return Path{}, err
			}
		} else if strings.HasPrefix(ss, "weight:") {
			var err error
			path.Weight, err = strconv.ParseFloat(strings.TrimPrefix(ss, "weight:"), 64)
			if err != nil {
				return Path{}, err
			}
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

func (p *Path) draw(gc *draw2dimg.GraphicContext, trans *transformer) {
	if len(p.Positions) <= 1 {
		return
	}

	gc.SetStrokeColor(p.Color)
	gc.SetFillColor(p.FillColor)
	gc.SetLineWidth(p.Weight)

	gc.MoveTo(trans.ll2p(p.Positions[0]))
	for _, ll := range p.Positions[1:] {
		gc.LineTo(trans.ll2p(ll))
	}

	if p.IsFilled {
		gc.Close()
		gc.FillStroke()
	} else {
		gc.Stroke()
	}
}
