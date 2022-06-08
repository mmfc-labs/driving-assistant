package main

import (
	"github.com/mmfc-labs/driving-assistant/pkg/apis"
	"github.com/mmfc-labs/driving-assistant/pkg/config"
	"github.com/mmfc-labs/driving-assistant/pkg/geo"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs"
	"github.com/mmfc-labs/driving-assistant/pkg/probe"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoute(t *testing.T) {
	err := config.LoadConfig("./config.yaml", func(setting config.Setting, probeManager probe.Manager) {
		probeManager.CalculateToward()
		lbsClient := lbs.NewLBS(setting, probeManager)
		tests := []struct {
			name            string      // 测试用例名称
			from            geo.Coord   // from 坐标
			to              geo.Coord   // to 坐标
			wantAvoidProbes []geo.Coord // 应该要避让的摄像头（存在多次路线规划，很难在设计用例时知道哪些探头是需要避让的）
			wantError       error       // 路线规划错误， error == nil 成功规划 错误类型：https://github.com/mmfc-labs/driving-assistant/blob/main/pkg/apis/error.go
		}{
			{
				name:            "简单路线",
				from:            geo.NewCoord(22.577781, 113.910683),
				to:              geo.NewCoord(22.576752, 113.914866),
				wantAvoidProbes: []geo.Coord{geo.NewCoord(22.577590, 113.914101)},
				wantError:       nil,
			},
			{
				name:            "路线必须经过避让区，导致第三方地图无法规划路线",
				from:            geo.NewCoord(22.578005, 113.913589),
				to:              geo.NewCoord(22.577347, 113.914361),
				wantAvoidProbes: nil,
				wantError:       apis.ErrorRouteFailed,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, avoidProbes, _, err := lbsClient.Route(tt.from, tt.to)

				// err == nil 说明规划路线成功
				assert.Equal(t, tt.wantError, err)

				// 判断需要避让的摄像头是否匹配
				assert.ElementsMatch(t, tt.wantAvoidProbes, avoidProbes)
			})
		}
	})
	if err != nil {
		return
	}
}
