CREATE KEYSPACE IF NOT EXISTS log_analysis
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

CREATE TABLE IF NOT EXISTS log_analysis.event_raw (
    project_id UUID,
    event_name TEXT,
    event_time TIMESTAMP,
    inserted_time TIMESTAMP,
    payload MAP<TEXT, TEXT>,
    PRIMARY KEY (project_id, event_time)
);