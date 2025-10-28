package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.ClearUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to clear the users table: %w\n", err)
	}

	fmt.Printf("Cleared the users table\n")
	return nil
}
