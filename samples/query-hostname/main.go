package main

import (
	"context"
	"time"
	"log"

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
	username, password, dbName, timeoutDuration, serverLocations, err := parseCommandParam()
	if nil != err {
		log.Fatalf("failed on parsing command parameter: %v", err)
		return
	}
	log.Printf("username = %v, password = %v, db-name = %v, timeout = %v, %d server locations.",
		username, password, dbName, timeoutDuration, len(serverLocations))
	hostnameText, err := queryHostname(username, password, dbName, timeoutDuration, serverLocations)
	if nil != err {
		log.Fatalf("failed on query hostname: %v", err)
		return
	}
	log.Printf("Hostname: %v", hostnameText)
}
