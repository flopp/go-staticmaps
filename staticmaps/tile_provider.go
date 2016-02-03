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
	t.Attribution = "Maps © MapQuest; Data © OSM and contributors, ODbL"
	t.TileSize = 256
	t.URLPattern = "http://otile%[1]s.mqcdn.com/tiles/1.0.0/osm/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"1", "2", "3", "4"}
	return t
}

func newTileProviderThunderforest(name string) *TileProvider {
	t := new(TileProvider)
	t.Name = fmt.Sprintf("thundeforest-%s", name)
	t.Attribution = "Maps © Thundeforest; Data © OSM and contributors, ODbL"
	t.TileSize = 256
	t.URLPattern = "https://%[1]s.tile.thunderforest.com/landscape/%[2]d/%[3]d/%[4]d.png"
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

func GetTileProviders() map[string]*TileProvider {
	m := make(map[string]*TileProvider)

	t1 := NewTileProviderMapQuest()
	t2 := NewTileProviderThunderforestLandscape()
	t3 := NewTileProviderThunderforestOutdoors()
	t4 := NewTileProviderThunderforestTransport()

	m[t1.Name] = t1
	m[t2.Name] = t2
	m[t3.Name] = t3
	m[t4.Name] = t4

	return m
}
