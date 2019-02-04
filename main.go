package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
)

func simg(img image.Image, width, heigh int) error {
	resized := imaging.Resize(img, width, heigh, imaging.Lanczos)
	err := imaging.Save(resized, "/tmp/t.png", imaging.PNGCompressionLevel(png.BestCompression))
	if err != nil {
		panic(err)
	}
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [ option ... ] dimension [ file ... ]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()
	if len(flag.Args()) < 2 {
		flag.Usage()
	}
	w, h, err := parseDementions(flag.Arg(0))
	fmt.Println(w, h,err)
	return
	if len(flag.Args()) > 2 {

	} else {
		// TODO(wgr): implement a reader that understands multiple images comming
		// in from a data stream.
		img, _, err := image.Decode(os.Stdin)
		_ = img
		_ = err
	}
	r, err := imaging.Open("/tmp/test.png")
	if err != nil {
		panic(err)
	}
	simg(r, w, h)
}

func parseDementions(dim string) (width int, height int, err error) {
	dims := strings.SplitN(dim, "x", 2)
	if len(dims) < 2 {
		dims = strings.SplitN(dim, "X", 2)
		if len(dims) < 2 {
			err = fmt.Errorf("dimention: format should be widthxheight")
			return
		}
	}
	width, err = strconv.Atoi(dims[0])
	if err != nil {
		return
	}
	height, err = strconv.Atoi(dims[1])
	if err != nil {
		return
	}
	return
}
