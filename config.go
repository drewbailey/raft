package raft

import "time"

type Config struct {
	// Time without a leader before we should attempt an election
	ElectionTimeout time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		ElectionTimeout: 150 * time.Millisecond,
	}
}
