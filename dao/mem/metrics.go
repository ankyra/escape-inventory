package mem

import (
	. "github.com/ankyra/escape-inventory/dao/types"
)

func (a *dao) GetUserMetrics(username string) (*Metrics, error) {
	metrics, found := a.metrics[username]
	if !found {
		a.metrics[username] = NewMetrics(0)
		return a.metrics[username], nil
	}
	return metrics, nil
}

func (a *dao) SetUserMetrics(username string, previous, new *Metrics) error {
	_, found := a.metrics[username]
	if !found {
		return NotFound
	}
	a.metrics[username] = new
	return nil
}
