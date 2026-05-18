package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lynndixon/gator/internal/database"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>\n", cmd.Name)
	}

	feedURL := cmd.Args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("Could not retrieve a feed with that URL: %w", err)
	}

	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Could not retrieve current user: %w", err)
	}

	feedInfo, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Could not retrieve feed info: %w", err)
	}

	fmt.Printf("Feed name: %v", feedInfo.FeedName)
	fmt.Printf("Current user: %v", currentUser)

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Could not retrieve current user: %w", err)
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.Name)
	if err != nil {
		return fmt.Errorf("Could not retrieve followed feeds for user: %w", err)
	}

	for _, feed := range follows {
		fmt.Printf("- %v\n", feed.FeedName)
	}

	return nil
}
