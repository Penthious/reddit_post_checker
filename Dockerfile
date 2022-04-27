FROM golang:1.16.2-alpine3.12 as base
RUN apk update && \
    apk add --update make && \
    apk add --update git && \
    apk add build-base
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
COPY ./config_docker.json ./config.json
RUN make build

FROM base as dev
RUN apk add --update npm
RUN npm install -g nodemon

FROM golang:1.16.2-alpine3.12 as prod
# FROM gcr.io/distroless/base as PROD // pull this in when I figure out why its breaking
# COPY --from=BASE /src/config_docker.json ./config.json
COPY --from=base /src/app .
CMD ["./app"]