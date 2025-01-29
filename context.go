// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package sm (~ static maps) renders static map images from OSM tiles with markers, paths, and filled areas.
package sm

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"strings"
	"sync"

	"github.com/fogleman/gg"
	"github.com/golang/geo/r2"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// Context holds all information about the map image that is to be rendered
type Context struct {
	width  int
	height int

	hasZoom bool
	zoom    int
	maxZoom int

	hasCenter bool
	center    s2.LatLng

	hasBoundingBox bool
	boundingBox    s2.Rect

	background color.Color

	objects  []MapObject
	overlays []*TileProvider

	userAgent    string
	online       bool
	tileProvider *TileProvider
	cache        TileCache

	overrideAttribution *string
}

// NewContext creates a new instance of Context
func NewContext() *Context {
	t := new(Context)
	t.width = 512
	t.height = 512
	t.hasZoom = false
	t.maxZoom = 30
	t.hasCenter = false
	t.hasBoundingBox = false
	t.background = nil
	t.userAgent = ""
	t.online = true
	t.tileProvider = NewTileProviderOpenStreetMaps()
	t.cache = NewTileCacheFromUserCache(0777)
	return t
}

// SetTileProvider sets the TileProvider to be used
func (m *Context) SetTileProvider(t *TileProvider) {
	m.tileProvider = t
}

// SetCache takes a nil argument to disable caching
func (m *Context) SetCache(cache TileCache) {
	m.cache = cache
}

// SetOnline enables/disables online
// TileFetcher will only fetch tiles from cache if online = false
func (m *Context) SetOnline(online bool) {
	m.online = online
}

