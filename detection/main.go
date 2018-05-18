package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"

	"github.com/marianina8/go-computer-vision/detection/visual"
	"github.com/marianina8/go-computer-vision/detection/visual/haar"
	"github.com/marianina8/go-computer-vision/detection/visual/microsoft"
	"gocv.io/x/gocv"
)

func main() {
	inputImage := flag.String("img", "", "image to check for objects")
	inputVideo := flag.String("vid", "", "video to check for objects")
	webDeviceID := flag.Int("webcam", 0, "check webcam stream")
	detection := flag.String("detection", "haar", "type of detection (options: haar, azure)")
	flag.Parse()

	var detector visual.Eyes
	var err error

	switch *detection {
	case "haar":
		detector, _ = haar.New()
	case "azure":
		detector, err = microsoft.New()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if *inputImage != "" {
		displayImage(detector, *inputImage)
		return
	}
	if *inputVideo != "" {
		displayVideo(detector, *inputVideo)
		return
	}
	displayWebStream(detector, *webDeviceID)
}

func displayVideo(detector visual.Eyes, inputVideo string) {
	stream, err := gocv.VideoCaptureFile(inputVideo)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stream.Close()

	// open display window
	window := gocv.NewWindow("Detect")
	defer window.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()
	for {
		if ok := stream.Read(&img); !ok {
			fmt.Printf("cannot read device\n")
			return
		}
		if img.Empty() {
			continue
		}
		b, _ := gocv.IMEncode(".jpg", img)
		info := visual.ImageInfo{
			Mat:   img,
			Bytes: &b,
		}
		objects := detector.Detect(info)
		blue := color.RGBA{0, 0, 255, 0}
		for _, objects := range objects {
			for objType, detection := range objects {
				r := toRectangle(detection)
				pt := image.Pt(r.Min.X, r.Max.Y-3)
				gocv.PutText(&img, objType, pt, gocv.FontHersheyPlain, 1, blue, 2)
				gocv.Rectangle(&img, r, blue, 2)
			}
		}
		// show the image in the window, and wait 1 millisecond
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}

func displayWebStream(detector visual.Eyes, deviceID int) {
	// open webcam
	webcam, err := gocv.VideoCaptureDevice(int(deviceID))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	// open display window
	window := gocv.NewWindow("Detect")
	defer window.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	fmt.Printf("start reading camera device: %v\n", deviceID)
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %d\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}
		b, _ := gocv.IMEncode(".jpg", img)
		info := visual.ImageInfo{
			Mat:   img,
			Bytes: &b,
		}
		objects := detector.Detect(info)
		blue := color.RGBA{0, 0, 255, 0}
		for _, objects := range objects {
			for objType, detection := range objects {
				r := toRectangle(detection)
				pt := image.Pt(r.Min.X, r.Max.Y-3)
				gocv.PutText(&img, objType, pt, gocv.FontHersheyPlain, 1, blue, 2)
				gocv.Rectangle(&img, r, blue, 2)
			}
		}
		// show the image in the window, and wait 1 millisecond
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}

func displayImage(detector visual.Eyes, inputImage string) {
	window := gocv.NewWindow("Detect")
	defer window.Close()

	img := gocv.IMRead(inputImage, gocv.IMReadColor)
	defer img.Close()

	info := visual.ImageInfo{
		Path: inputImage,
	}
	objects := detector.Detect(info)
	blue := color.RGBA{0, 0, 255, 0}
	for _, objects := range objects {
		for objType, detection := range objects {
			r := toRectangle(detection)
			pt := image.Pt(r.Min.X, r.Max.Y-3)
			gocv.PutText(&img, objType, pt, gocv.FontHersheyPlain, 1, blue, 2)
			gocv.Rectangle(&img, r, blue, 2)
		}
	}
	for {
		window.IMShow(img)
		if window.WaitKey(10) >= 0 {
			break
		}
	}
}

func toRectangle(detection visual.Detection) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: detection.Rectangle.Left,
			Y: detection.Rectangle.Top - detection.Rectangle.Height,
		},
		Max: image.Point{
			X: detection.Rectangle.Left + detection.Rectangle.Width,
			Y: detection.Rectangle.Top,
		},
	}
}
