package main

import (
	"context"
	"encoding/json"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {
	// Hard coded url for now
	url := "https://www.wagslane.dev/index.xml"

	rssFeed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Failed to fetch RSS feed: %w\n", err)
	}

	//fmt.Println(rssFeed)

	// print formatted JSON
	formattedJSON, err := json.MarshalIndent(rssFeed, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(formattedJSON))

	return nil
}
