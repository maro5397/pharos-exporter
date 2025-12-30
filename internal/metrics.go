package internal

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricsOnce sync.Once

	ProposeTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "validator_propose_total",
		Help: "Total number of propose attempts observed in logs.",
	})
	LastProposeTimestamp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "validator_last_propose_timestamp",
		Help: "Unix timestamp of the last propose event observed in logs.",
	})
	EndorseTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "validator_endorse_total",
		Help: "Total number of endorse events observed in logs.",
	})
	LastEndorseTimestamp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "validator_last_endorse_timestamp",
		Help: "Unix timestamp of the last endorse event observed in logs.",
	})

	VoteInclusionTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "validator_vote_inclusion_total",
		Help: "Total number of blocks where the validator vote was included.",
	})
	VoteInclusionTimestamp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "validator_vote_inclusion_timestamp",
		Help: "Unix timestamp when the validator vote was last included.",
	})
	ActiveTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "validator_active_total",
		Help: "Total number of blocks where the validator was active in the validator set.",
	})
	ActiveTimestamp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "validator_active_timestamp",
		Help: "Unix timestamp when validator active status was last observed.",
	})
)

func RegisterMetrics() {
	metricsOnce.Do(func() {
		prometheus.MustRegister(
			ProposeTotal,
			LastProposeTimestamp,
			EndorseTotal,
			LastEndorseTimestamp,
			VoteInclusionTotal,
			VoteInclusionTimestamp,
			ActiveTotal,
			ActiveTimestamp,
		)
	})
}
