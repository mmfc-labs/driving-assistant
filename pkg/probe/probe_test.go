package probe

import (
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
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

	yamlStr := `
  points:
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

	p2 := ProbeManager{}
	if err := yaml.Unmarshal([]byte(yamlStr), &p2); err != nil {
		t.Error(err)
	}
	fmt.Println(p2)

}
