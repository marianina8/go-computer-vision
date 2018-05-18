package microsoft

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/marianina8/go-computer-vision/detection/visual"
)

const (
	visionKeyEnvVariable = "azure_vision_key"
)

type azureDetection struct {
	visionKey string
	client    *http.Client
}

// New makes a new detector using haar cascades.
func New() (visual.Eyes, error) {
	if os.Getenv(visionKeyEnvVariable) == "" {
		return nil, errors.New("missing " + visionKeyEnvVariable)
	}
	return &azureDetection{
		visionKey: os.Getenv(visionKeyEnvVariable),
		client:    &http.Client{},
	}, nil
}

func (a *azureDetection) Detect(file visual.ImageInfo) visual.Objects {
	if file.Path != "" {
		return a.detectFromFile(file.Path)
	}
	return a.detectFromBytes(*file.Bytes)
}

func (a azureDetection) detectFromFile(filePath string) visual.Objects {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return visual.Objects{}
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
		return visual.Objects{}
	}
	objects, _ := a.request(fileBytes)
	return objects
}

func (a azureDetection) detectFromBytes(fileBytes []byte) visual.Objects {
	objects, _ := a.request(fileBytes)
	return objects
}

func (a azureDetection) request(fileBytes []byte) (visual.Objects, error) {
	req, err := a.visionRequest(fileBytes)
	if err != nil {
		return nil, err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	var resBody io.Reader
	peeker := bufio.NewReader(resp.Body)
	if head, err := peeker.Peek(1); err != nil {
		fmt.Println(err)
		return nil, err
	} else if head[0] == '{' {
		var body bytes.Buffer
		err = json.NewDecoder(io.TeeReader(peeker, &body)).Decode(&res)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		resBody = &body
	} else {
		resBody = peeker
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status:" + resp.Status)
	}
	var result response
	err = json.NewDecoder(resBody).Decode(&result)
	resJson, _ := json.Marshal(result)
	fmt.Println("result: ", string(resJson))
	fmt.Println(err)
	objects := visual.Objects{}
	for _, category := range result.Categories {
		for _, celeb := range category.Detail.Celebrities {
			detection := visual.Detection{
				Confidence: celeb.Confidence,
				Rectangle: visual.Rectangle{
					Top:    celeb.FaceRectangle.Top,
					Left:   celeb.FaceRectangle.Left,
					Width:  celeb.FaceRectangle.Width,
					Height: celeb.FaceRectangle.Height,
				},
			}
			object := visual.Object{}
			object[celeb.Name] = detection
			objects = append(objects, object)
		}
	}
	if len(objects) > 0 {
		return objects, nil
	}
	for _, face := range result.Faces {
		detection := visual.Detection{
			Rectangle: visual.Rectangle{
				Top:    face.FaceRectangle.Top,
				Left:   face.FaceRectangle.Left,
				Width:  face.FaceRectangle.Width,
				Height: face.FaceRectangle.Height,
			},
		}
		object := visual.Object{}
		objType := fmt.Sprintf("Face [%d, %s]", face.Age, face.Gender)
		object[objType] = detection
		objects = append(objects, object)
	}
	return objects, err
}

func (a azureDetection) visionRequest(fileBytes []byte) (*http.Request, error) {
	params := url.Values{}
	visualFeatures := []string{"Faces"}
	details := []string{"Celebrities"}
	if len(visualFeatures) == 0 {
		fmt.Println("missing param features")
		return nil, errors.New("missing param features")
	}
	params.Set("details", strings.Join(details, ","))
	params.Set("visualFeatures", strings.Join(visualFeatures, ","))
	apiURL := fmt.Sprintf("https://westus.api.cognitive.microsoft.com/vision/v2.0/analyze?%s", params.Encode())
	fmt.Println(apiURL)
	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(fileBytes))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Ocp-Apim-Subscription-Key", a.visionKey)
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("Content-Length", strconv.Itoa(len(fileBytes)))
	return req, nil
}

type response struct {
	Categories []struct {
		Detail struct {
			Celebrities []struct {
				Confidence    float64 `json:"confidence"`
				FaceRectangle struct {
					Height int `json:"height"`
					Left   int `json:"left"`
					Top    int `json:"top"`
					Width  int `json:"width"`
				} `json:"faceRectangle"`
				Name string `json:"name"`
			} `json:"celebrities"`
		} `json:"detail"`
		Name  string  `json:"name"`
		Score float64 `json:"score"`
	} `json:"categories"`
	Faces []struct {
		Age           int `json:"age"`
		FaceRectangle struct {
			Height int `json:"height"`
			Left   int `json:"left"`
			Top    int `json:"top"`
			Width  int `json:"width"`
		} `json:"faceRectangle"`
		Gender string `json:"gender"`
	} `json:"faces"`
	Metadata struct {
		Format string `json:"format"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"metadata"`
	RequestID string `json:"requestId"`
}
