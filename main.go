package main

import (
	"flag"
	"fmt"
	"image"
	"io"
	"os"
	"strconv"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/tiff"
)

type simage struct {
	i             image.Image
	width, height int
	enc           imaging.Format
	filter        imaging.ResampleFilter
}

func newsimage(img image.Image, enc string, w, h int, f imaging.ResampleFilter) simage {
	var e imaging.Format
	switch enc {
	case "jpeg":
		e = imaging.JPEG
	case "png":
		e = imaging.PNG
	case "gif":
		e = imaging.GIF
	case "tiff":
		e = imaging.TIFF
	default:
		e = imaging.JPEG
	}
	return simage{
		i:      img,
		width:  w,
		height: h,
		enc:    e,
		filter: f,
	}
}

func simg(w io.Writer, img simage) error {
	resized := imaging.Resize(img.i, img.width, img.height, img.filter)
	err := imaging.Encode(w, resized, img.enc)
	if err != nil {
		return err
	}
	return nil
}

var write = flag.Bool("w", false, "write to source instead of stdout")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [ option ... ] dimension [ file ... ]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.Usage()
	}
	w, h, err := parseDementions(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "simg: %v\n", err)
		os.Exit(1)
	}
	if len(flag.Args()) >= 2 {
		for _, v := range flag.Args()[1:] {
			r, err := os.OpenFile(v, os.O_RDWR, 0)
			if err != nil {
				fmt.Fprintf(os.Stderr, "simg: %v\n", err)
				os.Exit(1)
			}
			img, f, err := image.Decode(r)
			if err != nil {
				fmt.Fprintf(os.Stderr, "simg: %v\n", err)
				os.Exit(1)
			}
			var s = newsimage(img, f, w, h, imaging.Lanczos)
			var dest io.Writer
			if *write {
				dest = r
			} else {
				dest = os.Stdout
			}
			err = simg(dest, s)
			if err != nil {
				fmt.Fprintf(os.Stderr, "simg: %v\n", err)
				os.Exit(1)
			}
			r.Close()
		}
	} else {
		if *write {
			// filename is "<stdin>"
			fmt.Fprintf(os.Stderr, "simg: can't use -w on stdin\n")
			os.Exit(1)
		}
		// TODO(wgr): implement a reader that understands multiple images comming
		// in from a data stream.
		img, f, err := image.Decode(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "simg: %v\n", err)
			os.Exit(1)
		}
		s := newsimage(img, f, w, h, imaging.Lanczos)
		err = simg(os.Stdout, s)
	}
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
