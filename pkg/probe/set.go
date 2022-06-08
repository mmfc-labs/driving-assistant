package probe

import "github.com/mmfc-labs/driving-assistant/pkg/geo"

type Set struct {
	m     map[geo.Coord]struct{}
	slice []geo.Coord
}

func NewProbeSet() *Set {
	return &Set{
		m:     make(map[geo.Coord]struct{}),
		slice: make([]geo.Coord, 0),
	}
}

func (s *Set) Add(c geo.Coord) {
	if _, ok := s.m[c]; !ok {
		s.m[c] = struct{}{}
		s.slice = append(s.slice, c)
	}
}

func (s *Set) Len() int {
	return len(s.slice)
}

func (s *Set) ToSlice() []geo.Coord {
	return s.slice
}

func (s *Set) ToAvoidArea(offset float64) [][]geo.Coord {
	avoidAreas := make([][]geo.Coord, 0)
	for _, p := range s.slice {
		avoidAreas = append(avoidAreas, p.ToAvoidArea(offset))
	}
	return avoidAreas
}
