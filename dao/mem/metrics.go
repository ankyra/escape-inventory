package mem

import (
	. "github.com/ankyra/escape-inventory/dao/types"
)

func (a *dao) GetUserMetrics(userID string) (*Metrics, error) {
	metrics, found := a.metrics[userID]
	if !found {
		a.metrics[userID] = NewMetrics(0)
		return a.metrics[userID], nil
	}
	return metrics, nil
}

func (a *dao) SetUserMetrics(userID string, previous, new *Metrics) error {
	_, found := a.metrics[userID]
	if !found {
		return NotFound
	}
	a.metrics[userID] = new
	return nil
}
