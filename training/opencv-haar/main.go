package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/marianina8/go-computer-vision/training/opencv-haar/images"
	"github.com/marianina8/go-computer-vision/training/opencv-haar/samples"
	"github.com/marianina8/go-computer-vision/training/opencv-haar/training"
)

func main() {
	var (
		err error

		// Help
		help = flag.Bool("help", false, "help")

		// Negative Sample parameters
		getNegatives       = flag.Bool("getNegatives", false, "download negative files 100x100 to /negatives folder and generates bg file")
		imageLinksFile     = flag.String("link", "imagelinks.txt", "file containing image-net links to download negative files")
		grayScaleNegatives = flag.Bool("grayscale", true, "save negative files as grayscale")
		picNumStart        = flag.Int("picNum", 1, "file number to start negative file download")
		numNegDownloads    = flag.Int("numNegDownloads", 4000, "number of negative files to download")
		defaultNegHeight   = flag.Int("negHeight", 100, "height of negative image file")
		defaultNegWidth    = flag.Int("negWidth", 100, "width of negative image file")
		defaultNegFolder   = flag.String("negFolder", "negatives", "folder containing negative images")
		bgFile             = flag.String("bgFile", "bg.txt", "filepath of bg (negatives) file")

		// Positive Sample parameters
		multiplePos         = flag.Bool("multiplePos", false, "use multiple positive images and generate pos file")
		defaultPosFolder    = flag.String("posFolder", "positives", "folder containing positives images")
		posFile             = flag.String("posFile", "pos.txt", "filepath of bg (positives) file")
		posImage            = flag.String("posImage", "positive.jpg", "filepath of positive image (smaller than negative height x width)")
		numGeneratedSamples = flag.Int("numSamples", 3900, "number of positive samples to generate")
		posVecFile          = flag.String("posVecFile", "positives.vec", "filepath of positive vector file")
		maxAngle            = flag.Float64("maxAngle", 0.5, "max X, Y, Z angle used for generating positive samples")

		// Training Cascade parameters
		dataFolder     = flag.String("dataFolder", "data", "data folder that stores the haar cascade training data")
		detectHeight   = flag.Int("detectHeight", 20, "detection height")
		detectWidth    = flag.Int("detectWidth", 20, "detection width")
		trainNumPos    = flag.Int("trainNumPos", 3600, "number of positive samples to use for training")
		trainNumNeg    = flag.Int("trainNumNeg", 1800, "number of negatives to use for training")
		trainNumStages = flag.Int("trainNumStages", 13, "number of training stages")
	)
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	if *help {
		flag.Usage()
		return
	}
	if *multiplePos {
		err = generatePosFile(*defaultPosFolder, *posFile)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	if *getNegatives {
		// Download negative image files
		err = images.Get(*imageLinksFile, *defaultNegFolder, *grayScaleNegatives, *picNumStart, *numNegDownloads, *defaultNegHeight, *defaultNegWidth)
		if err != nil {
			log.Fatal(err)
			return
		}
		// Generate bg file from negative (background) files
		err = generateBGFile(*defaultNegFolder, *bgFile)
		if err != nil {
			log.Fatal(err)
			return
		}

	}
	// Create positive samples
	err = samples.CreateSamples(*posImage, *bgFile, *numGeneratedSamples, *maxAngle)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Create positive vector file
	err = samples.CreatePositiveVectorFile(*numGeneratedSamples, *posVecFile, *detectHeight, *detectWidth)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Train haar cascade file
	err = training.HaarCascade(*dataFolder, *posVecFile, *bgFile, *trainNumPos, *trainNumNeg, *trainNumStages, *detectHeight, *detectWidth)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func generateBGFile(negFolder, bgFile string) error {
	negFiles, err := ioutil.ReadDir(negFolder)
	if err != nil {
		return err
	}
	var negFilePaths string
	for _, negFilePath := range negFiles {
		negFilePaths += negFolder + "/" + negFilePath.Name() + "\n"
	}
	err = ioutil.WriteFile(bgFile, []byte(negFilePaths), 0666)
	if err != nil {
		return err
	}
	return nil
}

func generatePosFile(posFolder, posFile string) error {
	posFiles, err := ioutil.ReadDir(posFolder)
	if err != nil {
		return err
	}
	var posFilePaths string
	for _, posFilePath := range posFiles {
		posFilePaths += posFolder + "/" + posFilePath.Name() + "\n"
	}
	err = ioutil.WriteFile(posFile, []byte(posFilePaths), 0666)
	if err != nil {
		return err
	}
	return nil
}
