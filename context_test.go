package sm

import (
	"image/color"
	"testing"

	"github.com/golang/geo/s2"
)

func TestRenderEverything(t *testing.T) {
	width := 640
	height := 480

	ctx := NewContext()
	ctx.SetSize(width, height)
	ctx.SetTileProvider(NewTileProviderNone())
	ctx.SetBackground(color.RGBA{255, 255, 255, 255})

	coords1 := s2.LatLngFromDegrees(48.123, 7.0)
	coords2 := s2.LatLngFromDegrees(48.987, 8.0)
	coords3 := s2.LatLngFromDegrees(47.987, 7.5)
	coords4 := s2.LatLngFromDegrees(48.123, 9.0)

	p := make([]s2.LatLng, 0, 3)
	p = append(p, coords1)
	p = append(p, coords2)
	p = append(p, coords3)
	path := NewPath(p, color.RGBA{0, 0, 255, 255}, 4.0)

	a := make([]s2.LatLng, 0, 3)
	a = append(a, coords1)
	a = append(a, coords3)
	a = append(a, coords4)
	area := NewArea(a, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 255, 0, 50}, 3.0)

	marker := NewMarker(coords1, color.RGBA{255, 0, 0, 255}, 16.0)
	circle := NewCircle(coords2, color.RGBA{0, 255, 0, 255}, color.RGBA{0, 255, 0, 100}, 10000.0, 2.0)

	ctx.AddObject(area)
	ctx.AddObject(path)
	ctx.AddObject(marker)
	ctx.AddObject(circle)

	img, err := ctx.Render()
	if err != nil {
		t.Errorf("failed to render: %v", err)
	}

	if img.Bounds().Dx() != width || img.Bounds().Dy() != height {
		t.Errorf("unexpected image size: %d x %d; expected %d x %d", img.Bounds().Dx(), img.Bounds().Dy(), width, height)
	}
}
