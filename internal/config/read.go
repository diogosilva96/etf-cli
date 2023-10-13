package config

import (
	"fmt"
	"log"
	"os"
)

func Read(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(data))
}
