package tencent

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mmfc-labs/driving-assistant/lbs"
)

type LBS struct {
	key string
}

func NewLBSClient(key string) *LBS {
	c := &LBS{}
	c.key = key
	return c
}

func (c *LBS) GetRoute(from1, from2, to1, to2 float64) ([]lbs.Route, error) {
	client := resty.New()

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"from":     fmt.Sprintf("%f,%f", from1, from2),
			"to":       fmt.Sprintf("%f,%f", to1, to2),
			"output":   "json",
			"callback": "cb",
			"key":      c.key,
		}).
		Get("https://apis.map.qq.com/ws/direction/v1/driving/")

	if err != nil {
		return nil, err
	}

	p := struct {
		Status int
		Result struct {
			Routes []struct {
				Polyline []float64
			}
		}
	}{}

	err = json.Unmarshal(resp.Body(), &p)
	if err != nil {
		return nil, err
	}

	routes := make([]lbs.Route, 0, len(p.Result.Routes))

	for _, route := range p.Result.Routes {
		for i := 2; i < len(route.Polyline); i++ {
			route.Polyline[i] = route.Polyline[i-2] + route.Polyline[i]/1000000
		}
		points := make([]lbs.Coord, 0, len(route.Polyline)/2)
		for i := 0; i < len(route.Polyline); i = i + 2 {
			points = append(points, lbs.Coord{
				Lat: route.Polyline[i],
				Lon: route.Polyline[i+1],
			})

		}
		routes = append(routes, lbs.Route{Points: points})

	}
	return routes, nil
}
