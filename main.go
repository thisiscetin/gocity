package gocity

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// City ..
type City struct {
	Name  string
	I, J  float64
	Roads map[*City]float64
}

// NewCity ..
func NewCity(name string, i, j float64) *City {
	return &City{
		name,
		i,
		j,
		make(map[*City]float64, 0),
	}
}

// Map ..
type Map struct {
	Cities       []*City
	RoadsBuilt   bool
	AllReachable bool
	sync.Mutex
}

// NewMap ..
func NewMap() *Map {
	return &Map{
		Cities: make([]*City, 0, 1),
	}
}

// DistanceTo ..
func (c *City) DistanceTo(c1 *City) float64 {
	return math.Sqrt(math.Pow((c.I-c1.I), 2) + math.Pow((c.J-c1.J), 2))
}

// AddCity ..
func (cm *Map) AddCity(name string, i, j float64) error {
	if i < 0 || j < 0 {
		return errors.New("i, j coordinates should be positive")
	}
	cm.Lock()
	defer cm.Unlock()

	for _, c := range cm.Cities {
		if c.I == i && c.J == j {
			return fmt.Errorf("there is a city named %s on %f, %f", c.Name,
				i, j)
		}
	}

	cm.Cities = append(cm.Cities, NewCity(name, i, j))
	return nil
}

// FindCity ..
func (cm *Map) FindCity(name string) *City {
	for _, c := range cm.Cities {
		if c.Name == name {
			return c
		}
	}
	return nil
}

type roads struct {
	cap       int
	cities    []*City
	distances []float64
}

func newRoads(cap int) *roads {
	return &roads{
		cap,
		make([]*City, 0, cap),
		make([]float64, 0, cap),
	}
}

func (r *roads) addCity(c *City, d float64) {
	for i, dis := range r.distances {
		if dis > d && len(r.cities) == r.cap {
			r.cities = append(r.cities[:i], r.cities[i+1:]...)
			r.distances = append(r.distances[:i], r.distances[i+1:]...)

			break
		}
	}

	if len(r.cities) < r.cap {
		r.cities = append(r.cities, c)
		r.distances = append(r.distances, d)
	}
}

// BuildRoads .. number of cities to build roads
func (cm *Map) BuildRoads(closest int) {
	cm.Lock()
	defer cm.Unlock()

	for i, c0 := range cm.Cities {
		tempRoads := newRoads(closest)

		for j, c1 := range cm.Cities {
			if i == j {
				continue
			}
			tempRoads.addCity(c1, c0.DistanceTo(c1))
		}

		// make sure we are adding clearing roads before adding new
		c0.Roads = make(map[*City]float64, 0)

		for i := 0; i < tempRoads.cap; i++ {
			c0.Roads[tempRoads.cities[i]] = tempRoads.distances[i]
		}
	}

	// build vertices in reverse order to simulate roads flowing on both ways
	// record connected cities in the mean time to check all cities are reachable
	connected := make(map[*City]bool, 0)
	for _, c := range cm.Cities {
		for k, v := range c.Roads {
			k.Roads[c] = v
			connected[k] = true
		}
	}

	cm.RoadsBuilt = true
	cm.AllReachable = len(cm.Cities) == len(connected)
}

// Path ..
type Path struct {
	Starting    *City
	Route       []*City
	Distance    float64
	TimeElapsed time.Duration
	TimeFound   time.Time
}

func (p *Path) String() string {
	return fmt.Sprintf("Started and ended at %s, made %d visits, distance %f, took %v to discover",
		p.Starting.Name, len(p.Route), p.Distance, p.TimeElapsed)
}

// FindRandomPath ..
func (cm *Map) FindRandomPath(starting *City, distanceCap float64) (*Path, error) {
	if !cm.RoadsBuilt {
		return nil, fmt.Errorf("first build roads")
	}
	if !cm.AllReachable {
		return nil, fmt.Errorf("all cities are not reachable, build more roads")
	}

	timeStart, path := time.Now(), &Path{
		Starting: starting,
		Route:    []*City{starting},
	}

	cityCount := len(cm.Cities)
	seen := make(map[*City]bool, cityCount)

	rand.Seed(timeStart.UnixNano())
	for {
		i := rand.Intn(len(starting.Roads))

		for k, v := range starting.Roads {
			if i == 0 {
				starting = k
				seen[k] = true
				path.Distance += v
				path.Route = append(path.Route, k)

				break
			}
			i--
		}

		if path.Distance > distanceCap && distanceCap > 0 {
			return nil, nil
		}
		if len(seen) == cityCount {
			path.TimeFound = time.Now()
			path.TimeElapsed = time.Since(timeStart)

			return path, nil
		}
	}
}
