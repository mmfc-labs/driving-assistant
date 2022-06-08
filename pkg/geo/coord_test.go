package geo

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestBearing(t *testing.T) {
	a1 := NewCoord(23.265333, 113.354559)
	a2 := NewCoord(23.265136, 113.354565)
	gap := math.Abs(a1.BearingTo(a2) - a2.BearingTo(a1))
	assert.Equal(t, true, gap > 179 && gap < 181)

}
