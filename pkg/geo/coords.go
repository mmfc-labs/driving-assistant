package geo

import (
	"fmt"
	"github.com/xyctruth/stream"
	"strings"
)

type Coords []Coord

func (c Coords) String() string {
	s := stream.NewSliceByMapping[Coord, string, string](c).Map(func(p Coord) string {
		return p.String()
	}).ToSlice()
	return fmt.Sprintf(strings.Join(s, ";"))
}
