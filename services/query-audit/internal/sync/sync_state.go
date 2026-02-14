package sync

import "time"

// SyncState tracks sync progress
type SyncState struct {
	ServiceName      string    `json:"service_name"`
	LastSyncedCursor string    `json:"last_synced_cursor"`
	LastSyncTime     time.Time `json:"last_sync_time"`
	SyncStatus       string    `json:"sync_status"`
}

// NewSyncState creates a new sync state
func NewSyncState(serviceName string) *SyncState {
	return &SyncState{
		ServiceName:      serviceName,
		LastSyncedCursor: "",
		LastSyncTime:     time.Now(),
		SyncStatus:       "initialized",
	}
}

// UpdateCursor updates the sync cursor
func (ss *SyncState) UpdateCursor(cursor string) {
	ss.LastSyncedCursor = cursor
	ss.LastSyncTime = time.Now()
	ss.SyncStatus = "synced"
}

// MarkFailed marks the sync as failed
func (ss *SyncState) MarkFailed() {
	ss.SyncStatus = "failed"
}
