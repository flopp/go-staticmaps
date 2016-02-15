// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package staticmaps

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg" // to be able to decode jpegs
	_ "image/png"  // to be able to decode pngs
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Wessie/appdirs"
)

type TileFetcher struct {
	tileProvider *TileProvider
	cacheDir     string
	useCaching   bool
}

func NewTileFetcher(tileProvider *TileProvider) *TileFetcher {
	t := new(TileFetcher)
	t.tileProvider = tileProvider
	app := appdirs.New("go-staticmaps", "flopp.net", "0.1")
	t.cacheDir = fmt.Sprintf("%s/%s", app.UserCache(), tileProvider.Name)
	t.useCaching = true
	return t
}

func (t *TileFetcher) url(zoom, x, y int) string {
	shard := ""
	ss := len(t.tileProvider.Shards)
	if len(t.tileProvider.Shards) > 0 {
		shard = t.tileProvider.Shards[(x+y)%ss]
	}
	return t.tileProvider.GetURL(shard, zoom, x, y)
}

func (t *TileFetcher) cacheFileName(zoom int, x, y int) string {
	return fmt.Sprintf("%s/%d-%d-%d", t.cacheDir, zoom, x, y)
}

func (t *TileFetcher) ToggleCaching(enabled bool) {
	t.useCaching = enabled
}

func (t *TileFetcher) Fetch(zoom, x, y int) (image.Image, error) {
	if t.useCaching {
		fileName := t.cacheFileName(zoom, x, y)
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

	if t.useCaching {
		fileName := t.cacheFileName(zoom, x, y)
		err = t.storeCache(fileName, data)
		if err != nil {
			fmt.Println("Failed to store image as", fileName)
			fmt.Println(err)
		}
	}

	return img, nil
}

func (t *TileFetcher) download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

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

func (t *TileFetcher) createCacheDir() error {
	src, err := os.Stat(t.cacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(t.cacheDir, 0777)
		}
		return err
	}
	if src.IsDir() {
		return nil
	}
	return fmt.Errorf("File exists but is not a directory: %s", t.cacheDir)
}

func (t *TileFetcher) storeCache(fileName string, data []byte) error {
	err := t.createCacheDir()
	if err != nil {
		return err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	return nil
}
