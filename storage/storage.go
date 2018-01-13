package storage

// API is the interface for a storage API
type API interface {
	// Get will get a result from the map by using the ID as the key
	Get(id string) (*Result, error)
	// Set will store a new result in the map
	Set(id string, res Result) error
}

// Result contains a routing response
type Result struct {
	Status        string     `json:"status"`
	Error         string     `json:"error,omitempty"`
	Paths         [][]string `json:"path,omitempty"`
	TotalDistance int        `json:"total_distance,omitempty"`
	TotalTime     int        `json:"total_time,omitempty"`
}
