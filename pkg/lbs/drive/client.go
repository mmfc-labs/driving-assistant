package drive

import (
	"fmt"
	geo "github.com/kellydunn/golang-geo"
)

type Client interface {
	GetRoutes(from, to Coord, avoids []Coord, avoidAreaOffset float64) ([]Route, error)
	GetDistanceMatrix(froms, tos []Coord) ([]int, error)
}

type Route struct {
	Points []Coord
}

type Coord struct {
	Lat float64 `json:"lat" yaml:"lat"`
	Lon float64 `json:"lon" yaml:"lon"`
}

func (c Coord) String() string {
	return fmt.Sprintf("%f,%f", c.Lat, c.Lon)
}

func (c Coord) GeoPoint() *geo.Point {
	return geo.NewPoint(c.Lat, c.Lon)
}
