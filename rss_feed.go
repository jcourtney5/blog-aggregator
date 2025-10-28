package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	rssFeed := &RSSFeed{}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return rssFeed, fmt.Errorf("Error creating request: %w\n", err)
	}
	req.Header.Set("User-Agent", "gator")

	// Create HTTP client
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return rssFeed, fmt.Errorf("Error making request: %w\n", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return rssFeed, fmt.Errorf("Unexpected status code: %v\n", resp.Status)
	}

	// Read the data from the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rssFeed, fmt.Errorf("Error reading response body: %w\n", err)
	}

	// Decode the xml into an RSSFeed
	if err := xml.Unmarshal(body, &rssFeed); err != nil {
		return rssFeed, fmt.Errorf("Error decoding XML: %w\n", err)
	}

	// html Unescape all Title and Description fields
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(rssFeed.Channel.Item[i].Title)
		rssFeed.Channel.Item[i].Description = html.UnescapeString(rssFeed.Channel.Item[i].Description)
	}

	return rssFeed, nil
}
