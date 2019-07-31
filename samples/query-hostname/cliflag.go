package main

import (
	"errors"
	"flag"
	"time"

	mysqlroundrobinconnector "github.com/yinyin/go-mysql-round-robin-connector"
)

func parseCommandParam() (username, password, dbName string, timeoutDuration time.Duration, serverLocations []mysqlroundrobinconnector.Location, err error) {
	flag.StringVar(&username, "user", "", "user name of database instance")
	flag.StringVar(&password, "pass", "", "password of database instance")
	flag.StringVar(&dbName, "db", "mysql", "database name")
	flag.DurationVar(&timeoutDuration, "timeout", 10*time.Second, "timeout duration for dial")
	flag.Parse()
	serverAddresses := flag.Args()
	for _, addr := range serverAddresses {
		if "" == addr {
			continue
		}
		var loc mysqlroundrobinconnector.Location
		if '/' == addr[0] {
			loc.Network = "unix"
		} else {
			loc.Network = "tcp"
		}
		loc.Address = addr
		serverLocations = append(serverLocations, loc)
	}
	if len(serverLocations) == 0 {
		err = errors.New("server addresses is required")
	}
	return
}
