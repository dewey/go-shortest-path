package geo

// API is the interface for a geo API
type API interface {
	CalculateDirections(locations []Location) (Directions, error)
}

// Location defines a geo location
type Location struct {
	Latitude  string
	Longitude string
}
