package img

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"
	"tweet-img/pkg/utils"

	"github.com/disintegration/gift"
)

type Image struct {
	FileLocation string
	Data         image.Image
	Prefix       string
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
