package main

import (
	"database/sql"
	"log"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"

	mysqlroundrobinconnector "github.com/yinyin/go-mysql-round-robin-connector"
)

const defaultLocationName = "location-1"
const dialerName = "rr-conn"

func connectMySQL(username, password, dbName, extraAddrPath string, timeoutDuration time.Duration, serverLocations []mysqlroundrobinconnector.Location) (conn *sql.DB, err error) {
	mysqldriver.RegisterDialContext(dialerName, mysqlroundrobinconnector.RoundRobinDialContext) // TODO: go-sql-driver/mysql 1.4.3+
	// mysqldriver.RegisterDial(dialerName, mysqlroundrobinconnector.RoundRobinDial)
	if err = mysqlroundrobinconnector.RegisterLocations(defaultLocationName, serverLocations); nil != err {
		return
	}
	cfg := mysqldriver.NewConfig()
	cfg.User = username
	cfg.Passwd = password
	cfg.Net = dialerName
	if extraAddrPath == "" {
		cfg.Addr = defaultLocationName
	} else {
		cfg.Addr = defaultLocationName + "/" + extraAddrPath
	}
	cfg.DBName = dbName
	cfg.Timeout = timeoutDuration
	cfg.ReadTimeout = timeoutDuration
	cfg.WriteTimeout = timeoutDuration
	// cfg.ParseTime = true
	dsn := cfg.FormatDSN()
	log.Printf("INFO: connecting with DSN: %s", dsn)
	return sql.Open("mysql", dsn)
}
