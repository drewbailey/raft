package raft

// StableStore is used to privde stable storage
// of key configurations
type StableStore interface {
	Set(key []byte, val []byte) error
	Get(key []byte) ([]byte, error)

	SetUint64(key []byte, val uint64) error
	GetUint64(key []byte) (uint64, error)
}
