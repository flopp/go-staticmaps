// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package staticmaps

import (
	"errors"
	"image"
	"image/draw"
	"math"
)

// MapCreator class
type MapCreator struct {
	width  int
	height int

	hasZoom bool
	zoom    uint

	hasCenter bool
	center    LatLng
}

// NewMapCreator creates a new instance of MapCreator
func NewMapCreator() *MapCreator {
	t := new(MapCreator)
	t.width = 512
	t.height = 512
	t.hasZoom = false
	t.hasCenter = false
	return t
}

// SetSize sets the size of the generated image
func (m *MapCreator) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetZoom sets the zoom level
func (m *MapCreator) SetZoom(zoom uint) {
	m.zoom = zoom
	m.hasZoom = true
}

// SetCenter sets the center coordinates
func (m *MapCreator) SetCenter(center LatLng) {
	m.center = center
	m.hasCenter = true
}

func ll2xy(ll LatLng, zoom uint) (float64, float64) {
	tiles := math.Exp2(float64(zoom))
	x := tiles * (ll.Lng() + 180.0) / 360.0
	y := tiles * (1 - math.Log(math.Tan(ll.LatRadians())+(1.0/math.Cos(ll.LatRadians())))/math.Pi) / 2.0
	return x, y
}

// Create actually creates the image
func (m *MapCreator) Create() (image.Image, error) {
	if !m.hasCenter {
		return nil, errors.New("No center coordinates specified")
	}
	if !m.hasZoom {
		return nil, errors.New("No zoom specified")
	}

	tile_size := 256
	tiles_x := int(math.Ceil(float64(m.width)/float64(tile_size))) + 2
	tiles_y := int(math.Ceil(float64(m.height)/float64(tile_size))) + 2

	x_offset := -int(math.Floor(float64(tiles_x) * 0.5))
	y_offset := -int(math.Floor(float64(tiles_y) * 0.5))

	center_x, center_y := ll2xy(m.center, m.zoom)
	center_tile_x := int(center_x)
	center_tile_y := int(center_y)
	//origin_x := center_tile_x + x_offset
	//origin_y := center_tile_y + y_offset

	imageWidth := tiles_x * tile_size
	imageHeight := tiles_y * tile_size
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	t := NewTileFetcher("http://otile1.mqcdn.com/tiles/1.0.0/osm/%[1]d/%[2]d/%[3]d.png", "cache")

	for xx := 0; xx < tiles_x; xx++ {
		x := center_tile_x + xx + x_offset
		// if x < 0 {
		// 	x = x + math.Exp2(float64(m.zoom))
		// }
		for yy := 0; yy < tiles_y; yy++ {
			y := center_tile_y + yy + y_offset

			tileImg, err := t.Fetch(m.zoom, x, y)

			if err == nil {
				draw.Draw(img, image.Rect(xx*tile_size, yy*tile_size, (xx+1)*tile_size, (yy+1)*tile_size),
					tileImg, image.ZP, draw.Src)
			}
		}
	}

	center_pixel_x := int((float64(x_offset) + center_x - float64(int(center_x))) * float64(tile_size))
	center_pixel_y := int((float64(y_offset) + center_y - float64(int(center_y))) * float64(tile_size))

	croppedImg := image.NewRGBA(image.Rect(0, 0, m.width, m.height))
	draw.Draw(croppedImg, image.Rect(0, 0, m.width, m.height),
		img, image.Point{m.width/2 - center_pixel_x, m.height/2 - center_pixel_y},
		draw.Src)

	return croppedImg, nil
}
