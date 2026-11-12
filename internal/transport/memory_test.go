package transport

import (
	"context"
	"errors"
	"testing"
)

type echo struct{ term uint64 }

func (e echo) HandleRequestVote(_ context.Context, r *RequestVoteRequest) (*RequestVoteResponse, error) {
	return &RequestVoteResponse{Term: e.term, VoteGranted: true}, nil
}
func (e echo) HandleAppendEntries(_ context.Context, r *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	return &AppendEntriesResponse{Term: e.term, Success: true}, nil
}
func (e echo) HandleInstallSnapshot(_ context.Context, r *InstallSnapshotRequest) (*InstallSnapshotResponse, error) {
	return &InstallSnapshotResponse{Term: e.term}, nil
}

func TestMemoryNetwork(t *testing.T) {
	n := NewMemoryNetwork()
	n.Register("a", echo{1})
	n.Register("b", echo{2})
	a, ctx := n.Endpoint("a"), context.Background()

	if r, err := a.RequestVote(ctx, "b", &RequestVoteRequest{}); err != nil || r.Term != 2 {
		t.Fatalf("got %v, %v", r, err)
	}
	n.Isolate("b")
	if _, err := a.AppendEntries(ctx, "b", &AppendEntriesRequest{}); !errors.Is(err, ErrUnreachable) {
		t.Fatalf("want ErrUnreachable, got %v", err)
	}
	n.Heal()
	if _, err := a.InstallSnapshot(ctx, "b", &InstallSnapshotRequest{}); err != nil {
		t.Fatal(err)
	}
	n.Unregister("b")
	if _, err := a.RequestVote(ctx, "b", &RequestVoteRequest{}); !errors.Is(err, ErrUnreachable) {
		t.Fatalf("want ErrUnreachable after unregister, got %v", err)
	}
}
