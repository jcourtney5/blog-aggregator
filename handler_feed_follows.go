package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jcourtney5/blog-aggregator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	// Make sure there is enough args
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %v <url>", cmd.name)
	}

	// Get the args
	url := cmd.args[0]

	// Get the feed from the DB using the url
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Feed with url %s not found\n", url)
	}

	// Create the feed follow in the db
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Failed to create the feed_follow in the db: %w\n", err)
	}

	fmt.Println("Feed follow has been created:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("=========================================")

	return nil
}

func handlerListFeedFollows(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Failed to get the list of feed follows: %w\n", err)
	}

	if len(feedFollows) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	fmt.Printf("The current user %s is following these feeds:\n", s.cfg.CurrentUserName)
	for _, feedFollow := range feedFollows {
		fmt.Printf("* Feed URL:      %s\n", feedFollow.FeedName)
	}

	return nil
}

func handlerRemoveFollow(s *state, cmd command, user database.User) error {
	// Make sure there is enough args
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %v <url>", cmd.name)
	}

	// Get the args
	url := cmd.args[0]

	// Get the feed from the DB using the url
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Feed with url %s not found\n", url)
	}

	// Remove the feed follow in the db
	params := database.RemoveFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	}
	err = s.db.RemoveFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Failed to remove the feed_follow in the db: %w\n", err)
	}

	fmt.Printf("Successfully removed the feed follow\n")
	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}
