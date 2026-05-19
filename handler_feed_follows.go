package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lynndixon/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>\n", cmd.Name)
	}

	feedURL := cmd.Args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("Could not retrieve a feed with that URL: %w", err)
	}

	feedInfo, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Could not retrieve feed info: %w", err)
	}

	printFeedFollow(user.Name, feedInfo.FeedName)

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>\n", cmd.Name)
	}

	feedUrl := cmd.Args[0]

	unfollowedFeed, err := s.db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		Name: user.Name,
		Url:  feedUrl,
	})
	if err != nil {
		return fmt.Errorf("Could not find that url in current user's followed feeds: %w", err)
	}

	fmt.Printf("%v is no longer following %v", unfollowedFeed.UserName, unfollowedFeed.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("Could not retrieve followed feeds for user: %w", err)
	}

	if len(follows) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	fmt.Printf("Feed follows for user %s:\n", user.Name)

	for _, feed := range follows {
		fmt.Printf("- %v\n", feed.FeedName)
	}

	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}
