package port

import "errors"

var ErrInvalidID = errors.New("port id is required")

type Port struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Province    string    `json:"province"`
	Coordinates []float64 `json:"coordinates"`
	Timezone    string    `json:"timezone"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Code        string    `json:"code"`
}

func (p Port) Validate() error {
	if p.ID == "" {
		return ErrInvalidID
	}
	return nil
}
