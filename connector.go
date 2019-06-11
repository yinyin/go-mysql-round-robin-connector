package mysqlroundrobinconnector

import (
	"context"
	"net"
)

// RoundRobinDial implements `DialContextFunc` of `github.com/go-sql-driver/mysql` driver.
func RoundRobinDial(ctx context.Context, addr string) (net.Conn, error) {
	locSet := getLocationSet(addr)
	if nil == locSet {
		return nil, &UnknownLocationsErr{Name: addr}
	}
	// TODO: connect to locations
	return nil, nil
}
