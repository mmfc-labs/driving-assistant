package lbs

import (
	"github.com/mmfc-labs/driving-assistant/pkg/config"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	"github.com/mmfc-labs/driving-assistant/pkg/probe"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsAvoid(t *testing.T) {
	s := config.Setting{
		LBSKey:      "KN6BZ-G526D-JAI4V-PGSJ2-6L5U6-YYFBV",
		Offset:      6,
		TowardRange: 90,
	}

	lbs := NewLBS(s, probe.ProbeManager{})

	tests := []struct {
		name  string
		cur   drive.Coord
		next  drive.Coord
		probe probe.Probe
		want  bool
	}{
		{
			name: "case",
			cur:  drive.Coord{Lat: 22.578057, Lon: 113.913546},
			next: drive.Coord{Lat: 22.577495, Lon: 113.914098},
			probe: probe.Probe{
				Coord:   drive.Coord{Lat: 22.577590, Lon: 113.914101},
				Towards: []probe.Toward{{Coord: drive.Coord{Lat: 22.577651, Lon: 113.914054}, Value: 0}}},
			want: false,
		},
		{
			name: "case",
			cur:  drive.Coord{Lat: 22.578057, Lon: 113.913546},
			next: drive.Coord{Lat: 22.577495, Lon: 113.914098},
			probe: probe.Probe{
				Coord:   drive.Coord{Lat: 22.577590, Lon: 113.914101},
				Towards: []probe.Toward{{Coord: drive.Coord{Lat: 22.577566, Lon: 113.914139}, Value: 0}}},
			want: true,
		},
		{
			name: "case",
			cur:  drive.Coord{Lat: 22.578057, Lon: 113.913546},
			next: drive.Coord{Lat: 22.577495, Lon: 113.914098},
			probe: probe.Probe{
				Coord:   drive.Coord{Lat: 22.577590, Lon: 113.914101},
				Towards: []probe.Toward{{Coord: drive.Coord{Lat: 22.577651, Lon: 113.914054}, Value: 0}, {Coord: drive.Coord{Lat: 22.577566, Lon: 113.914139}, Value: 0}}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, lbs.isAvoid(tt.cur, tt.next, tt.probe))
		})
	}

}
