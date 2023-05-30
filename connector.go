package mysqlroundrobinconnector

import (
	"context"
	"net"
	"time"
)

const defaultConnectionTimeout = 10 * time.Second

func getContextRemainNanoseconds(ctx context.Context) (baseTime time.Time, remainNanoseconds int64, err error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		baseTime = time.Now()
		return
	}
	baseTime = time.Now()
	remainNanoseconds = deadline.Sub(baseTime).Nanoseconds()
	if remainNanoseconds <= 0 {
		err = &TimeoutErr{
			ReferenceTime: baseTime,
			DeadlineTime:  deadline,
		}
	}
	return
}

func portionOfTimeoutDuration(remainNanoseconds, timeoutWeight, totalTimeoutWeight int64) (timeoutDuration time.Duration) {
	timeoutNanosec := (remainNanoseconds * timeoutWeight) / totalTimeoutWeight
	timeoutDuration = time.Nanosecond * time.Duration(timeoutNanosec)
	return
}

// RoundRobinDialContext implements `DialContextFunc` of `github.com/go-sql-driver/mysql` driver.
func RoundRobinDialContext(ctx context.Context, addr string) (net.Conn, error) {
	locSet := getLocationSet(addr)
	if nil == locSet {
		return nil, &UnknownLocationsErr{Name: addr}
	}
	locSet.lck.Lock()
	defer locSet.lck.Unlock()
	if locSet.shuffleOnNextDial {
		var skipCount int
		if locSet.preferFirstLocation {
			skipCount = 1
		}
		locSet.shuffleLocationsInPlace(skipCount)
		locSet.shuffleOnNextDial = false
	}
	// locations := locSet.shuffledLocations(parsedAddr.orderedCount)
	baseTime, remainNanoseconds, err := getContextRemainNanoseconds(ctx)
	if nil != err {
		return nil, err
	}
	dailsErr := &DialsErr{}
	var targetTimeoutWeight int64
	locSize := len(locSet.locations)
	for targetOffset := 0; targetOffset < locSize; targetOffset++ {
		targetIndex := (targetOffset + locSet.nextDialLocation) % locSize
		targetLoc := locSet.locations[targetIndex]
		targetTimeoutWeight += targetLoc.TimeoutWeight
		var timeoutDuration time.Duration
		if remainNanoseconds == 0 {
			timeoutDuration = defaultConnectionTimeout
		} else if timeoutDuration = portionOfTimeoutDuration(
			remainNanoseconds,
			targetTimeoutWeight, locSet.totalTimeoutWeight); timeoutDuration <= 0 {
			continue
		}
		if netConn, err := targetLoc.dialContext(ctx, timeoutDuration, baseTime); nil != err {
			dailsErr.append(err)
		} else {
			if (!locSet.preferFirstLocation) || (targetIndex != 0) {
				nextDialLoc := (targetIndex + 1) % locSize
				locSet.shuffleOnNextDial = locSet.shuffleLocations && (nextDialLoc <= locSet.nextDialLocation)
				locSet.nextDialLocation = nextDialLoc
			}
			return netConn, nil
		}
	}
	return nil, dailsErr
}

// RoundRobinDial implements `DialFunc` of `github.com/go-sql-driver/mysql` driver.
func RoundRobinDial(addr string) (net.Conn, error) {
	backgroundCtx := context.Background()
	timeoutCtx, cancel := context.WithTimeout(backgroundCtx, time.Second*10)
	defer cancel()
	return RoundRobinDialContext(timeoutCtx, addr)
}
