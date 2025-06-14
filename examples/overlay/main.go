// This is an example on how to use a map overlay.

package main

import (
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

func main() {
	ctx := sm.NewContext()
	ctx.SetSize(1600, 1200)

	ctx.SetCenter(s2.LatLngFromDegrees(48.78110, -3.59638))
	ctx.SetZoom(15)

	// base map
	ctx.SetTileProvider(sm.NewTileProviderOpenStreetMaps())
	ctx.AddOverlay(sm.NewTileProviderOpenSeaMap())

	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}

	if err := gg.SavePNG("overlay.png", img); err != nil {
		panic(err)
	}
}
