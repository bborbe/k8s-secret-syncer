// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import (
	"context"

	"github.com/bborbe/errors"
	libk8s "github.com/bborbe/k8s"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// SecretManager implements libk8s.SecretEventProcessor for syncing secrets
// from source namespace to target namespace.
//
//counterfeiter:generate -o ../mocks/secret-manager.go --fake-name SecretManager . SecretManager
type SecretManager = libk8s.SecretEventProcessor

func NewSecretManager(
	clientset kubernetes.Interface,
	namespace libk8s.Namespace,
) libk8s.SecretEventProcessor {
	return &secretManager{
		clientset: clientset,
		namespace: namespace,
	}
}

type secretManager struct {
	clientset kubernetes.Interface
	namespace libk8s.Namespace
}

func (s *secretManager) OnUpdate(ctx context.Context, secret corev1.Secret) error {
	secret = corev1.Secret{
		TypeMeta: secret.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name:        secret.Name,
			Namespace:   s.namespace.String(),
			Annotations: secret.Annotations,
			Labels:      secret.Labels,
		},
		Immutable:  secret.Immutable,
		Data:       secret.Data,
		StringData: secret.StringData,
		Type:       secret.Type,
	}

	currentSecret, err := s.clientset.CoreV1().
		Secrets(secret.Namespace).
		Get(ctx, secret.Name, metav1.GetOptions{})
	if err != nil {
		_, err = s.clientset.CoreV1().
			Secrets(secret.Namespace).
			Create(ctx, &secret, metav1.CreateOptions{})
		if err != nil {
			return errors.Wrap(ctx, err, "create secret failed")
		}
		glog.V(2).Infof("secret %s/%s created successful", secret.Namespace, secret.Name)
		return nil
	}
	secret.ResourceVersion = currentSecret.ResourceVersion
	_, err = s.clientset.CoreV1().
		Secrets(secret.Namespace).
		Update(ctx, &secret, metav1.UpdateOptions{})
	if err != nil {
		return errors.Wrap(ctx, err, "update secret failed")
	}
	glog.V(2).Infof("secret %s/%s updated successful", secret.Namespace, secret.Name)
	return nil
}

func (s *secretManager) OnDelete(ctx context.Context, secret corev1.Secret) error {
	if err := s.clientset.CoreV1().Secrets(s.namespace.String()).Delete(ctx, secret.Name, metav1.DeleteOptions{}); err != nil {
		return errors.Wrap(ctx, err, "delete secret failed")
	}
	glog.V(2).Infof("secret %s/%s deleted successful", s.namespace, secret.Name)
	return nil
}
