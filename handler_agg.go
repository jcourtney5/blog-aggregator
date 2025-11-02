package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jcourtney5/blog-aggregator/internal/database"
	"github.com/lib/pq"
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

	// Save the posts
	for _, item := range rssFeed.Channel.Item {
		// parse the publishedAt time
		publishedAt := sql.NullTime{}
		parsedTime, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err == nil {
			publishedAt = sql.NullTime{Time: parsedTime, Valid: true}
		} else {
			log.Printf("Failed to parse published at %s: %v\n", item.PubDate, err)
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			// Only log an error if didn't fail unique constraint
			if pgerr, ok := err.(*pq.Error); ok {
				if pgerr.Code != "23505" {
					log.Printf("Failed to save post: %v\n", err)
				}
			}
		}
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
