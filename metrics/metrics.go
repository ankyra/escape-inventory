/*
Copyright 2017 Ankyra

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	DownloadCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "escape_downloads",
			Help: "Number of downloads",
		},
	)
	UploadCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "escape_uploads",
			Help: "Number of uploads",
		},
	)
	ResponsesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "escape_http_responses_total",
			Help: "The count of http responses issued, classified by code and method",
		},
		[]string{"code", "method"},
	)
	ResponsesLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "escape_http_responses_latency",
			Help: "The latency of http responses issued, classified by code and method",
		},
		[]string{"code", "method"},
	)
)

func init() {
	metrics := []prometheus.Collector{
		DownloadCounter,
		UploadCounter,
		ResponsesTotal,
		ResponsesLatency,
	}
	for _, metric := range metrics {
		prometheus.MustRegister(metric)
	}
}
