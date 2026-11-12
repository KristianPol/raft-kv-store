// Package raft implements the Raft consensus algorithm. It only knows about
// opaque []byte commands, never about the key-value layer or the transport.
package raft

// Command is an opaque state-machine command.
type Command []byte

// ApplyMsg is delivered on the apply channel, in order, once per node lifetime.
// Either CommandValid or SnapshotValid is set.
type ApplyMsg struct {
	CommandValid bool
	Command      Command
	Index, Term  uint64

	SnapshotValid               bool
	Snapshot                    []byte
	SnapshotIndex, SnapshotTerm uint64
}

// Node is a single Raft peer as seen by the layer above it.
type Node interface {
	// Propose appends cmd to the log if this node is leader.
	Propose(cmd Command) (index, term uint64, isLeader bool)
	// State reports the current term, whether this node is leader, and the known leader.
	State() (term uint64, isLeader bool, leaderID string)
	// Snapshot is called by the state machine to compact the log up to index.
	Snapshot(index uint64, data []byte)
	// Stop shuts the node down.
	Stop()
}
