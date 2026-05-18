package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lynndixon/gator/internal/database"
)

const URL = "https://www.wagslane.dev/index.xml"

func handlerAgg(s *state, cmd command) error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return fmt.Errorf("Failed to create new HTTP request: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP error: %w", err)
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response body: %w", err)
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal XML: %w", err)
	}

	fmt.Print(rssFeed, "\n")

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	currentUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("could not retrieve current user: %w", err)
	}

	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feed name> <url>\n", cmd.Name)
	}

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed)
	fmt.Println()
	fmt.Println("==========================")

	return nil
}

// Helper function to print the fields of a feed
func printFeed(feed database.Feed) {
	fmt.Printf(" * ID:		%v\n", feed.ID)
	fmt.Printf(" * Name:	%v\n", feed.Name)
	fmt.Printf(" * URL:		%v\n", feed.Url)
}

func handlerGetFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s\n", cmd.Name)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Could not retrieve feeds: %w", err)
	}

	if len(feeds) == 0 {
		return errors.New("No feeds found in database.")
	}

	for _, feed := range feeds {
		fmt.Printf("%v, %v, created by: %v\n", feed.Name, feed.Url, feed.User)
	}

	return nil
}
