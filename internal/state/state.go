package state

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/yourorg/portwatch/internal/scanner"
)

// Snapshot holds a persisted port state with metadata.
type Snapshot struct {
	Ports     []scanner.Port `json:"ports"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// Store persists and retrieves port state snapshots to/from disk.
type Store struct {
	mu   sync.RWMutex
	path string
}

// NewStore creates a Store that persists state to the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the current port set to disk as a JSON snapshot.
func (s *Store) Save(ports scanner.PortSet) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	slice := make([]scanner.Port, 0, len(ports))
	for _, p := range ports {
		slice = append(slice, p)
	}

	snap := Snapshot{
		Ports:     slice,
		UpdatedAt: time.Now().UTC(),
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}

// Load reads a previously saved snapshot from disk.
// Returns an empty Snapshot and no error if the file does not exist.
func (s *Store) Load() (Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return Snapshot{}, nil
	}
	if err != nil {
		return Snapshot{}, err
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}

// ToPortSet converts a Snapshot back into a PortSet for diffing.
func (s Snapshot) ToPortSet() scanner.PortSet {
	ports := make([]scanner.Port, len(s.Ports))
	copy(ports, s.Ports)
	return scanner.PortSetFromSlice(ports)
}
