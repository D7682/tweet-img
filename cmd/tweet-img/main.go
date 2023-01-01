package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
	"tweet-img/pkg/client"
	"tweet-img/pkg/img"
)

func RandomSeed() int {
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := 1

	randomInt := rand.Intn(max-min+1) + min
	return randomInt
}

/* func UI() {
	a := app.New()
	w := a.NewWindow("Hello")

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	w.ShowAndRun()
} */

// Create a save option for a function so that a person is able to easily decide
// if they want to save an image as well or not.

func main() {
	// UI()
	start := time.Now()
	urls := []string{
		"https://picsum.photos/2560/1440",
		// "https://picsum.photos/2560/1440",
		"https://thecatapi.com/api/images/get?format=src&type=image&mime_types=jpg&size=full",
		//"https://cataas.com/cat?width=2560&height=1440",
		// "https://picsum.photos/1920/1080",
	}

	c, err := client.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	var images [2]img.Image

	for i := 0; i <= 1; i++ {
		wg.Add(1)
		go func(i int) {
			seed := RandomSeed()
			currentURL := urls[seed]
			newImage, err := c.NewImg(currentURL)
			if err != nil {
				log.Fatal(err)
			}
			images[i] = newImage
			defer wg.Done()
		}(i)
	}
	wg.Wait()

	joined, err := img.Join(images[0], images[1])
	if err != nil {
		log.Fatal(err)
	}

	err = c.Send(*joined)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Elapsed: %v\n", time.Since(start))
}
