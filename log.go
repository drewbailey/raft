package raft

type LogType uint8

const (
	LogCommand LogType = iota
)

// Log entries are replicated to all members of the Raft cluster
// form the heart of the replicated state machine.
type Log struct {
	Index uint64
	Term  uint64
	Type  LogType
	Data  []byte
}

// LogStore is used to provide an interface for storing
// and retrieving logs
type LogStore interface {
	// LastIndex returns the last index written, default 0
	LastIndex() (uint64, error)

	// GetLog gets a log entry at a given index
	GetLog(index uint64, log *Log) error

	// StoreLog stores a log entry
	StoreLog(log *Log) error

	// DeleteRange deletes a range of log entries. Inclusive
	DeleteRange(min, max uint64) error
}
