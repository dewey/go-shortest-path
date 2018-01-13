# go-shortest-path

Simple API based on the Google Maps Directions API, nothing to see here. Please move along!

![](https://i.imgur.com/It5nzO8.png)

## Missing

- Logging
- Persistant storage backend (currently just in-memory, but easily extendable by implementing the storage interface)
- Export Prometheus metrics
- Add tests

## Build

```zsh
GOOS=linux go build api.go
docker build -f Dockerfile -t go-shortest-path:latest .
docker tag 5bd59077e93d tehwey/go-shortest-path
docker push tehwey/go-shortest-path
```

## Run

`$ GOOGLE_MAPS_API_TOKEN=TOKEN_HERE docker-compose -f docker-compose.yml up -d`
