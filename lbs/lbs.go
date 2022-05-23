package lbs

type LBS interface {
	GetRoutes(from1, from2, to1, to2 float64) []Route
}

type Plan struct {
	Status int
	Result struct {
		Routes []Route
	}
}

type Route struct {
	Points []Coord
}

type Coord struct {
	Lat float64
	Lon float64
}
