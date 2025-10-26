package main

import (
	"fmt"

	"github.com/jcourtney5/blog-aggregator/internal/config"
)

type state struct {
	cfg *config.Config
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

	// Set user which should update and save config file
	err := s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("Failed to set user: %w\n", err)
	}

	fmt.Printf("Username '%s' has been set\n", username)

	return nil
}
