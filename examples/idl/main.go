// This is an example on how to use a multiline attribution string.

package main

import (
	"image/color"

	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

func main() {
	ctx := sm.NewContext()
	ctx.SetSize(1920, 1080)

	newyork := sm.NewMarker(s2.LatLngFromDegrees(40.641766, -73.780968), color.RGBA{255, 0, 0, 255}, 16.0)
	hongkong := sm.NewMarker(s2.LatLngFromDegrees(22.308046, 113.918480), color.RGBA{0, 0, 255, 255}, 16.0)
	ctx.AddObject(newyork)
	ctx.AddObject(hongkong)
	path := make([]s2.LatLng, 0, 2)
	path = append(path, newyork.Position)
	path = append(path, hongkong.Position)
	ctx.AddObject(sm.NewPath(path, color.RGBA{0, 255, 0, 255}, 4.0))

	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}

	if err := gg.SavePNG("idl.png", img); err != nil {
		panic(err)
	}
}
