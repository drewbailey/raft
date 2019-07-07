package raft

// FSM provides an interface that can be implemented by
// clients to make use of the replicated log
type FSM interface {
	// ApplyLog is invoked once a log entry is committed
	ApplyLog([]byte)
}
