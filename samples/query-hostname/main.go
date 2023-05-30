package main

import (
	"context"
	"database/sql"
	"log"
	"time"
)

func queryHostname(dbConn *sql.DB) (hostnameText string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err = dbConn.QueryRowContext(ctx, "SELECT @@hostname").Scan(&hostnameText); nil != err {
		log.Printf("ERROR: failed on query hostname: %v", err)
		return
	}
	return
}

func main() {
	username, password, dbName, extraAddrPath, timeoutDuration, serverLocations, connOpts, loopCount, err := parseCommandParam()
	if nil != err {
		log.Fatalf("failed on parsing command parameter: %v", err)
		return
	}
	log.Printf("username = %v, password = %v, db-name = %v, looping = %d, timeout = %v, %d server locations.",
		username, password, dbName, loopCount,
		timeoutDuration, len(serverLocations))
	dbConn, err := connectMySQL(username, password, dbName, extraAddrPath, timeoutDuration, serverLocations, connOpts)
	if nil != err {
		log.Fatalf("ERROR: failed on connecting to database instance: %v", err)
		return
	}
	for loopCount > 0 {
		loopCount--
		if hostnameText, err := queryHostname(dbConn); nil != err {
			log.Printf("query hostname failed (remain=%d): %v", loopCount, err)
		} else {
			log.Printf("Hostname (remain=%d): %v", loopCount, hostnameText)
		}
		if loopCount > 0 {
			time.Sleep(time.Second * 2)
		}
	}
}
