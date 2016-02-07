// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package staticmaps

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"

	"github.com/cheggaaa/pb"
	"github.com/golang/freetype/truetype"
	"github.com/golang/geo/s2"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/mitchellh/go-homedir"
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

func ll2xy(ll s2.LatLng, zoom int) (float64, float64) {
	tiles := math.Exp2(float64(zoom))
	x := tiles * (ll.Lng.Degrees() + 180.0) / 360.0
	y := tiles * (1 - math.Log(math.Tan(ll.Lat.Radians())+(1.0/math.Cos(ll.Lat.Radians())))/math.Pi) / 2.0
	return x, y
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

func LoadFont() {
	fontData, err := Asset("assets/luxisr.ttf")
	if err != nil {
		log.Panic(err)
	}
	font, err := truetype.Parse(fontData)
	if err != nil {
		log.Panic(err)
	}
	draw2d.RegisterFont(draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilySans, Style: draw2d.FontStyleNormal}, font)
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

	center_x, center_y := ll2xy(center, zoom)

	tileSize := m.tileProvider.TileSize
	ww := float64(m.width) / float64(tileSize)
	hh := float64(m.height) / float64(tileSize)
	imgTileOriginX := int(center_x - 0.5*ww)
	imgTileOriginY := int(center_y - 0.5*hh)
	tileCountX := 1 + int(center_x+0.5*ww) - imgTileOriginX
	tileCountY := 1 + int(center_y+0.5*hh) - imgTileOriginY

	imageWidth := tileCountX * tileSize
	imageHeight := tileCountY * tileSize
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	t := NewTileFetcher(m.tileProvider)
	cacheBaseDir, err := homedir.Expand("~/.cache/go-staticmaps")
	if err == nil {
		t.SetCacheBaseDir(cacheBaseDir)
	} else {
		fmt.Println("Unable to determine user's home directory => no caching of downloaded tiles")
	}

	bar := pb.StartNew(tileCountX * tileCountY).Prefix("Fetching tiles: ")
	for xx := 0; xx < tileCountX; xx++ {
		x := imgTileOriginX + xx
		if x < 0 {
			x = x + (1 << uint(zoom))
		}
		for yy := 0; yy < tileCountY; yy++ {
			y := imgTileOriginY + yy
			bar.Increment()
			tileImg, err := t.Fetch(zoom, x, y)

			if err == nil {
				rect := image.Rect(xx*tileSize, yy*tileSize, (xx+1)*tileSize, (yy+1)*tileSize)
				draw.Draw(img, rect, tileImg, image.ZP, draw.Src)
			}
		}
	}
	bar.Finish()

	imgCenterPixelX := int((center_x - float64(imgTileOriginX)) * float64(tileSize))
	imgCenterPixelY := int((center_y - float64(imgTileOriginY)) * float64(tileSize))

	gc := draw2dimg.NewGraphicContext(img)

	for _, path := range m.paths {
		if len(path.Positions) <= 1 {
			break
		}

		gc.SetStrokeColor(path.Color)
		gc.SetFillColor(path.FillColor)
		gc.SetLineWidth(path.Weight)

		for i, ll := range path.Positions {
			x, y := ll2xy(ll, zoom)
			x = float64(imgCenterPixelX) + (x-center_x)*float64(tileSize)
			y = float64(imgCenterPixelY) + (y-center_y)*float64(tileSize)
			if i == 0 {
				gc.MoveTo(x, y)
			} else {
				gc.LineTo(x, y)
			}
		}

		if path.IsFilled {
			gc.Close()
			gc.FillStroke()
		} else {
			gc.Stroke()
		}
	}

	for i := range m.markers {
		marker := m.markers[i]
		gc.SetStrokeColor(color.RGBA{0, 0, 0, 0xff})
		gc.SetFillColor(marker.Color)
		gc.SetLineWidth(1.0)
		x, y := ll2xy(marker.Position, zoom)
		x = float64(imgCenterPixelX) + (x-center_x)*float64(tileSize)
		y = float64(imgCenterPixelY) + (y-center_y)*float64(tileSize) - marker.Size
		radius := 0.5 * marker.Size
		gc.ArcTo(x, y, radius, radius, 150.0*math.Pi/180.0, 240.0*math.Pi/180.0)
		gc.LineTo(x, y+marker.Size)
		gc.Close()
		gc.FillStroke()
	}

	croppedImg := image.NewRGBA(image.Rect(0, 0, int(m.width), int(m.height)))
	draw.Draw(croppedImg, image.Rect(0, 0, int(m.width), int(m.height)),
		img, image.Point{imgCenterPixelX - int(m.width)/2, imgCenterPixelY - int(m.height)/2},
		draw.Src)

	// draw attribution box
	gc = draw2dimg.NewGraphicContext(croppedImg)

	gc.SetFillColor(color.RGBA{0, 0, 0, 0x7f})
	gc.MoveTo(0, float64(m.height)-14.0)
	gc.LineTo(float64(m.width), float64(m.height)-14.0)
	gc.LineTo(float64(m.width), float64(m.height))
	gc.LineTo(0, float64(m.height))
	gc.Close()
	gc.Fill()

	// draw attribution
	gc.SetFontData(draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilySans, Style: draw2d.FontStyleNormal})
	gc.SetStrokeColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.SetFontSize(8)
	gc.FillStringAt(m.tileProvider.Attribution, 4.0, float64(m.height)-4.0)

	return croppedImg, nil
}
