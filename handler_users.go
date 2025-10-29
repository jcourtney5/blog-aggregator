package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jcourtney5/blog-aggregator/internal/database"
	"github.com/lib/pq"
)

func handlerLogin(s *state, cmd command) error {
	// Make sure there is enough args
	if len(cmd.args) == 0 {
		return fmt.Errorf("login missing the <username> arguement")
	}

	// Get username from first arg
	username := cmd.args[0]

	// Check if the user exists in the db first
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("User %s not found\n", username)
	}

	// Set user which should update and save config file
	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Failed to set user: %w\n", err)
	}

	fmt.Printf("Username '%s' has been set\n", username)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	// Make sure there is enough args
	if len(cmd.args) == 0 {
		return fmt.Errorf("register missing the <username> arguement")
	}

	// Get username from first arg
	username := cmd.args[0]

	// Create the user in the db
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      username,
	}
	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Failed to create the user in the db: %w\n", err)
	}

	// Set user which should update and save config file
	err = s.cfg.SetUser(username)
	if err != nil {
		// check for and exit on unique constraint violation
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == "23505" {
				return fmt.Errorf("Username '%s' already exists: %w\n", username, err)
			}
		} else {
			return fmt.Errorf("Failed to set user: %w\n", err)
		}
	}

	fmt.Println("User has been created:")
	printUser(user)
	fmt.Println("=========================================")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	// Get all the users from the DB
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get all the users from the users table: %w\n", err)
	}

	fmt.Printf("All the current users:\n")
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
