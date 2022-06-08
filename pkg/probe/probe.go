package probe

import (
	"github.com/mmfc-labs/driving-assistant/pkg/geo"
)

type Probe struct {
	geo.Coord `mapstructure:",squash"`
	Towards   []Toward `yaml:"towards" json:"towards"` // 探头朝向的坐标
}

type Toward struct {
	geo.Coord `mapstructure:",squash"`
	Value     float64 `json:"value" yaml:"value"` // 探头朝向的角度值
}
