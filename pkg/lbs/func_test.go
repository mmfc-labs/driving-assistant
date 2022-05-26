package lbs

import (
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvCoordToQuadrilateral(t *testing.T) {
	tests := []struct {
		name        string
		inputCoord  drive.Coord
		inputOffset float64
		want        []drive.Coord
	}{
		{
			name:        "case",
			inputCoord:  drive.Coord{Lat: 1, Lon: 1},
			inputOffset: 0.1,
			want: []drive.Coord{
				{Lat: 1.1, Lon: 0.9},
				{Lat: 1.1, Lon: 1.1},
				{Lat: 0.9, Lon: 0.9},
				{Lat: 0.9, Lon: 1.1},
			},
		},
		{
			name:        "case",
			inputCoord:  drive.Coord{Lat: 1.1, Lon: 1.1},
			inputOffset: 0.1,
			want: []drive.Coord{
				{Lat: 1.2, Lon: 1},
				{Lat: 1.2, Lon: 1.2},
				{Lat: 1, Lon: 1},
				{Lat: 1, Lon: 1.2},
			},
		},
		{
			name:        "case",
			inputCoord:  drive.Coord{Lat: 1.000001, Lon: 1.000001},
			inputOffset: 0.000001,
			want: []drive.Coord{
				{Lat: 1.000002, Lon: 1.000000},
				{Lat: 1.000002, Lon: 1.000002},
				{Lat: 1.000000, Lon: 1.000000},
				{Lat: 1.000000, Lon: 1.000002},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvCoordToQuadrilateral(tt.inputCoord, tt.inputOffset)
			assert.Equal(t, tt.want, got)
		})
	}

}
