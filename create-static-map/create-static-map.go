// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/flopp/go-coordsparser"
	"github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/jessevdk/go-flags"
)

func getTileProviderOrExit(name string) *sm.TileProvider {
	tileProviders := sm.GetTileProviders()
	tp := tileProviders[name]
	if tp != nil {
		return tp
	}

	if name != "list" {
		fmt.Println("Bad map type:", name)
	}
	fmt.Println("Possible map types (to be used with --type/-t):")
	// print sorted keys
	keys := make([]string, 0, len(tileProviders))
	for k := range tileProviders {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println(k)
	}
	os.Exit(0)

	return nil
}

func main() {
	var opts struct {
		//		ClearCache bool     `long:"clear-cache" description:"Clears the tile cache"`
		Width   int      `long:"width" description:"Width of the generated static map image" value-name:"PIXELS" default:"512"`
		Height  int      `long:"height" description:"Height of the generated static map image" value-name:"PIXELS" default:"512"`
		Output  string   `short:"o" long:"output" description:"Output file name" value-name:"FILENAME" default:"map.png"`
		Type    string   `short:"t" long:"type" description:"Select the map type; list possible map types with '--type list'" value-name:"MAPTYPE"`
		Center  string   `short:"c" long:"center" description:"Center coordinates (lat,lng) of the static map" value-name:"LATLNG"`
		Zoom    int      `short:"z" long:"zoom" description:"Zoom factor" value-name:"ZOOMLEVEL"`
		BBox    string   `short:"b" long:"bbox" description:"Bounding box of the static map" value-name:"LATLNG|LATLNG"`
		Markers []string `short:"m" long:"marker" description:"Add a marker to the static map" value-name:"MARKER"`
		Paths   []string `short:"p" long:"path" description:"Add a path to the static map" value-name:"PATH"`
		Areas   []string `short:"a" long:"area" description:"Add an area to the static map" value-name:"AREA"`
	}

	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	parser.LongDescription = `Creates a static map`
	_, err := parser.Parse()

	if parser.FindOptionByLongName("help").IsSet() {
		parser.WriteHelp(os.Stdout)
		os.Exit(0)
	}

	ctx := sm.NewContext()

	if parser.FindOptionByLongName("type").IsSet() {
		tp := getTileProviderOrExit(opts.Type)
		if tp != nil {
			ctx.SetTileProvider(tp)
		}
	}

	ctx.SetSize(opts.Width, opts.Height)

	if parser.FindOptionByLongName("zoom").IsSet() {
		ctx.SetZoom(opts.Zoom)
	}

	if parser.FindOptionByLongName("center").IsSet() {
		lat, lng, err := coordsparser.Parse(opts.Center)
		if err != nil {
			log.Fatal(err)
		} else {
			ctx.SetCenter(s2.LatLngFromDegrees(lat, lng))
		}
	}

	if parser.FindOptionByLongName("bbox").IsSet() {
		pair := strings.Split(opts.BBox, "|")
		if len(pair) != 2 {
			log.Fatalf("Bad LATLNG|LATLNG pair: %s", opts.BBox)
		}

		bbox := s2.EmptyRect()

		lat, lng, err := coordsparser.Parse(pair[0])
		if err != nil {
			log.Fatal(err)
		} else {
			bbox = bbox.AddPoint(s2.LatLngFromDegrees(lat, lng))
		}

		lat, lng, err = coordsparser.Parse(pair[1])
		if err != nil {
			log.Fatal(err)
		} else {
			bbox = bbox.AddPoint(s2.LatLngFromDegrees(lat, lng))
		}

		ctx.SetBoundingBox(bbox)
	}

	for _, markerString := range opts.Markers {
		markers, err := sm.ParseMarkerString(markerString)
		if err != nil {
			log.Fatal(err)
		} else {
			for _, marker := range markers {
				ctx.AddMarker(marker)
			}
		}
	}

	for _, pathString := range opts.Paths {
		paths, err := sm.ParsePathString(pathString)
		if err != nil {
			log.Fatal(err)
		} else {
			for _, path := range paths {
				ctx.AddPath(path)
			}
		}
	}

	for _, areaString := range opts.Areas {
		area, err := sm.ParseAreaString(areaString)
		if err != nil {
			log.Fatal(err)
		} else {
			ctx.AddArea(area)
		}
	}

	img, err := ctx.Render()
	if err != nil {
		log.Fatal(err)
		return
	}

	if err = gg.SavePNG(opts.Output, img); err != nil {
		log.Fatal(err)
		return
	}
}
