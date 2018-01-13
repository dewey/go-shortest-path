# go-shortest-path

Simple API based on the Google Maps Directions API, nothing to see here. Please move along!

![](https://i.imgur.com/It5nzO8.png)

## Missing

- Logging
- Persistant storage backend (currently just in-memory, but easily extendable by implementing the storage interface)
- Compile Go binary and copy it into scratch image for smaller image size
- Export Prometheus metrics

## Build

```zsh
$ docker build -f Dockerfile -t go-shortest-path:latest .
$ docker tag 5bd59077e93d tehwey/go-shortest-path
$ docker push tehwey/go-shortest-path

```

## Run

`$ docker-compose -f docker-compose.yml up -d`
