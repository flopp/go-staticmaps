// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package staticmaps

import (
	"image/color"
)

type Marker struct {
	Position LatLng
	Color    color.RGBA
	Size     float64
}
