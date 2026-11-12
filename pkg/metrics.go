// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	"github.com/prometheus/client_golang/prometheus"
)

// NewLastUpdateTimestampGauge returns the gauge that records, per target
// namespace, when a secret was last delivered successfully.
//
// It is deliberately not pre-initialised with a zero value. A zero would make
// a "no delivery for N minutes" alert fire on a freshly started pod before its
// first sync, and it would also hide the strongest signal this metric carries:
// an absent series means the process itself is gone. Alerts should therefore
// read `time() - k8s_secret_syncer_last_update_timestamp_seconds > 600 and
// k8s_secret_syncer_last_update_timestamp_seconds > 0`, which fires only once
// the syncer has delivered at least once and then stopped.
func NewLastUpdateTimestampGauge() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "k8s_secret_syncer_last_update_timestamp_seconds",
			Help: "Unix timestamp of the last successful secret delivery, by target namespace.",
		},
		[]string{"namespace"},
	)
}
