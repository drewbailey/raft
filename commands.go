package raft

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
