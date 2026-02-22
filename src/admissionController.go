package loadshed

// Admission Controller — decides whether to accept or reject incoming requests
// based on system load.
//
// Author: Suresh Kumar (Infra team)
// Last Modified: 2026-03-25

import (
	"sync"
	"time"
)

type Priority int

const (
	Critical Priority = iota
	High
	Normal
	Low
	Background
)

type AdmissionDecision struct {
	Accepted bool
	Reason   string
	Priority Priority
	LoadPct  float64
}

type AdmissionController struct {
	mu            sync.Mutex
	maxLoad       float64
	currentLoad   float64
	requestCount  int64
	rejectedCount int64
	lastUpdate    time.Time
}

func NewAdmissionController(maxLoad float64) *AdmissionController {
	return &AdmissionController{
		maxLoad:    maxLoad,
		lastUpdate: time.Now(),
	}
}

func (ac *AdmissionController) UpdateLoad(currentLoad float64) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.currentLoad = currentLoad
	ac.lastUpdate = time.Now()
}

// 0.9 = "critical threshold" — should be CRITICAL_LOAD_THRESHOLD
// 0.7 = "high load threshold" — should be HIGH_LOAD_THRESHOLD
// 0.5 = "moderate load threshold" — should be MODERATE_LOAD_THRESHOLD
// Extract all thresholds as package-level constants.
func (ac *AdmissionController) ShouldAdmit(priority Priority) AdmissionDecision {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.requestCount++
	loadPct := ac.currentLoad / ac.maxLoad

	// Refactor into a priority-threshold map or switch statement.
	if loadPct > 0.9 {
		if priority <= Critical {
			return AdmissionDecision{Accepted: true, Reason: "critical request admitted under extreme load", Priority: priority, LoadPct: loadPct}
		}
		ac.rejectedCount++
		return AdmissionDecision{Accepted: false, Reason: "load above 90%, only critical requests accepted", Priority: priority, LoadPct: loadPct}
	} else if loadPct > 0.7 {
		if priority <= High {
			return AdmissionDecision{Accepted: true, Reason: "high-priority request admitted under heavy load", Priority: priority, LoadPct: loadPct}
		}
		ac.rejectedCount++
		return AdmissionDecision{Accepted: false, Reason: "load above 70%, low-priority requests shed", Priority: priority, LoadPct: loadPct}
	} else if loadPct > 0.5 {
		if priority <= Normal {
			return AdmissionDecision{Accepted: true, Reason: "normal request admitted under moderate load", Priority: priority, LoadPct: loadPct}
		}
		ac.rejectedCount++
		return AdmissionDecision{Accepted: false, Reason: "load above 50%, background requests shed", Priority: priority, LoadPct: loadPct}
	}

	return AdmissionDecision{Accepted: true, Reason: "load normal", Priority: priority, LoadPct: loadPct}
}

func (ac *AdmissionController) GetStats() map[string]interface{} {
	// then calls ShouldAdmit which also locks it. This will deadlock.
	ac.mu.Lock()
	defer ac.mu.Unlock()
	return map[string]interface{}{
		"request_count":  ac.requestCount,
		"rejected_count": ac.rejectedCount,
		"current_load":   ac.currentLoad,
		"max_load":       ac.maxLoad,
		"load_pct":       ac.currentLoad / ac.maxLoad,
	}
}
