package apis

import (
	"github.com/mmfc-labs/driving-assistant/pkg/probe"
)

type ProbeReq struct {
	Lat  float64 `form:"lat" json:"lat"`   // 当前位置经度
	Lon  float64 `form:"lon" json:"lon"`   // 当前位置纬度
	Near float64 `form:"near" json:"near"` // 获取附近多少公里的探头，0为获取所有探头
}

type ProbeResp struct {
	Probes []probe.Probe `json:"probes"`
}
