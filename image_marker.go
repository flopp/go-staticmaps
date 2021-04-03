// Copyright 2021 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sm

import (
	"fmt"
	"image"
	_ "image/jpeg" // to be able to decode jpegs
	_ "image/png"  // to be able to decode pngs
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/flopp/go-coordsparser"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

// ImageMarker represents an image marker on the map
type ImageMarker struct {
	MapObject
	Position s2.LatLng
	Img      image.Image
	OffsetX  float64
	OffsetY  float64
}

// NewImageMarker creates a new ImageMarker
func NewImageMarker(pos s2.LatLng, img image.Image, offsetX, offsetY float64) *ImageMarker {
	m := new(ImageMarker)
	m.Position = pos
	m.Img = img
	m.OffsetX = offsetX
	m.OffsetY = offsetY

	return m
}

// ParseImageMarkerString parses a string and returns an array of image markers
func ParseImageMarkerString(s string) ([]*ImageMarker, error) {
	markers := make([]*ImageMarker, 0)

	var img image.Image = nil
	offsetX := 0.0
	offsetY := 0.0

	for _, ss := range strings.Split(s, "|") {
		if ok, suffix := hasPrefix(ss, "image:"); ok {
			file, err := os.Open(suffix)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			img, _, err = image.Decode(file)
			if err != nil {
				return nil, err
			}
		} else if ok, suffix := hasPrefix(ss, "offsetx:"); ok {
			var err error
			offsetX, err = strconv.ParseFloat(suffix, 64)
			if err != nil {
				return nil, err
			}
		} else if ok, suffix := hasPrefix(ss, "offsety:"); ok {
			var err error
			offsetY, err = strconv.ParseFloat(suffix, 64)
			if err != nil {
				return nil, err
			}
		} else {
			lat, lng, err := coordsparser.Parse(ss)
			if err != nil {
				return nil, err
			}
			if img == nil {
				return nil, fmt.Errorf("cannot create an ImageMarker without an image: %s", s)
			}
			m := NewImageMarker(s2.LatLngFromDegrees(lat, lng), img, offsetX, offsetY)
			markers = append(markers, m)
		}
	}
	return markers, nil
}

// SetImage sets the marker's image
func (m *ImageMarker) SetImage(img image.Image) {
	m.Img = img
}

// SetOffsetX sets the marker's x offset
func (m *ImageMarker) SetOffsetX(offset float64) {
	m.OffsetX = offset
}

// SetOffsetY sets the marker's y offset
func (m *ImageMarker) SetOffsetY(offset float64) {
	m.OffsetY = offset
}

// ExtraMarginPixels return the marker's left, top, right, bottom pixel extent.
func (m *ImageMarker) ExtraMarginPixels() (float64, float64, float64, float64) {
	size := m.Img.Bounds().Size()
	return m.OffsetX, m.OffsetY, float64(size.X) - m.OffsetX, float64(size.Y) - m.OffsetY
}

// Bounds returns single point rect containing the marker's geographical position.
func (m *ImageMarker) Bounds() s2.Rect {
	r := s2.EmptyRect()
	r = r.AddPoint(m.Position)
	return r
}

// Draw draws the object in the given graphical context.
func (m *ImageMarker) Draw(gc *gg.Context, trans *Transformer) {
	if !CanDisplay(m.Position) {
		log.Printf("ImageMarker coordinates not displayable: %f/%f", m.Position.Lat.Degrees(), m.Position.Lng.Degrees())
		return
	}

	x, y := trans.LatLngToXY(m.Position)
	gc.DrawImage(m.Img, int(x-m.OffsetX), int(y-m.OffsetY))
}
