package samples

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// CreatePositiveVectorFile
// opencv_createsamples -info {infoPath} -num {numSamples} -w {detectWidth} -h {detectHeight} -vec {posVecFile}
func CreatePositiveVectorFile(numSamples int, posVecFile string, detectHeight, detectWidth int) error {
	mode := int64(0777)
	if _, err := os.Stat(infoFolder); os.IsNotExist(err) {
		os.MkdirAll(infoFolder, os.FileMode(mode))
	}
	cmdName, err := exec.LookPath(cvCreateSamplesCmd)
	if err != nil {
		return err
	}
	cmdArgs := []string{
		"-info", infoPath,
		"-num", strconv.Itoa(numSamples),
		"-w", strconv.Itoa(detectWidth),
		"-h", strconv.Itoa(detectHeight),
		"-vec", posVecFile,
	}
	fmt.Println(cvCreateSamplesCmd, strings.Join(cmdArgs, " "))
	_, err = exec.Command(cmdName, cmdArgs...).CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
