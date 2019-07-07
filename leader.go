package raft

type LeaderState struct {
	followers map[string]*FollowerState
}

type FollowerState struct {
	// next index to send
	nextIndex uint64

	// last known replicated index
	replicatedIndex uint64
}
