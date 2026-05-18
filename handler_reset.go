package main

import (
	"context"
	"fmt"

	"github.com/pressly/goose/v3"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't delete users: %w", err)
	}

	// Initialize goose and run migrations
	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Assuming your migrations are in "sql/schema"
	err = goose.Up(s.rawDB, "sql/schema")
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	fmt.Println("Database reset successfully!")
	return nil

}
