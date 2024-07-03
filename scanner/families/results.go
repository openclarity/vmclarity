package families

import (
	"errors"
	"sync"
)

// Results stores results from all families. Safe for concurrent usage.
type Results struct {
	mu      sync.RWMutex
	results []any
}

func NewResults() *Results {
	return &Results{}
}

func (r *Results) SetFamilyResult(result any) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.results = append(r.results, result)
}

// GetFamilyResult returns results for a specific family from the given results.
func GetFamilyResult[familyType any](r *Results) (familyType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, result := range r.results {
		res, ok := result.(familyType)
		if ok {
			return res, nil
		}
	}

	var res familyType
	return res, errors.New("missing result")
}
