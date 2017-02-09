package gocity

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistanceTo(t *testing.T) {
	c0, c1 := &City{
		I: 0,
		J: 0,
	}, &City{
		I: 3,
		J: 4,
	}

	assert.Equal(t, float64(5), c0.DistanceTo(c1))
	assert.Equal(t, float64(5), c1.DistanceTo(c0))
}

func TestBuildRoads(t *testing.T) {
	cm := NewMap()

	cm.AddCity("a", 0.0, 1.0)
	cm.AddCity("b", 0.0, 2.0)

	cm.BuildRoads(1)
	// expected roads to be built for the closest distance 1
	// a -> b
	// b -> a
	// c -> a || c-> b
	cityA := cm.FindCity("a")
	cityB := cm.FindCity("b")

	assert.Equal(t, 1, len(cityA.Roads))
	assert.Equal(t, float64(1), cityA.Roads[cityB])
	assert.Equal(t, float64(1), cityA.Roads[cityB])

	cm.AddCity("c", 1.0, 0.0)
	cm.BuildRoads(1)

	cityC := cm.FindCity("c")

	assert.Equal(t, 1, len(cityC.Roads))
	assert.Equal(t, math.Sqrt(2.0), cityC.Roads[cityA])

	// we are adding road to cityA which is coming from cityC
	assert.Equal(t, 1, len(cityC.Roads))
	assert.Equal(t, math.Sqrt(2.0), cityC.Roads[cityA])

	cm.BuildRoads(2)

	assert.Equal(t, 2, len(cityB.Roads))
	assert.Equal(t, 2, len(cityA.Roads))
	assert.Equal(t, 2, len(cityC.Roads))
}

func TestFindRandomPath(t *testing.T) {
	cm := NewMap()

	cm.AddCity("a", 0.0, 1.0)
	cm.AddCity("b", 0.0, 2.0)
	cm.AddCity("c", 1.0, 2.0)
	cm.AddCity("d", 2.0, 2.0)

	cm.BuildRoads(2)

	assert.Equal(t, true, cm.AllReachable)

	cityA := cm.FindCity("a")

	p, err := cm.FindRandomPath(cityA, 0)

	assert.NoError(t, err)
	assert.Equal(t, cityA, p.Starting)
	assert.NotNil(t, p.TimeElapsed)
	assert.NotNil(t, p.TimeFound)
	assert.True(t, len(p.Route) >= 4)
}
