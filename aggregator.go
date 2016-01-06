package main

import "sync"

type aggregator struct {
	countersLock sync.RWMutex
	guagesLock   sync.RWMutex
	timingsLock  sync.RWMutex
	setsLock     sync.RWMutex

	counters map[string]int64
	guages   map[string]int64
	timings  map[string]int64
	sets     map[string]bool
}

func newAggregator() *aggregator {
	result := &aggregator{}
	result.counters = make(map[string]int64)
	result.guages = make(map[string]int64)
	result.timings = make(map[string]int64)
	result.sets = make(map[string]bool)

	return result
}

func (agg *aggregator) processCounter(m *Metric) {
	agg.countersLock.Lock()
	defer agg.countersLock.Unlock()

	// Counters are only 64bit ints.
	agg.counters[string(m.Bucket)] += int64(m.Value)
}

func (agg *aggregator) processGuage(m *Metric) {
	agg.guagesLock.Lock()
	defer agg.guagesLock.Unlock()

	// Guages are only 64bit ints.
	agg.guages[string(m.Bucket)] += int64(m.Value * (1.0 / m.SamplingRate))
}

func (agg *aggregator) processTiming(m *Metric) {
	agg.timingsLock.Lock()
	defer agg.timingsLock.Unlock()

	// Timings are only 64bit ints.
	agg.timings[string(m.Bucket)] += int64(m.Value * (1.0 / m.SamplingRate))
}

func (agg *aggregator) processSet(m *Metric) {
	agg.setsLock.Lock()
	defer agg.setsLock.Unlock()

	agg.sets[string(m.Bucket)] = true
}
