package drive

import (
	"fmt"
	"github.com/mmfc-labs/driving-assistant/pkg/geo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvCoordToQuadrilateral(t *testing.T) {
	tests := []struct {
		name        string
		inputCoord  geo.Coord
		inputOffset float64
		want        []geo.Coord
	}{
		{
			name:        "case",
			inputCoord:  geo.NewCoord(1, 1),
			inputOffset: 0.1,
			want: []geo.Coord{
				geo.NewCoord(1.1, 0.9),
				geo.NewCoord(1.1, 1.1),
				geo.NewCoord(0.9, 0.9),
				geo.NewCoord(0.9, 1.1),
			},
		},
		{
			name:        "case",
			inputCoord:  geo.NewCoord(1.1, 1.1),
			inputOffset: 0.1,
			want: []geo.Coord{
				geo.NewCoord(1.2, 1),
				geo.NewCoord(1.2, 1.2),
				geo.NewCoord(1, 1),
				geo.NewCoord(1, 1.2),
			},
		},
		{
			name:        "case",
			inputCoord:  geo.NewCoord(1.000001, 1.000001),
			inputOffset: 0.000001,
			want: []geo.Coord{
				geo.NewCoord(1.000002, 1.000000),
				geo.NewCoord(1.000002, 1.000002),
				geo.NewCoord(1.000000, 1.000000),
				geo.NewCoord(1.000000, 1.000002),
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

func TestDecimal(t *testing.T) {
	c := geo.NewCoord(22.560413, 113.874613)
	cs := c.ToAvoidArea(0.000100)
	fmt.Println(geo.Coords(cs).String())
}
