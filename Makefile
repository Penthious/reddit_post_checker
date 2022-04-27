# Description: The Makefile is to automate the common build tasks
# 
# Example: make build
# Note: Use modern make to get additional help context from the comments (https://github.com/tj/mmake)

# Build the application for deployment
VERSION := $(shell git rev-parse --short HEAD)
NOW := $(shell date +'%Y-%m-%d_%T')

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

build:
	@env GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(NOW)'" -o app .
.PHONY: build

# Run appilication in docker
run:
	@docker-compose up
.PHONY: run

# Tags the docker image for prod
tag:
	@docker build --target=prod -t penthious/reddit_post_checker:latest .

# Pushes up docker image to docker hub (requires login before this command)
prod: tag
	@docker push penthious/reddit_post_checker:latest
