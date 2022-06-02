package apis

import "errors"

var (
	// RouteOutOfRange 超出路线规划次数
	RouteOutOfRange = errors.New("ROUTE_OUT_OF_RANGE")
	// AvoidOutOfRange 超出避让探头数量
	AvoidOutOfRange = errors.New("AVOID_OUT_OF_RANGE")
)
