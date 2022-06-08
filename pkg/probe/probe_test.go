package probe

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJson(t *testing.T) {
	jsonStr := `
{
	"lat":1,
	"lon":2
}
`
	p := Probe{}
	if err := json.Unmarshal([]byte(jsonStr), &p); err != nil {
		t.Error(err)
	}

	assert.Equal(t, float64(1), p.Lat)
	assert.Equal(t, float64(2), p.Lon)

	yamlStr := `
  probes:
    - lat: 22.57682
      lon: 113.913137
    - lat: 22.57759
      lon: 113.914101
      towards: # 探头朝向的坐标，有可能多个
        - lat: 22.577651
          lon: 113.914054
    - lat: 22.576349
      lon: 113.914133
    - lat: 22.576952
      lon: 113.914656
    - lat: 22.576974
      lon: 113.914728
`

	p2 := Manager{}
	if err := yaml.Unmarshal([]byte(yamlStr), &p2); err != nil {
		t.Error(err)
	}
	assert.Equal(t, 5, len(p2.Probes))
	assert.Equal(t, float64(22.57682), p2.Probes[0].Lat)
	assert.Equal(t, float64(113.913137), p2.Probes[0].Lon)
	assert.Equal(t, 1, len(p2.Probes[1].Towards))
	assert.Equal(t, float64(22.577651), p2.Probes[1].Towards[0].Lat)
	assert.Equal(t, float64(113.914054), p2.Probes[1].Towards[0].Lon)

}
