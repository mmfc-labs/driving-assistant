package main

import (
	"github.com/mmfc-labs/driving-assistant/pkg/config"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	"github.com/mmfc-labs/driving-assistant/pkg/probe"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoute(t *testing.T) {
	err := config.LoadConfig("./config.yaml", func(setting config.Setting, probeManager probe.ProbeManager) {
		probeManager.CalculateToward()
		lbsClient := lbs.NewLBS(setting, probeManager)
		tests := []struct {
			name string        // 测试用例名称
			from drive.Coord   // from 坐标
			to   drive.Coord   // to 坐标
			want []drive.Coord // 应该要屏蔽的摄像头
		}{
			{
				name: "case1",
				from: drive.Coord{Lat: 22.577781, Lon: 113.910683},
				to:   drive.Coord{Lat: 22.576752, Lon: 113.914866},
				want: []drive.Coord{{Lat: 22.576952, Lon: 113.914656}, {Lat: 22.576974, Lon: 113.914728}, {Lat: 22.57759, Lon: 113.914101}}},
			{
				name: "case2",
				from: drive.Coord{Lat: 22.577781, Lon: 113.910683},
				to:   drive.Coord{Lat: 22.576752, Lon: 113.914866},
				want: []drive.Coord{{Lat: 22.576952, Lon: 113.914656}, {Lat: 22.576974, Lon: 113.914728}, {Lat: 22.57759, Lon: 113.914101}}},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, avoidProbes, _, err := lbsClient.Route(tt.from, tt.to)

				if err != nil {
					assert.Error(t, err)
				}
				assert.Equal(t, len(tt.want), len(avoidProbes))
				assert.ElementsMatch(t, tt.want, avoidProbes)
			})
		}
	})
	if err != nil {
		return
	}
}