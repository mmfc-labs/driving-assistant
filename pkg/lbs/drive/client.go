package drive

type Client interface {
	GetRoutes(from, to Coord) ([]Route, error)
	GetDistanceMatrix(froms, tos []Coord) ([]int, error)
}

type Route struct {
	Points []Coord
}

type Coord struct {
	Lat float64
	Lon float64
}
