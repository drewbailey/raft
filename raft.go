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

	logs LogStore

	commitIndex uint64

	lastApplied uint64

	leader *LeaderStore
}
