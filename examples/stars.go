// This is an example on how to use custom object types (here: 5-pointed stars).

package main

import (
	"image/color"
	"math"
	"math/rand"

	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

// Star represents a 5-pointed star on the map
type Star struct {
	sm.MapObject
	Position s2.LatLng
	Size     float64
}

// NewStar creates a new Star
func NewStar(pos s2.LatLng, size float64) *Star {
	s := new(Star)
	s.Position = pos
	s.Size = size
	return s
}

func (s *Star) ExtraMarginPixels() (float64, float64, float64, float64) {
	return s.Size * 0.5, s.Size * 0.5, s.Size * 0.5, s.Size * 0.5
}

func (s *Star) Bounds() s2.Rect {
	r := s2.EmptyRect()
	r = r.AddPoint(s.Position)
	return r
}

func (s *Star) Draw(gc *gg.Context, trans *sm.Transformer) {
	if !sm.CanDisplay(s.Position) {
		return
	}

	x, y := trans.LatLngToXY(s.Position)
	gc.ClearPath()
	gc.SetLineWidth(1)
	gc.SetLineCap(gg.LineCapRound)
	gc.SetLineJoin(gg.LineJoinRound)
	for i := 0; i <= 10; i++ {
		a := float64(i) * 2 * math.Pi / 10.0
		if i%2 == 0 {
			gc.LineTo(x+s.Size*math.Cos(a), y+s.Size*math.Sin(a))
		} else {
			gc.LineTo(x+s.Size*0.5*math.Cos(a), y+s.Size*0.5*math.Sin(a))
		}
	}
	gc.SetColor(color.RGBA{0xff, 0xff, 0x00, 0xff})
	gc.FillPreserve()
	gc.SetColor(color.RGBA{0xff, 0x00, 0x00, 0xff})
	gc.Stroke()
}

func main() {
	ctx := sm.NewContext()
	ctx.SetSize(400, 300)

	for i := 0; i < 10; i++ {
		star := NewStar(
			s2.LatLngFromDegrees(40+rand.Float64()*10, rand.Float64()*10),
			10+rand.Float64()*10,
		)
		ctx.AddObject(star)
	}

	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}

	if err := gg.SavePNG("stars.png", img); err != nil {
		panic(err)
	}
}
