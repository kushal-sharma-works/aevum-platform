package storage

// EventMapping defines the ES mapping for events
const EventMapping = `{
  "mappings": {
    "properties": {
      "event_id": {"type": "keyword"},
      "stream_id": {"type": "keyword"},
      "sequence_number": {"type": "long"},
      "event_type": {"type": "keyword"},
      "payload": {"type": "object", "enabled": true},
      "metadata": {"type": "object", "enabled": true},
      "occurred_at": {"type": "date"},
      "ingested_at": {"type": "date"},
      "schema_version": {"type": "keyword"}
    }
  }
}`

// DecisionMapping defines the ES mapping for decisions
const DecisionMapping = `{
  "mappings": {
    "properties": {
      "decision_id": {"type": "keyword"},
      "event_id": {"type": "keyword"},
      "stream_id": {"type": "keyword"},
      "rule_id": {"type": "keyword"},
      "rule_version": {"type": "keyword"},
      "status": {"type": "keyword"},
      "deterministic_hash": {"type": "keyword"},
      "input": {"type": "object", "enabled": true},
      "output": {"type": "object", "enabled": true},
      "trace": {
        "type": "nested",
        "properties": {
          "step": {"type": "integer"},
          "condition": {"type": "text"},
          "result": {"type": "boolean"},
          "message": {"type": "text"},
          "timestamp": {"type": "date"}
        }
      },
      "evaluated_at": {"type": "date"},
      "event_occurred_at": {"type": "date"}
    }
  }
}`

// SyncStateMapping defines the ES mapping for sync state
const SyncStateMapping = `{
  "mappings": {
    "properties": {
      "service_name": {"type": "keyword"},
      "last_synced_cursor": {"type": "keyword"},
      "last_sync_time": {"type": "date"},
      "sync_status": {"type": "keyword"}
    }
  }
}`
