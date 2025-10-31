package main

import (
	"context"
	"fmt"

	"github.com/jcourtney5/blog-aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		// Get the current user from the DB
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Current User %s not found\n", s.cfg.CurrentUserName)
		}

		return handler(s, cmd, user)
	}
}
