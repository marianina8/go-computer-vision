package gophercon

import (
	"fmt"

	yolo "../../github.com/ZanLabs/go-yolo"
)

func main() {
	fmt.Println("Hi")
	// not save
	yolo.VideoDetector(
		"./cfg/coco.data", // datacfg
		"./cfg/yolo.cfg",  // cfgfile
		"./yolo.weights",  // weightfile
		"/path/video.mp4", // video that you want recognize
		0.24,              // thresh default: 0.24
		0.5)               // hierThresh default: 0.5

	// save video
	yolo.VideoDetector(
		"./cfg/coco.data", // datacfg
		"./cfg/yolo.cfg",  // cfgfile
		"./yolo.weights",  // weightfile
		"/path/video.mp4", // video that you want recognize
		0.24,              // thresh default: 0.24
		0.5,               // hierThresh default: 0.5
		"/PATHTO/A_VIDEO_NAME") // ignore the suffix
}
