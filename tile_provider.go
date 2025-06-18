// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sm

import "fmt"

// TileProvider encapsulates all infos about a map tile provider service (name, url scheme, attribution, etc.)
type TileProvider struct {
	Name           string
	Attribution    string
	IgnoreNotFound bool
	TileSize       int
	URLPattern     string // "%[1]s" => shard, "%[2]d" => zoom, "%[3]d" => x, "%[4]d" => y, "%[5]s" => API key
	Shards         []string
	APIKey         string
}

// IsNone returns true if t is an empyt TileProvider (e.g. no configured Url)
func (t TileProvider) IsNone() bool {
	return len(t.URLPattern) == 0
}

func (t *TileProvider) getURL(shard string, zoom, x, y int, apikey string) string {
	if t.IsNone() {
		return ""
	}
	return fmt.Sprintf(t.URLPattern, shard, zoom, x, y, apikey)
}

// NewTileProviderOpenStreetMaps creates a TileProvider struct for OSM's tile service
func NewTileProviderOpenStreetMaps() *TileProvider {
	t := new(TileProvider)
	t.Name = "osm"
	t.Attribution = "Maps and Data (c) openstreetmap.org and contributors, ODbL"
	t.TileSize = 256
	t.URLPattern = "https://%[1]s.tile.openstreetmap.org/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"a", "b", "c"}
	return t
}

func newTileProviderThunderforest(name string, apikey string) *TileProvider {
	t := new(TileProvider)
	t.Name = fmt.Sprintf("thunderforest-%s", name)
	t.Attribution = "Maps (c) Thundeforest; Data (c) OSM and contributors, ODbL"
	t.TileSize = 256
	t.APIKey = apikey
	t.URLPattern = "https://%[1]s.tile.thunderforest.com/" + name + "/%[2]d/%[3]d/%[4]d.png?apikey=%[5]s"
	t.Shards = []string{"a", "b", "c"}
	return t
}

// NewTileProviderThunderforestLandscape creates a TileProvider struct for thundeforests's 'landscape' tile service
func NewTileProviderThunderforestLandscape(thunderforestApiKey string) *TileProvider {
	return newTileProviderThunderforest("landscape", thunderforestApiKey)
}

// NewTileProviderThunderforestOutdoors creates a TileProvider struct for thundeforests's 'outdoors' tile service
func NewTileProviderThunderforestOutdoors(thunderforestApiKey string) *TileProvider {
	return newTileProviderThunderforest("outdoors", thunderforestApiKey)
}

