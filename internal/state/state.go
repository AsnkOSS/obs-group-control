// Package state holds the shared application state and notifies observers.
package state

import "sync"

// Phase is the recording lifecycle phase of the whole device group.
type Phase int

const (
	Idle Phase = iota
	Starting
	Recording
	Stopping
	Failed
)

func (p Phase) String() string {
	switch p {
	case Idle:
		return "Idle"
	case Starting:
		return "Starting..."
	case Recording:
		return "Recording"
	case Stopping:
		return "Stopping..."
	case Failed:
		return "Error"
	default:
		return "Unknown"
	}
}

// Snapshot is an immutable view of the current state.
type Snapshot struct {
	Phase       Phase
	DeviceCount int
	Message     string
}

// Store is a thread-safe state container with change notification.
type Store struct {
	mu       sync.Mutex
	snapshot Snapshot
	onChange func(Snapshot)
}

// NewStore creates a store; onChange is called after every update.
func NewStore(deviceCount int, onChange func(Snapshot)) *Store {
	return &Store{
		snapshot: Snapshot{Phase: Idle, DeviceCount: deviceCount},
		onChange: onChange,
	}
}

// Get returns the current snapshot.
func (s *Store) Get() Snapshot {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.snapshot
}

// Set updates the phase and message, then notifies the observer.
func (s *Store) Set(phase Phase, message string) {
	s.mu.Lock()
	s.snapshot.Phase = phase
	s.snapshot.Message = message
	snap := s.snapshot
	cb := s.onChange
	s.mu.Unlock()
	if cb != nil {
		cb(snap)
	}
}
