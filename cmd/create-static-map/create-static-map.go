// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"sort"

	"github.com/flopp/go-coordsparser"
	"github.com/flopp/go-staticmaps/staticmaps"
	"github.com/golang/geo/s2"
	"github.com/jessevdk/go-flags"
)

func main() {
	var opts struct {
		//		ClearCache bool     `long:"clear-cache" description:"Clears the tile cache"`
		Width   int      `long:"width" description:"Width of the generated static map image" value-name:"PIXELS" default:"512"`
		Height  int      `long:"height" description:"Height of the generated static map image" value-name:"PIXELS" default:"512"`
		Output  string   `short:"o" long:"output" description:"Output file name" value-name:"FILENAME" default:"map.png"`
		Type    string   `short:"t" long:"type" description:"Select the map type; list possible map types with '--type list'" value-name:"MAPTYPE"`
		Center  string   `short:"c" long:"center" description:"Center coordinates (lat,lng) of the static map" value-name:"LATLNG"`
		Zoom    int      `short:"z" long:"zoom" description:"Zoom factor" value-name:"ZOOMLEVEL"`
		Markers []string `short:"m" long:"marker" description:"Add a marker to the static map" value-name:"MARKER"`
		Paths   []string `short:"p" long:"path" description:"Add a path to the static map" value-name:"PATH"`
	}

	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	parser.LongDescription = `Creates a static map`
	_, err := parser.Parse()

	if parser.FindOptionByLongName("help").IsSet() {
		parser.WriteHelp(os.Stdout)
		os.Exit(0)
	}

	m := staticmaps.NewMapCreator()

	tileProviders := staticmaps.GetTileProviders()
	if parser.FindOptionByLongName("type").IsSet() && (opts.Type == "list" || tileProviders[opts.Type] == nil) {
		if opts.Type != "list" {
			fmt.Println("Bad map type:", opts.Type)
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
	} else if parser.FindOptionByLongName("type").IsSet() {
		m.SetTileProvider(tileProviders[opts.Type])
	}

	m.SetSize(opts.Width, opts.Height)

	if parser.FindOptionByLongName("zoom").IsSet() {
		m.SetZoom(opts.Zoom)
	}

	if parser.FindOptionByLongName("center").IsSet() {
		lat, lng, err := coordsparser.Parse(opts.Center)
		if err == nil {
			m.SetCenter(s2.LatLngFromDegrees(lat, lng))
		} else {
			log.Fatal(err)
		}
	}

	for _, markerString := range opts.Markers {
		markers, err := staticmaps.ParseMarkerString(markerString)
		if err != nil {
			log.Fatal(err)
		}
		for _, marker := range markers {
			m.AddMarker(marker)
		}
	}

	for _, pathString := range opts.Paths {
		path, err := staticmaps.ParsePathString(pathString)
		if err != nil {
			log.Fatal(err)
		}
		m.AddPath(path)
	}

	img, err := m.Create()
	if err != nil {
		log.Fatal(err)
		return
	}

	file, err := os.Create(opts.Output)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	png.Encode(file, img)
}
