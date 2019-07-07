package raft

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestLevelDB_StableStore(t *testing.T) {
	var l interface{} = &LevelDBStore{}
	_, ok := l.(StableStore)
	if !ok {
		t.Fatalf("LevelDBStore is not StableStore")
	}
}

func TestLevelDB_SetGet(t *testing.T) {
	// Create test dir
	dir, err := ioutil.TempDir("", "raft")
	if err != nil {
		t.Fatalf("err %v", err)
	}
	defer os.RemoveAll(dir)

	// New leveldb
	l, err := NewLevelDBStore(dir)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer l.Close()

	// Get a non existent key
	key := []byte("learningasigo")
	_, err = l.Get(key)
	if err.Error() != "not found" {
		t.Fatalf("err: %v", err)
	}

	val := []byte("test value")
	if err := l.Set(key, val); err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err := l.Get(key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if bytes.Compare(val, out) != 0 {
		t.Fatalf("Did not get resutl back: %v %v", val, out)
	}
}

func TestLevelDB_Logs(t *testing.T) {
	// Create a test dir
	dir, err := ioutil.TempDir("", "raft")
	if err != nil {
		t.Fatalf("err: %v ", err)
	}
	defer os.RemoveAll(dir)

	// New level
	l, err := NewLevelDBStore(dir)
	if err != nil {
		t.Fatalf("err: %v ", err)
	}
	defer l.Close()

	// Should be no last index
	idx, err := l.LastIndex()
	if err != nil {
		t.Fatalf("err: %v ", err)
	}
	if idx != 0 {
		t.Fatalf("bad idx: %d", idx)
	}

	// Try a filed fetch
	var out Log
	if err := l.GetLog(10, &out); err.Error() != "log not found" {
		t.Fatalf("err: %v ", err)
	}

	// Write out a log
	log := Log{
		Index: 10,
		Term:  3,
		Type:  LogCommand,
		Data:  []byte("test"),
	}
	if err := l.StoreLog(&log); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Try to fetch
	if err := l.GetLog(10, &out); err != nil {
		t.Fatalf("err: %v ", err)
	}

	// Check the highest index
	idx, err = l.LastIndex()
	if err != nil {
		t.Fatalf("err: %v ", err)
	}
	if idx != 10 {
		t.Fatalf("bad idx: %d", idx)
	}

	// Delete the entire range
	if err := l.DeleteRange(1, 10); err != nil {
		t.Fatalf("err: %v ", err)
	}

	// Index should be zero again
	idx, err = l.LastIndex()
	if err != nil {
		t.Fatalf("err: %v ", err)
	}
	if idx != 0 {
		t.Fatalf("bad idx: %d", idx)
	}

	// Should not be able to fetch
	if err := l.GetLog(10, &out); err.Error() != "log not found" {
		t.Fatalf("err: %v ", err)
	}
}
