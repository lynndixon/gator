package main

import "fmt"

// This is the function signature of all command handlers
func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>\n", cmd.Name)
	}

	username := cmd.Args[0]

	err := s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w\n", err)
	}

	fmt.Printf("Username has been set to: %v\n", username)

	return nil
}
