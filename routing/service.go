package routing

import (
	"errors"

	"github.com/dewey/go-shortest-path/storage"

	"github.com/dewey/go-shortest-path/geo"
	"github.com/google/uuid"
)

// Service provides access to the routing functions
type Service interface {
	CalculateShortestPath(locations []geo.Location) (string, error)
	GetShortestPath(token string) (storage.Result, error)
}

type service struct {
	api geo.API
	db  storage.API
}

// NewService initializes a new routing service
func NewService(api geo.API, db storage.API) Service {
	return &service{
		api: api,
		db:  db,
	}
}

// CalculateShortestPath calculates the shortest paths by using the Google Maps API
func (s *service) CalculateShortestPath(locations []geo.Location) (string, error) {
	token := uuid.New()
	if err := s.db.Set(token.String(), storage.Result{Status: "in progress"}); err != nil {
		return "", err
	}
	go s.calculate(token.String(), locations)
	return token.String(), nil
}

func (s *service) calculate(token string, locations []geo.Location) {
	dir, err := s.api.CalculateDirections(locations)
	if err != nil {
		if err = s.db.Set(token, storage.Result{Status: "failure", Error: err.Error()}); err != nil {
			return
		}
		return
	}
	var tDistance, tTime int
	var wpOrder []int
	for _, r := range dir.Routes {
		for i, l := range r.Legs {
			if i != len(r.Legs)-1 {
				tDistance = tDistance + l.Distance.Value
				tTime = tTime + l.Duration.Value
			}
		}
		// We can assign that directly as "alternatives=true" is not set in our request
		// and so only one route will ever be returned
		wpOrder = r.WaypointOrder
	}
	// Reorder waypoints based on waypoint_order calculated by the API,
	// we have to ignore the origin there as it'll always be first
	wpOrdered := make([]geo.Location, len(locations)-1)
	for i, o := range wpOrder {
		wpOrdered[i] = locations[o+1]
	}
	var paths [][]string
	paths = append(paths, []string{locations[0].Latitude, locations[0].Longitude})
	for _, p := range wpOrdered {
		paths = append(paths, []string{p.Latitude, p.Longitude})
	}
	err = s.db.Set(token, storage.Result{
		Status:        "success",
		Paths:         paths,
		TotalDistance: tDistance,
		TotalTime:     tTime,
	})
	if err != nil {
		return
	}
}

// GetShortestPath gets the shortest path calculation's result from the database
func (s *service) GetShortestPath(token string) (storage.Result, error) {
	res, err := s.db.Get(token)
	if err != nil {
		return storage.Result{}, err
	}
	if res != nil {
		return *res, nil
	}
	return storage.Result{}, errors.New("invalid token")
}
