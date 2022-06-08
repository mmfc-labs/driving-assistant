package geo

import (
	"fmt"
	geo "github.com/kellydunn/golang-geo"
	"github.com/shopspring/decimal"
)

type Coord struct {
	Lat float64 `json:"lat" yaml:"lat"`
	Lon float64 `json:"lon" yaml:"lon"`
}

func NewCoord(lat, lon float64) Coord {
	return Coord{
		Lat: lat,
		Lon: lon,
	}
}

func (c Coord) String() string {
	return fmt.Sprintf("%f,%f", c.Lat, c.Lon)
}

func (c Coord) Distance(c2 Coord) float64 {
	return c.geoPoint().GreatCircleDistance(c2.geoPoint())
}

func (c Coord) BearingTo(c2 Coord) float64 {
	return c.geoPoint().BearingTo(c2.geoPoint())
}

func (c Coord) geoPoint() *geo.Point {
	return geo.NewPoint(c.Lat, c.Lon)
}

// ToAvoidArea 将1个坐标转换成一个正方形（4个点坐标）
func (c Coord) ToAvoidArea(offset float64) []Coord {
	offsetD := decimal.NewFromFloat(offset)
	c1 := Coord{
		Lat: decimal.NewFromFloat(c.Lat).Add(offsetD).Round(6).InexactFloat64(),
		Lon: decimal.NewFromFloat(c.Lon).Sub(offsetD).Round(6).InexactFloat64(),
	}

	c2 := Coord{
		Lat: decimal.NewFromFloat(c.Lat).Add(offsetD).Round(6).InexactFloat64(),
		Lon: decimal.NewFromFloat(c.Lon).Add(offsetD).Round(6).InexactFloat64(),
	}

	c3 := Coord{
		Lat: decimal.NewFromFloat(c.Lat).Sub(offsetD).Round(6).InexactFloat64(),
		Lon: decimal.NewFromFloat(c.Lon).Sub(offsetD).Round(6).InexactFloat64(),
	}

	c4 := Coord{
		Lat: decimal.NewFromFloat(c.Lat).Sub(offsetD).Round(6).InexactFloat64(),
		Lon: decimal.NewFromFloat(c.Lon).Add(offsetD).Round(6).InexactFloat64(),
	}

	return []Coord{c1, c2, c3, c4}
}
