package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jcourtney5/blog-aggregator/internal/config"
	"github.com/jcourtney5/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	// Read the current config file
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to our SQL db
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	// Init state struct
	st := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	// Init commands struct
	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	// Register our commands
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)

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
