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

// Context holds all information about the map image that is to be rendered
type Context struct {
	width  int
	height int

	hasZoom bool
	zoom    int

	hasCenter bool
	center    s2.LatLng

	markers []*Marker
	paths   []*Path
	areas   []*Area

	tileProvider *TileProvider
}

// NewContext creates a new instance of Context
func NewContext() *Context {
	t := new(Context)
	t.width = 512
	t.height = 512
	t.hasZoom = false
	t.hasCenter = false
	t.tileProvider = NewTileProviderMapQuest()
	return t
}

// SetTileProvider sets the TileProvider to be used
func (m *Context) SetTileProvider(t *TileProvider) {
	m.tileProvider = t
}

// SetSize sets the size of the generated image
func (m *Context) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetZoom sets the zoom level
func (m *Context) SetZoom(zoom int) {
	m.zoom = zoom
	m.hasZoom = true
}

// SetCenter sets the center coordinates
func (m *Context) SetCenter(center s2.LatLng) {
	m.center = center
	m.hasCenter = true
}

// AddMarker adds a marker to the Context
func (m *Context) AddMarker(marker *Marker) {
	m.markers = append(m.markers, marker)
}

// ClearMarkers removes all markers from the Context
func (m *Context) ClearMarkers() {
	m.markers = nil
}

// AddPath adds a path to the Context
func (m *Context) AddPath(path *Path) {
	m.paths = append(m.paths, path)
}

// ClearPaths removes all paths from the Context
func (m *Context) ClearPaths() {
	m.paths = nil
}

// AddArea adds an area to the Context
func (m *Context) AddArea(area *Area) {
	m.areas = append(m.areas, area)
}

// ClearAreas removes all areas from the Context
func (m *Context) ClearAreas() {
	m.areas = nil
}

func (m *Context) determineBounds() s2.Rect {
	r := s2.EmptyRect()
	for _, marker := range m.markers {
		r = r.Union(marker.bounds())
	}
	for _, path := range m.paths {
		r = r.Union(path.bounds())
	}
	for _, area := range m.areas {
		r = r.Union(area.bounds())
	}
	return r
}

func (m *Context) determineExtraMarginPixels() float64 {
	p := 0.0
	for _, marker := range m.markers {
		if pp := marker.extraMarginPixels(); pp > p {
			p = pp
		}
	}
	for _, path := range m.paths {
		if pp := path.extraMarginPixels(); pp > p {
			p = pp
		}
	}
	for _, area := range m.areas {
		if pp := area.extraMarginPixels(); pp > p {
			p = pp
		}
	}
	return p
}

func (m *Context) determineZoom(bounds s2.Rect, center s2.LatLng) int {
	b := bounds.AddPoint(center)
	if b.IsEmpty() || b.IsPoint() {
		return 15
	}

	tileSize := m.tileProvider.TileSize
	margin := 4.0 + m.determineExtraMarginPixels()
	w := (float64(m.width) - 2.0*margin) / float64(tileSize)
	h := (float64(m.height) - 2.0*margin) / float64(tileSize)
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

// Render actually renders the map image including all map objects (markers, paths, areas)
func (m *Context) Render() (image.Image, error) {
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

	// draw map objects
	gc := gg.NewContextForRGBA(img)
	for _, area := range m.areas {
		area.draw(gc, trans)
	}
	for _, path := range m.paths {
		path.draw(gc, trans)
	}
	for _, marker := range m.markers {
		marker.draw(gc, trans)
	}

	// crop image
	croppedImg := image.NewRGBA(image.Rect(0, 0, int(m.width), int(m.height)))
	draw.Draw(croppedImg, image.Rect(0, 0, int(m.width), int(m.height)),
		img, image.Point{trans.pCenterX - int(m.width)/2, trans.pCenterY - int(m.height)/2},
		draw.Src)

	// draw attribution
	_, textHeight := gc.MeasureString(m.tileProvider.Attribution)
	boxHeight := textHeight + 4.0
	gc = gg.NewContextForRGBA(croppedImg)
	gc.SetRGBA(0.0, 0.0, 0.0, 0.5)
	gc.DrawRectangle(0.0, float64(m.height)-boxHeight, float64(m.width), boxHeight)
	gc.Fill()
	gc.SetRGBA(1.0, 1.0, 1.0, 0.75)
	gc.DrawString(m.tileProvider.Attribution, 4.0, float64(m.height)-4.0)

	return croppedImg, nil
}
