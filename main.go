package main

import (
	"fmt"
	"log"

	"github.com/jcourtney5/blog-aggregator/internal/config"
)

func main() {

	// Read the current config file
	configData, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Config before: %+v\n", configData)

	// Update current username
	err = configData.SetUser("Jesse")
	if err != nil {
		log.Fatal(err)
	}

	// Read updated config file data to see if it changed
	configData, err = config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Config after: %+v\n", configData)
}
