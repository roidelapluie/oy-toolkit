// Copyright 2022 The o11y toolkit Authors
// spdx-license-identifier: apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	webConfig              = webflag.AddFlags(kingpin.CommandLine)
	listenAddress          = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9101").String()
	metricsPath            = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	disableExporterMetrics = kingpin.Flag("web.disable-exporter-metrics", "Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).").Bool()
)

func Serve(logger log.Logger, r *prometheus.Registry) {
	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)

	if !*disableExporterMetrics {
		r.MustRegister(
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
			collectors.NewGoCollector(),
		)
	}

	http.Handle(*metricsPath, promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	srv := &http.Server{Addr: *listenAddress}
	if err := web.ListenAndServe(srv, *webConfig, logger); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
