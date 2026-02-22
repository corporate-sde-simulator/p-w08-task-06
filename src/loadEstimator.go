package loadshed

// Load Estimator — estimates system load from various metrics.
//
// Author: Suresh Kumar (Infra team)
// Last Modified: 2026-03-25

import (
	"math"
	"time"
)

type LoadSample struct {
	CPUPercent    float64
	MemoryPercent float64
	QueueDepth    int
	Latency99     time.Duration
	Timestamp     time.Time
}

type LoadEstimator struct {
	samples    []LoadSample
	maxSamples int
	weights    map[string]float64
}

func NewLoadEstimator(maxSamples int) *LoadEstimator {
	return &LoadEstimator{
		samples:    make([]LoadSample, 0),
		maxSamples: maxSamples,
		weights: map[string]float64{
			"cpu":     0.4,
			"memory":  0.2,
			"queue":   0.25,
			"latency": 0.15,
		},
	}
}

func (le *LoadEstimator) AddSample(sample LoadSample) {
	sample.Timestamp = time.Now()
	le.samples = append(le.samples, sample)
	if len(le.samples) > le.maxSamples {
		le.samples = le.samples[1:]
	}
}

func (le *LoadEstimator) EstimateLoad() float64 {
	if len(le.samples) == 0 {
		return 0
	}

	latest := le.samples[len(le.samples)-1]

	cpuLoad := latest.CPUPercent / 100.0
	memLoad := latest.MemoryPercent / 100.0
	queueLoad := math.Min(float64(latest.QueueDepth)/1000.0, 1.0)
	latencyLoad := math.Min(float64(latest.Latency99.Milliseconds())/5000.0, 1.0)

	weighted := cpuLoad*le.weights["cpu"] +
		memLoad*le.weights["memory"] +
		queueLoad*le.weights["queue"] +
		latencyLoad*le.weights["latency"]

	return math.Min(weighted, 1.0)
}

func (le *LoadEstimator) GetTrend(window time.Duration) string {
	if len(le.samples) < 2 {
		return "insufficient_data"
	}

	cutoff := time.Now().Add(-window)
	var recent, older []float64

	midpoint := time.Now().Add(-window / 2)
	for _, s := range le.samples {
		if s.Timestamp.After(cutoff) {
			load := s.CPUPercent / 100.0
			if s.Timestamp.After(midpoint) {
				recent = append(recent, load)
			} else {
				older = append(older, load)
			}
		}
	}

	if len(recent) == 0 || len(older) == 0 {
		return "insufficient_data"
	}

	avgRecent := avg(recent)
	avgOlder := avg(older)

	if avgRecent > avgOlder*1.2 {
		return "increasing"
	} else if avgRecent < avgOlder*0.8 {
		return "decreasing"
	}
	return "stable"
}

func avg(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}
