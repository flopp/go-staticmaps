// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package staticmaps renders static map images from OSM tiles with markers, paths, and filled areas.
package staticmaps

import (
	"errors"
	"image"
	"image/draw"
	"math"

	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

// MapCreator class
type MapCreator struct {
	width  int
	height int

	hasZoom bool
	zoom    int

	hasCenter bool
	center    s2.LatLng

	markers []Marker
	paths   []Path

	tileProvider *TileProvider
}

// NewMapCreator creates a new instance of MapCreator
func NewMapCreator() *MapCreator {
	t := new(MapCreator)
	t.width = 512
	t.height = 512
	t.hasZoom = false
	t.hasCenter = false
	t.tileProvider = NewTileProviderMapQuest()
	return t
}

// SetTileProvider sets the TileProvider to be used
func (m *MapCreator) SetTileProvider(t *TileProvider) {
	m.tileProvider = t
}

// SetSize sets the size of the generated image
func (m *MapCreator) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetZoom sets the zoom level
func (m *MapCreator) SetZoom(zoom int) {
	m.zoom = zoom
	m.hasZoom = true
}

// SetCenter sets the center coordinates
func (m *MapCreator) SetCenter(center s2.LatLng) {
	m.center = center
	m.hasCenter = true
}

// AddMarker adds a marker to the MapCreator
func (m *MapCreator) AddMarker(marker Marker) {
	m.markers = append(m.markers, marker)
}

// ClearMarkers removes all markers from the MapCreator
func (m *MapCreator) ClearMarkers() {
	m.markers = nil
}

// AddPath adds a path to the MapCreator
func (m *MapCreator) AddPath(path Path) {
	m.paths = append(m.paths, path)
}

// ClearPaths removes all paths from the MapCreator
func (m *MapCreator) ClearPaths() {
	m.paths = nil
}

func (m *MapCreator) determineBounds() s2.Rect {
	r := s2.EmptyRect()
	for _, marker := range m.markers {
		r = r.AddPoint(marker.Position)
	}
	for _, path := range m.paths {
		for _, position := range path.Positions {
			r = r.AddPoint(position)
		}
	}

	return r
}

func (m *MapCreator) determineZoom(bounds s2.Rect, center s2.LatLng) int {
	b := bounds.AddPoint(center)
	if b.IsEmpty() || b.IsPoint() {
		return 15
	}

	tileSize := m.tileProvider.TileSize
	margin := 16
	w := float64(m.width-2*margin) / float64(tileSize)
	h := float64(m.height-2*margin) / float64(tileSize)
	minX := (b.Lo().Lng.Degrees() + 180.0) / 360.0
	maxX := (b.Hi().Lng.Degrees() + 180.0) / 360.0
	minY := (1.0 - math.Log(math.Tan(b.Lo().Lat.Radians())+(1.0/math.Cos(b.Lo().Lat.Radians())))/math.Pi) / 2.0
	maxY := (1.0 - math.Log(math.Tan(b.Hi().Lat.Radians())+(1.0/math.Cos(b.Hi().Lat.Radians())))/math.Pi) / 2.0
	dx := math.Abs(maxX - minX)
	dy := math.Abs(maxY - minY)

	zoom := 1
	for zoom < 30 {
		tiles := float64(uint(1) << uint(zoom))
		if dx*tiles > w || dy*tiles > h {
			return zoom - 1
		}
		zoom = zoom + 1
	}

	return 15
}

type transformer struct {
	zoom               int
	tileSize           int
	pWidth, pHeight    int
	pCenterX, pCenterY int
	tCountX, tCountY   int
	tCenterX, tCenterY float64
	tOriginX, tOriginY int
}

func newTransformer(width int, height int, zoom int, llCenter s2.LatLng, tileSize int) *transformer {
	t := new(transformer)
	t.zoom = zoom
	t.tileSize = tileSize
	t.tCenterX, t.tCenterY = t.ll2t(llCenter)

	ww := float64(width) / float64(tileSize)
	hh := float64(height) / float64(tileSize)

	t.tOriginX = int(math.Floor(t.tCenterX - 0.5*ww))
	t.tOriginY = int(math.Floor(t.tCenterY - 0.5*hh))

	t.tCountX = 1 + int(math.Floor(t.tCenterX+0.5*ww)) - t.tOriginX
	t.tCountY = 1 + int(math.Floor(t.tCenterY+0.5*hh)) - t.tOriginY

	t.pWidth = t.tCountX * tileSize
	t.pHeight = t.tCountY * tileSize

	t.pCenterX = int((t.tCenterX - float64(t.tOriginX)) * float64(tileSize))
	t.pCenterY = int((t.tCenterY - float64(t.tOriginY)) * float64(tileSize))

	return t
}

func (t *transformer) ll2t(ll s2.LatLng) (float64, float64) {
	tiles := math.Exp2(float64(t.zoom))
	x := tiles * (ll.Lng.Degrees() + 180.0) / 360.0
	y := tiles * (1 - math.Log(math.Tan(ll.Lat.Radians())+(1.0/math.Cos(ll.Lat.Radians())))/math.Pi) / 2.0
	return x, y
}

func (t *transformer) ll2p(ll s2.LatLng) (float64, float64) {
	x, y := t.ll2t(ll)
	x = float64(t.pCenterX) + (x-t.tCenterX)*float64(t.tileSize)
	y = float64(t.pCenterY) + (y-t.tCenterY)*float64(t.tileSize)
	return x, y
}

// Create actually creates the image
func (m *MapCreator) Create() (image.Image, error) {
	bounds := m.determineBounds()

	center := m.center
	if !m.hasCenter {
		if bounds.IsEmpty() {
			return nil, errors.New("No center coordinates specified, cannot determine center from markers")
		}
		center = bounds.Center()
	}

	zoom := m.zoom
	if !m.hasZoom {
		zoom = m.determineZoom(bounds, center)
	}

	tileSize := m.tileProvider.TileSize
	trans := newTransformer(m.width, m.height, zoom, center, tileSize)
	img := image.NewRGBA(image.Rect(0, 0, trans.pWidth, trans.pHeight))

	// fetch and draw tiles to img
	t := NewTileFetcher(m.tileProvider)
	for xx := 0; xx < trans.tCountX; xx++ {
		x := trans.tOriginX + xx
		if x < 0 {
			x = x + (1 << uint(zoom))
		}
		for yy := 0; yy < trans.tCountY; yy++ {
			y := trans.tOriginY + yy
			tileImg, err := t.Fetch(zoom, x, y)

			if err == nil {
				rect := image.Rect(xx*tileSize, yy*tileSize, (xx+1)*tileSize, (yy+1)*tileSize)
				draw.Draw(img, rect, tileImg, image.ZP, draw.Src)
			}
		}
	}

	dc := gg.NewContextForRGBA(img)

	for _, path := range m.paths {
		path.draw(dc, trans)
	}

	for _, marker := range m.markers {
		marker.draw(dc, trans)
	}
	croppedImg := image.NewRGBA(image.Rect(0, 0, int(m.width), int(m.height)))
	draw.Draw(croppedImg, image.Rect(0, 0, int(m.width), int(m.height)),
		img, image.Point{trans.pCenterX - int(m.width)/2, trans.pCenterY - int(m.height)/2},
		draw.Src)

	// draw attribution
	_, textHeight := dc.MeasureString(m.tileProvider.Attribution)
	boxHeight := textHeight + 4.0
	dc = gg.NewContextForRGBA(croppedImg)
	dc.SetRGBA(0.0, 0.0, 0.0, 0.5)
	dc.DrawRectangle(0.0, float64(m.height)-boxHeight, float64(m.width), boxHeight)
	dc.Fill()
	dc.SetRGBA(1.0, 1.0, 1.0, 0.75)
	dc.DrawString(m.tileProvider.Attribution, 4.0, float64(m.height)-4.0)

	return croppedImg, nil
}
