package img

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/disintegration/gift"
	"github.com/h2non/bimg"
	"gopkg.in/yaml.v3"
)

type Image struct {
	file_location string
	data          image.Image
	prefix        string
}

// GenFile will generate a file with a random name for use
func GenFile(dir, imgPrefix string) (string, *os.File, error) {
	currentTime := time.Now().Format("Monday, 02 Jan 2006 15.04.05 MST")
	fileName := fmt.Sprintf("%v %v.jpg", imgPrefix, currentTime)

	f, err := os.Create(filepath.Join(dir, fileName))
	if err != nil {
		return "", nil, err
	}
	return filepath.Join(dir, fileName), f, nil
}

func NewImg(endpoint string) (Image, error) {
	var img Image
	client := http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := client.Get(endpoint)
	if err != nil {
		return img, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return img, err
	}

	newImage, err := bimg.NewImage(body).Resize(2560, 1440)
	if err != nil {
		return img, err
	}

	if bimg.NewImage(newImage).Type() == "jpg" {
		fmt.Println(os.Stderr, "The image was converted into jpg")
	}

	reader := bytes.NewReader(newImage)
	image, _, err := image.Decode(reader)
	if err != nil {
		return img, err
	}

	img.prefix = "img"
	img.data = image
	if strings.Contains(endpoint, "cataas") {
		img.prefix = "cat"
	}

	return img, nil
}

func (i *Image) Write() error {
	fileName, file, err := GenFile("examples", i.prefix)
	if err != nil {
		return err
	}
	defer file.Close()
	i.file_location = fileName

	var opt jpeg.Options
	opt.Quality = 100

	err = jpeg.Encode(file, i.data, &opt)
	if err != nil {
		return err
	}
	return nil
}

type Config struct {
	ApiKey       string `yaml:"api_key"`
	ApiSecret    string `yaml:"api_secret"`
	BearerToken  string `yaml:"bearer_token"`
	AccessToken  string `yaml:"access_token"`
	AccessSecret string `yaml:"access_secret"`
}

func (i Image) Send() error {
	fmt.Println("running")
	f, err := os.ReadFile(".yaml")
	if err != nil {
		return err
	}

	var c Config
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		return err
	}

	fmt.Println("running")

	v := url.Values{}

	api := anaconda.NewTwitterApiWithCredentials(c.AccessToken, c.AccessSecret, c.ApiKey, c.ApiSecret)
	t, err := api.PostTweet(i.file_location, nil)
	if err != nil {
		return err
	}
	fmt.Println(t)
	fmt.Println("running")
	data, err := os.ReadFile(i.file_location)
	if err != nil {
		return err
	}

	mediaResponse, err := api.UploadMedia(base64.StdEncoding.EncodeToString(data))
	if err != nil {
		return err
	}

	fmt.Println(mediaResponse.MediaID)

	v.Set("media_ids", strconv.FormatInt(mediaResponse.MediaID, 10))
	v.Set("in_reply_to_status_id", mediaResponse.MediaIDString)

	tweetString := fmt.Sprintf("@%s", t.User.ScreenName)
	_, err = api.PostTweet(tweetString, v)
	if err != nil {
		return err
	}
	return nil
}

func Join(img1, img2 Image) (*Image, error) {
	// starting position of the second image (bottom left)
	sp2 := image.Point{img1.data.Bounds().Dx(), 0}

	// new rectangle for the second image
	r2 := image.Rectangle{sp2, sp2.Add(img2.data.Bounds().Size())}

	// rectangle for the big image
	r := image.Rectangle{image.Point{0, 0}, r2.Max}

	rgba := image.NewRGBA(r)

	draw.Draw(rgba, img1.data.Bounds(), img1.data, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, r2, img2.data, image.Point{0, 0}, draw.Src)

	g := gift.New(
		gift.Invert(),
		gift.Gamma(0.5),
	)

	dstImage := image.NewRGBA(rgba.Bounds())
	g.Draw(dstImage, rgba)

	// Draw the fgImage over the dstImage at the (100, 100) position

	return &Image{data: dstImage, prefix: img1.prefix}, nil
}
