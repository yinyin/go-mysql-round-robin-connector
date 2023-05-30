package main

import (
	"errors"
	"flag"
	"time"

	mysqlroundrobinconnector "github.com/yinyin/go-mysql-round-robin-connector"
)

func parseCommandParam() (username, password, dbName, extraAddrPath string, timeoutDuration time.Duration, serverLocations []mysqlroundrobinconnector.Location, connOpts *mysqlroundrobinconnector.Options, loopCount int, err error) {
	var connPreferFirstLocation bool
	var connShuffleLocations bool
	flag.StringVar(&username, "user", "", "user name of database instance")
	flag.StringVar(&password, "pass", "", "password of database instance")
	flag.StringVar(&dbName, "db", "mysql", "database name")
	flag.StringVar(&extraAddrPath, "extraAddr", "", "additional address path of DSN (such as `-`)")
	flag.DurationVar(&timeoutDuration, "timeout", 10*time.Second, "timeout duration for dial")
	flag.BoolVar(&connPreferFirstLocation, "preferFirstLocation", false, "enable prefer first location option")
	flag.BoolVar(&connShuffleLocations, "shuffleLocations", false, "enable shuffle locations option")
	flag.IntVar(&loopCount, "loop", 1, "looping count")
	flag.Parse()
	serverAddresses := flag.Args()
	for _, addr := range serverAddresses {
		if addr == "" {
			continue
		}
		var loc mysqlroundrobinconnector.Location
		if addr[0] == '/' {
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
	if connPreferFirstLocation || connShuffleLocations {
		connOpts = &mysqlroundrobinconnector.Options{
			PreferFirstLocation: connPreferFirstLocation,
			ShuffleLocations:    connShuffleLocations,
		}
	}
	return
}
