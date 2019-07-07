package raft

// StableStore is used to privde stable storage
// of key configurations
type StableStore interface {
	// Returns current term
	CurrentTerm() (uint64, error)

	// Returns candidate voted for this term
	VotedFor() (string, error)

	// Sets the current term. Clears the current vote
	SetCurrentTerm(uint64) error

	// Sets a candidate vote for the current term
	SetVote(string) error

	// Returns our candidate ID.
	CandidateID() (string, error)
}
