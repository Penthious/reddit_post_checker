FROM golang:1.15.6-alpine3.12 as BASE
RUN apk update && \
    apk add --update make && \
    apk add --update git && \
    apk add build-base
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN make build

FROM BASE as DEV
RUN apk add --update npm
RUN npm install -g nodemon

# "Distroless" images contain only your application and its runtime dependencies. 
# They do not contain package managers, shells or any other programs you would expect to find in a standard Linux distribution.
FROM gcr.io/distroless/base-debian10 as PROD
COPY --from=BASE /main /
CMD ["./reddit_post_checker"]