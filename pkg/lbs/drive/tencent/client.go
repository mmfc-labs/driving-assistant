package tencent

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
)

type Client struct {
	key string
}

func NewClient(key string) *Client {
	c := &Client{}
	c.key = key
	return c
}

func (c *Client) GetRoutes(from, to drive.Coord) ([]drive.Route, error) {
	client := resty.New()

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"from":     fmt.Sprintf("%f,%f", from.Lat, from.Lon),
			"to":       fmt.Sprintf("%f,%f", to.Lat, to.Lon),
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

	routes := make([]drive.Route, 0, len(p.Result.Routes))

	for _, route := range p.Result.Routes {
		for i := 2; i < len(route.Polyline); i++ {
			route.Polyline[i] = route.Polyline[i-2] + route.Polyline[i]/1000000
		}
		points := make([]drive.Coord, 0, len(route.Polyline)/2)
		for i := 0; i < len(route.Polyline); i = i + 2 {
			points = append(points, drive.Coord{
				Lat: route.Polyline[i],
				Lon: route.Polyline[i+1],
			})

		}
		routes = append(routes, drive.Route{Points: points})

	}
	return routes, nil
}
