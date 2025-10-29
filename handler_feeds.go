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
	printFeed(feed)
	fmt.Println("=========================================")

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	// Get all the users from the DB
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get all the users from the users table: %w\n", err)
	}

	// Put users in a map so we can lookup id -> user
	userMap := make(map[uuid.UUID]database.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// Get all the feeds from the DB
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get all the feeds from the feeds table: %w\n", err)
	}

	fmt.Printf("All the current users:\n")
	for _, feed := range feeds {
		fmt.Println("Feed-------------------------------")
		fmt.Printf("* Name:          %s\n", feed.Name)
		fmt.Printf("* URL:           %s\n", feed.Url)
		fmt.Printf("* User:          %s\n", userMap[feed.UserID].Name)
	}

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}
