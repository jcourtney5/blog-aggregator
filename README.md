# Blog Aggregator 

Command line RSS feed post tracker and aggregator project

---

### Requirements
* Golang version 1.25.3 or higher
* Postgres version 16.10 or higher

---

### How To Build And Run
* Download code from github
* Run "go install ." which will build and copy the blog-aggregator binary to your go bin folder
* Create a ".gatorconfig.json" file in your root user directory (ex: ~/.gatorconfig.json)
* Put the postgres connection string in the json file
    ```json
    {
      "db_url": "<postgres_connection_string>?sslmode=disable"
    }
    ```
* Run the program with "blog-aggregator \<command> \<args>"

---

### List of commands
* register \<username>
  * Add the *username* to the system
* login \<username>
  * Set username as the current user (*many commands user the current user*)
* users
  * List all the users in the system
* reset
  * Reset all the data in the DB to start over
* addfeed \<name> <url>
  * Add an RSS feed to the system and have current user follow it
* feeds
  * List all the RSS feeds in the system
* follow \<url>
  * Follow the feed for the current user
* unfollow \<url>
  * Unfollow the feed for the current user
* following
  * List all the feeds the current user is following
* agg \<time_between_requests>
  * Start the fetch loop to get all the latest posts for each RSS feed the current user follows.
  * time_between_requests is the time gap between updating feeds (ex: 30s, 1m, 2m, 1h)
* browse \<limit>
  * Get the most recent posts for the current user up to *limit* count
