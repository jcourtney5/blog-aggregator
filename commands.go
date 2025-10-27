package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jcourtney5/blog-aggregator/internal/config"
	"github.com/jcourtney5/blog-aggregator/internal/database"
	"github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}

	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

	fmt.Printf("User has been created %v\n", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ClearUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to clear the users table: %w\n", err)
	}

	fmt.Printf("Cleared the users table\n")
	return nil
}
