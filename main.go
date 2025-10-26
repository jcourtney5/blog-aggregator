package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jcourtney5/blog-aggregator/internal/config"
)

func main() {

	// Read the current config file
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Init state struct
	st := &state{
		cfg: &cfg,
	}

	// Init commands struct
	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	// Register our commands
	cmds.register("login", handlerLogin)

	// Get the command line args (skip first one which is program name)
	args := os.Args[1:]

	// Make sure we at least have the command
	if len(args) == 0 {
		fmt.Printf("Not enough arguments passed, need at least one for the command\n")
		os.Exit(1)
	}

	// Create command struct
	command := command{
		name: args[0],
		args: args[1:],
	}

	// Run the command
	err = cmds.run(st, command)
	if err != nil {
		fmt.Printf("Command failed: %v\n", err)
		os.Exit(1)
	}
}
