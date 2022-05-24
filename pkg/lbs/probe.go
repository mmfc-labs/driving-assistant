package lbs

import (
	"encoding/json"
	"github.com/mmfc-labs/driving-assistant/pkg/lbs/drive"
	"io/ioutil"
)

func LoadProbe() *Probe {
	probes := &Probe{}

	bytes, err := ioutil.ReadFile("./probe.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, probes)
	if err != nil {
		panic(err)
	}

	return probes
}

type Probe struct {
	Points []drive.Coord
}
