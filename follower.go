package raft

func (r *Raft) followerAppendEntries(rpc RPC, a *AppendEntriesRequest) {}

func (r *Raft) followerRequestVote(rpc RPC, req *RequestVoteRequest) {}
