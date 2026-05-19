package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/lynndixon/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return errors.New("usage: browse <limit of number of posts>")
	}

	limit := 2

	if len(cmd.Args) == 1 {
		if limitNum, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = limitNum
		} else {
			return fmt.Errorf("could not parse optional limit argument into an integer. please enter a valid integer if you wish to limit the number of posts to browse: %w", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		ID:    user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		return fmt.Errorf("Could not fetch posts for user %v: %w", user.Name, err)
	}

	fmt.Printf("Found %v for user %v:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf(" * %v posted at %v:\n", post.Feedname, post.PublishedAt.Time.Format("Sat Mar 7 11:06:39 GMT 2015"))
		fmt.Printf(" ---\"%s\"--- \n", post.Title.String)
		fmt.Printf("	%s\n", post.Description.String)
		fmt.Printf(" Link: %s\n", post.Url.String)
		fmt.Println("----------------------------------------")
	}

	return nil
}

func printPost(post database.Post) {
	fmt.Printf(" * Title: \"%v\"\n", post.Title)
	fmt.Printf(" * Description:	%v\n", post.Description)
}
