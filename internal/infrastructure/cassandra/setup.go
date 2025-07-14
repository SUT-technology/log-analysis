package cassandra

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

// CassandraClient wraps a gocql.Session.
type CassandraClient struct {
	session *gocql.Session
}

func NewCassandraClient(hosts []string, keyspace string) (*CassandraClient, error) {
	// مرحله ۱: ساخت session موقتی بدون keyspace
	tempCluster := gocql.NewCluster(hosts...)
	tempCluster.Consistency = gocql.Quorum
	tempCluster.ConnectTimeout = 10 * time.Second

	tempSession, err := tempCluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("create temp session (no keyspace): %w", err)
	}
	defer tempSession.Close()

	// ایجاد keyspace
	createKeyspace := fmt.Sprintf(`
		CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {
			'class': 'NetworkTopologyStrategy',
			'datacenter1': 3
		};`, keyspace)

	if err := tempSession.Query(createKeyspace).Exec(); err != nil {
		return nil, fmt.Errorf("create keyspace: %w", err)
	}

	// مرحله ۲: اتصال اصلی با keyspace
	mainCluster := gocql.NewCluster(hosts...)
	mainCluster.Keyspace = keyspace
	mainCluster.Consistency = gocql.Quorum
	mainCluster.ConnectTimeout = 10 * time.Second

	mainSession, err := mainCluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("connect with keyspace: %w", err)
	}

	client := &CassandraClient{session: mainSession}

	// مرحله ۳: ساخت جدول
	if err := client.initSchema(); err != nil {
		return nil, fmt.Errorf("init cassandra schema: %w", err)
	}

	return client, nil
}

func (c *CassandraClient) initSchema() error {
	createTable := `
	CREATE TABLE IF NOT EXISTS events_raw (
		id UUID,
		project_id UUID,
		event_name TEXT,
		event_time TIMESTAMP,
		inserted_time TIMESTAMP,
		payload map<TEXT,TEXT>,
		PRIMARY KEY ((project_id), event_name, event_time)
	) WITH default_time_to_live = 0;`

	return c.session.Query(createTable).Exec()
}
