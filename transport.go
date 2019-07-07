package raft

import "net"

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
