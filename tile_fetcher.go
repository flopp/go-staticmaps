// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sm

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg" // to be able to decode jpegs
	_ "image/png"  // to be able to decode pngs
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

var errTileNotFound = errors.New("error 404: tile not found")

// TileFetcher downloads map tile images from a TileProvider
type TileFetcher struct {
	tileProvider *TileProvider
	cache        TileCache
	userAgent    string
}

// NewTileFetcher creates a new Tilefetcher struct
func NewTileFetcher(tileProvider *TileProvider, cache TileCache) *TileFetcher {
	t := new(TileFetcher)
	t.tileProvider = tileProvider
	t.cache = cache
	t.userAgent = "Mozilla/5.0+(compatible; go-staticmaps/0.1; https://github.com/flopp/go-staticmaps)"
	return t
}

// SetUserAgent sets the HTTP user agent string used when downloading map tiles
func (t *TileFetcher) SetUserAgent(a string) {
	t.userAgent = a
}

func (t *TileFetcher) url(zoom, x, y int) string {
	shard := ""
	ss := len(t.tileProvider.Shards)
	if len(t.tileProvider.Shards) > 0 {
		shard = t.tileProvider.Shards[(x+y)%ss]
	}
	return t.tileProvider.getURL(shard, zoom, x, y)
}

func cacheFileName(cache TileCache, providerName string, zoom, x, y int) string {
	return path.Join(
		cache.Path(),
		providerName,
		strconv.Itoa(zoom),
		strconv.Itoa(x),
		strconv.Itoa(y),
	)
}

// Fetch download (or retrieves from the cache) a tile image for the specified zoom level and tile coordinates
func (t *TileFetcher) Fetch(zoom, x, y int) (image.Image, error) {
	if t.cache != nil {
		fileName := cacheFileName(t.cache, t.tileProvider.Name, zoom, x, y)
		cachedImg, err := t.loadCache(fileName)
		if err == nil {
			return cachedImg, nil
		}
	}

	url := t.url(zoom, x, y)
	data, err := t.download(url)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	if t.cache != nil {
		fileName := cacheFileName(t.cache, t.tileProvider.Name, zoom, x, y)
		if err := t.storeCache(fileName, data); err != nil {
			log.Printf("Failed to store map tile as '%s': %s", fileName, err)
		}
	}

	return img, nil
}

func (t *TileFetcher) download(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", t.userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		// Great! Nothing to do.

	case http.StatusNotFound:
		return nil, errTileNotFound

	default:
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (t *TileFetcher) loadCache(fileName string) (image.Image, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (t *TileFetcher) createCacheDir(path string) error {
	src, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, t.cache.Perm())
		}
		return err
	}
	if src.IsDir() {
		return nil
	}

	return fmt.Errorf("file exists but is not a directory: %s", path)
}

func (t *TileFetcher) storeCache(fileName string, data []byte) error {
	dir, _ := filepath.Split(fileName)

	if err := t.createCacheDir(dir); err != nil {
		return err
	}

	// Create file using the configured directory create permission with the
	// 'x' bit removed.
	file, err := os.OpenFile(
		fileName,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		t.cache.Perm()&0666,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, bytes.NewBuffer(data)); err != nil {
		return err
	}

	return nil
}
