package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/mmfc-labs/driving-assistant/pkg/probe"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//const TencentKey = "RB7BZ-CGUKW-AU4RO-RIFUW-57GFS-L4BP7"
const TencentKey = "KN6BZ-G526D-JAI4V-PGSJ2-6L5U6-YYFBV"
const Offset = 5                 // 路面距离计算偏移量，单位米
const AvoidAreaOffset = 0.000100 // 生成四边形避让区偏移量, 单位经纬度

// LoadConfig watch configPath change, callback fn
func LoadConfig(configPath string, fn func(s Setting, probe probe.Probe)) error {
	var err error
	var setting Setting
	var p probe.Probe

	conf := viper.New()
	conf.SetConfigFile(configPath)
	conf.SetConfigType("yaml")

	if err = conf.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	if err = conf.UnmarshalKey("setting", &setting); err != nil {
		return fmt.Errorf("fatal error config setting: %w", err)
	}
	if err = conf.UnmarshalKey("probe", &p); err != nil {
		return fmt.Errorf("fatal error config probe: %w", err)
	}

	conf.OnConfigChange(func(in fsnotify.Event) {
		var setting Setting
		var p probe.Probe
		if err = conf.UnmarshalKey("setting", &setting); err != nil {
			log.Error("Fatal error config setting: %w", err)
			return
		}
		if err = conf.UnmarshalKey("probe", &p); err != nil {
			log.Error("Fatal error config probe: %w", err)
			return
		}
		fn(setting, p)
	})

	conf.WatchConfig()
	fn(setting, p)

	return nil
}

type Config struct {
	Setting Setting     `yaml:"setting"`
	Probe   probe.Probe `yaml:"probe"`
}

type Setting struct {
	LBSKey          string  `yaml:"lbsKey"`
	RoadOffset      int     `yaml:"roadOffset"`
	LineOffset      int     `yaml:"lineOffset"`
	AvoidAreaOffset float64 `yaml:"avoidAreaOffset"`
	MaxAvoid        int     `yaml:"maxAvoid"`
	MaxRoute        int     `yaml:"maxRoute"`
}
