package main

import (
	"fmt"
	"io"
	"sync"
)

type aggregator struct {
	countersLock sync.RWMutex
	guagesLock   sync.RWMutex
	timingsLock  sync.RWMutex
	setsLock     sync.RWMutex

	counters map[string]int64
	guages   map[string]int64
	timings  map[string][]int64
	sets     map[string]bool
}

func newAggregator() *aggregator {
	result := &aggregator{}
	result.counters = make(map[string]int64)
	result.guages = make(map[string]int64)
	result.timings = make(map[string][]int64)
	result.sets = make(map[string]bool)

	return result
}

func (agg *aggregator) processCounter(m *Metric) {
	agg.countersLock.Lock()
	defer agg.countersLock.Unlock()

	// Counters are only 64bit ints.
	agg.counters[string(m.Bucket)] += int64(m.Value * (1 / m.SamplingRate))
}

func (agg *aggregator) processGuage(m *Metric) {
	agg.guagesLock.Lock()
	defer agg.guagesLock.Unlock()

	// Guages are only 64bit ints.
	agg.guages[string(m.Bucket)] = int64(m.Value)
}

func (agg *aggregator) processTiming(m *Metric) {
	agg.timingsLock.Lock()
	defer agg.timingsLock.Unlock()

	// Timings are only 64bit ints.
	agg.timings[string(m.Bucket)] = append(agg.timings[string(m.Bucket)], int64(m.Value))
}

func (agg *aggregator) processSet(m *Metric) {
	agg.setsLock.Lock()
	defer agg.setsLock.Unlock()

	agg.sets[string(m.Bucket)] = true
}

func (agg *aggregator) writeCounters(w io.Writer) {
	agg.countersLock.RLock()
	defer agg.countersLock.RUnlock()
	fmt.Fprintf(w, "Dumping counters:\n--------------\n")
	for k, v := range agg.counters {
		fmt.Fprintf(w, "%s: %d\n", k, v)
	}

	fmt.Fprintf(w, "Total Counters: %d\n--------------\n", len(agg.counters))
}

func (agg *aggregator) writeGuages(w io.Writer) {
	agg.guagesLock.RLock()
	defer agg.guagesLock.RUnlock()

	fmt.Fprintf(w, "Dumping guages:\n--------------\n")
	for k, v := range agg.guages {
		fmt.Fprintf(w, "%s: %d\n", k, v)
	}

	fmt.Fprintf(w, "Total guages: %d\n--------------\n", len(agg.guages))
}

func (agg *aggregator) writeTimings(w io.Writer) {
	agg.timingsLock.RLock()
	defer agg.timingsLock.RUnlock()

	fmt.Fprintf(w, "Dumping timings:\n--------------\n")
	for k, v := range agg.timings {
		fmt.Fprintf(w, "%s:", k)
		for _, i := range v {
			fmt.Fprintf(w, " %d", i)
		}
		fmt.Fprintf(w, "\n")

	}

	fmt.Fprintf(w, "Total timings: %d\n--------------\n", len(agg.timings))
}
