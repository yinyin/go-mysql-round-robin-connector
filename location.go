package mysqlroundrobinconnector

import (
	"sync"
)

// Location is the network and address of MySQL server instance.
type Location struct {
	Network       string
	Address       string
	TimeoutWeight int
}

type locationSet struct {
	locations          []Location
	lastConnectedIndex int
	totalTimeoutWeight int
}

var (
	locationSetsLock sync.RWMutex
	locationSets     map[string]*locationSet
)

// RegisterLocations add a set of MySQL instance locations with name to reference it.
func RegisterLocations(name string, locations []*Location) error {
	if 0 == len(locations) {
		return &EmptyLocationsErr{
			Name: name,
		}
	}
	locationSetsLock.Lock()
	defer locationSetsLock.Unlock()
	if nil == locationSets {
		locationSets = make(map[string]*locationSet)
	}
	aux := &locationSet{
		locations: make([]Location, 0, len(locations)),
	}
	totalTimeoutWeight := 0
	for _, loc := range locations {
		var timeoutWeight int
		if timeoutWeight = loc.TimeoutWeight; timeoutWeight < 1 {
			timeoutWeight = 1
		}
		nLoc := Location{
			Network:       loc.Network,
			Address:       loc.Address,
			TimeoutWeight: timeoutWeight,
		}
		aux.locations = append(aux.locations, nLoc)
		totalTimeoutWeight = totalTimeoutWeight + timeoutWeight
	}
	aux.totalTimeoutWeight = totalTimeoutWeight
	locationSets[name] = aux
	return nil
}

// UnregisterLocations remove a set of MySQL instance locations with previous register name.
func UnregisterLocations(name string) {
	locationSetsLock.Lock()
	defer locationSetsLock.Unlock()
	if nil == locationSets {
		return
	}
	delete(locationSets, name)
}

func getLocationSet(name string) (locSet *locationSet) {
	locationSetsLock.RLock()
	defer locationSetsLock.RUnlock()
	locSet = locationSets[name]
	return
}
