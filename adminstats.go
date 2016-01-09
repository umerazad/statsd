package main

import (
	"encoding/json"
	"log"
	"time"
)

type adminStats struct {
	// counters
	TotalRecords          int           `json:"total_records"`
	BadRecords            int           `json:"bad_records"`
	CounterRecords        int           `json:"counters_received"`
	GuageRecords          int           `json:"guages_received"`
	TimingRecords         int           `json:"timers_received"`
	SetRecords            int           `json:"sets_received"`
	DumpCounterRequests   int           `json:"dump_counters_requests"`
	DumpTimingRequests    int           `json:"dump_timers_requests"`
	DumpGuageRequests     int           `json:"dump_guages_requests"`
	DumpSetRequests       int           `json:"dump_sets_requests"`
	DelCountersRequests   int           `json:"delcounters_requests"`
	DelGuagesRequests     int           `json:"delguages_requests"`
	DelTimersRequests     int           `json:"deltimers_requests"`
	HealthRequests        int           `json:"health_requests"`
	LastFlushTimestamp    time.Duration `json:"last_flush_timestamp"`
	TotalAdminConnections int           `json:"total_admin_connections"`
}

func (s *adminStats) String() string {
	result, err := json.Marshal(stats)
	if err != nil {
		log.Printf("Failed to convert stats to json: %v", err)
	}
	return string(result)
}

func newAdminStats() *adminStats {
	return &adminStats{}
}
