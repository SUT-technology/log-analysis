package db

import (
    "github.com/gocql/gocql"
    "log"
)

func InitCassandra() *gocql.Session {
    cluster := gocql.NewCluster("localhost")
    cluster.Keyspace = "logsystem"
    cluster.Consistency = gocql.Quorum

    session, err := cluster.CreateSession()
    if err != nil {
        log.Fatalf("Cassandra connection failed: %v", err)
    }

    // Create keyspace (if needed)
    if err := session.Query(`
        CREATE KEYSPACE IF NOT EXISTS logsystem
        WITH replication = {
            'class': 'SimpleStrategy',
            'replication_factor': 3
        }`).Exec(); err != nil {
        log.Fatalf("Failed to create keyspace: %v", err)
    }

    // Create table
    if err := session.Query(`
        CREATE TABLE IF NOT EXISTS events (
            project_id UUID,
            event_name TEXT,
            event_id UUID,
            timestamp TIMESTAMP,
            payload MAP<TEXT, TEXT>,
            created_at TIMESTAMP,
            PRIMARY KEY ((project_id, event_name), timestamp, event_id)
        ) WITH default_time_to_live = 0;
    `).Exec(); err != nil {
        log.Fatalf("Failed to create events table: %v", err)
    }

    return session
}
