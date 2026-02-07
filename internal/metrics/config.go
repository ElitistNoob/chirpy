package metrics

import "sync/atomic"

type Metrics struct {
	fileserverHits atomic.Int32
}

func New() *Metrics {
	return &Metrics{}
}
