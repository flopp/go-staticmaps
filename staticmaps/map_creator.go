// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package staticmaps

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/cheggaaa/pb"
	"github.com/llgcode/draw2d/draw2dimg"
)

// MapCreator class
type MapCreator struct {
	width  uint
	height uint

	hasZoom bool
	zoom    uint

	hasCenter bool
	center    LatLng

	markers []Marker
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
func (m *MapCreator) SetSize(width, height uint) {
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

func (m *MapCreator) AddMarker(marker Marker) {
	n := len(m.markers)
	if n == cap(m.markers) {
		// Grow. We double its size and add 1, so if the size is zero we still grow.
		newSlice := make([]Marker, n, 2*n+1)
		copy(newSlice, m.markers)
		m.markers = newSlice
	}
	m.markers = m.markers[0 : n+1]
	m.markers[n] = marker
}

func (m *MapCreator) ClearMarkers() {
	m.markers = m.markers[0:0]
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

	center_x, center_y := ll2xy(m.center, m.zoom)

	tile_size := 256
	ww := float64(m.width) / float64(tile_size)
	hh := float64(m.height) / float64(tile_size)
	imgTileOriginX := int(center_x - 0.5*ww)
	imgTileOriginY := int(center_y - 0.5*hh)
	tileCountX := 1 + int(center_x+0.5*ww) - imgTileOriginX
	tileCountY := 1 + int(center_y+0.5*hh) - imgTileOriginY

	imageWidth := tileCountX * tile_size
	imageHeight := tileCountY * tile_size
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	t := NewTileFetcher("http://otile1.mqcdn.com/tiles/1.0.0/osm/%[1]d/%[2]d/%[3]d.png", "cache")

	bar := pb.StartNew(tileCountX * tileCountY).Prefix("Fetching tiles: ")
	for xx := 0; xx < tileCountX; xx++ {
		x := imgTileOriginX + xx
		if x < 0 {
			x = x + (1 << m.zoom)
		}
		for yy := 0; yy < tileCountY; yy++ {
			y := imgTileOriginY + yy
			bar.Increment()
			tileImg, err := t.Fetch(m.zoom, x, y)

			if err == nil {
				rect := image.Rect(xx*tile_size, yy*tile_size, (xx+1)*tile_size, (yy+1)*tile_size)
				draw.Draw(img, rect, tileImg, image.ZP, draw.Src)
			}
		}
	}
	bar.Finish()

	imgCenterPixelX := int((center_x - float64(imgTileOriginX)) * float64(tile_size))
	imgCenterPixelY := int((center_y - float64(imgTileOriginY)) * float64(tile_size))

	gc := draw2dimg.NewGraphicContext(img)

	for i := range m.markers {
		marker := m.markers[i]
		gc.SetStrokeColor(color.RGBA{0, 0, 0, 0xff})
		gc.SetFillColor(marker.Color)
		x, y := ll2xy(marker.Position, m.zoom)
		x = float64(imgCenterPixelX) + (x-center_x)*float64(tile_size)
		y = float64(imgCenterPixelY) + (y-center_y)*float64(tile_size)
		radius := 0.5 * marker.Size
		gc.ArcTo(x, y, radius, radius, 0, 2*math.Pi)
		gc.Close()
		gc.FillStroke()
	}

	croppedImg := image.NewRGBA(image.Rect(0, 0, int(m.width), int(m.height)))
	draw.Draw(croppedImg, image.Rect(0, 0, int(m.width), int(m.height)),
		img, image.Point{imgCenterPixelX - int(m.width)/2, imgCenterPixelY - int(m.height)/2},
		draw.Src)

	return croppedImg, nil
}
