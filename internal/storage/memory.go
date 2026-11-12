package storage

import "sync"

// Memory is an in-memory Storage for tests. Its contents survive a simulated
// restart as long as the same value is handed to the new node.
type Memory struct {
	mu       sync.Mutex
	term     uint64
	votedFor string
	log      []LogEntry
	meta     SnapshotMeta
	snap     []byte
}

// NewMemory returns an empty Memory store.
func NewMemory() *Memory { return &Memory{} }

// SaveState stores the current term and vote.
func (m *Memory) SaveState(term uint64, votedFor string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.term, m.votedFor = term, votedFor
	return nil
}

// LoadState returns the stored term and vote.
func (m *Memory) LoadState() (uint64, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.term, m.votedFor, nil
}

// AppendEntries appends copies of entries to the log.
func (m *Memory) AppendEntries(entries []LogEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range entries {
		e.Command = append([]byte(nil), e.Command...)
		m.log = append(m.log, e)
	}
	return nil
}

// TruncateFrom deletes every entry with Index >= index.
func (m *Memory) TruncateFrom(index uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, e := range m.log {
		if e.Index >= index {
			m.log = m.log[:i]
			break
		}
	}
	return nil
}

// Entries returns copies of entries with lo <= Index < hi.
func (m *Memory) Entries(lo, hi uint64) ([]LogEntry, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var out []LogEntry
	for _, e := range m.log {
		if e.Index >= lo && e.Index < hi {
			e.Command = append([]byte(nil), e.Command...)
			out = append(out, e)
		}
	}
	return out, nil
}

// SaveSnapshot stores the snapshot and drops log entries it covers.
func (m *Memory) SaveSnapshot(meta SnapshotMeta, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.meta, m.snap = meta, append([]byte(nil), data...)
	keep := m.log[:0:0]
	for _, e := range m.log {
		if e.Index > meta.LastIncludedIndex {
			keep = append(keep, e)
		}
	}
	m.log = keep
	return nil
}

// LoadSnapshot returns the stored snapshot, or a zero meta and nil data if none.
func (m *Memory) LoadSnapshot() (SnapshotMeta, []byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.meta, append([]byte(nil), m.snap...), nil
}
