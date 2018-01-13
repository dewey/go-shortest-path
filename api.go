package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dewey/go-shortest-path/geo"
	"github.com/dewey/go-shortest-path/routing"
	"github.com/dewey/go-shortest-path/storage"
	"github.com/go-chi/chi"
)

func main() {
	gmr, err := geo.NewGoogleMapsRepository(os.Getenv("GOOGLE_MAPS_API_TOKEN"))
	if err != nil {
		log.Fatal(err.Error())
	}
	mem, err := storage.NewInMemoryRepository()
	if err != nil {
		log.Fatal(err.Error())
	}
	var routingService = routing.NewService(gmr, mem)

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("go-shortest-path"))
	})
	r.Mount("/", routing.NewHandler(routingService))
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err.Error())
	}
}
