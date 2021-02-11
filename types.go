package main

import "github.com/turnage/graw/reddit"

type config struct {
	Reddit struct {
		UserAgent    string `json:"user_agent"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Username     string `json:"username"`
		Password     string `json:"password"`
	}
	Discord struct {
		Webhook string `json:"webhook"`
		Enabled bool   `json:"enabled"`
	} `json:"discord"`
	Browser struct {
		Enabled bool `json:"enabled"`
	}
	Keywords struct {
		Terms   []string `json:"terms"`
		Enabled bool     `json:"enabled"`
	}
	Debug      bool     `json:"debug"`
	Subreddits []string `json:"subreddits"`
}

type redditBot struct {
	bot    reddit.Bot
	config config
}

type notifier struct{}
