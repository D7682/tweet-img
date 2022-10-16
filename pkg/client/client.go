package client

import (
	"image"
	"image-manipulation/pkg/img"
	"net/http"
	"strings"
	"time"

	"github.com/disintegration/gift"
)

type ImageClient struct {
	C *http.Client
}

func NewClient() ImageClient {
	c := &http.Client{
		Timeout: time.Second * 5,
	}

	return ImageClient{
		C: c,
	}
}

func (i ImageClient) NewImg(endpoint string) (img.Image, error) {
	var img img.Image

	resp, err := i.C.Get(endpoint)
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
