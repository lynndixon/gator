package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/lynndixon/gator/internal/config"
	"github.com/lynndixon/gator/internal/database"
)

type state struct {
	db    *database.Queries
	rawDB *sql.DB
	cfg   *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("could not read config file: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("Could not open a connection to database at: %v", cfg.DbUrl)
	}
	defer db.Close()
	dbQueries := database.New(db)

	programState := &state{
		db:    dbQueries,
		rawDB: db,
		cfg:   &cfg,
	}

	cmds := commands{
		RegisteredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("addfeed", handlerAddFeed)
	cmds.register("agg", handlerAgg)
	cmds.register("login", handlerLogin)
	cmds.register("feeds", handlerGetFeeds)
	cmds.register("follow", handlerFollow)
	cmds.register("following", handlerFollowing)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)

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
