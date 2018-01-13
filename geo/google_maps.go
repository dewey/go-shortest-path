package geo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type googleMapsRepository struct {
	Token string
}

// Directions contains the API response from the Google Maps API for a direction request
type Directions struct {
	GeocodedWaypoints []struct {
		GeocoderStatus string   `json:"geocoder_status"`
		PlaceID        string   `json:"place_id"`
		Types          []string `json:"types"`
	} `json:"geocoded_waypoints"`
	Routes []struct {
		Bounds struct {
			Northeast struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"northeast"`
			Southwest struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"southwest"`
		} `json:"bounds"`
		Copyrights string `json:"copyrights"`
		Legs       []struct {
			Distance struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"distance"`
			Duration struct {
				Text  string `json:"text"`
				Value int    `json:"value"`
			} `json:"duration"`
			EndAddress  string `json:"end_address"`
			EndLocation struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"end_location"`
			StartAddress  string `json:"start_address"`
			StartLocation struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"start_location"`
			Steps []struct {
				Distance struct {
					Text  string `json:"text"`
					Value int    `json:"value"`
				} `json:"distance"`
				Duration struct {
					Text  string `json:"text"`
					Value int    `json:"value"`
				} `json:"duration"`
				EndLocation struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"end_location"`
				HTMLInstructions string `json:"html_instructions"`
				Polyline         struct {
					Points string `json:"points"`
				} `json:"polyline"`
				StartLocation struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"start_location"`
				TravelMode string `json:"travel_mode"`
				Maneuver   string `json:"maneuver,omitempty"`
			} `json:"steps"`
			TrafficSpeedEntry []interface{} `json:"traffic_speed_entry"`
			ViaWaypoint       []interface{} `json:"via_waypoint"`
		} `json:"legs"`
		OverviewPolyline struct {
			Points string `json:"points"`
		} `json:"overview_polyline"`
		Summary       string        `json:"summary"`
		Warnings      []interface{} `json:"warnings"`
		WaypointOrder []int         `json:"waypoint_order"`
	} `json:"routes"`
	Status string `json:"status"`
}

// NewGoogleMapsRepository returns a newly initialized Google Maps repository
func NewGoogleMapsRepository(token string) (API, error) {
	if token == "" {
		return nil, errors.New("GOOGLE_MAPS_API_TOKEN not allowed to be empty")
	}
	return &googleMapsRepository{
		Token: token,
	}, nil
}

func (r *googleMapsRepository) CalculateDirections(locations []Location) (Directions, error) {
	client := &http.Client{}
	u, err := url.Parse("https://maps.googleapis.com/maps/api/directions/json")
	if err != nil {
		return Directions{}, err
	}
	q := u.Query()
	q.Add("key", r.Token)
	q.Add("mode", "driving")
	// Origin and destination are the same in this case as we are visiting the waypoints and then going back home
	q.Add("origin", fmt.Sprintf("%s,%s", locations[0].Latitude, locations[0].Longitude))
	q.Add("destination", fmt.Sprintf("%s,%s", locations[0].Latitude, locations[0].Longitude))

	// We have to build the waypoint list by removing the first location (origin) and prefixing the parameter with
	// optimize:true to enable waypoint optimization in the Google Maps API
	var wp bytes.Buffer
	for i, loc := range locations {
		if i > 0 {
			wp.WriteString(fmt.Sprintf("%s,%s", loc.Latitude, loc.Longitude))
		}
		if i != (len(locations) - 1) {
			wp.WriteString("|")
		}
	}
	q.Add("waypoints", fmt.Sprintf("optimize:true|%s", wp.String()))
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return Directions{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return Directions{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Directions{}, fmt.Errorf("unexpected status code from google maps api: %d", resp.StatusCode)
	}

	var dir Directions
	if err := json.NewDecoder(resp.Body).Decode(&dir); err != nil {
		return Directions{}, err
	}

	// If the routing is impossible or a coordinate doesn't exist Google's API returns empty routes and status "ZERO_RESULTS"
	if len(dir.Routes) < 1 {
		return Directions{}, errors.New("no routes returned from API")
	}

	return dir, nil
}
