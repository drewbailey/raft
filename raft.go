package raft

type RaftState uint8

const (
	Follower RaftState = iota
	Candidate
	Leader
)

type Raft struct {
	// Configuration
	conf *Config

	// Current state
	state RaftState

	// stable is a StableStore implementation
	stable StableStore

	// logs is a LogStore implemeation to keep our logs
	logs LogStore

	// Highest committed log entry
	commitIndex uint64

	// last applied log to the FSM
	lastApplied uint64

	// if leader, extra state
	leader *LeaderState
}
