// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sm

import "fmt"

// TileProvider encapsulates all infos about a map tile provider service (name, url scheme, attribution, etc.)
type TileProvider struct {
	Name        string
	Attribution string
	TileSize    int
	URLPattern  string // "%[1]s" => shard, "%[2]d" => zoom, "%[3]d" => x, "%[4]d" => y
	Shards      []string
}

func (t *TileProvider) getURL(shard string, zoom, x, y int) string {
	return fmt.Sprintf(t.URLPattern, shard, zoom, x, y)
}

// NewTileProviderMapQuest creates a TileProvider struct for mapquest's tile service
func NewTileProviderMapQuest() *TileProvider {
	t := new(TileProvider)
	t.Name = "mapquest"
	t.Attribution = "Maps (c) MapQuest; Data (c) OSM and contributors, ODbL"
	t.TileSize = 256
	t.URLPattern = "http://otile%[1]s.mqcdn.com/tiles/1.0.0/osm/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"1", "2", "3", "4"}
	return t
}

func newTileProviderThunderforest(name string) *TileProvider {
	t := new(TileProvider)
	t.Name = fmt.Sprintf("thunderforest-%s", name)
	t.Attribution = "Maps (c) Thundeforest; Data (c) OSM and contributors, ODbL"
	t.TileSize = 256
	t.URLPattern = "https://%[1]s.tile.thunderforest.com/" + name + "/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"a", "b", "c"}
	return t
}

// NewTileProviderThunderforestLandscape creates a TileProvider struct for thundeforests's 'landscape' tile service
func NewTileProviderThunderforestLandscape() *TileProvider {
	return newTileProviderThunderforest("landscape")
}

// NewTileProviderThunderforestOutdoors creates a TileProvider struct for thundeforests's 'outdoors' tile service
func NewTileProviderThunderforestOutdoors() *TileProvider {
	return newTileProviderThunderforest("outdoors")
}

// NewTileProviderThunderforestTransport creates a TileProvider struct for thundeforests's 'transport' tile service
func NewTileProviderThunderforestTransport() *TileProvider {
	return newTileProviderThunderforest("transport")
}

// NewTileProviderStamenToner creates a TileProvider struct for stamens' 'toner' tile service
func NewTileProviderStamenToner() *TileProvider {
	t := new(TileProvider)
	t.Name = "stamen-toner"
	t.Attribution = "Maps (c) Stamen; Data (c) OSM and contributors, ODbL"
	t.TileSize = 256
	t.URLPattern = "http://%[1]s.tile.stamen.com/toner/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"a", "b", "c", "d"}
	return t
}

// NewTileProviderOpenTopoMap creates a TileProvider struct for opentopomaps's tile service
func NewTileProviderOpenTopoMap() *TileProvider {
	t := new(TileProvider)
	t.Name = "opentopomap"
	t.Attribution = "Maps (c) OpenTopoMap [CC-BY-SA]; Data (c) OSM and contributors [ODbL]; Data (c) SRTM"
	t.TileSize = 256
	t.URLPattern = "http://%[1].tile.opentopomap.org/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"a", "b", "c"}
	return t
}

// GetTileProviders returns a map of all available TileProviders
func GetTileProviders() map[string]*TileProvider {
	m := make(map[string]*TileProvider)

	list := []*TileProvider{
		NewTileProviderMapQuest(),
		NewTileProviderThunderforestLandscape(),
		NewTileProviderThunderforestOutdoors(),
		NewTileProviderThunderforestTransport(),
		NewTileProviderStamenToner(),
        NewTileProviderOpenTopoMap()}

	for _, tp := range list {
		m[tp.Name] = tp
	}

	return m
}
