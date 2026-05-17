package main

import (
	"log"
	"os"

	"github.com/lynndixon/gator/internal/config"
	"github.com/lib/pq"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("could not read config file: %v", err)
	}

	programState := &state{
		cfg: &cfg,
	}

	cmds := commands{
		RegisteredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	args := os.Args

	if len(args) < 2 {
		log.Fatal("Usage: cli <command> [args...]\n")
		return
	}

	cmdName := args[1]
	cmdArgs := args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}

}
