package raft

import "net"

// AppendEntriesRequest ...
type AppendEntriesRequest struct {
	// Provide the current term and leader
	Term     uint64
	LeaderID string

	// Proivde the previous entries for integrity checking
	PrevLogEntry uint64
	PrevLogTerm  uint64

	// New entries to commit
	Entries []*Log

	// Commit index on the leader
	LeaderCommitIndex uint64
}

// AppendEntriesResponse ...
type AppendEntriesResponse struct {
	Term uint64

	Success bool
}

// RequestVoteRequest ...
type RequestVoteRequest struct {
	// Provide the term and our id
	Term        uint64
	CandidateId string

	// Used to ensure safety
	LastLogIndex uint64
	LastLogTerm  uint64
}

// RequestVoteResponse ...
type RequestVoteResponse struct {
	Term    uint64
	Granted bool
}

// RPCResponse ...
type RPCResponse struct {
	Response interface{}
	Error    error
}

// RPC ...
type RPC struct {
	Command  interface{}
	RespChan chan<- RPCResponse
}

// Respond is used to respond with a response, error or both
func (r *RPC) Respond(resp interface{}, err error) {
	r.RespChan <- RPCResponse{resp, err}
}

// Transport interface
type Transport interface {
	// Consumer returns a channel that can be used to
	// consume and respond to RPC requests
	Consumer() <-chan RPC

	AppendEntries(target net.Addr, args *AppendEntriesRequest, resp *AppendEntriesResponse) error

	RequestVote(target net.Addr, args *RequestVoteRequest, resp *RequestVoteResponse) error
}
