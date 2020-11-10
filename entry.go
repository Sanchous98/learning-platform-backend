package confucius

import "sync"

const (
	Undefined Status = iota
	Inactive
	Ok
	Serving
	Stopping
	Stopped
)

type Status uint8

type containerEntry struct {
	entryLock sync.Mutex
	service   Service
	name      string
	status    Status
}

func (e *containerEntry) getStatus() Status {
	e.entryLock.Lock()
	defer e.entryLock.Unlock()

	return e.status
}

func (e *containerEntry) setStatus(status Status) {
	e.entryLock.Lock()
	e.status = status
	e.entryLock.Unlock()
}

func (e *containerEntry) hasStatus(status Status) bool {
	return status == e.getStatus()
}
