# gator

## Gator is an RSS aggreGATOR (Get it?) that I wrote as part of a guided project on Boot.dev, a backend programming course.

Please forgive any funky markdown style - I've only written a few ReadMes!

### To use gator, you'll need two things first: Postgres and Go.

macOS with brew:

`brew install postgresql@15`

Linux / WSL (Debian):

```
sudo apt update
sudo apt install postgresql postgresql-contrib
```

### Next, install the gator CLI:

`go install gator`


### Gator is only configured to run locally. It needs a config file on your computer to work.


If there isn't one already, create a gatorconfig.json file in gator's root directory, like so:

```
{
   "db_url": "connection_string_goes_here",
   "current_user_name": "username_goes_here"
}
```

This keeps track of who is currently logged in, and the connection credentials for the Postgres database.

A connection string looks like this: protocol://username:password@host:port/database

Some examples:

macOS (no password, your username): `postgres://wagslane:@localhost:5432/gator`
Linux (password set to postgres, user also set to postgres): `postgres://postgres:postgres@localhost:5432/gator`

You can test your connection string by running

`psql "postgres://wagslane:@localhost:5432/gator"`

Note: There's no user-based authentication for this app. If someone has the database credentials, they can act as any user. 


### JSON configured? Good!

### Gator Commands:


Gator can take a number of commands, but you must first `register` a username and `login` to set that name as the active user, so you can follow feeds.

`addfeed` <feed name> <url> - add a feed to the database's list of known feeds.

`agg` <time duration string: 1h, 1m, 1s etc> - begin aggregating a users feeds for browsing later. CAUTION - do not set the refresh duration to be too short, lest you DOS a server!

`browse` <limit of number of posts> - browse a user's post by a given number of posts. Defaults to two posts per browse if not given the optional argument.

`feeds` - retrieves a list of feeds on the server.

`follow` <url> - follow a feed at a given URL.

`following` - retrieves a list of the current user's followed feeds.

`login` <username> - log in as the provided user.

`register` <username> - register a new user.

`reset` - resets the database - mainly helpful for testing.

`unfollow` <url> unfollow a feed at the given URL.

`users` - retrieves a list of registered users on the database.


### Happy RSS Chomping!