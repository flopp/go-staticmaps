// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package staticmaps

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type TileFetcher struct {
	url_schema string
	cache_dir  string
}

func NewTileFetcher(url_schema, cache_dir string) *TileFetcher {
	t := new(TileFetcher)
	t.url_schema = url_schema
	t.cache_dir = cache_dir
	return t
}

func (t *TileFetcher) url(zoom uint, x, y int) string {
	return fmt.Sprintf(t.url_schema, zoom, x, y)
}

func (t *TileFetcher) cache_file_name(zoom uint, x, y int) string {
	return fmt.Sprintf("%s/%d-%d-%d", t.cache_dir, zoom, x, y)
}

func (t *TileFetcher) Fetch(zoom uint, x, y int) (image.Image, error) {
	file_name := t.cache_file_name(zoom, x, y)
	cached_img, err := t.load_cache(file_name)
	if err == nil {
		return cached_img, nil
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

	err = t.store_cache(file_name, data)
	if err != nil {
		fmt.Println("Failed to store image as", file_name)
        fmt.Println(err)
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

func (t *TileFetcher) load_cache(file_name string) (image.Image, error) {
	file, err := os.Open(file_name)
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

func (t *TileFetcher) create_cache_dir() error {
	src, err := os.Stat(t.cache_dir)
	if err != nil {
        if os.IsNotExist(err) {
            return os.Mkdir(t.cache_dir, 0777)
        } else {
            return err
        }
	}
	if src.IsDir() {
		return nil
	}
	return errors.New(fmt.Sprintf("File exists but is not a directory: %s", t.cache_dir))
}

func (t *TileFetcher) store_cache(file_name string, data []byte) error {
	err := t.create_cache_dir()
	if err != nil {
		return err
	}

	file, err := os.Create(file_name)
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
