package raft

import (
	"fmt"
	"log"
	"sync"
)

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

	// FSM is the state machine that can handle the logs
	fsm FSM

	// The transport layer we use
	trans Transport

	// if leader, extra state
	leader *LeaderState

	// Shutdown channel to exit
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

func NewRaft(
	conf *Config,
	stable StableStore,
	logs LogStore,
	fsm FSM,
	trans Transport) (*Raft, error) {
	r := &Raft{
		conf:        conf,
		state:       Follower,
		stable:      stable,
		logs:        logs,
		commitIndex: 0,
		lastApplied: 0,
		fsm:         fsm,
		trans:       trans,
		shutdownCh:  make(chan struct{}),
	}

	go r.run()
	return r, nil
}

func (r *Raft) run() {
	for {
		// Check if we are doing a shutdown
		select {
		case <-r.shutdownCh:
			return
		default:
		}

		switch r.state {
		case Follower:
			r.runFollower()
		case Candidate:
			r.runCandidate()
		case Leader:
			r.runLeader()
		}
	}
}

func (r *Raft) runFollower() {
	ch := r.trans.Consumer()
	for {
		select {
		case rpc := <-ch:
			// Handle the command
			switch cmd := rpc.Command.(type) {
			case *AppendEntriesRequest:
				r.followerAppendEntries(rpc, cmd)
			case *RequestVoteRequest:
				r.followerRequestVote(rpc, cmd)
			default:
				log.Printf("[ERR] in Follower state, got unexpected command: %#v", rpc.Command)
				rpc.Respond(nil, fmt.Errorf("Unexpected command"))
			}
		case <-randomTimeout(r.conf.HeartbeatTimeout):
			// Hearbeat failed! Go to the candidate sttate
			r.state = Candidate
			return
		case <-r.shutdownCh:
			return
		}
	}
}

func (r *Raft) runCandidate() {
	for {
		select {
		case <-r.shutdownCh:
			return
		}
	}
}

func (r *Raft) runLeader() {
	for {
		select {
		case <-r.shutdownCh:
			return
		}
	}
}

func (r *Raft) Shutdown() {
	r.shutdownLock.Lock()
	defer r.shutdownLock.Lock()

	if r.shutdownCh != nil {
		close(r.shutdownCh)
		r.shutdownCh = nil
	}
}
