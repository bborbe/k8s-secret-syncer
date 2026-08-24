// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"os"
	"time"

	libhttp "github.com/bborbe/http"
	libk8s "github.com/bborbe/k8s"
	"github.com/bborbe/log"
	libmetrics "github.com/bborbe/metrics"
	"github.com/bborbe/run"
	libsentry "github.com/bborbe/sentry"
	"github.com/bborbe/service"
	libtime "github.com/bborbe/time"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/bborbe/k8s-secret-syncer/pkg/factory"
)

func main() {
	app := &application{}
	os.Exit(service.Main(context.Background(), app, &app.SentryDSN, &app.SentryProxy))
}

type application struct {
	Listen          string            `required:"true"  arg:"listen"            env:"LISTEN"            usage:"address to listen to"`
	SentryDSN       string            `required:"false" arg:"sentry-dsn"        env:"SENTRY_DSN"        usage:"Sentry DSN"                                               display:"length"`
	SentryProxy     string            `required:"false" arg:"sentry-proxy"      env:"SENTRY_PROXY"      usage:"Sentry Proxy"`
	SourceNamespace libk8s.Namespace  `required:"true"  arg:"source-namespace"  env:"SOURCE_NAMESPACE"  usage:"source namespace"`
	TargetNamespace libk8s.Namespace  `required:"true"  arg:"target-namespace"  env:"TARGET_NAMESPACE"  usage:"target namespace"`
	SecretPrefix    string            `required:"true"  arg:"secret-prefix"     env:"SECRET_PREFIX"     usage:"only secret with this prefix are synced"`
	Kubeconfig      string            `required:"false" arg:"kubeconfig"        env:"KUBECONFIG"        usage:"Path to k8s config"`
	BuildGitVersion string            `required:"false" arg:"build-git-version" env:"BUILD_GIT_VERSION" usage:"Build Git version (git describe --tags --always --dirty)"                  default:"dev"`
	BuildGitCommit  string            `required:"false" arg:"build-git-commit"  env:"BUILD_GIT_COMMIT"  usage:"Build Git commit hash"                                                     default:"none"`
	BuildDate       *libtime.DateTime `required:"false" arg:"build-date"        env:"BUILD_DATE"        usage:"Build timestamp (RFC3339)"`
}

func (a *application) Run(ctx context.Context, sentryClient libsentry.Client) error {
	libmetrics.NewBuildInfoMetrics().SetBuildInfo(a.BuildGitVersion, a.BuildGitCommit, a.BuildDate)

	return service.Run(
		ctx,
		a.createSecretSyncer(),
		a.createHTTPServer(),
	)
}

func (a *application) createSecretSyncer() run.Func {
	return factory.CreateSecretSyncer(
		a.Kubeconfig,
		a.SecretPrefix,
		libk8s.Namespace(a.SourceNamespace),
		libk8s.Namespace(a.TargetNamespace),
	)
}

func (a *application) createHTTPServer() run.Func {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		router := mux.NewRouter()
		router.Path("/healthz").Handler(libhttp.NewPrintHandler("OK"))
		router.Path("/readiness").Handler(libhttp.NewPrintHandler("OK"))
		router.Path("/metrics").Handler(promhttp.Handler())
		router.Path("/setloglevel/{level}").
			Handler(log.NewSetLoglevelHandler(ctx, log.NewLogLevelSetter(2, 5*time.Minute)))
		router.Path("/gc").Handler(libhttp.NewGarbageCollectorHandler())

		glog.V(2).Infof("starting http server listen on %s", a.Listen)
		return libhttp.NewServer(
			a.Listen,
			router,
		).Run(ctx)
	}
}
