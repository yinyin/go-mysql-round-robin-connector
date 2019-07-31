package main

import (
	"context"
	"log"
	"time"

	mysqlroundrobinconnector "github.com/yinyin/go-mysql-round-robin-connector"
)

func queryHostname(username, password, dbName string, timeoutDuration time.Duration, serverLocations []mysqlroundrobinconnector.Location) (hostnameText string, err error) {
	dbConn, err := connectMySQL(username, password, dbName, timeoutDuration, serverLocations)
	if nil != err {
		log.Printf("ERROR: failed on connecting to database instance: %v", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err = dbConn.QueryRowContext(ctx, "SELECT @@hostname").Scan(&hostnameText); nil != err {
		log.Printf("ERROR: failed on query hostname: %v", err)
		return
	}
	return
}

func main() {
	username, password, dbName, timeoutDuration, serverLocations, loopCount, err := parseCommandParam()
	if nil != err {
		log.Fatalf("failed on parsing command parameter: %v", err)
		return
	}
	log.Printf("username = %v, password = %v, db-name = %v, looping = %d, timeout = %v, %d server locations.",
		username, password, dbName, loopCount,
		timeoutDuration, len(serverLocations))
	for loopCount > 0 {
		loopCount--
		if hostnameText, err := queryHostname(username, password, dbName, timeoutDuration, serverLocations); nil != err {
			log.Printf("query hostname failed (remain=%d): %v", loopCount, err)
		} else {
			log.Printf("Hostname (remain=%d): %v", loopCount, hostnameText)
		}
		if loopCount > 0 {
			time.Sleep(time.Second * 2)
		}
	}
}
