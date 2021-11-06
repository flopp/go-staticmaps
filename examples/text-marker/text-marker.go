// This is an example on how to create a custom text marker.

package main

import (
	"image/color"
	"math"

	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	sm "github.com/shanghuiyang/go-staticmaps"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// TextMarker
type TextMarker struct {
	sm.MapObject
	Position   s2.LatLng
	Text       string
	TextWidth  float64
	TextHeight float64
	TipSize    float64
}

// NewTextMarker creates a new TextMarker
func NewTextMarker(pos s2.LatLng, text string) *TextMarker {
	s := new(TextMarker)
	s.Position = pos
	s.Text = text
	s.TipSize = 8.0

	d := &font.Drawer{
		Face: basicfont.Face7x13,
	}
	s.TextWidth = float64(d.MeasureString(s.Text) >> 6)
	s.TextHeight = 13.0
	return s
}

func (s *TextMarker) ExtraMarginPixels() (float64, float64, float64, float64) {
	w := math.Max(4.0+s.TextWidth, 2*s.TipSize)
	h := s.TipSize + s.TextHeight + 4.0
	return w * 0.5, h, w * 0.5, 0.0
}

func (s *TextMarker) Bounds() s2.Rect {
	r := s2.EmptyRect()
	r = r.AddPoint(s.Position)
	return r
}

func (s *TextMarker) Draw(gc *gg.Context, trans *sm.Transformer) {
	if !sm.CanDisplay(s.Position) {
		return
	}

	w := math.Max(4.0+s.TextWidth, 2*s.TipSize)
	h := s.TextHeight + 4.0
	x, y := trans.LatLngToXY(s.Position)
	gc.ClearPath()
	gc.SetLineWidth(1)
	gc.SetLineCap(gg.LineCapRound)
	gc.SetLineJoin(gg.LineJoinRound)
	gc.LineTo(x, y)
	gc.LineTo(x-s.TipSize, y-s.TipSize)
	gc.LineTo(x-w*0.5, y-s.TipSize)
	gc.LineTo(x-w*0.5, y-s.TipSize-h)
	gc.LineTo(x+w*0.5, y-s.TipSize-h)
	gc.LineTo(x+w*0.5, y-s.TipSize)
	gc.LineTo(x+s.TipSize, y-s.TipSize)
	gc.LineTo(x, y)
	gc.SetColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.FillPreserve()
	gc.SetColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	gc.Stroke()

	gc.SetRGBA(0.0, 0.0, 0.0, 1.0)
	gc.DrawString(s.Text, x-s.TextWidth*0.5, y-s.TipSize-4.0)
}

func main() {
	ctx := sm.NewContext()
	ctx.SetSize(400, 300)

	berlin := NewTextMarker(s2.LatLngFromDegrees(52.517037, 13.388860), "Berlin")
	london := NewTextMarker(s2.LatLngFromDegrees(51.507322, 0.127647), "London")
	paris := NewTextMarker(s2.LatLngFromDegrees(48.856697, 2.351462), "Paris")
	ctx.AddObject(berlin)
	ctx.AddObject(london)
	ctx.AddObject(paris)

	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}

	if err := gg.SavePNG("text-markers.png", img); err != nil {
		panic(err)
	}
}
