package config

import (
	"fmt"
	"github.com/mmfc-labs/driving-assistant/pkg/geo"
	"github.com/mmfc-labs/driving-assistant/pkg/probe"
	"testing"
)

func TestPrintProbes(t *testing.T) {
	err := LoadConfig("../../config.yaml", func(setting Setting, probeManager probe.Manager) {
		points := make([]geo.Coord, 0)
		for _, p := range probeManager.Probes {
			points = append(points, p.Coord)
		}
		fmt.Println(geo.Coords(points).String())
	})

	if err != nil {
		t.Error(err)
	}
}
