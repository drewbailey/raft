package raft

import "time"

type Config struct {
	// Time in follower state without a leader before we attempt election
	HeartbeatTimeout time.Duration

	// Time without a leader before we should attempt an election
	ElectionTimeout time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		HeartbeatTimeout: 100 * time.Millisecond,
		ElectionTimeout:  150 * time.Millisecond,
	}
}
