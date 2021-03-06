/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package apisrv

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// OpenCensus Measures. These are exported as metrics to your monitoring system
// https://godoc.org/go.opencensus.io/stats
//
// When making opencensus stats, the 'name' param, with forward slashes changed
// to underscores, is appended to the 'namespace' value passed to the
// prometheus exporter to become the Prometheus metric name. You can also look
// into having Prometheus rewrite your metric names on scrape.
//
//  For example:
//   - defining the promethus export namespace "open_match" when instanciating the exporter:
//			pe, err := promethus.NewExporter(promethus.Options{Namespace: "open_match"})
//   - and naming the request counter "backend/requests_total":
//			MGrpcRequests := stats.Int64("backendapi/requests_total", ...
//   - results in the prometheus metric name:
//			open_match_backendapi_requests_total
//   - [note] when using opencensus views to aggregate the metrics into
//     distribution buckets and such, multiple metrics
//     will be generated with appended types ("<metric>_bucket",
//     "<metric>_count", "<metric>_sum", for example)
//
// In addition, OpenCensus stats propogated to Prometheus have the following
// auto-populated labels pulled from kubernetes, which we should avoid to
// prevent overloading and having to use the HonorLabels param in Prometheus.
//
// - Information about the k8s pod being monitored:
//		"pod" (name of the monitored k8s pod)
//		"namespace" (k8s namespace of the monitored pod)
// - Information about how promethus is gathering the metrics:
//		"instance" (IP and port number being scraped by prometheus)
//		"job" (name of the k8s service being scraped by prometheus)
//		"endpoint" (name of the k8s port in the k8s service being scraped by prometheus)
//
var (
	// Logging instrumentation
	// There's no need to record this measurement directly if you use
	// the logrus hook provided in metrics/helper.go after instantiating the
	// logrus instance in your application code.
	// https://godoc.org/github.com/sirupsen/logrus#LevelHooks
	BeLogLines = stats.Int64("backendapi/logs_total", "Number of Backend API lines logged", "1")

	// Failure instrumentation
	BeFailures = stats.Int64("backendapi/failures_total", "Number of Backend API failures", "1")

	// Counting operations
	BeAssignments                = stats.Int64("backendapi/assignments_total", "Number of players assigned to matches", "1")
	BeAssignmentFailures         = stats.Int64("backendapi/assignment/failures_total", "Number of player match assigment failures", "1")
	BeAssignmentDeletions        = stats.Int64("backendapi/assignment/deletions_total", "Number of player match assigment deletions", "1")
	BeAssignmentDeletionFailures = stats.Int64("backendapi/assignment/deletions/failures_total", "Number of player match assigment deletion failures", "1")
)

var (
	// KeyMethod is used to tag a measure with the currently running API method.
	KeyMethod, _ = tag.NewKey("method")
	// KeySeverity is used to tag a the severity of a log message.
	KeySeverity, _ = tag.NewKey("severity")
)

// Package metrics provides some convience views.
// You need to register the views for the data to actually be collected.
// Note: The OpenCensus View 'Description' is exported to Prometheus as the HELP string.
// Note: If you get a "Failed to export to Prometheus: inconsistent label
// cardinality" error, chances are you forgot to set the tags specified in the
// view for a given measure when you tried to do a stats.Record()
var (
	BeLogCountView = &view.View{
		Name:        "log_lines/total",
		Measure:     BeLogLines,
		Description: "The number of lines logged",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{KeySeverity},
	}

	BeFailureCountView = &view.View{
		Name:        "failures",
		Measure:     BeFailures,
		Description: "The number of failures",
		Aggregation: view.Count(),
	}

	BeAssignmentCountView = &view.View{
		Name:        "backend/assignments",
		Measure:     BeAssignments,
		Description: "The number of successful player match assignments",
		Aggregation: view.Count(),
	}

	BeAssignmentFailureCountView = &view.View{
		Name:        "backend/assignments/failures",
		Measure:     BeAssignmentFailures,
		Description: "The number of player match assignment failures",
		Aggregation: view.Count(),
	}

	BeAssignmentDeletionCountView = &view.View{
		Name:        "backend/assignments/deletions",
		Measure:     BeAssignmentDeletions,
		Description: "The number of successful player match assignments",
		Aggregation: view.Count(),
	}

	BeAssignmentDeletionFailureCountView = &view.View{
		Name:        "backend/assignments/deletions/failures",
		Measure:     BeAssignmentDeletionFailures,
		Description: "The number of player match assignment failures",
		Aggregation: view.Count(),
	}
)

// DefaultBackendAPIViews are the default backend API OpenCensus measure views.
var DefaultBackendAPIViews = []*view.View{
	BeLogCountView,
	BeFailureCountView,
	BeAssignmentCountView,
	BeAssignmentFailureCountView,
	BeAssignmentDeletionCountView,
	BeAssignmentDeletionFailureCountView,
}
