package drive

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/xyctruth/stream"
	"strings"
)

// ConvCoordToAvoidArea 将1个坐标转换成一个正方形（4个点坐标）
func ConvCoordToAvoidArea(c Coord, offset float64) []Coord {
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

func FmtCoord(points ...Coord) string {
	s := stream.NewSliceByMapping[Coord, string, string](points).Map(func(p Coord) string {
		return p.String()
	}).ToSlice()
	return fmt.Sprintf(strings.Join(s, ";"))
}
