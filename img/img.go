package img

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Image struct {
	data   image.Image
	prefix string
}

// GenFile will generate a file with a random name for use
func GenFile(dir, imgPrefix string) (*os.File, error) {
	currentTime := time.Now().Format("Monday, 02 Jan 2006 15.04.05 MST")
	fileName := fmt.Sprintf("%v %v.jpg", imgPrefix, currentTime)

	f, err := os.Create(filepath.Join(dir, fileName))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func NewImg(endpoint string) (Image, error) {
	var img Image
	resp, err := http.Get(endpoint)
	if err != nil {
		return img, err
	}
	defer resp.Body.Close()

	image, _, err := image.Decode(resp.Body)
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

func (i Image) Write() error {
	file, err := GenFile("examples", i.prefix)
	if err != nil {
		return err
	}
	defer file.Close()

	var opt jpeg.Options
	opt.Quality = 100

	err = jpeg.Encode(file, i.data, &opt)
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

	return &Image{data: rgba, prefix: img1.prefix}, nil
}
