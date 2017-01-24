// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

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
