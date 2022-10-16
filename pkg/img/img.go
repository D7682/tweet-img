package img

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image-manipulation/pkg/utils"
	"image/draw"
	"image/jpeg"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/disintegration/gift"
	"gopkg.in/yaml.v3"
)

type Image struct {
	FileLocation string
	Data         image.Image
	Prefix       string
}

type Config struct {
	ApiKey       string `yaml:"api_key"`
	ApiSecret    string `yaml:"api_secret"`
	BearerToken  string `yaml:"bearer_token"`
	AccessToken  string `yaml:"access_token"`
	AccessSecret string `yaml:"access_secret"`
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

func (i *Image) Write() error {
	utils.Examples()
	fileName, file, err := GenFile("examples", i.Prefix)
	if err != nil {
		return err
	}
	defer file.Close()
	i.FileLocation = fileName

	var opt jpeg.Options
	opt.Quality = 100

	err = jpeg.Encode(file, i.Data, &opt)
	if err != nil {
		return err
	}
	return nil
}

func (i Image) Send() error {
	f, err := os.ReadFile(".yaml")
	if err != nil {
		return err
	}

	var c Config
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		return err
	}

	v := url.Values{}

	api := anaconda.NewTwitterApiWithCredentials(c.AccessToken, c.AccessSecret, c.ApiKey, c.ApiSecret)
	/* t, err := api.PostTweet("", nil)
	if err != nil {
		return err
	} */

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, i.Data, nil)
	if err != nil {
		return err
	}
	data := buf.Bytes()

	mediaResponse, err := api.UploadMedia(base64.StdEncoding.EncodeToString(data))
	if err != nil {
		return err
	}

	v.Set("media_ids", strconv.FormatInt(mediaResponse.MediaID, 10))
	v.Set("in_reply_to_status_id", mediaResponse.MediaIDString)

	// tweetString := fmt.Sprintf("@%s", t.User.ScreenName)
	// _, err = api.PostTweet(tweetString, v)
	// if err != nil {
	// 	return err
	// }
	_, err = api.PostTweet("", v)
	if err != nil {
		return err
	}
	return nil
}

func Join(img1, img2 Image) (*Image, error) {
	// rectangle for view
	r := image.Rectangle{image.Point{0, 0}, image.Point{img1.Data.Bounds().Dx() + img2.Data.Bounds().Dx(), img1.Data.Bounds().Dy()}}
	rgba := image.NewRGBA(r)

	// locations on the rgba rectangle that each image will be drawn at.
	sp1 := image.Rectangle{rgba.Bounds().Min, image.Point{rgba.Bounds().Max.X / 2, rgba.Bounds().Max.Y}}
	sp2 := image.Rectangle{image.Point{rgba.Bounds().Max.X / 2, 0}, rgba.Bounds().Max}

	draw.Draw(rgba, sp1, img1.Data, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, sp2, img2.Data, image.Point{0, 0}, draw.Src)

	// Add The Image Filter, etc. with the gift library.
	g := gift.New(
		gift.Invert(),
		gift.Gamma(0.5),
	)

	dstImage := image.NewRGBA(g.Bounds(rgba.Bounds()))
	g.Draw(dstImage, rgba)

	return &Image{Data: dstImage, Prefix: img1.Prefix}, nil
}
