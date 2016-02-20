// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package staticmaps

import "fmt"

type TileProvider struct {
	Name        string
	Attribution string
	TileSize    int
	URLPattern  string // "%[1]s" => shard, "%[2]d" => zoom, "%[3]d" => x, "%[4]d" => y
	Shards      []string
}

func (t *TileProvider) GetURL(shard string, zoom, x, y int) string {
	return fmt.Sprintf(t.URLPattern, shard, zoom, x, y)
}

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

func NewTileProviderThunderforestLandscape() *TileProvider {
	return newTileProviderThunderforest("landscape")
}

func NewTileProviderThunderforestOutdoors() *TileProvider {
	return newTileProviderThunderforest("outdoors")
}

func NewTileProviderThunderforestTransport() *TileProvider {
	return newTileProviderThunderforest("transport")
}

func NewTileProviderStamenToner() *TileProvider {
	t := new(TileProvider)
	t.Name = "stamen-toner"
	t.Attribution = "Maps (c) Stamen; Data (c) OSM and contributors, ODbL"
	t.TileSize = 256
	t.URLPattern = "http://%[1]s.tile.stamen.com/toner/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"a", "b", "c", "d"}
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
		NewTileProviderStamenToner()}

	for _, tp := range list {
		m[tp.Name] = tp
	}

	return m
}
