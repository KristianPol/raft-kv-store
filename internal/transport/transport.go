// Package transport defines the node-to-node RPC interface and its messages.
package transport

import (
	"context"
	"errors"
)

// ErrUnreachable is returned when the peer cannot be contacted.
var ErrUnreachable = errors.New("transport: peer unreachable")

// Entry is a replicated log entry on the wire.
type Entry struct {
	Index, Term uint64
	Command     []byte
}

// RequestVoteRequest is the RequestVote RPC argument (Figure 2).
type RequestVoteRequest struct {
	Term         uint64
	CandidateID  string
	LastLogIndex uint64
	LastLogTerm  uint64
}

// RequestVoteResponse is the RequestVote RPC result.
type RequestVoteResponse struct {
	Term        uint64
	VoteGranted bool
}

// AppendEntriesRequest is the AppendEntries RPC argument (Figure 2).
type AppendEntriesRequest struct {
	Term         uint64
	LeaderID     string
	PrevLogIndex uint64
	PrevLogTerm  uint64
	Entries      []Entry
	LeaderCommit uint64
}

// AppendEntriesResponse is the AppendEntries RPC result, with the §5.3
// fast-backtracking hints.
type AppendEntriesResponse struct {
	Term          uint64
	Success       bool
	ConflictTerm  uint64
	ConflictIndex uint64
}

// InstallSnapshotRequest is the InstallSnapshot RPC argument.
type InstallSnapshotRequest struct {
	Term              uint64
	LeaderID          string
	LastIncludedIndex uint64
	LastIncludedTerm  uint64
	Data              []byte
}

// InstallSnapshotResponse is the InstallSnapshot RPC result.
type InstallSnapshotResponse struct {
	Term uint64
}

// Transport sends RPCs to peers.
type Transport interface {
	RequestVote(ctx context.Context, peer string, req *RequestVoteRequest) (*RequestVoteResponse, error)
	AppendEntries(ctx context.Context, peer string, req *AppendEntriesRequest) (*AppendEntriesResponse, error)
	InstallSnapshot(ctx context.Context, peer string, req *InstallSnapshotRequest) (*InstallSnapshotResponse, error)
}

// Handler receives RPCs on behalf of a node.
type Handler interface {
	HandleRequestVote(ctx context.Context, req *RequestVoteRequest) (*RequestVoteResponse, error)
	HandleAppendEntries(ctx context.Context, req *AppendEntriesRequest) (*AppendEntriesResponse, error)
	HandleInstallSnapshot(ctx context.Context, req *InstallSnapshotRequest) (*InstallSnapshotResponse, error)
}
