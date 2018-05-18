package images

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
)

/*
 *
 *  Name: store_images.go
 *  Author: Marian Montagnino
 *  Description: storeImages is a function that given specified parameters will download negative images from www.image-net.org
 *  Parameters: takes in links to images from http://www.image-net.org/, folder name to store images and limit (int) the number of downloaded images.
 *      . links ([]string) - links contains a slice of strings representing a single link to image-net containing links to negative images
 *      example link: http://www.image-net.org/api/text/imagenet.synset.geturls?wnid=n12102133
 *      . folderName (string) - folder name to store all negative images (negative images are saved as numbers, ie. 1.jpg, 2.jpg, ...)
 *      . grayscale (bool) - if set to true, converts image to grayscale
 *      . start (int) - number to start saving file.  If function was run earlier and the last negative image saved is 20.jpg, pass in 21 for start so
 *      you do not overwrite any files saved in previous run
 *      . limit (int) - limit number of negative files saved
 *      . height (int) - resize to height (if 0 preserve height)
 *      . width (int) - resize to width (if 0 preserve width)
 *
 */

var (
	errInsuficientNegImageFiles = errors.New("Insufficient files.  Add more links to imagelinks.txt file to download more images: Visit http://image-net.org/download-imageurls for more details.")
	client                      = &http.Client{
		Timeout: 60 * time.Second,
	}
)

func Get(imageLinkFile string, folderName string, grayscale bool, start, numDownload, height, width int) error {
	links, err := getLinksFromFile(imageLinkFile)
	if err != nil {
		return err
	}
	picNum := start
	for _, imageLink := range links {
		fmt.Println("\nimage-net.org link:", imageLink)
		// create a request to image-net
		req, err := http.NewRequest(http.MethodGet, imageLink, nil)
		if err != nil {
			fmt.Println("bad request:", err.Error())
			continue
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("unable to complete request: err=%s resp=%v\n", err, resp)
			continue
		}
		defer resp.Body.Close()
		// create folder if does not exist
		mode := int64(0777)
		if _, err := os.Stat(folderName); os.IsNotExist(err) {
			os.MkdirAll(folderName, os.FileMode(mode))
		}
		// scan through each line of the response body
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			fmt.Println("\n\tdownloading:", scanner.Text())
			if picNum > numDownload {
				return nil
			}
			// get the image associated with the link
			resp, err := http.Get(scanner.Text())
			if err != nil || resp.StatusCode != http.StatusOK {
				if err != nil {
					fmt.Printf("\t...error getting file: err=%s skipping!\n", err.Error())
					continue
				}
				fmt.Printf("\t...error getting file: resp=%d skipping!\n", resp.StatusCode)
				continue
			}
			defer resp.Body.Close()
			//open a file for writing
			filePath := filepath.Join(folderName, strconv.Itoa(picNum)+".jpg")
			file, err := os.Create(filePath)
			if err != nil {
				fmt.Println("\t...error creating:", err, "skipping!")
				continue
			}
			// Use io.Copy to just dump the response body to the file. This supports huge files
			n, err := io.Copy(file, resp.Body)
			if err != nil || n < 3000 {
				_ = os.Remove(filePath)
				fmt.Printf("\t...error copying bytes (n=%d) skipping!\n", n)
				continue
			}
			file.Close()
			// open the file for image manipulation
			srcImg, err := imaging.Open(filePath)
			if srcImg == nil || err != nil {
				fmt.Println("\t...error opening:", err, "skipping!")
				continue
			}
			// resize image to 100x100
			if height > 0 || width > 0 {
				srcImg = imaging.Resize(srcImg, width, height, imaging.Lanczos)
			}
			// convert image to grayscale
			if grayscale {
				srcImg = imaging.Grayscale(srcImg)
			}
			err = imaging.Save(srcImg, filePath)
			if err != nil {
				fmt.Println("\t...error saving:", err, "skipping!")
				continue
			}
			picNum++
			fmt.Println("\t...success!")
		}
	}
	if picNum < numDownload {
		return errInsuficientNegImageFiles
	}
	return nil
}

func getLinksFromFile(imageLinkFile string) ([]string, error) {
	links := []string{}
	file, err := os.Open(imageLinkFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		links = append(links, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return links, nil
}
