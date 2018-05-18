package samples

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	cvCreateSamplesCmd = "opencv_createsamples"
	infoFolder         = "info"
	infoFile           = "info.lst"
	infoPath           = infoFolder + "/" + infoFile
)

// CreateSamples
// opencv_createsamples -img {posFile} -bg {bgFile} -info {infoPath} -pngoutput {infoFolder} -maxxangle {maxAngle} -maxyangle {maxAngle} -maxzangle {maxAngle} -num {numSamples}
func CreateSamples(posFile string, bgFile string, numSamples int, maxAngle float64) error {
	mode := int64(0777)
	if _, err := os.Stat(infoFolder); os.IsNotExist(err) {
		os.MkdirAll(infoFolder, os.FileMode(mode))
	}
	cmdName, err := exec.LookPath(cvCreateSamplesCmd)
	if err != nil {
		return err
	}
	cmdArgs := []string{
		"-img", posFile,
		"-bg", bgFile,
		"-info", infoPath,
		// only run this to view the created sample png files
		// "-pngoutput", infoFolder,
		"-bgcolor 255 -bgthresh 8",
		"-maxxangle", fmt.Sprintf("%.1f", maxAngle),
		"-maxyangle", fmt.Sprintf("-%.1f", maxAngle),
		"-maxzangle", fmt.Sprintf("%.1f", maxAngle),
		"-num", strconv.Itoa(numSamples),
	}
	fmt.Println(cvCreateSamplesCmd, strings.Join(cmdArgs, " "))
	_, err = exec.Command(cmdName, cmdArgs...).CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
