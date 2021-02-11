# Reddit post checker
This script will monitor a list of subreddits for new posts. It can send messages via discord or terminal, if enabled it will also auto open your default browser with the current post url.

# Setup
You then need to configure copy the `config.json.template` file to `config.json` and update it with the api keys and settings you want. 

With docker:
> `docker-compose up -d`

Without Docker: (for when you want to open with browser as docker does not allow you to open browser)
> `go run .`

Building the binary and running in the background:
> `make build && nohup ./reddit_post_checker &`

### Reddit 
Create your bot / mod account and go to: https://www.reddit.com/prefs/apps/. Create a script app and it will give you the credentials you need. 

### Discord
In your server settings create a webhook: https://support.discordapp.com/hc/en-us/articles/228383668-Intro-to-Webhooks

### Browser
Opens your default browser with the current url of the reddit post, Only works if you enable it in config and are not using docker.