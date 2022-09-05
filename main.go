package main

import (
	"fmt"
	"image-manipulation/img"
	"image-manipulation/utils"
	"log"
	"math/rand"
	"time"
)

func RandomSeed() int {
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := 1

	randomInt := rand.Intn(max-min+1) + min
	return randomInt
}

func main() {
	urls := []string{
		"https://picsum.photos/2560/1440",
		// "https://cataas.com/cat?width=2560&height=1440",
		"https://picsum.photos/1920/1080",
	}

	seed := RandomSeed()
	currentURL := urls[seed]

	utils.Examples()

	img1, err := img.NewImg(currentURL)
	if err != nil {
		log.Fatal(err)
	}

	img2, err := img.NewImg(currentURL)
	if err != nil {
		log.Fatal(err)
	}

	joined, err := img.Join(img1, img2)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 1)
	fmt.Println("running")
	fmt.Println("running")
	err = joined.Write()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("running")
	err = joined.Send()
	if err != nil {
		log.Fatal(err)
	}
}
