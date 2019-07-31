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
	baseTime, remainNanoseconds, err := getContextRemainNanoseconds(ctx)
	if nil != err {
		return nil, err
	}
	dailsErr := &DialsErr{}
	var targetTimeoutWeight int64
	for _, targetLoc := range locSet.locations {
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
