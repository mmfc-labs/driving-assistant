package drive

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvCoordToQuadrilateral(t *testing.T) {
	tests := []struct {
		name        string
		inputCoord  Coord
		inputOffset float64
		want        []Coord
	}{
		{
			name:        "case",
			inputCoord:  Coord{Lat: 1, Lon: 1},
			inputOffset: 0.1,
			want: []Coord{
				{Lat: 1.1, Lon: 0.9},
				{Lat: 1.1, Lon: 1.1},
				{Lat: 0.9, Lon: 0.9},
				{Lat: 0.9, Lon: 1.1},
			},
		},
		{
			name:        "case",
			inputCoord:  Coord{Lat: 1.1, Lon: 1.1},
			inputOffset: 0.1,
			want: []Coord{
				{Lat: 1.2, Lon: 1},
				{Lat: 1.2, Lon: 1.2},
				{Lat: 1, Lon: 1},
				{Lat: 1, Lon: 1.2},
			},
		},
		{
			name:        "case",
			inputCoord:  Coord{Lat: 1.000001, Lon: 1.000001},
			inputOffset: 0.000001,
			want: []Coord{
				{Lat: 1.000002, Lon: 1.000000},
				{Lat: 1.000002, Lon: 1.000002},
				{Lat: 1.000000, Lon: 1.000000},
				{Lat: 1.000000, Lon: 1.000002},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvCoordToAvoidArea(tt.inputCoord, tt.inputOffset)
			assert.Equal(t, tt.want, got)
		})
	}

}

func TestDecimal(t *testing.T) {
	c := Coord{Lat: 22.560413, Lon: 113.874613}
	cs := ConvCoordToAvoidArea(c, 0.000100)
	fmt.Println(FmtCoord(cs...))
}
