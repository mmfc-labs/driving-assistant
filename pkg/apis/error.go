package apis

import "errors"

var (
	// ErrorRouteOutOfRange 超出路线规划次数
	ErrorRouteOutOfRange = errors.New("ROUTE_OUT_OF_RANGE")
	// ErrorAvoidOutOfRange 超出避让探头数量
	ErrorAvoidOutOfRange = errors.New("AVOID_OUT_OF_RANGE")
	// ErrorRouteFailed 第三方地图规划路线失败，一般是因为路线必须经过避让区
	ErrorRouteFailed = errors.New("ROUTE_FAILED")
)
