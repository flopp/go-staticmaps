// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package staticmaps

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

type LatLng struct {
	lat float64
	lng float64
}

func LatLngFromDegrees(lat, lng float64) LatLng {
	return LatLng{lat, lng}
}

func LatLngFromRadians(lat, lng float64) LatLng {
	return LatLng{lat * 180.0 / math.Pi, lng * 180.0 / math.Pi}
}

func LatLngFromString(s string) (LatLng, error) {
	re := regexp.MustCompile(`^\s*([+-]?\d+\.?\d*)\s*,\s*([+-]?\d+\.?\d*)\s*$`)

	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return LatLng{}, fmt.Errorf("Cannot parse lat,lng string: %s", s)
	}

	lat, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return LatLng{}, fmt.Errorf("Cannot parse lat,lng string: %s", s)
	}

	lng, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return LatLng{}, fmt.Errorf("Cannot parse lat,lng string: %s", s)
	}

	return LatLngFromDegrees(lat, lng), nil
}

func (ll *LatLng) Lat() float64 {
	return ll.lat
}

func (ll *LatLng) LatRadians() float64 {
	return ll.Lat() * math.Pi / 180.0
}

func (ll *LatLng) Lng() float64 {
	return ll.lng
}

func (ll *LatLng) LngRadians() float64 {
	return ll.Lng() * math.Pi / 180.0
}
