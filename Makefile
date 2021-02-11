# Description: The Makefile is to automate the common build tasks
# 
# Example: make build
# Note: Use modern make to get additional help context from the comments (https://github.com/tj/mmake)

# Build the application for deployment
build:
	@go build .
.PHONY: build

# Run appilication in docker
run:
	@docker-compose up
.PHONY: run