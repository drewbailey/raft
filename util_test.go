package raft

import (
	"regexp"
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	prev := generateUUID()
	for i := 0; i < 100; i++ {
		id := generateUUID()
		if prev == id {
			t.Fatalf("Should get a new ID!")
		}

		matched, err := regexp.MatchString(
			"[\\da-f]{8}-[\\da-f]{4}-[\\da-f]{4}-[\\da-f]{4}-[\\da-f]{12}", id)
		if !matched || err != nil {
			t.Fatalf("expected match %s %v %s", id, matched, err)
		}
	}
}

func TestAsyncNtoify(t *testing.T) {
	chs := []chan struct{}{
		make(chan struct{}),
		make(chan struct{}, 1),
		make(chan struct{}, 2),
	}

	asyncNotify(chs)
	asyncNotify(chs)
	asyncNotify(chs)

	// Try to read
	select {
	case <-chs[0]:
		t.Fatalf("should not have message!")
	default:
	}
	select {
	case <-chs[1]:
	default:
		t.Fatalf("should have message!")
	}
	select {
	case <-chs[2]:
	default:
		t.Fatalf("should have message!")
	}
	select {
	case <-chs[2]:
	default:
		t.Fatalf("should have message!")
	}
}
