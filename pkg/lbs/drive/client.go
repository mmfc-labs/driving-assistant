package drive

import (
	"github.com/mmfc-labs/driving-assistant/pkg/geo"
)

type Client interface {
	GetRoutes(from, to geo.Coord, avoids []geo.Coord, avoidAreaOffset float64) ([]Route, error)
	GetDistanceMatrix(froms, tos []geo.Coord) ([]int, error)
}

type Route struct {
	Points []geo.Coord
}
