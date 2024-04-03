// This is an example on how to use a multiline attribution string.

package main

import (
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

func main() {
	ctx := sm.NewContext()
	ctx.SetSize(400, 300)
	ctx.OverrideAttribution("This is a\nmulti-line\nattribution string.")
	ctx.SetCenter(s2.LatLngFromDegrees(48, 7.9))
	ctx.SetZoom(13)

	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}

	if err := gg.SavePNG("multiline-attribution.png", img); err != nil {
		panic(err)
	}
}
