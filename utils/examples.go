package utils

import (
	"log"
	"os"
)

// Examples is the function that checks if the "examples" folder already exists
func Examples() {
	// Check if the examples folder already exists.
	if _, err := os.Stat("examples"); os.IsNotExist(err) {
		err = os.Mkdir("examples", 0777)
		if err != nil {
			log.Fatal(err)
		}
	}
}
