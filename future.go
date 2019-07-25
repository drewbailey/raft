package raft

// ApplyFuture is used to represent an application that may occur in the future
type ApplyFuture interface {
	Error() error
}

// errorFuture is used to return a static error
type errorFuture struct {
	err error
}

func (e errorFuture) Error() error {
	return e.err
}

// logFuture is used to apply a log entry and waits until
// the log is considered committed
type logFuture struct {
	log   Log
	err   error
	errCh chan error
}

func (l *logFuture) Error() error {
	if l.err != nil {
		return l.err
	}
	l.err = <-l.errCh

	return l.err
}

func (l *logFuture) respond(err error) {
	l.errCh <- err
	close(l.errCh)
}