// NewTileProviderThunderforestTransport creates a TileProvider struct for thundeforests's 'transport' tile service
func NewTileProviderThunderforestTransport(thunderforestApiKey string) *TileProvider {
	return newTileProviderThunderforest("transport", thunderforestApiKey)
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

// NewTileProviderStamenTerrain creates a TileProvider struct for stamens' 'terrain' tile service
func NewTileProviderStamenTerrain() *TileProvider {
	t := new(TileProvider)
	t.Name = "stamen-terrain"
	t.Attribution = "Maps (c) Stamen; Data (c) OSM and contributors, ODbL"
	t.TileSize = 256
	t.URLPattern = "http://%[1]s.tile.stamen.com/terrain/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"a", "b", "c", "d"}
	return t
}

// NewTileProviderOpenTopoMap creates a TileProvider struct for opentopomap's tile service
func NewTileProviderOpenTopoMap() *TileProvider {
	t := new(TileProvider)
	t.Name = "opentopomap"
	t.Attribution = "Maps (c) OpenTopoMap [CC-BY-SA]; Data (c) OSM and contributors [ODbL]; Data (c) SRTM"
	t.TileSize = 256
	t.URLPattern = "http://%[1]s.tile.opentopomap.org/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"a", "b", "c"}
	return t
}

// NewTileProviderWikimedia creates a TileProvider struct for Wikimedia's tile service
func NewTileProviderWikimedia() *TileProvider {
	t := new(TileProvider)
	t.Name = "wikimedia"
	t.Attribution = "Map (c) Wikimedia; Data (c) OSM and contributors, ODbL."
	t.TileSize = 256
	t.URLPattern = "https://maps.wikimedia.org/osm-intl/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{}
	return t
}

// NewTileProviderOpenCycleMap creates a TileProvider struct for OpenCycleMap's tile service
func NewTileProviderOpenCycleMap() *TileProvider {
	t := new(TileProvider)
	t.Name = "cycle"
	t.Attribution = "Maps and Data (c) openstreetmaps.org and contributors, ODbL"
	t.TileSize = 256
	t.URLPattern = "http://%[1]s.tile.opencyclemap.org/cycle/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"a", "b"}
	return t
}

// NewTileProviderOpenSeaMap creates a TileProvider struct for OpenSeaMap's tile service
func NewTileProviderOpenSeaMap() *TileProvider {
	t := new(TileProvider)
	t.Name = "sea"
	t.Attribution = "Maps and Data (c) openstreetmaps.org and contributors, ODbL"
	t.TileSize = 256
	t.URLPattern = "http://t1.openseamap.org/seamark/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{}
	return t
}

// NewTileProviderCarto creates a TileProvider struct for Carto's tile service
// See https://github.com/CartoDB/basemap-styles?tab=readme-ov-file#1-web-raster-basemaps for available names
func NewTileProviderCarto(name string) *TileProvider {
	t := new(TileProvider)
	t.Name = fmt.Sprintf("carto-%s", name)
	t.Attribution = "Map (c) Carto [CC BY 3.0] Data (c) OSM and contributors, ODbL."
	t.TileSize = 256
	t.URLPattern = "https://cartodb-basemaps-%[1]s.global.ssl.fastly.net/" + name + "/%[2]d/%[3]d/%[4]d.png"
	t.Shards = []string{"a", "b", "c", "d"}
	return t
}

// NewTileProviderCartoLight creates a TileProvider struct for Carto's tile service (light variant)
func NewTileProviderCartoLight() *TileProvider {
	return NewTileProviderCarto("light_all")
}

// NewTileProviderCartoDark creates a TileProvider struct for Carto's tile service (dark variant)
func NewTileProviderCartoDark() *TileProvider {
	return NewTileProviderCarto("dark_all")
}

// NewTileProviderArcgisWorldImagery creates a TileProvider struct for Arcgis' WorldImagery tiles
func NewTileProviderArcgisWorldImagery() *TileProvider {
	t := new(TileProvider)
	t.Name = "arcgis-worldimagery"
	t.Attribution = "Source: Esri, Maxar, GeoEye, Earthstar Geographics, CNES/Airbus DS, USDA, USGS, AeroGRID, IGN, and the GIS User Community"
	t.TileSize = 256
	t.URLPattern = "https://server.arcgisonline.com/arcgis/rest/services/World_Imagery/MapServer/tile/%[2]d/%[4]d/%[3]d"
	t.Shards = []string{}
	return t
}

// NewTileProviderNone creates a TileProvider struct that does not provide any tiles
func NewTileProviderNone() *TileProvider {
	t := new(TileProvider)
	t.Name = "none"
	t.Attribution = ""
	t.TileSize = 256
	t.URLPattern = ""
	t.Shards = []string{}
	return t
}

// GetTileProviders returns a map of all available TileProviders
func GetTileProviders(thunderforestApiKey string) map[string]*TileProvider {
	m := make(map[string]*TileProvider)

	list := []*TileProvider{
		NewTileProviderThunderforestLandscape(thunderforestApiKey),
		NewTileProviderThunderforestOutdoors(thunderforestApiKey),
		NewTileProviderThunderforestTransport(thunderforestApiKey),
		NewTileProviderStamenToner(),
		NewTileProviderStamenTerrain(),
		NewTileProviderOpenTopoMap(),
		NewTileProviderOpenStreetMaps(),
		NewTileProviderOpenCycleMap(),
		NewTileProviderOpenSeaMap(),
		NewTileProviderCarto("rastertiles/voyager"),
		NewTileProviderCartoLight(),
		NewTileProviderCartoDark(),
		NewTileProviderArcgisWorldImagery(),
		NewTileProviderWikimedia(),
		NewTileProviderNone(),
	}

	for _, tp := range list {
		m[tp.Name] = tp
	}

	return m
}
