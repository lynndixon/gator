package main

import (
	"errors"
)

// command struct
// contains a name and slice of string arguments
// in the case of login, name login, handler expects arguments
// slice to contain one string, ther useername.

type command struct {
	Name string
	Args []string
}

type commands struct {
	RegisteredCommands map[string]func(*state, command) error
}

// This method registers a new handler function for a command name.
func (c *commands) register(name string, f func(*state, command) error) {
	c.RegisteredCommands[name] = f
}

// This method runs a given command with the provided state if it exists.
func (c *commands) run(s *state, cmd command) error {
	f, ok := c.RegisteredCommands[cmd.Name]
	if !ok {
		return errors.New("Command not registered.\n")
	}
	return f(s, cmd)
}
