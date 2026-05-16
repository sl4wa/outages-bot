package outage

// SnapshotStore stores normalized outages for deduplication.
type SnapshotStore interface {
	Load() ([]*Outage, error)
	Save(outages []*Outage) error
}
