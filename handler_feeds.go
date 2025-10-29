package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jcourtney5/blog-aggregator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	// Make sure there is enough args
	if len(cmd.args) < 2 {
		return fmt.Errorf("usage: %v <name> <url>", cmd.name)
	}

	// Get the args
	name := cmd.args[0]
	url := cmd.args[1]

	// Get the current user from the DB
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Current User %s not found\n", s.cfg.CurrentUserName)
	}

	// Create the feed in the db
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Failed to create the feed in the db: %w\n", err)
	}

	fmt.Println("Feed has been created:")
	printFeed(feed, user)
	fmt.Println("=========================================")

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	// Get all the feeds from the DB
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get all the feeds from the feeds table: %w\n", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Println("All the current feeds:")
	fmt.Println("=========================================")
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get user: %w", err)
		}
		printFeed(feed, user)
		fmt.Println("=========================================")
	}

	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
}
