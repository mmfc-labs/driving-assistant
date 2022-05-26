package lbs

import (
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	"github.com/shopspring/decimal"
)

// ConvCoordToQuadrilateral 将1个坐标转换成一个正方形（4个点坐标）
func ConvCoordToQuadrilateral(c drive.Coord, offset float64) []drive.Coord {
	offsetD := decimal.NewFromFloat(offset)
	c1 := drive.Coord{
		Lat: decimal.NewFromFloat(c.Lat).Add(offsetD).InexactFloat64(),
		Lon: decimal.NewFromFloat(c.Lon).Sub(offsetD).InexactFloat64(),
	}

	c2 := drive.Coord{
		Lat: decimal.NewFromFloat(c.Lat).Add(offsetD).InexactFloat64(),
		Lon: decimal.NewFromFloat(c.Lon).Add(offsetD).InexactFloat64(),
	}

	c3 := drive.Coord{
		Lat: decimal.NewFromFloat(c.Lat).Sub(offsetD).InexactFloat64(),
		Lon: decimal.NewFromFloat(c.Lon).Sub(offsetD).InexactFloat64(),
	}

	c4 := drive.Coord{
		Lat: decimal.NewFromFloat(c.Lat).Sub(offsetD).InexactFloat64(),
		Lon: decimal.NewFromFloat(c.Lon).Add(offsetD).InexactFloat64(),
	}

	return []drive.Coord{c1, c2, c3, c4}
}
