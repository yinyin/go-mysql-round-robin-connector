package mysqlroundrobinconnector

import (
	"context"
	"math/rand"
	"net"
	"sync"
	"time"
)

// Location is the network and address of MySQL server instance.
type Location struct {
	Network       string
	Address       string
	TimeoutWeight int64
}

func (loc *Location) dialContext(ctx context.Context, timeout time.Duration, baseTime time.Time) (conn net.Conn, err error) {
	subCtx, ctxCancel := context.WithDeadline(ctx, baseTime.Add(timeout))
	defer ctxCancel()
	nd := net.Dialer{Timeout: timeout}
	return nd.DialContext(subCtx, loc.Network, loc.Address)
}

type locationSet struct {
	locations          []Location
	totalTimeoutWeight int64
}

func (s *locationSet) shuffledLocations(orderedCount int) (result []Location) {
	if orderedCount == -1 {
		return s.locations
	}
	totalCnt := len(s.locations)
	if orderedCount > totalCnt {
		orderedCount = totalCnt
	}
	result = append(make([]Location, 0, totalCnt), s.locations...)
	for idx := orderedCount; idx < (totalCnt - 1); idx++ {
		w := totalCnt - idx
		t := rand.Intn(w)
		if t == 0 {
			continue
		}
		result[idx], result[idx+t] = result[idx+t], result[idx]
	}
	return
}

var (
	locationSetsLock sync.RWMutex
	locationSets     map[string]*locationSet
)

// RegisterLocations add a set of MySQL instance locations with name to reference it.
func RegisterLocations(name string, locations []Location) error {
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
	var totalTimeoutWeight int64
	for _, loc := range locations {
		var timeoutWeight int64
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
	if 0 == totalTimeoutWeight {
		return &EmptyLocationsErr{
			Name: name,
		}
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
