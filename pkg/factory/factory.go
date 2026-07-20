// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package factory

import (
	"context"

	"github.com/bborbe/errors"
	libk8s "github.com/bborbe/k8s"
	"github.com/bborbe/run"
	libtime "github.com/bborbe/time"
	k8s_kubernetes "k8s.io/client-go/kubernetes"

	"github.com/bborbe/k8s-secret-syncer/pkg"
)

func CreateSecretSyncer(
	kubeconfig string,
	secretPrefix string,
	sourceNamespace libk8s.Namespace,
	targetNamespace libk8s.Namespace,
) run.Func {
	return func(ctx context.Context) error {
		clientset, err := libk8s.CreateClientset(kubeconfig)
		if err != nil {
			return errors.Wrap(ctx, err, "create clientset failed")
		}
		return CreateSecretWatcher(
			clientset,
			secretPrefix,
			sourceNamespace,
			targetNamespace,
		).Watch(ctx)
	}
}

func CreateSecretWatcher(
	clientset k8s_kubernetes.Interface,
	secretPrefix string,
	sourceNamespace libk8s.Namespace,
	targetNamespace libk8s.Namespace,
) libk8s.SecretWatcher {
	return libk8s.NewSecretWatcherRetry(
		libk8s.NewSecretWatcher(
			clientset,
			pkg.NewSecretFilter(
				libk8s.NewSecretEventProcessorSkipError(
					pkg.NewSecretManager(
						clientset,
						targetNamespace,
					),
				),
				secretPrefix,
			),
			sourceNamespace,
		),
		libtime.NewWaiterDuration(),
		1*libtime.Second,
	)
}
