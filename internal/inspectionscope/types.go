package inspectionscope

import (
	"errors"
	"sync"
)

var ErrBusy = errors.New("campaign inspection scope is busy")

// Scope serializes inspection-style work so that two unrelated campaigns never
// trip over each other. A campaign holds the scope while its inspection or
// overlap check runs; once every in-flight holder for that campaign reports
// done, the scope is released and a different campaign may acquire it.
//
// The depth counter preserves the original behavior of permitting concurrent
// same-campaign holders while still freeing the scope the moment the last one
// finishes, which is what prevents state from leaking across campaigns.
type Scope struct {
	mu     sync.Mutex
	owner  string
	depth  int
}

func New() *Scope { return &Scope{} }

// release decrements the holder count for campaignID and clears the owner
// once the last holder finishes. A mismatched campaignID is ignored rather
// than reported as busy: the Done methods run from deferred cleanup paths
// where returning an error would be dropped anyway, and silently no-oping
// keeps a late, mismatched done from corrupting an unrelated holder's count.
func (s *Scope) release(campaignID string) {
	if s.owner != campaignID || s.depth == 0 {
		return
	}
	s.depth--
	if s.depth <= 0 {
		s.depth = 0
		s.owner = ""
	}
}
