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
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/jessevdk/go-flags"
)

func handleTypeOption(ctx *sm.Context, parameter string) {
	tileProviders := sm.GetTileProviders()
	tp := tileProviders[parameter]
	if tp != nil {
		ctx.SetTileProvider(tp)
		return
	}

	if parameter != "list" {
		fmt.Println("Bad map type:", parameter)
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
}

func handleCenterOption(ctx *sm.Context, parameter string) {
	lat, lng, err := coordsparser.Parse(parameter)
	if err != nil {
		log.Fatal(err)
	} else {
		ctx.SetCenter(s2.LatLngFromDegrees(lat, lng))
	}
}

func handleBboxOption(ctx *sm.Context, parameter string) {
	pair := strings.Split(parameter, "|")
	if len(pair) != 2 {
		log.Fatalf("Bad NW|SE coordinates pair: %s", parameter)
	}

	var err error
	var nwlat float64
	var nwlng float64
	var selat float64
	var selng float64
	nwlat, nwlng, err = coordsparser.Parse(pair[0])
	if err != nil {
		log.Fatal(err)
	}
	selat, selng, err = coordsparser.Parse(pair[1])
	if err != nil {
		log.Fatal(err)
	}

	var bbox *s2.Rect
	bbox, err = sm.CreateBBox(nwlat, nwlng, selat, selng)
	if err != nil {
		log.Fatal(err)
	}

	ctx.SetBoundingBox(*bbox)
}

func handleBackgroundOption(ctx *sm.Context, parameter string) {
	color, err := sm.ParseColorString(parameter)
	if err != nil {
		log.Fatal(err)
	}

	ctx.SetBackground(color)
}

func handleMarkersOption(ctx *sm.Context, parameters []string) {
	for _, s := range parameters {
		markers, err := sm.ParseMarkerString(s)
		if err != nil {
			log.Fatal(err)
		} else {
			for _, marker := range markers {
				ctx.AddMarker(marker)
			}
		}
	}
}

func handlePathsOption(ctx *sm.Context, parameters []string) {
	for _, s := range parameters {
		paths, err := sm.ParsePathString(s)
		if err != nil {
			log.Fatal(err)
		} else {
			for _, path := range paths {
				ctx.AddPath(path)
			}
		}
	}
}

func handleAreasOption(ctx *sm.Context, parameters []string) {
	for _, s := range parameters {
		area, err := sm.ParseAreaString(s)
		if err != nil {
			log.Fatal(err)
		} else {
			ctx.AddArea(area)
		}
	}
}

func handleCirclesOption(ctx *sm.Context, parameters []string) {
	for _, s := range parameters {
		circles, err := sm.ParseCircleString(s)
		if err != nil {
			log.Fatal(err)
		} else {
			for _, circle := range circles {
				ctx.AddCircle(circle)
			}
		}
	}
}

func main() {
	var opts struct {
		//		ClearCache bool     `long:"clear-cache" description:"Clears the tile cache"`
		Width      int      `long:"width" description:"Width of the generated static map image" value-name:"PIXELS" default:"512"`
		Height     int      `long:"height" description:"Height of the generated static map image" value-name:"PIXELS" default:"512"`
		Output     string   `short:"o" long:"output" description:"Output file name" value-name:"FILENAME" default:"map.png"`
		Type       string   `short:"t" long:"type" description:"Select the map type; list possible map types with '--type list'" value-name:"MAPTYPE"`
		Center     string   `short:"c" long:"center" description:"Center coordinates (lat,lng) of the static map" value-name:"LATLNG"`
		Zoom       int      `short:"z" long:"zoom" description:"Zoom factor" value-name:"ZOOMLEVEL"`
		BBox       string   `short:"b" long:"bbox" description:"Bounding box of the static map" value-name:"nwLATLNG|seLATLNG"`
		Background string   `long:"background" description:"Background color" value-name:"COLOR" default:"transparent"`
		UserAgent  string   `short:"u" long:"useragent" description:"Overwrite the default HTTP user agent string" value-name:"USERAGENT"`
		Markers    []string `short:"m" long:"marker" description:"Add a marker to the static map" value-name:"MARKER"`
		Paths      []string `short:"p" long:"path" description:"Add a path to the static map" value-name:"PATH"`
		Areas      []string `short:"a" long:"area" description:"Add an area to the static map" value-name:"AREA"`
		Circles    []string `short:"C" long:"circle" description:"Add a circle to the static map" value-name:"CIRCLE"`
	}

	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	parser.LongDescription = `Creates a static map`
	_, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}

	if parser.FindOptionByLongName("help").IsSet() {
		parser.WriteHelp(os.Stdout)
		os.Exit(0)
	}

	ctx := sm.NewContext()

	if parser.FindOptionByLongName("type").IsSet() {
		handleTypeOption(ctx, opts.Type)
	}

	ctx.SetSize(opts.Width, opts.Height)

	if parser.FindOptionByLongName("zoom").IsSet() {
		ctx.SetZoom(opts.Zoom)
	}

	if parser.FindOptionByLongName("center").IsSet() {
		handleCenterOption(ctx, opts.Center)
	}

	if parser.FindOptionByLongName("bbox").IsSet() {
		handleBboxOption(ctx, opts.BBox)
	}

	if parser.FindOptionByLongName("background").IsSet() {
		handleBackgroundOption(ctx, opts.Background)
	}

	if parser.FindOptionByLongName("useragent").IsSet() {
		ctx.SetUserAgent(opts.UserAgent)
	}

	handleMarkersOption(ctx, opts.Markers)
	handlePathsOption(ctx, opts.Paths)
	handleAreasOption(ctx, opts.Areas)
	handleCirclesOption(ctx, opts.Circles)

	img, err := ctx.Render()
	if err != nil {
		log.Fatal(err)
	}

	if err = gg.SavePNG(opts.Output, img); err != nil {
		log.Fatal(err)
	}
}
