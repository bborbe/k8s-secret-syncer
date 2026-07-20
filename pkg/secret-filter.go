// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	"context"
	"strings"

	libk8s "github.com/bborbe/k8s"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
)

func NewSecretFilter(
	secretManager libk8s.SecretEventProcessor,
	prefix string,
) libk8s.SecretEventProcessor {
	return &secretFilter{
		secretManager: secretManager,
		prefix:        prefix,
	}
}

type secretFilter struct {
	secretManager libk8s.SecretEventProcessor
	prefix        string
}

func (s *secretFilter) OnUpdate(ctx context.Context, secret corev1.Secret) error {
	if !strings.HasPrefix(secret.Name, s.prefix) {
		glog.V(3).Infof("secret '%s' has no prefix '%s' => skip update", secret.Name, s.prefix)
		return nil
	}
	return s.secretManager.OnUpdate(ctx, secret)
}

func (s *secretFilter) OnDelete(ctx context.Context, secret corev1.Secret) error {
	if !strings.HasPrefix(secret.Name, s.prefix) {
		glog.V(3).Infof("secret '%s' has no prefix '%s' => skip delete", secret.Name, s.prefix)
		return nil
	}
	return s.secretManager.OnDelete(ctx, secret)
}
