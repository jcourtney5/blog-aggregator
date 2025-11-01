package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jcourtney5/blog-aggregator/internal/database"
)

// example feeds
// go run . addfeed 'TechCrunch' https://techcrunch.com/feed/
// go run . addfeed 'Hacker News' https://news.ycombinator.com/rss
// go run . addfeed 'Boot.dev Blog' https://blog.boot.dev/index.xml

func handlerAgg(s *state, cmd command) error {
	// Make sure there is enough args
	if len(cmd.args) < 1 || len(cmd.args) > 2 {
		return fmt.Errorf("usage: %v <time_between_requests>", cmd.name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("unable to parse duration: %w\n", err)
	}

	fmt.Printf("Collecting Feeds Every %v\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	return nil
}

func scrapeFeeds(s *state) {
	// get the next feed to fetch
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Printf("Failed to get next feed to fetch: %v\n", err)
		return
	}

	log.Printf("Next feed to fetch is:\n")
	log.Printf("* Name:          %s\n", feed.Name)
	log.Printf("* URL:           %s\n", feed.Url)
	log.Println("=========================================")

	// fetch the feed
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Failed to fetch RSS feed %s: %v\n", feed.Name, err)
		return
	}

	//err = print_rss_feed(rssFeed)
	//err = print_rss_feed_json_formatted(rssFeed)
	err = print_rss_feed_titles(rssFeed)
	if err != nil {
		log.Printf("Failed to print RSS feed %s: %v\n", feed.Name, err)
	}

	// mark feed as fetched
	_, err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:            feed.ID,
		UpdatedAt:     time.Now().UTC(),
		LastFetchedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		log.Printf("Failed to mark feed %s as fetched: %v\n", feed.Name, err)
	}

	log.Printf("Feed %s collected, %v posts found\n", feed.Name, len(rssFeed.Channel.Item))
}

func print_rss_feed(rssFeed *RSSFeed) error {
	fmt.Printf("%+v\n", rssFeed)
	return nil
}

func print_rss_feed_json_formatted(rssFeed *RSSFeed) error {
	// print formatted JSON
	formattedJSON, err := json.MarshalIndent(rssFeed, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(formattedJSON))
	return nil
}

func print_rss_feed_titles(rssFeed *RSSFeed) error {
	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("* Title: %s\n", item.Title)
	}
	return nil
}
