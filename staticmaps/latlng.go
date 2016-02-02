// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package staticmaps

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/golang/geo/s2"
)

func ParseLatLngFromString(s string) (s2.LatLng, error) {
	re := regexp.MustCompile(`^\s*([+-]?\d+\.?\d*)\s*,\s*([+-]?\d+\.?\d*)\s*$`)

	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return s2.LatLng{}, fmt.Errorf("Cannot parse lat,lng string: %s", s)
	}

	lat, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return s2.LatLng{}, fmt.Errorf("Cannot parse lat,lng string: %s", s)
	}

	lng, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return s2.LatLng{}, fmt.Errorf("Cannot parse lat,lng string: %s", s)
	}

	return s2.LatLngFromDegrees(lat, lng), nil
}
