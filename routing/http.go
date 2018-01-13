package routing

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dewey/go-shortest-path/geo"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// NewHandler initializes a new routing http Handler
func NewHandler(s Service) *chi.Mux {
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Post("/route", calculateShortestPathHandler(s))
	r.Get("/route/{token}", getShortestPathResultHandler(s))
	return r
}

type payload struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func calculateShortestPathHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pl payload
		var req [][]string
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Converting request format into our internal one which is nicer to work with
		var locations []geo.Location
		for _, loc := range req {
			if len(loc) == 2 {
				locations = append(locations, geo.Location{Latitude: loc[0], Longitude: loc[1]})
			} else {
				pl.Error = errors.New("a geo location can only contain 2 elements").Error()
				b, err := json.Marshal(pl)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(b)
				return
			}
		}
		if len(locations) < 2 {
			pl.Error = errors.New("not enough values to calculate paths").Error()
			b, err := json.Marshal(pl)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(b)
			return
		}
		token, err := s.CalculateShortestPath(locations)
		if err != nil {
			pl.Error = err.Error()
			b, err := json.Marshal(pl)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(b)
			return
		}
		pl.Token = token
		b, err := json.Marshal(pl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func getShortestPathResultHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pl payload
		token := chi.URLParam(r, "token")
		res, err := s.GetShortestPath(token)
		if err != nil {
			pl.Error = err.Error()
			b, err := json.Marshal(pl)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}
		b, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if res.Status == "in progress" {
			w.WriteHeader(http.StatusAccepted)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		w.Write(b)
	}
}
