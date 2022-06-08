package probe

import (
	"github.com/mmfc-labs/driving-assistant/pkg/geo"
	"github.com/xyctruth/stream"
)

type Manager struct {
	Probes []Probe `yaml:"probes"`
}

// Near 获取附近的探头
// distanceFilter 单位:KM
func (p *Manager) Near(cur geo.Coord, distanceFilter float64) []Probe {
	return stream.NewSlice(p.Probes).Filter(func(point Probe) bool {
		distance := cur.Distance(point.Coord)
		return distance < distanceFilter
	}).ToSlice()
}

// All 全部探头
func (p *Manager) All() []Probe {
	return p.Probes
}

func (p *Manager) CalculateToward() {
	for _, probe := range p.Probes {
		for i, _ := range probe.Towards {
			probe.Towards[i].Value = probe.BearingTo(probe.Towards[i].Coord)
		}
	}
}
