// Package storage defines Raft's durable state interface.
package storage

// LogEntry is one persisted log entry.
type LogEntry struct {
	Index, Term uint64
	Command     []byte
}

// SnapshotMeta describes the last log entry covered by a snapshot.
type SnapshotMeta struct {
	LastIncludedIndex, LastIncludedTerm uint64
}

// Storage persists Raft state. Implementations must make each call durable
// before returning.
type Storage interface {
	SaveState(term uint64, votedFor string) error
	LoadState() (term uint64, votedFor string, err error)
	// AppendEntries appends entries after the current last entry.
	AppendEntries(entries []LogEntry) error
	// TruncateFrom deletes every entry with Index >= index.
	TruncateFrom(index uint64) error
	// Entries returns entries with lo <= Index < hi.
	Entries(lo, hi uint64) ([]LogEntry, error)
	SaveSnapshot(meta SnapshotMeta, data []byte) error
	LoadSnapshot() (SnapshotMeta, []byte, error)
}
