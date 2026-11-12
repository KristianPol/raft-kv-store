package storage

import "testing"

func TestMemory(t *testing.T) {
	var s Storage = NewMemory()
	s.SaveState(3, "n2")
	if term, v, _ := s.LoadState(); term != 3 || v != "n2" {
		t.Fatalf("state: %d %q", term, v)
	}
	s.AppendEntries([]LogEntry{{1, 1, []byte("a")}, {2, 1, []byte("b")}, {3, 2, []byte("c")}})
	if got, _ := s.Entries(2, 4); len(got) != 2 || string(got[0].Command) != "b" {
		t.Fatalf("entries: %v", got)
	}
	s.TruncateFrom(3)
	if got, _ := s.Entries(0, 10); len(got) != 2 {
		t.Fatalf("after truncate: %v", got)
	}
	s.SaveSnapshot(SnapshotMeta{1, 1}, []byte("snap"))
	if got, _ := s.Entries(0, 10); len(got) != 1 || got[0].Index != 2 {
		t.Fatalf("after snapshot: %v", got)
	}
	if m, d, _ := s.LoadSnapshot(); m.LastIncludedIndex != 1 || string(d) != "snap" {
		t.Fatalf("snapshot: %v %q", m, d)
	}
}
