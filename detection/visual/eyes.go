package visual

import "gocv.io/x/gocv"

type Eyes interface {
	Detect(info ImageInfo) Objects
}

type Objects []Object
type Object map[string]Detection

type Detection struct {
	Confidence float64
	Rectangle  Rectangle
}

type Rectangle struct {
	Top    int
	Left   int
	Height int
	Width  int
}

type ImageInfo struct {
	Path  string
	Bytes *[]byte
	Mat   gocv.Mat
}
