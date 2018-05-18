package training

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	cvTrainCascadeCmd = "opencv_traincascade"
)

// HaarCascade trains a haar cascade file based on the input parameters
// Opencv_traincascade -data {dataFolder} -vec {posVecFile} -bg {bgFile} -numPos {numPositive} -numNeg {numNegative} -numStages {numStages} -w {detectWidth} -h {detectHeight}
func HaarCascade(dataFolder string, posVecFile string, bgFile string, numPositive, numNegative, numStages, detectHeight, detectWidth int) error {
	mode := int64(0777)
	if _, err := os.Stat(dataFolder); os.IsNotExist(err) {
		os.MkdirAll(dataFolder, os.FileMode(mode))
	}
	cmdName, err := exec.LookPath(cvTrainCascadeCmd)
	if err != nil {
		return err
	}
	cmdArgs := []string{
		"-data", dataFolder,
		"-vec", posVecFile,
		"-bg", bgFile,
		"-numPos", strconv.Itoa(numPositive),
		"-numNeg", strconv.Itoa(numNegative),
		"-numStages", strconv.Itoa(numStages),
		"-w", strconv.Itoa(detectWidth),
		"-h", strconv.Itoa(detectHeight),
	}
	fmt.Println(cvTrainCascadeCmd, strings.Join(cmdArgs, " "))
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stdout = os.Stdout
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	cmd.Start()
	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
	return nil
}
