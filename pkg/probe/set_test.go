package probe

import (
	"github.com/mmfc-labs/driving-assistant/pkg/geo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSet(t *testing.T) {
	s := NewProbeSet()
	s.Add(geo.NewCoord(1, 1))
	s.Add(geo.NewCoord(1, 1))
	assert.Equal(t, 1, s.Len())
	s.Add(geo.NewCoord(1, 2))
	assert.Equal(t, 2, s.Len())
	assert.ElementsMatch(t, geo.Coords{geo.NewCoord(1, 1), geo.NewCoord(1, 2)}, s.ToSlice())
}
