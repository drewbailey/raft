package raft

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/ugorji/go/codec"
)

const logPath = "logs"
const confPath = "conf"
const maxOpenFiles = 128

// LevelDBStore provides an implementation of LogStore
// as well as StableStore, allowing it to be used as
// the primary backing store
type LevelDBStore struct {
	logs *leveldb.DB
	conf *leveldb.DB
}

func NewLevelDBStore(base string) (*LevelDBStore, error) {
	// Get the paths

	logLoc := filepath.Join(base, logPath)
	confLoc := filepath.Join(base, confPath)

	// create the struct
	ldb := &LevelDBStore{}

	// LevelDB opts
	opts := &opt.Options{
		Compression:            opt.SnappyCompression,
		OpenFilesCacheCapacity: maxOpenFiles,
		Strict:                 opt.StrictAll,
	}

	// Open the DBs
	db, err := leveldb.OpenFile(logLoc, opts)
	if err != nil {
		log.Printf("[ERR] Failed to open logs leveldb: %v", err)
		return nil, err
	}
	ldb.logs = db

	db, err = leveldb.OpenFile(logLoc, opts)
	if err != nil {
		log.Printf("[ERR] Failed to open conf leveldb: %v", err)
		return nil, err
	}

	ldb.conf = db

	return ldb, nil
}

func (l *LevelDBStore) Close() error {
	err1 := l.logs.Close()
	err2 := l.conf.Close()
	if err1 == nil && err2 == nil {
		return nil
	} else if err1 == nil && err2 != nil {
		return err2
	} else if err1 != nil && err2 == nil {
		return err1
	} else {
		return fmt.Errorf("Failed to close DB: Got: %v and %v",
			err1, err2)
	}
}

func (l *LevelDBStore) LastIndex() (uint64, error) {
	// Get an iterator
	it := l.logs.NewIterator(nil, nil)
	defer it.Release()

	// Seek to the last value
	it.Last()

	// Check if there is a key
	key := it.Key()
	if key == nil {
		// Nothing written yet
		return 0, it.Error()
	}

	// Convert the key to the index
	return bytesToUint64(key), it.Error()
}

// Gets a log entry at a given index
func (l *LevelDBStore) GetLog(index uint64, logOut *Log) error {
	key := uint64ToBytes(index)

	// Get an iterator
	snap, err := l.logs.GetSnapshot()
	if err != nil {
		return err
	}
	defer snap.Release()

	// Look for the key
	val, err := snap.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return fmt.Errorf("log not found")
	} else if err != nil {
		return err
	}

	// Convert the value to a log
	return decode(val, logOut)
}

func (l *LevelDBStore) StoreLog(log *Log) error {
	// Concert to an on-disk format
	key := uint64ToBytes(log.Index)
	val, err := encode(log)
	if err != nil {
		return err
	}

	// write
	opts := &opt.WriteOptions{Sync: true}
	return l.logs.Put(key, val.Bytes(), opts)
}

func (l *LevelDBStore) DeleteRange(min, max uint64) error {
	// Create a batch operation
	batch := &leveldb.Batch{}
	for i := min; i <= max; i++ {
		key := uint64ToBytes(i)
		batch.Delete(key)
	}

	// Apply
	opts := &opt.WriteOptions{Sync: true}
	return l.logs.Write(batch, opts)
}

func (l *LevelDBStore) Set(key, val []byte) error {
	opts := &opt.WriteOptions{Sync: true}
	return l.conf.Put(key, val, opts)
}

func (l *LevelDBStore) Get(key []byte) ([]byte, error) {
	// Get snapshot view
	snap, err := l.conf.GetSnapshot()
	if err != nil {
		return nil, err
	}
	defer snap.Release()

	// Find key
	val, err := snap.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return nil, fmt.Errorf("not found")
	} else if err != nil {
		return nil, err
	}

	// Copy to new buffer
	buf := make([]byte, len(val))
	copy(buf, val)
	return buf, nil
}

func (l *LevelDBStore) SetUint64(key []byte, val uint64) error {
	return l.Set(key, uint64ToBytes(val))
}

func (l *LevelDBStore) GetUint64(key []byte) (uint64, error) {
	buf, err := l.Get(key)
	if err != nil {
		return 0, err
	}
	return bytesToUint64(buf), nil
}

// Convert uint to a byte slice
func uint64ToBytes(u uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, u)
	return buf
}

func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// Decode reverses the encode operation on a byte slice input
func decode(buf []byte, out interface{}) error {
	r := bytes.NewBuffer(buf)
	hd := codec.MsgpackHandle{}
	dec := codec.NewDecoder(r, &hd)
	return dec.Decode(out)
}

// Encode writes an encoded object to a new bytes buffer
func encode(in interface{}) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	hd := codec.MsgpackHandle{}
	enc := codec.NewEncoder(buf, &hd)
	err := enc.Encode(in)
	return buf, err
}
