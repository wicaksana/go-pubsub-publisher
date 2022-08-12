# syntax=docker/dockerfile:1

### Build
FROM golang:1.18.4-bullseye AS build

RUN mkdir /app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o publish .


### Deploy
FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /app/publish /publish
USER nonroot:nonroot
ENTRYPOINT ["/publish"]