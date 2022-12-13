package client

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/jpeg"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"tweet-img/pkg/config"
	"tweet-img/pkg/img"

	"github.com/ChimeraCoder/anaconda"
	"github.com/disintegration/gift"
	"gopkg.in/yaml.v3"
)

type ImageClient struct {
	http   *http.Client
	Config *config.Config
}

func NewClient() (*ImageClient, error) {
	f, err := os.ReadFile("../../.yaml")
	if err != nil {
		return nil, errors.New("Error in reading the configuration file")
	}

	var conf config.Config
	err = yaml.Unmarshal(f, &conf)
	if err != nil {
		return nil, errors.New("Error in parsing the configuration file")
	}

	c := &http.Client{
		Timeout: time.Second * 5,
	}

	return &ImageClient{
		http:   c,
		Config: &conf,
	}, nil
}

func (i ImageClient) NewImg(endpoint string) (img.Image, error) {
	var img img.Image

	resp, err := i.http.Get(endpoint)
	if err != nil {
		return img, err
	}
	defer resp.Body.Close()

	src, _, err := image.Decode(resp.Body)
	if err != nil {
		return img, err
	}

	giftFilter := gift.Resize(2560, 1440, gift.LanczosResampling)
	dst := image.NewRGBA(giftFilter.Bounds(src.Bounds()))
	giftFilter.Draw(dst, src, &gift.Options{
		Parallelization: true,
	})

	img.Prefix = "img"
	img.Data = dst
	if strings.Contains(endpoint, "cataas") {
		img.Prefix = "cat"
	}

	return img, nil
}

func (i ImageClient) Send(img img.Image) error {
	v := url.Values{}
	api := anaconda.NewTwitterApiWithCredentials(i.Config.AccessToken, i.Config.AccessSecret, i.Config.ApiKey, i.Config.ApiSecret)
	/* t, err := api.PostTweet("", nil)
	if err != nil {
		return err
	} */

	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img.Data, nil)
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
