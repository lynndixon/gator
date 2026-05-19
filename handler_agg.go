package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lynndixon/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %v <time duration string: 1h, 1m, 1s etc>", cmd.Name)
	}

	durationString := cmd.Args[0]
	timeBetweenRequests, err := time.ParseDuration(durationString)
	if err != nil {
		return fmt.Errorf("invalid duration string: %w", err)
	}

	fmt.Printf("Collecting feeds every %v", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			log.Printf("Error scraping feeds: %v\n", err)
		}
	}
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Could not get next feed to fetch: %w", err)
	}
	log.Println("Found next feed to fetch!")
	return scrapeFeed(s.db, nextFeed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) error {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("Could not mark feed as fetched: %w", err)
	}

	fetchedFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("Could not fetch %v feed: %w", feed.Name, err)
	}

	for _, post := range fetchedFeed.Channel.Item {
		createdPost, err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       sql.NullString{String: post.Title, Valid: post.Title != ""},
			Url:         sql.NullString{String: post.Link, Valid: post.Link != ""},
			Description: sql.NullString{String: post.Description, Valid: post.Description != ""},
			PublishedAt: sql.NullTime{Time: parsePubDate(post.PubDate), Valid: post.PubDate != ""},
			FeedID:      feed.ID,
		})
		if err != nil {
			return fmt.Errorf("Error saving post to database: %w", err)
		}
		if createdPost.ID == uuid.Nil && createdPost.Title.String == "" {
			log.Printf("Duplicate post URL detected; skipping: %s", post.Link)
		}
		log.Printf("New post saved: ID=%s, Title=%s", createdPost.ID, createdPost.Title.String)
	}

	return nil
}

// Helper function for parsing published date
func parsePubDate(pubDate string) time.Time {
	// Attempt standard RSS datetime format first
	const rssDateFormat = "Mon, 02 Jan 2006 15:04:05 MST"
	t, err := time.Parse(rssDateFormat, pubDate)
	if err == nil {
		return t
	}

	// Fallback format(s) in case of failure
	const alternativeFormat = "2006-01-02T15:04:05Z" // ISO-8601 as an example
	t, err = time.Parse(alternativeFormat, pubDate)
	if err == nil {
		return t
	}

	// If parsing fails, log the error and return a zero time
	log.Printf("Failed to parse pubDate: %s, error: %v", pubDate, err)
	return time.Time{} // A zero time to indicate failure
}
