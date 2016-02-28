package sm

import (
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

// MapObject is the interface for all objects on the map
type MapObject interface {
	bounds() s2.Rect
	extraMarginPixels() float64
	draw(dc *gg.Context, trans *transformer)
}
