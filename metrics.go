package main

import (
	"bytes"
	"fmt"
	"strconv"
)

// MetricType represents the type of a telemetry sample.
type MetricType int

const (
	// UNKNOWN represents an unrecognized metric type.
	UNKNOWN MetricType = iota

	// COUNTER represents counting metrics.
	COUNTER = 1 << iota

	// TIMING represents timings type metric.
	TIMING

	// GUAGE represents guages type metrics.
	GUAGE

	// SET type represents unique value counter.
	SET
)

// Metric is used to represent all types of statsd telemetry samples.
type Metric struct {
	Type         MetricType
	Bucket       []byte
	ValueStr     []byte
	Value        float64
	SamplingRate float64
}

func (m *Metric) String() string {
	return fmt.Sprintf("{%s, %s, %s, %f}", m.Type, m.Bucket, m.ValueStr, m.SamplingRate)
}

func (m MetricType) String() string {
	switch m {
	case COUNTER:
		return "counter"
	case TIMING:
		return "timing"
	case GUAGE:
		return "guage"
	case SET:
		return "set"
	}

	return "unknown"
}

// ParseMetric reads a statsd metric from the given string.
func parseMetric(record []byte) (*Metric, error) {
	record = bytes.TrimSpace(record)

	if len(record) == 0 {
		return nil, fmt.Errorf("Parse error: empty record")
	}

	bucket, rest := tokenize(':', record)
	if len(bucket) == 0 {
		return nil, fmt.Errorf("Malformed record: No bucket name.")
	}

	if len(rest) == 0 {
		// Statsd spec is fuzzy about it and there are Statsd implementations
		// that treat "a" as "a:1|c" but we'll drop these requests.
		return nil, fmt.Errorf("Malformed record: No value/type found.")
	}

	value, rest := tokenize('|', rest)
	if len(value) == 0 {
		return nil, fmt.Errorf("Malformed record: No value found after '|'")
	}

	mtype, rest := tokenize('|', rest)
	if len(mtype) == 0 {
		return nil, fmt.Errorf("Malformed record: No 'type' found in %q", record)
	}

	_, samplingRate := tokenize('@', rest)

	return newMetric(bucket, value, mtype, samplingRate)
}

func newMetric(bucket, value, mtype, rate []byte) (*Metric, error) {
	var metricType MetricType
	var samplingRate float64 = 1
	var floatValue float64
	var err error

	switch string(mtype) {
	case "g":
		metricType = GUAGE
	case "c":
		metricType = COUNTER
	case "ms":
		metricType = TIMING
	case "s":
		metricType = SET
	default:
		return nil, fmt.Errorf("Malformed metric. Failed to parse type: %s", mtype)
	}

	if metricType != SET {
		floatValue, err = strconv.ParseFloat(string(value), 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse value: %s", value)
		}

		if len(rate) > 0 {
			samplingRate, err = strconv.ParseFloat(string(rate), 64)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse sampling rate: %s: %v", rate, err)
			}
		}

	}

	return &Metric{metricType, bucket, value, floatValue, samplingRate}, nil
}

func tokenize(needle byte, source []byte) ([]byte, []byte) {
	index := bytes.IndexByte(source, needle)
	if index == -1 {
		return source, nil
	}

	return bytes.TrimSpace(source[:index]), source[index+1:]
}
