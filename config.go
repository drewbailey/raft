package raft

import "time"

type Config struct {
	// Time in follower state without a leader before we attempt election
	HeartbeatTimeout time.Duration

	// Time without a leader before we should attempt an election
	ElectionTimeout time.Duration

	// Time without an Apply() operation before we heartbeat to ensure
	// a timely commit. Should be far less than HeartbeatTimeout to ensure
	// we don't lose leadership.
	CommitTimeout time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		HeartbeatTimeout: 100 * time.Millisecond,
		ElectionTimeout:  150 * time.Millisecond,
		CommitTimeout:    5 * time.Millisecond,
	}
}
