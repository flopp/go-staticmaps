// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import "flag"
import "fmt"
import "github.com/flopp/go-staticmaps/staticmaps"
import (
	"image/png"
	"log"
	"os"
)

func main() {
	output := flag.String("output", "output.png", "name of the generated image file")
	width := flag.Int("width", 400, "width of the generated static map image")
	height := flag.Int("height", 300, "height of the generated static map image")
	centerString := flag.String("center", "", `center of the map ("lat,lng", e.g. "47.123,7.567")`)
	flag.Parse()

	fmt.Println("output:", *output)
	fmt.Println("width:", *width)
	fmt.Println("height:", *height)
	fmt.Println("center:", *centerString)

	m := staticmaps.NewMapCreator()
	m.SetSize(*width, *height)
	m.SetZoom(14)

	if *centerString != "" {
		center, err := staticmaps.LatLngFromString(*centerString)
		if err == nil {
			m.SetCenter(center)
		} else {
			log.Fatal(err)
		}
	}

	img, err := m.Create()
	if err != nil {
		log.Fatal(err)
		return
	}

	file, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	png.Encode(file, img)
}
