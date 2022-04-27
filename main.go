package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	_ "embed"

	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

var (
	Version   string
	BuildTime string
	Logger zerolog.Logger
)

//go:embed config.json
var f []byte

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	Logger = zerolog.New(io.MultiWriter(os.Stdout)).With().Timestamp().
		Interface("BuildTime", BuildTime).
		Interface("Version", Version).
		Logger()
}
func loadConfig() (config, error) {
	var c config
	err := json.Unmarshal(f, &c)

	if err != nil {
		return config{}, fmt.Errorf("error unmarshalling json: %v", err)
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
		Logger.Info().Msg("test thing")
		fmt.Println("==========================================")
		fmt.Println("Title: \n", message)
		fmt.Println("User: \n", p.Author)
		fmt.Println("Body: \n", p.SelfText)
		fmt.Println("==========================================")
	}
	if c.Discord.Enabled {
		if err := n.notifyDiscord(message, c.Discord.Webhook); err != nil {
			return fmt.Errorf("error notifying discord: %v", err)
		}
	}
	if c.Browser.Enabled {
		if err := n.openBrowser(p.URL); err != nil {
			return fmt.Errorf("error notifying browser: %v", err)
		}
	}

	return nil
}

func (n notifier) notifyDiscord(message string, webhook string) error {
	values := map[string]string{"content": message}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("error marshalling discord json: %v", err)
	}

	req, err := http.NewRequest("POST", webhook, bytes.NewBuffer(jsonValue))
	if err != nil {
		return fmt.Errorf("error setting up request to discord webhook: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request to discord webhook: %v", err)
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
	Logger.Info().Msg("started")

	for {
		fmt.Printf("Version: %s, Build Time: %s\n", Version, BuildTime)

		c, err := loadConfig()
		if err != nil {
			fmt.Println("Failed to load config: ", err)
		}

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
}
