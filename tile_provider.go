// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sm

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func fetchURL(url string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	// req.Header.Set("User-Agent", t.userAgent)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// TileProvider encapsulates all infos about a map tile provider service (name, url scheme, attribution, etc.)

type TileProvider interface {
	Name() string
	Attribution() string
	TileSize() int
	URLPattern() string
	TileURL(int, int, int) string
	Shards() []string
	FetchTile(int, int, int) ([]byte, error)
}

type DefaultTileProvider struct {
	TileProvider
	name        string
	attribution string
	tileSize    int
	urlPattern  string // "%[1]s" => shard, "%[2]d" => zoom, "%[3]d" => x, "%[4]d" => y
	shards      []string
}

func (t *DefaultTileProvider) Name() string {
	return t.name
}

func (t *DefaultTileProvider) Attribution() string {
	return t.attribution
}

func (t *DefaultTileProvider) TileSize() int {
	return t.tileSize
}

func (t *DefaultTileProvider) URLPattern() string {
	return t.urlPattern
}

func (t *DefaultTileProvider) Shards() []string {
	return t.shards
}

func (t *DefaultTileProvider) TileURL(zoom int, x int, y int) string {

	possible_shards := t.Shards()
	shard := ""

	ss := len(possible_shards)

	if len(possible_shards) > 0 {
		shard = possible_shards[(x+y)%ss]
	}

	return fmt.Sprintf(t.URLPattern(), shard, zoom, x, y)
}

func (t *DefaultTileProvider) FetchTile(z int, x int, y int) ([]byte, error) {

	url := t.TileURL(z, x, y)
	return fetchURL(url)
}

// NewTileProviderOpenStreetMaps creates a TileProvider struct for OSM's tile service

func NewTileProviderOpenStreetMaps() TileProvider {
	t := &DefaultTileProvider{
		name:        "osm",
		attribution: "Maps and Data (c) openstreetmap.org and contributors, ODbL",
		tileSize:    256,
		urlPattern:  "http://%[1]s.tile.openstreetmap.org/%[2]d/%[3]d/%[4]d.png",
		shards:      []string{"a", "b", "c"},
	}
	return t
}

func newTileProviderThunderforest(name string) TileProvider {
	t := &DefaultTileProvider{
		name:        fmt.Sprintf("thunderforest-%s", name),
		attribution: "Maps (c) Thundeforest; Data (c) OSM and contributors, ODbL",
		tileSize:    256,
		urlPattern:  "https://%[1]s.tile.thunderforest.com/" + name + "/%[2]d/%[3]d/%[4]d.png",
		shards:      []string{"a", "b", "c"},
	}
	return t
}

// NewTileProviderThunderforestLandscape creates a TileProvider struct for thundeforests's 'landscape' tile service
func NewTileProviderThunderforestLandscape() TileProvider {
	return newTileProviderThunderforest("landscape")
}

// NewTileProviderThunderforestOutdoors creates a TileProvider struct for thundeforests's 'outdoors' tile service
func NewTileProviderThunderforestOutdoors() TileProvider {
	return newTileProviderThunderforest("outdoors")
}

// NewTileProviderThunderforestTransport creates a TileProvider struct for thundeforests's 'transport' tile service
func NewTileProviderThunderforestTransport() TileProvider {
	return newTileProviderThunderforest("transport")
}

// NewTileProviderStamenToner creates a TileProvider struct for stamens' 'toner' tile service
func NewTileProviderStamenToner() TileProvider {
	t := &DefaultTileProvider{
		name:        "stamen-toner",
		attribution: "Maps (c) Stamen; Data (c) OSM and contributors, ODbL",
		tileSize:    256,
		urlPattern:  "http://%[1]s.tile.stamen.com/toner/%[2]d/%[3]d/%[4]d.png",
		shards:      []string{"a", "b", "c", "d"},
	}
	return t
}

// NewTileProviderStamenTerrain creates a TileProvider struct for stamens' 'terrain' tile service
func NewTileProviderStamenTerrain() TileProvider {
	t := &DefaultTileProvider{
		name:        "stamen-terrain",
		attribution: "Maps (c) Stamen; Data (c) OSM and contributors, ODbL",
		tileSize:    256,
		urlPattern:  "http://%[1]s.tile.stamen.com/terrain/%[2]d/%[3]d/%[4]d.png",
		shards:      []string{"a", "b", "c", "d"},
	}
	return t
}

// NewTileProviderOpenTopoMap creates a TileProvider struct for opentopomap's tile service
func NewTileProviderOpenTopoMap() TileProvider {
	t := &DefaultTileProvider{
		name:        "opentopomap",
		attribution: "Maps (c) OpenTopoMap [CC-BY-SA]; Data (c) OSM and contributors [ODbL]; Data (c) SRTM",
		tileSize:    256,
		urlPattern:  "http://%[1]s.tile.opentopomap.org/%[2]d/%[3]d/%[4]d.png",
		shards:      []string{"a", "b", "c"},
	}
	return t
}

// NewTileProviderWikimedia creates a TileProvider struct for Wikimedia's tile service
func NewTileProviderWikimedia() TileProvider {
	t := &DefaultTileProvider{
		name:        "wikimedia",
		attribution: "Map (c) Wikimedia; Data (c) OSM and contributors, ODbL.",
		tileSize:    256,
		urlPattern:  "https://maps.wikimedia.org/osm-intl/%[2]d/%[3]d/%[4]d.png",
		shards:      []string{},
	}
	return t
}

// NewTileProviderOpenCycleMap creates a TileProvider struct for OpenCycleMap's tile service
func NewTileProviderOpenCycleMap() TileProvider {
	t := &DefaultTileProvider{
		name:        "cycle",
		attribution: "Maps and Data (c) openstreetmaps.org and contributors, ODbL",
		tileSize:    256,
		urlPattern:  "http://%[1]s.tile.opencyclemap.org/cycle/%[2]d/%[3]d/%[4]d.png",
		shards:      []string{"a", "b"},
	}
	return t
}

func newTileProviderCarto(name string) TileProvider {
	t := &DefaultTileProvider{
		name:        fmt.Sprintf("carto-%s", name),
		attribution: "Map (c) Carto [CC BY 3.0] Data (c) OSM and contributors, ODbL.",
		tileSize:    256,
		urlPattern:  "https://cartodb-basemaps-%[1]s.global.ssl.fastly.net/" + name + "_all/%[2]d/%[3]d/%[4]d.png",
		shards:      []string{"a", "b", "c", "d"},
	}
	return t
}

// NewTileProviderCartoLight creates a TileProvider struct for Carto's tile service (light variant)
func NewTileProviderCartoLight() TileProvider {
	return newTileProviderCarto("light")
}

// NewTileProviderCartoDark creates a TileProvider struct for Carto's tile service (dark variant)
func NewTileProviderCartoDark() TileProvider {
	return newTileProviderCarto("dark")
}

// GetTileProviders returns a map of all available TileProviders
func GetTileProviders() map[string]TileProvider {
	m := make(map[string]TileProvider)

	list := []TileProvider{
		NewTileProviderOpenStreetMaps(),
		NewTileProviderOpenCycleMap(),
		NewTileProviderThunderforestLandscape(),
		NewTileProviderThunderforestOutdoors(),
		NewTileProviderThunderforestTransport(),
		NewTileProviderStamenToner(),
		NewTileProviderStamenTerrain(),
		NewTileProviderOpenTopoMap(),
		NewTileProviderOpenStreetMaps(),
		NewTileProviderOpenCycleMap(),
		NewTileProviderCartoLight(),
		NewTileProviderCartoDark(),
	}

	for _, tp := range list {
		m[tp.Name()] = tp
	}

	return m
}
