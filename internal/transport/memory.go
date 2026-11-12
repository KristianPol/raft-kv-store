package transport

import (
	"context"
	"sync"
)

// MemoryNetwork is an in-process network for tests. Links can be cut to
// simulate crashes and partitions. ponytail: no delay/reorder yet, add in phase 6.
type MemoryNetwork struct {
	mu       sync.Mutex
	handlers map[string]Handler
	cut      map[[2]string]bool // directed: {from, to}
}

// NewMemoryNetwork returns an empty network.
func NewMemoryNetwork() *MemoryNetwork {
	return &MemoryNetwork{handlers: map[string]Handler{}, cut: map[[2]string]bool{}}
}

// Register attaches h as the node at addr (replacing any previous one).
func (n *MemoryNetwork) Register(addr string, h Handler) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.handlers[addr] = h
}

// Unregister removes the node at addr, as if it crashed.
func (n *MemoryNetwork) Unregister(addr string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	delete(n.handlers, addr)
}

// SetLink enables or disables traffic between a and b in both directions.
func (n *MemoryNetwork) SetLink(a, b string, up bool) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, k := range [][2]string{{a, b}, {b, a}} {
		if up {
			delete(n.cut, k)
		} else {
			n.cut[k] = true
		}
	}
}

// Isolate cuts addr off from every other registered node.
func (n *MemoryNetwork) Isolate(addr string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for other := range n.handlers {
		n.cut[[2]string{addr, other}] = true
		n.cut[[2]string{other, addr}] = true
	}
}

// Heal restores every link.
func (n *MemoryNetwork) Heal() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.cut = map[[2]string]bool{}
}

// Endpoint returns the Transport that node self uses to reach its peers.
func (n *MemoryNetwork) Endpoint(self string) Transport { return &memEndpoint{n, self} }

type memEndpoint struct {
	net  *MemoryNetwork
	self string
}

// route returns the handler for peer, or ErrUnreachable. The lock is not held on return.
func (e *memEndpoint) route(ctx context.Context, peer string) (Handler, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	e.net.mu.Lock()
	defer e.net.mu.Unlock()
	h, ok := e.net.handlers[peer]
	if !ok || e.net.cut[[2]string{e.self, peer}] {
		return nil, ErrUnreachable
	}
	return h, nil
}

func (e *memEndpoint) RequestVote(ctx context.Context, peer string, req *RequestVoteRequest) (*RequestVoteResponse, error) {
	h, err := e.route(ctx, peer)
	if err != nil {
		return nil, err
	}
	return h.HandleRequestVote(ctx, req)
}

func (e *memEndpoint) AppendEntries(ctx context.Context, peer string, req *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	h, err := e.route(ctx, peer)
	if err != nil {
		return nil, err
	}
	return h.HandleAppendEntries(ctx, req)
}

func (e *memEndpoint) InstallSnapshot(ctx context.Context, peer string, req *InstallSnapshotRequest) (*InstallSnapshotResponse, error) {
	h, err := e.route(ctx, peer)
	if err != nil {
		return nil, err
	}
	return h.HandleInstallSnapshot(ctx, req)
}
