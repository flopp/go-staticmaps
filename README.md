# go-staticmaps
Render static map images with go

## Installation

    go get -u github.com/flopp/go-staticmaps
    go get -u github.com/cheggaaa/pb
    go get -u github.com/llgcode/draw2d/draw2dimg
    go get -u github.com/jessevdk/go-flags

## Usage

Create a static map image "map1.png" with size 800x600, centered on "N 48 E 7.8" with zoom level 14:

    cd $GOPATH/src/github.com/flopp/go-staticmaps
    go run cmd/create-static-map.go --center "48,7.8" --zoom 14 --width 800 --height 600 --output "map1.png"

Add some markers (one red, two green)...

    go run cmd/create-static-map.go --center "48,7.8" --zoom 14 --width 800 --height 600 --marker "color:red|48,7.8" --marker "color:green|47.99,7.8|48.01,7.8" --output "map2.png"



## License
Copyright 2016 Florian Pigorsch. All rights reserved.

Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
