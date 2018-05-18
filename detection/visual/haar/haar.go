package haar

import (
	"encoding/json"
	"fmt"
	"image"

	"github.com/marianina8/go-computer-vision/detection/visual"
	"gocv.io/x/gocv"
)

type haarDetection struct {
	cascadeFrontalFace string
}

// New makes a new detector using haar cascades.
func New() (visual.Eyes, error) {
	return &haarDetection{
		cascadeFrontalFace: "./visual/haar/cascades/haarcascade_frontalface_default.xml",
	}, nil
}

func (h haarDetection) Detect(file visual.ImageInfo) visual.Objects {
	if file.Path != "" {
		return detectFromFile("face", h.cascadeFrontalFace, file.Path)
	}
	return detectFromMat("face", h.cascadeFrontalFace, file.Mat)
}

func detectFromMat(objType string, classifierXML string, img gocv.Mat) visual.Objects {
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()
	classifier.Load(classifierXML)
	rects := classifier.DetectMultiScale(img)
	return toObjects(objType, rects)
}

func detectFromFile(objType, classifierXML, filepath string) visual.Objects {
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()
	classifier.Load(classifierXML)
	img := gocv.IMRead(filepath, gocv.IMReadColor)
	if img.Empty() {
		fmt.Println("Error reading image from:", filepath)
		return visual.Objects{}
	}
	rects := classifier.DetectMultiScale(img)
	return toObjects(objType, rects)
}

func toObjects(objType string, rects []image.Rectangle) visual.Objects {
	objs := visual.Objects{}
	for _, r := range rects {
		detection := visual.Detection{
			Rectangle: visual.Rectangle{
				Top:    r.Max.Y,
				Left:   r.Min.X,
				Height: r.Max.Y - r.Min.Y,
				Width:  r.Max.X - r.Min.X,
			},
		}
		f, _ := json.Marshal(detection.Rectangle)
		fmt.Println("Detect:", string(f))
		obj := visual.Object{}
		obj[objType] = detection
		objs = append(objs, obj)
	}
	return objs
}
