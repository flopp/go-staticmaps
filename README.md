[![GoDoc](https://godoc.org/github.com/flopp/go-staticmaps?status.svg)](https://godoc.org/github.com/flopp/go-staticmaps)
[![Go Report Card](http://goreportcard.com/badge/flopp/go-staticmaps)](http://goreportcard.com/report/flopp/go-staticmaps)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/flopp/go-staticmaps/)

# go-staticmaps
A go (golang) library and command line tool to render static map images using OpenStreetMap tiles.

## What?
go-staticmaps is a golang library that allows you to create nice static map images from OpenStreetMap tiles, along with markers of different size and color, as well as paths and colored areas.

go-staticmaps comes with a command line tool called `create-static-map` for use in shell scripts, etc.

![Static map of the Berlin Marathon](https://raw.githubusercontent.com/flopp/flopp.github.io/master/go-staticmaps/berlin-marathon.png)

## How?

### Installation

Installing go-staticmaps is as easy as

```bash
go get -u github.com/flopp/go-staticmaps
```

### Library Usage

Create a 400x300 pixel map with a red marker:

```go

import (
  "image/color"

  "github.com/flopp/go-staticmaps"
  "github.com/fogleman/gg"
  "github.com/golang/geo/s2"
)

func main() {
  ctx := sm.NewContext()
  ctx.SetSize(400, 300)
  ctx.AddMarker(sm.NewMarker(s2.LatLng{52.514536, 13.350151}, color.RGBA{0xff, 0, 0, 0xff}, 16.0))

  img, err := ctx.Render()
  if err != nil {
    panic(err)
  }

  if err := gg.SavePNG("my-map.png", img); err != nil {
    panic(err)
  }
}
```


See [GoDoc](https://godoc.org/github.com/flopp/go-staticmaps/staticmaps) for a complete documentation and the source code of the [command line tool](https://github.com/flopp/go-staticmaps/blob/master/cmd/create-static-map/create-static-map.go) for an example how to use the package.


### Command Line Usage

    Usage:
      create-static-map [OPTIONS]

    Creates a static map

    Application Options:
          --width=PIXELS       Width of the generated static map image (default: 512)
          --height=PIXELS      Height of the generated static map image (default: 512)
      -o, --output=FILENAME    Output file name (default: map.png)
      -t, --type=MAPTYPE       Select the map type; list possible map types with '--type list'
      -c, --center=LATLNG      Center coordinates (lat,lng) of the static map
      -z, --zoom=ZOOMLEVEL     Zoom factor
      -m, --marker=MARKER      Add a marker to the static map
      -p, --path=PATH          Add a path to the static map
      -a, --area=AREA          Add an area to the static map

    Help Options:
      -h, --help               Show this help message

### General
The command line interface tries to resemble [Google's Static Maps API](https://developers.google.com/maps/documentation/static-maps/intro).
If `--center` or `--zoom` are not given, *good* values are determined from the specified markers and paths.

### Markers
The `--marker` option defines one or more map markers of the same style. Use multiple `--marker` options to add markers of different styles.

    --marker MARKER_STYLES|LATLNG|LATLNG|...

`LATLNG` is a comma separated pair of latitude and longitude, e.g. `52.5153,13.3564`.

`MARKER_STYLES` consists of a set of style descriptors separated by the pipe character `|`:

- `color:COLOR` - where `COLOR` is either of the form `0xRRGGBB`, `0xRRGGBBAA`, or one of `black`, `blue`, `brown`, `green`, `orange`, `purple`, `red`, `yellow`, `white` (default: `red`)
- `size:SIZE` - where `SIZE` is one of `mid`, `small`, `tiny`, or some number > 0 (default: `mid`)
- `label:LABEL` - where `LABEL` is an alpha numeric character, i.e. `A`-`Z`, `a`-`z`, `0`-`9`; (default: no label)

### Paths
The `--path` option defines a path on the map. Use multiple `--path` options to add multiple paths to the map.

    --path PATH_STYLES|LATLNG|LATLNG|...

`PATH_STYLES` consists of a set of style descriptors separated by the pipe character `|`:

- `color:COLOR` - where `COLOR` is either of the form `0xRRGGBB`, `0xRRGGBBAA`, or one of `black`, `blue`, `brown`, `green`, `orange`, `purple`, `red`, `yellow`, `white` (default: `red`)
- `weight:WEIGHT` - where `WEIGHT` is the line width in pixels (defaut: `5`)

### Areas
The `--area` option defines a closed area on the map. Use multiple `--area` options to add multiple areas to the map.

    --area AREA_STYLES|LATLNG|LATLNG|...

`AREA_STYLES` consists of a set of style descriptors separated by the pipe character `|`:

- `color:COLOR` - where `COLOR` is either of the form `0xRRGGBB`, `0xRRGGBBAA`, or one of `black`, `blue`, `brown`, `green`, `orange`, `purple`, `red`, `yellow`, `white` (default: `red`)
- `weight:WEIGHT` - where `WEIGHT` is the line width in pixels (defaut: `5`)
- `fill:COLOR` - where `COLOR` is either of the form `0xRRGGBB`, `0xRRGGBBAA`, or one of `black`, `blue`, `brown`, `green`, `orange`, `purple`, `red`, `yellow`, `white` (default: none)


## Examples

### Basic Maps

Centered at "N 52.514536 E 13.350151" with zoom level 10:

```bash
$ create-static-map --width 600 --height 400 -o map1.png -c "52.514536,13.350151" -z 10
```
![Example 1](https://raw.githubusercontent.com/flopp/flopp.github.io/master/go-staticmaps/map1.png)

A map with a marker at "N 52.514536 E 13.350151" with zoom level 14 (no need to specify the map's center - it is automatically computed from the marker(s)):

```bash
$ create-static-map --width 600 --height 400 -o map2.png -z 14 -m "52.514536,13.350151"
```

![Example 2](https://raw.githubusercontent.com/flopp/flopp.github.io/master/go-staticmaps/map2.png)

A map with two markers (red and green). If there are more than two markers in the map, a *good* zoom level can be determined automatically:

```bash
$ create-static-map --width 600 --height 400 -o map3.png -m "red|52.514536,13.350151" -m "green|52.516285,13.377746"
```

![Example 3](https://raw.githubusercontent.com/flopp/flopp.github.io/master/go-staticmaps/map3.png)




### Create a map of the Berlin Marathon

    create-static-map --width 800 --height 600 \
      --marker "color:green|52.5153,13.3564" \
      --marker "color:red|52.5160,13.3711" \
      --output "berlin-marathon.png" \
      --path "color:blue|weight:2|52.5153,13.3564|52.5146,13.3519|52.5143,13.3511|52.5139,13.3502|\
        52.5139,13.3496|52.5143,13.3484|52.5129,13.3280|52.5128,13.3234|52.5128,13.3230|52.5138,13.3226|\
        52.5146,13.3225|52.5170,13.3244|52.5220,13.3286|52.5223,13.3285|52.5238,13.3297|52.5246,13.3346|\
        52.5223,13.3675|52.5221,13.3685|52.5209,13.3739|52.5217,13.3754|52.5221,13.3764|52.5272,13.3872|\
        52.5294,13.3976|52.5283,13.4114|52.5274,13.4145|52.5249,13.4201|52.5226,13.4176|52.5222,13.4169|\
        52.5206,13.4216|52.5189,13.4277|52.5189,13.4282|52.5188,13.4288|52.5182,13.4289|52.5180,13.4282|\
        52.5142,13.4252|52.5131,13.4238|52.5098,13.4212|52.5110,13.4165|52.5037,13.4104|52.5034,13.4105|\
        52.4992,13.4179|52.4989,13.4178|52.4988,13.4183|52.4955,13.4204|52.4880,13.4251|52.4865,13.4241|\
        52.4874,13.4209|52.4895,13.4065|52.4938,13.3836|52.4935,13.3672|52.4942,13.3626|52.4914,13.3622|\
        52.4910,13.3607|52.4905,13.3602|52.4890,13.3451|52.4857,13.3452|52.4831,13.3451|52.4815,13.3449|\
        52.4787,13.3440|52.4724,13.3361|52.4710,13.3295|52.4715,13.3291|52.4712,13.3283|52.4716,13.3194|\
        52.4706,13.3175|52.4674,13.3088|52.4681,13.3077|52.4677,13.3063|52.4691,13.2979|52.4707,13.2898|\
        52.4707,13.2893|52.4768,13.2811|52.4801,13.2863|52.4802,13.2861|52.4885,13.3021|52.4884,13.3055|\
        52.4905,13.3142|52.4927,13.3111|52.4971,13.3116|52.4995,13.3128|52.5007,13.3132|52.5026,13.3253|\
        52.5045,13.3347|52.5022,13.3420|52.5020,13.3432|52.5001,13.3515|52.4999,13.3539|52.4980,13.3621|\
        52.4998,13.3628|52.5040,13.3664|52.5053,13.3678|52.5084,13.3695|52.5096,13.3763|52.5096,13.3781|\
        52.5107,13.3928|52.5110,13.3968|52.5123,13.3934|52.5159,13.3929|52.5170,13.3907|52.5160,13.3711"

![Static map of the Berlin Marathon](https://raw.githubusercontent.com/flopp/flopp.github.io/master/go-staticmaps/berlin-marathon.png)

### Create a map of the US capitals

    create-static-map --width 800 --height 400 \
      --output "us-capitals.png" \
      --marker "color:blue|size:tiny|32.3754,-86.2996|58.3637,-134.5721|33.4483,-112.0738|34.7244,-92.2789|\
        38.5737,-121.4871|39.7551,-104.9881|41.7665,-72.6732|39.1615,-75.5136|30.4382,-84.2806|33.7545,-84.3897|\
        21.2920,-157.8219|43.6021,-116.2125|39.8018,-89.6533|39.7670,-86.1563|41.5888,-93.6203|39.0474,-95.6815|\
        38.1894,-84.8715|30.4493,-91.1882|44.3294,-69.7323|38.9693,-76.5197|42.3589,-71.0568|42.7336,-84.5466|\
        44.9446,-93.1027|32.3122,-90.1780|38.5698,-92.1941|46.5911,-112.0205|40.8136,-96.7026|39.1501,-119.7519|\
        43.2314,-71.5597|40.2202,-74.7642|35.6816,-105.9381|42.6517,-73.7551|35.7797,-78.6434|46.8084,-100.7694|\
        39.9622,-83.0007|35.4931,-97.4591|44.9370,-123.0272|40.2740,-76.8849|41.8270,-71.4087|34.0007,-81.0353|\
        44.3776,-100.3177|36.1589,-86.7821|30.2687,-97.7452|40.7716,-111.8882|44.2627,-72.5716|37.5408,-77.4339|\
        47.0449,-122.9016|38.3533,-81.6354|43.0632,-89.4007|41.1389,-104.8165"

![Static map of the US capitals](https://raw.githubusercontent.com/flopp/flopp.github.io/master/go-staticmaps/us-capitals.png)

### Create a map of Australia
...where the Northern Territory is highlighted and the capital Canberra is marked.

    create-static-map --width 800 --height 600 \
      --center="-26.284973,134.303764" \
      --output "australia.png" \
      --marker "color:blue|-35.305200,149.121574" \
      --area "color:0x00FF00|fill:0x00FF007F|weight:2|-25.994024,129.013847|-25.994024,137.989677|-16.537670,138.011649|\
        -14.834820,135.385917|-12.293236,137.033866|-11.174554,130.398124|-12.925791,130.167411|-14.866678,129.002860"

![Static map of Australia](https://raw.githubusercontent.com/flopp/flopp.github.io/master/go-staticmaps/australia.png)

## Acknowledgements
Besides the go standard library, go-staticmaps uses

- MapQuest (https://developer.mapquest.com/), Thunderforest (http://www.thunderforest.com/), and Stamen (http://maps.stamen.com/) as map tile providers
- Go Graphics (https://github.com/fogleman/gg) for 2D drawing
- S2 geometry library (https://github.com/golang/geo) for spherical geometry calculations
- appdirs (https://github.com/Wessie/appdirs) for platform specific system directories
- go-coordsparser (https://github.com/flopp/go-coordsparser) for parsing geo coordinates

## License
Copyright 2016 Florian Pigorsch. All rights reserved.

Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
