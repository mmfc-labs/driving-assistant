package lbs

import (
	"github.com/mmfc-labs/driving-assistant/pkg/config"
	"github.com/mmfc-labs/driving-assistant/pkg/geo"
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

	lbs := NewLBS(s, probe.Manager{})

	tests := []struct {
		name  string
		cur   geo.Coord
		next  geo.Coord
		probe probe.Probe
		want  bool
	}{
		{
			name: "不同方向探头",
			cur:  geo.NewCoord(22.578057, 113.913546),
			next: geo.NewCoord(22.577495, 113.914098),
			probe: probe.Probe{
				Coord:   geo.NewCoord(22.577590, 113.914101),
				Towards: []probe.Toward{{Coord: geo.NewCoord(22.577651, 113.914054), Value: 0}}},
			want: false,
		},
		{
			name: "相同方向探头",
			cur:  geo.NewCoord(22.578057, 113.913546),
			next: geo.NewCoord(22.577495, 113.914098),
			probe: probe.Probe{
				Coord:   geo.NewCoord(22.577590, 113.914101),
				Towards: []probe.Toward{{Coord: geo.NewCoord(22.577566, 113.914139), Value: 0}}},
			want: true,
		},
		{
			name: "双方向探头",
			cur:  geo.NewCoord(22.578057, 113.913546),
			next: geo.NewCoord(22.577495, 113.914098),
			probe: probe.Probe{
				Coord:   geo.NewCoord(22.577590, 113.914101),
				Towards: []probe.Toward{{Coord: geo.NewCoord(22.577651, 113.914054), Value: 0}, {Coord: geo.NewCoord(22.577566, 113.914139), Value: 0}}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, lbs.isAvoid(tt.cur, tt.next, tt.probe))
		})
	}

}
