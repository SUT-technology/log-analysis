package main

import (
	// "backend/kafka"
	"fmt"
	"os"

	"github.com/SUT-technology/log-analysis/db"
	// etc.
)

func main() {
    cockroach := db.InitCockroachDB(os.Getenv("COCKROACH_DSN"))
    cassandra := db.InitCassandra()
    clickhouse := db.InitClickHouse()
	
	fmt.Printf("%v\n%v\n%v",cockroach,cassandra,clickhouse)
    // kafka.EnsureLogTopic()

    // Start your server or background consumers...
}