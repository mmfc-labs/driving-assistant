package geo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
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
			inputCoord:  NewCoord(1, 1),
			inputOffset: 0.1,
			want: []Coord{
				NewCoord(1.1, 0.9),
				NewCoord(1.1, 1.1),
				NewCoord(0.9, 0.9),
				NewCoord(0.9, 1.1),
			},
		},
		{
			name:        "case",
			inputCoord:  NewCoord(1.1, 1.1),
			inputOffset: 0.1,
			want: []Coord{
				NewCoord(1.2, 1),
				NewCoord(1.2, 1.2),
				NewCoord(1, 1),
				NewCoord(1, 1.2),
			},
		},
		{
			name:        "case",
			inputCoord:  NewCoord(1.000001, 1.000001),
			inputOffset: 0.000001,
			want: []Coord{
				NewCoord(1.000002, 1.000000),
				NewCoord(1.000002, 1.000002),
				NewCoord(1.000000, 1.000000),
				NewCoord(1.000000, 1.000002),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.inputCoord.ToAvoidArea(tt.inputOffset)
			assert.Equal(t, tt.want, got)
		})
	}

}

func TestBearing(t *testing.T) {
	a1 := NewCoord(23.265333, 113.354559)
	a2 := NewCoord(23.265136, 113.354565)
	gap := math.Abs(a1.BearingTo(a2) - a2.BearingTo(a1))
	assert.Equal(t, true, gap > 179 && gap < 181)
}

func TestPrintDist(t *testing.T) {
	fmt.Println(NewCoord(22.590208, 113.890855).Distance(NewCoord(22.590322, 113.890853)))
}

func TestPrintQuadrilateral(t *testing.T) {
	c := NewCoord(22.590208, 113.890855)
	fmt.Println(Coords(c.ToAvoidArea(0.000100)).String())
}