// SetUserAgent sets the HTTP user agent string used when downloading map tiles
func (m *Context) SetUserAgent(a string) {
	m.userAgent = a
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

// SetMaxZoom sets the upper zoom level limit when using dynamic zoom
func (m *Context) SetMaxZoom(n int) {
	m.maxZoom = n
}

// SetCenter sets the center coordinates
func (m *Context) SetCenter(center s2.LatLng) {
	m.center = center
	m.hasCenter = true
}

// SetBoundingBox sets the bounding box
func (m *Context) SetBoundingBox(bbox s2.Rect) {
	m.boundingBox = bbox
	m.hasBoundingBox = true
}

// SetBackground sets the background color (used as a fallback for areas without map tiles)
func (m *Context) SetBackground(col color.Color) {
	m.background = col
}

// AddMarker adds a marker to the Context
//
// Deprecated: AddMarker is deprecated. Use the more general AddObject.
func (m *Context) AddMarker(marker *Marker) {
	m.AddObject(marker)
}

// ClearMarkers removes all markers from the Context
func (m *Context) ClearMarkers() {
	filtered := []MapObject{}
	for _, object := range m.objects {
		switch object.(type) {
		case *Marker:
			// skip
		default:
			filtered = append(filtered, object)
		}
	}
	m.objects = filtered
}

// AddPath adds a path to the Context
//
// Deprecated: AddPath is deprecated. Use the more general AddObject.
func (m *Context) AddPath(path *Path) {
	m.AddObject(path)
}

// ClearPaths removes all paths from the Context
func (m *Context) ClearPaths() {
	filtered := []MapObject{}
	for _, object := range m.objects {
		switch object.(type) {
		case *Path:
			// skip
		default:
			filtered = append(filtered, object)
		}
	}
	m.objects = filtered
}

// AddArea adds an area to the Context
//
// Deprecated: AddArea is deprecated. Use the more general AddObject.
func (m *Context) AddArea(area *Area) {
	m.AddObject(area)
}

// ClearAreas removes all areas from the Context
func (m *Context) ClearAreas() {
	filtered := []MapObject{}
	for _, object := range m.objects {
		switch object.(type) {
		case *Area:
			// skip
		default:
			filtered = append(filtered, object)
		}
	}
	m.objects = filtered
}

// AddCircle adds an circle to the Context
//
// Deprecated: AddCircle is deprecated. Use the more general AddObject.
func (m *Context) AddCircle(circle *Circle) {
	m.AddObject(circle)
}

// ClearCircles removes all circles from the Context
func (m *Context) ClearCircles() {
	filtered := []MapObject{}
	for _, object := range m.objects {
		switch object.(type) {
		case *Circle:
			// skip
		default:
			filtered = append(filtered, object)
		}
	}
	m.objects = filtered
}

// AddObject adds an object to the Context
func (m *Context) AddObject(object MapObject) {
	m.objects = append(m.objects, object)
}

// ClearObjects removes all objects from the Context
func (m *Context) ClearObjects() {
	m.objects = nil
}

// AddOverlay adds an overlay to the Context
func (m *Context) AddOverlay(overlay *TileProvider) {
	m.overlays = append(m.overlays, overlay)
}

// ClearOverlays removes all overlays from the Context
func (m *Context) ClearOverlays() {
	m.overlays = nil
}

// OverrideAttribution sets a custom attribution string (or none if empty)
//
// If the attribution string contains newline characters ("\n") it will printed across multiple lines.
// Pay attention: you might be violating the terms of usage for the
// selected map provider - only use the function if you are aware of this!
func (m *Context) OverrideAttribution(attribution string) {
	m.overrideAttribution = &attribution
}

// Attribution returns the current attribution string - either the overridden
// version (using OverrideAttribution) or the one set by the selected
// TileProvider.
func (m *Context) Attribution() string {
	if m.overrideAttribution != nil {
		return *m.overrideAttribution
	}
	return m.tileProvider.Attribution
}

func (m *Context) determineBounds() s2.Rect {
	r := s2.EmptyRect()
	for _, object := range m.objects {
		r = r.Union(object.Bounds())
	}
	return r
}

func (m *Context) determineExtraMarginPixels() (float64, float64, float64, float64) {
	maxL := 0.0
	maxT := 0.0
	maxR := 0.0
	maxB := 0.0
	if m.Attribution() != "" {
		maxB = 12.0
	}
	for _, object := range m.objects {
		l, t, r, b := object.ExtraMarginPixels()
		maxL = math.Max(maxL, l)
		maxT = math.Max(maxT, t)
		maxR = math.Max(maxR, r)
		maxB = math.Max(maxB, b)
	}
	return maxL, maxT, maxR, maxB
}

func (m *Context) determineZoom(bounds s2.Rect, center s2.LatLng) int {
	b := bounds.AddPoint(center)
	if b.IsEmpty() || b.IsPoint() {
		return 15
	}

	tileSize := m.tileProvider.TileSize
	marginL, marginT, marginR, marginB := m.determineExtraMarginPixels()
	w := (float64(m.width) - marginL - marginR) / float64(tileSize)
	h := (float64(m.height) - marginT - marginB) / float64(tileSize)
	if w <= 0 || h <= 0 {
		log.Printf("Object margins are bigger than the target image size => ignoring object margins for calculation of the zoom level")
		w = float64(m.width) / float64(tileSize)
		h = float64(m.height) / float64(tileSize)
	}
	minX := (b.Lo().Lng.Degrees() + 180.0) / 360.0
	maxX := (b.Hi().Lng.Degrees() + 180.0) / 360.0
	minY := (1.0 - math.Log(math.Tan(b.Lo().Lat.Radians())+(1.0/math.Cos(b.Lo().Lat.Radians())))/math.Pi) / 2.0
	maxY := (1.0 - math.Log(math.Tan(b.Hi().Lat.Radians())+(1.0/math.Cos(b.Hi().Lat.Radians())))/math.Pi) / 2.0

	dx := maxX - minX
	for dx < 0 {
		dx = dx + 1
	}
	for dx > 1 {
		dx = dx - 1
	}
	dy := math.Abs(maxY - minY)

	zoom := 1
	for zoom < m.maxZoom {
		tiles := float64(uint(1) << uint(zoom))
		if dx*tiles > w || dy*tiles > h {
			return zoom - 1
		}
		zoom = zoom + 1
	}

	return m.maxZoom
}

// determineCenter computes a point that is visually centered in Mercator projection
func (m *Context) determineCenter(bounds s2.Rect) s2.LatLng {
	latLo := bounds.Lo().Lat.Radians()
	latHi := bounds.Hi().Lat.Radians()
	yLo := math.Log((1+math.Sin(latLo))/(1-math.Sin(latLo))) / 2
	yHi := math.Log((1+math.Sin(latHi))/(1-math.Sin(latHi))) / 2
	lat := s1.Angle(math.Atan(math.Sinh((yLo + yHi) / 2)))
	lng := bounds.Center().Lng
	return s2.LatLng{Lat: lat, Lng: lng}
}

// adjustCenter adjust the center such that the map objects are properly centerd in the view wrt. their pixel margins.
func (m *Context) adjustCenter(center s2.LatLng, zoom int) s2.LatLng {
	if len(m.objects) == 0 {
		return center
	}

	transformer := newTransformer(m.width, m.height, zoom, center, m.tileProvider.TileSize)

	first := true
	minX := 0.0
	maxX := 0.0
	minY := 0.0
	maxY := 0.0
	for _, object := range m.objects {
		bounds := object.Bounds()
		nwX, nwY := transformer.LatLngToXY(bounds.Vertex(3))
		seX, seY := transformer.LatLngToXY(bounds.Vertex(1))
		l, t, r, b := object.ExtraMarginPixels()
		if first {
			minX = nwX - l
			maxX = seX + r
			minY = nwY - t
			maxY = seY + b
			first = false
		} else {
			minX = math.Min(minX, nwX-l)
			maxX = math.Max(maxX, seX+r)
			minY = math.Min(minY, nwY-t)
			maxY = math.Max(maxY, seY+b)
		}
	}

	if (maxX-minX) > float64(m.width) || (maxY-minY) > float64(m.height) {
		log.Printf("Object margins are bigger than the target image size => ignoring object margins for adjusting the center")
		return center
	}

	centerX := (maxX + minX) * 0.5
	centerY := (maxY + minY) * 0.5

	return transformer.XYToLatLng(centerX, centerY)
}

func (m *Context) determineZoomCenter() (int, s2.LatLng, error) {
	if m.hasBoundingBox && !m.boundingBox.IsEmpty() {
		center := m.determineCenter(m.boundingBox)
		return m.determineZoom(m.boundingBox, center), center, nil
	}

	if m.hasCenter {
		if m.hasZoom {
			return m.zoom, m.center, nil
		}
		return m.determineZoom(m.determineBounds(), m.center), m.center, nil
	}

	bounds := m.determineBounds()
	if !bounds.IsEmpty() {
		center := m.determineCenter(bounds)
		zoom := m.zoom
		if !m.hasZoom {
			zoom = m.determineZoom(bounds, center)
		}
		return zoom, m.adjustCenter(center, zoom), nil
	}

	return 0, s2.LatLngFromDegrees(0, 0), errors.New("cannot determine map extent: no center coordinates given, no bounding box given, no content (markers, paths, areas) given")
}

// Transformer implements coordinate transformation from latitude longitude to image pixel coordinates.
type Transformer struct {
	zoom               int
	numTiles           float64 // number of tiles per dimension at this zoom level
	tileSize           int     // tile size in pixels from this provider
	pWidth, pHeight    int     // pixel size of returned set of tiles
	pCenterX, pCenterY int     // pixel location of requested center in set of tiles
	tCountX, tCountY   int     // download area in tile units
	tCenterX, tCenterY float64 // tile index to requested center
	tOriginX, tOriginY int     // bottom left tile to download
	pMinX, pMaxX       int
	proj               s2.Projection
}

// Transformer returns an initialized Transformer instance.
func (m *Context) Transformer() (*Transformer, error) {
	zoom, center, err := m.determineZoomCenter()
	if err != nil {
		return nil, err
	}

	return newTransformer(m.width, m.height, zoom, center, m.tileProvider.TileSize), nil
}

func newTransformer(width int, height int, zoom int, llCenter s2.LatLng, tileSize int) *Transformer {
	t := new(Transformer)

	t.zoom = zoom
	t.numTiles = math.Exp2(float64(t.zoom))
	t.tileSize = tileSize
	// mercator projection from -0.5 to 0.5
	t.proj = s2.NewMercatorProjection(0.5)

	// fractional tile index to center of requested area
	t.tCenterX, t.tCenterY = t.ll2t(llCenter)

	ww := float64(width) / float64(tileSize)
	hh := float64(height) / float64(tileSize)

	// origin tile to fulfill request
	t.tOriginX = int(math.Floor(t.tCenterX - 0.5*ww))
	t.tOriginY = int(math.Floor(t.tCenterY - 0.5*hh))

	// tiles in each axis to fulfill request
	t.tCountX = 1 + int(math.Floor(t.tCenterX+0.5*ww)) - t.tOriginX
	t.tCountY = 1 + int(math.Floor(t.tCenterY+0.5*hh)) - t.tOriginY

	// final pixel dimensions of area returned
	t.pWidth = t.tCountX * tileSize
	t.pHeight = t.tCountY * tileSize

	// Pixel location in returned image for center of requested area
	t.pCenterX = int((t.tCenterX - float64(t.tOriginX)) * float64(tileSize))
	t.pCenterY = int((t.tCenterY - float64(t.tOriginY)) * float64(tileSize))

	t.pMinX = t.pCenterX - width/2
	t.pMaxX = t.pMinX + width

	return t
}

// ll2t returns fractional tile index for a lat/lng points
func (t *Transformer) ll2t(ll s2.LatLng) (float64, float64) {
	p := t.proj.FromLatLng(ll)
	return t.numTiles * (p.X + 0.5), t.numTiles * (1 - (p.Y + 0.5))
}

// LatLngToXY transforms a latitude longitude pair into image x, y coordinates.
func (t *Transformer) LatLngToXY(ll s2.LatLng) (float64, float64) {
	x, y := t.ll2t(ll)
	x = float64(t.pCenterX) + (x-t.tCenterX)*float64(t.tileSize)
	y = float64(t.pCenterY) + (y-t.tCenterY)*float64(t.tileSize)

	offset := t.numTiles * float64(t.tileSize)
	if x < float64(t.pMinX) {
		for x < float64(t.pMinX) {
			x = x + offset
		}
	} else if x >= float64(t.pMaxX) {
		for x >= float64(t.pMaxX) {
			x = x - offset
		}
	}
	return x, y
}

// XYToLatLng transforms image x, y coordinates to  a latitude longitude pair.
func (t *Transformer) XYToLatLng(x float64, y float64) s2.LatLng {
	xx := ((((x - float64(t.pCenterX)) / float64(t.tileSize)) + t.tCenterX) / t.numTiles) - 0.5
	yy := 0.5 - (((y-float64(t.pCenterY))/float64(t.tileSize))+t.tCenterY)/t.numTiles
	return t.proj.ToLatLng(r2.Point{X: xx, Y: yy})
}

// Rect returns an s2.Rect bounding box around the set of tiles described by Transformer.
func (t *Transformer) Rect() (bbox s2.Rect) {
	// transform from https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Go
	invNumTiles := 1.0 / t.numTiles
	// Get latitude bounds
	n := math.Pi - 2.0*math.Pi*float64(t.tOriginY)*invNumTiles
	bbox.Lat.Hi = math.Atan(0.5 * (math.Exp(n) - math.Exp(-n)))
	n = math.Pi - 2.0*math.Pi*float64(t.tOriginY+t.tCountY)*invNumTiles
	bbox.Lat.Lo = math.Atan(0.5 * (math.Exp(n) - math.Exp(-n)))
	// Get longtitude bounds, much easier
	bbox.Lng.Lo = float64(t.tOriginX)*invNumTiles*2.0*math.Pi - math.Pi
	bbox.Lng.Hi = float64(t.tOriginX+t.tCountX)*invNumTiles*2.0*math.Pi - math.Pi
	return bbox
}

// Render actually renders the map image including all map objects (markers, paths, areas)
func (m *Context) Render() (image.Image, error) {
	zoom, center, err := m.determineZoomCenter()
	if err != nil {
		return nil, err
	}

	tileSize := m.tileProvider.TileSize
	trans := newTransformer(m.width, m.height, zoom, center, tileSize)
	img := image.NewRGBA(image.Rect(0, 0, trans.pWidth, trans.pHeight))
	gc := gg.NewContextForRGBA(img)
	if m.background != nil {
		draw.Draw(img, img.Bounds(), &image.Uniform{m.background}, image.Point{}, draw.Src)
	}

	// fetch and draw tiles to img
	layers := []*TileProvider{m.tileProvider}
	if m.overlays != nil {
		layers = append(layers, m.overlays...)
	}

	for _, layer := range layers {
		if err := m.renderLayer(gc, zoom, trans, tileSize, layer); err != nil {
			return nil, err
		}
	}

	// draw map objects
	for _, object := range m.objects {
		object.Draw(gc, trans)
	}

	// crop image
	croppedImg := image.NewRGBA(image.Rect(0, 0, int(m.width), int(m.height)))
	draw.Draw(croppedImg, image.Rect(0, 0, int(m.width), int(m.height)),
		img, image.Point{trans.pCenterX - int(m.width)/2, trans.pCenterY - int(m.height)/2},
		draw.Src)

	// draw attribution
	attribution := m.Attribution()
	if attribution == "" {
		return croppedImg, nil
	}
	lines := strings.Split(attribution, "\n")
	lineHeight := 0.0
	for _, line := range lines {
		_, h := gc.MeasureString(line)
		if h > float64(lineHeight) {
			lineHeight = h
		}
	}
	margin := 2.0
	spacing := 2.0
	boxHeight := lineHeight*float64(len(lines)) + 2*margin + spacing*float64(len(lines)-1)
	gc = gg.NewContextForRGBA(croppedImg)
	gc.SetRGBA(0.0, 0.0, 0.0, 0.5)
	gc.DrawRectangle(0.0, float64(m.height)-boxHeight, float64(m.width), boxHeight)
	gc.Fill()
	gc.SetRGBA(1.0, 1.0, 1.0, 0.75)
	y := float64(m.height) - boxHeight
	for _, line := range lines {
		gc.DrawStringAnchored(line, margin, y, 0, 1)
		y += spacing + lineHeight
	}

	return croppedImg, nil
}

// RenderWithTransformer actually renders the map image including all map objects (markers, paths, areas).
// The returned image covers requested area as well as any tiles necessary to cover that area, which may
// be larger than the request.
//
// A Transformer is returned to support image registration with other data.
func (m *Context) RenderWithTransformer() (image.Image, *Transformer, error) {
	zoom, center, err := m.determineZoomCenter()
	if err != nil {
		return nil, nil, err
	}

	tileSize := m.tileProvider.TileSize
	trans := newTransformer(m.width, m.height, zoom, center, tileSize)
	img := image.NewRGBA(image.Rect(0, 0, trans.pWidth, trans.pHeight))
	gc := gg.NewContextForRGBA(img)
	if m.background != nil {
		draw.Draw(img, img.Bounds(), &image.Uniform{m.background}, image.Point{}, draw.Src)
	}

	// fetch and draw tiles to img
	layers := []*TileProvider{m.tileProvider}
	if m.overlays != nil {
		layers = append(layers, m.overlays...)
	}

	for _, layer := range layers {
		if err := m.renderLayer(gc, zoom, trans, tileSize, layer); err != nil {
			return nil, nil, err
		}
	}

	// draw map objects
	for _, object := range m.objects {
		object.Draw(gc, trans)
	}

	// draw attribution
	if m.tileProvider.Attribution == "" {
		return img, trans, nil
	}
	_, textHeight := gc.MeasureString(m.tileProvider.Attribution)
	boxHeight := textHeight + 4.0
	gc.SetRGBA(0.0, 0.0, 0.0, 0.5)
	gc.DrawRectangle(0.0, float64(trans.pHeight)-boxHeight, float64(trans.pWidth), boxHeight)
	gc.Fill()
	gc.SetRGBA(1.0, 1.0, 1.0, 0.75)
	gc.DrawString(m.tileProvider.Attribution, 4.0, float64(m.height)-4.0)

	return img, trans, nil
}

// RenderWithBounds actually renders the map image including all map objects (markers, paths, areas).
// The returned image covers requested area as well as any tiles necessary to cover that area, which may
// be larger than the request.
//
// Specific bounding box of returned image is provided to support image registration with other data
func (m *Context) RenderWithBounds() (image.Image, s2.Rect, error) {
	img, trans, err := m.RenderWithTransformer()
	if err != nil {
		return nil, s2.Rect{}, err

	}
	return img, trans.Rect(), nil
}

func (m *Context) renderLayer(gc *gg.Context, zoom int, trans *Transformer, tileSize int, provider *TileProvider) error {
	if provider.IsNone() {
		return nil
	}

	var wg sync.WaitGroup
	tiles := (1 << uint(zoom))
	fetchedTiles := make(chan *Tile)
	t := NewTileFetcher(provider, m.cache, m.online)
	if m.userAgent != "" {
		t.SetUserAgent(m.userAgent)
	}

	go func() {
		for xx := 0; xx < trans.tCountX; xx++ {
			x := trans.tOriginX + xx
			if x < 0 {
				x = x + tiles
			} else if x >= tiles {
				x = x - tiles
			}
			if x < 0 || x >= tiles {
				log.Printf("Skipping out of bounds tile column %d/?", x)
				continue
			}
			for yy := 0; yy < trans.tCountY; yy++ {
				y := trans.tOriginY + yy
				if y < 0 || y >= tiles {
					log.Printf("Skipping out of bounds tile %d/%d", x, y)
					continue
				}
				wg.Add(1)
				tile := &Tile{Zoom: zoom, X: x, Y: y}
				go func(wg *sync.WaitGroup, tile *Tile, xx, yy int) {
					defer wg.Done()
					if err := t.Fetch(tile); err == nil {
						tile.X = xx * tileSize
						tile.Y = yy * tileSize
						fetchedTiles <- tile
					} else if err == errTileNotFound && provider.IgnoreNotFound {
						log.Printf("Error downloading tile file: %s (Ignored)", err)
					} else {
						log.Printf("Error downloading tile file: %s", err)
					}
				}(&wg, tile, xx, yy)
			}
		}
		wg.Wait()
		close(fetchedTiles)
	}()

	for tile := range fetchedTiles {
		gc.DrawImage(tile.Img, tile.X, tile.Y)
	}

	return nil
}
