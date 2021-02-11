package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

func loadConfig() (config, error) {
	var c config

	cf, err := os.Open("config.json")
	if err != nil {
		return config{}, fmt.Errorf("Error loading config: %v", err)
	}

	defer cf.Close()

	byteValue, err := ioutil.ReadAll(cf)
	if err != nil {
		return config{}, fmt.Errorf("Error sending to bytes: %v", err)
	}

	err = json.Unmarshal(byteValue, &c)
	if err != nil {
		return config{}, fmt.Errorf("Error unmarshalling config: %v", err)
	}

	return c, nil
}

type notifierInterface interface {
	notify(p *reddit.Post, c config) error
	notifyDiscord(message string, webhook string) error
	openBrowser(url string) error
}

func newNotifier() notifierInterface {
	return notifier{}
}

func (n notifier) notify(p *reddit.Post, c config) error {
	message := fmt.Sprintf("%s | <%s>", p.Title, p.URL)

	if c.Debug {
		fmt.Println(message)
	}
	if c.Discord.Enabled {
		if err := n.notifyDiscord(message, c.Discord.Webhook); err != nil {
			return fmt.Errorf("Error notifing discord: %v\n", err)
		}
	}
	if c.Browser.Enabled {
		if err := n.openBrowser(p.URL); err != nil {
			return fmt.Errorf("Error notifing discord: %v\n", err)
		}
	}

	return nil
}

func (n notifier) notifyDiscord(message string, webhook string) error {
	values := map[string]string{"content": message}
	jsonValue, err := json.Marshal(values)
	req, err := http.NewRequest("POST", webhook, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return nil
	}

	return nil
}

// Does not work when using docker, as docker does not have context of host machine
func (n notifier) openBrowser(url string) error {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start()
}

func (r *redditBot) Post(p *reddit.Post) error {
	n := newNotifier()
	if r.config.Keywords.Enabled {
		found := false
		for _, keyword := range r.config.Keywords.Terms {
			if strings.Contains(strings.ToLower(p.Title), strings.ToLower(keyword)) {
				found = true
				break
			}
		}
		if found {
			return n.notify(p, r.config)
		}
	} else {
		return n.notify(p, r.config)
	}
	return nil
}

func main() {
	c, err := loadConfig()

	bot, err := reddit.NewBot(reddit.BotConfig{
		Agent: c.Reddit.UserAgent,
		App: reddit.App{
			ID:       c.Reddit.ClientID,
			Secret:   c.Reddit.ClientSecret,
			Username: c.Reddit.Username,
			Password: c.Reddit.Password,
		},
		Rate: 0,
	})
	if err != nil {
		fmt.Println("Failed to create bot handler: ", err)
		return
	}

	rc := graw.Config{Subreddits: c.Subreddits}
	handler := &redditBot{bot: bot, config: c}

	if c.Keywords.Enabled {
		fmt.Printf("Listening for new post in subreddits: %v for keywords with %v\n", c.Subreddits, c.Keywords.Terms)
	} else {
		fmt.Printf("Listening for new post in subreddits: %v\n", c.Subreddits)
	}
	if _, wait, err := graw.Run(handler, bot, rc); err != nil {
		fmt.Println("Failed to start graw run: ", err)
	} else {
		fmt.Println("graw run failed: ", wait())
	}
}
