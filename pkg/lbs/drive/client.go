package drive

import "fmt"

type Client interface {
	GetRoutes(from, to Coord, avoids []Coord, avoidAreaOffset float64) ([]Route, error)
	GetDistanceMatrix(froms, tos []Coord) ([]int, error)
}

type Route struct {
	Points []Coord
}

type Coord struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (c Coord) String() string {
	return fmt.Sprintf("%f,%f", c.Lat, c.Lon)
}
