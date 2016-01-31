# go-staticmaps
Render static map images with go

## Installation

    go get -u github.com/flopp/go-staticmaps
    go get -u github.com/cheggaaa/pb
    go get -u github.com/llgcode/draw2d/draw2dimg

## Usage

Create a static map image "staticmap.png" with size 800x600, centered on "N 48 E 7.8" with zoom level 14:

    cd $GOPATH/src/github.com/flopp/go-staticmaps
    go run main/main.go -center "48,7.8" -zoom 14 -width 800 -height 600 -output "staticmap.png"


## License
Copyright 2016 Florian Pigorsch. All rights reserved.

Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
