// Copyright 2016 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package staticmaps

import (
	"fmt"
	"image/color"
	"regexp"
	"strconv"
)

func ParseColorString(s string) (*color.RGBA, error) {
	re := regexp.MustCompile(`^\s*0x([A-Fa-f0-9]{2})([A-Fa-f0-9]{2})([A-Fa-f0-9]{2})\s*$`)
	matches := re.FindStringSubmatch(s)
	if matches != nil {
		r, errr := strconv.ParseInt(matches[1], 16, 8)
		g, errg := strconv.ParseInt(matches[2], 16, 8)
		b, errb := strconv.ParseInt(matches[3], 16, 8)
		if errr != nil || errg != nil || errb != nil {
			return nil, fmt.Errorf("Cannot parse color string: %s", s)
		}
		return &color.RGBA{uint8(r), uint8(g), uint8(b), 0xff}, nil
	}

	re = regexp.MustCompile(`^\s*0x([A-Fa-f0-9]{2})([A-Fa-f0-9]{2})([A-Fa-f0-9]{2})([A-Fa-f0-9]{2})\s*$`)
	matches = re.FindStringSubmatch(s)
	if matches != nil {
		r, errr := strconv.ParseInt(matches[1], 16, 8)
		g, errg := strconv.ParseInt(matches[2], 16, 8)
		b, errb := strconv.ParseInt(matches[3], 16, 8)
		a, erra := strconv.ParseInt(matches[4], 16, 8)
		if errr != nil || errg != nil || errb != nil || erra != nil {
			return nil, fmt.Errorf("Cannot parse color string: %s", s)
		}
		return &color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}, nil
	}

	if s == "black" {
		return &color.RGBA{0x00, 0x00, 0x00, 0xff}, nil
	} else if s == "blue" {
		return &color.RGBA{0x00, 0x00, 0xff, 0xff}, nil
	} else if s == "brown" {
		return &color.RGBA{0x96, 0x4b, 0x00, 0xff}, nil
	} else if s == "green" {
		return &color.RGBA{0x00, 0xff, 0x00, 0xff}, nil
	} else if s == "orange" {
		return &color.RGBA{0xff, 0x7f, 0x00, 0xff}, nil
	} else if s == "purple" {
		return &color.RGBA{0x7f, 0x00, 0x7f, 0xff}, nil
	} else if s == "red" {
		return &color.RGBA{0xff, 0x00, 0, 0xff}, nil
	} else if s == "yellow" {
		return &color.RGBA{0xff, 0xff, 0x00, 0xff}, nil
	} else if s == "white" {
		return &color.RGBA{0xff, 0xff, 0xff, 0xff}, nil
	} else {
		return nil, fmt.Errorf("Cannot parse color string: %s", s)
	}
}
