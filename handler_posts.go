package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jcourtney5/blog-aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2

	// Get the optional limit arg as an int
	if len(cmd.args) > 0 {
		if limitInt, err := strconv.Atoi(cmd.args[0]); err == nil {
			limit = limitInt
		} else {
			log.Printf("Failed to parse limit %s, using default of 2: %v\n", cmd.args[0], err)
		}
	}

	// Get the posts
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("Failed to get the recent posts: %w\n", err)
	}

	if len(posts) == 0 {
		fmt.Println("No posts found for this user.")
		return nil
	}

	fmt.Printf("Here are the %d most recent posts for the feeds the user %s follows:\n", len(posts), s.cfg.CurrentUserName)
	for _, post := range posts {
		printPost(&post)
		fmt.Println("=========================================")
	}

	return nil
}

func printPost(post *database.GetPostsForUserRow) {
	fmt.Printf("* Feed:          %s\n", post.FeedName)
	fmt.Printf("* Published At:  %v\n", post.PublishedAt.Time.Format(time.RFC822))
	fmt.Printf("* Title:         %s\n", post.Title)
	fmt.Printf("* URL:           %s\n", post.Url)
	fmt.Printf("* Description:   %s\n", post.Description.String)
}
