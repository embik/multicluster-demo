/*
Copyright 2025 The KCP Authors.

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

package controller

import (
	"context"
	"fmt"

	mcbuilder "sigs.k8s.io/multicluster-runtime/pkg/builder"
	mcmanager "sigs.k8s.io/multicluster-runtime/pkg/manager"
	mcreconcile "sigs.k8s.io/multicluster-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type ConfigMapReconciler struct {
	Mgr mcmanager.Manager
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigMapReconciler) SetupWithManager(mgr mcmanager.Manager) error {
	return mcbuilder.ControllerManagedBy(mgr).
		Named("multicluster-configmap-controller").
		For(&corev1.ConfigMap{}).
		Complete(r)
}

func (r *ConfigMapReconciler) Reconcile(ctx context.Context, req mcreconcile.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("cluster", req.ClusterName)

	cl, err := r.Mgr.GetCluster(ctx, req.ClusterName)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get cluster: %w", err)
	}
	client := cl.GetClient()

	// Retrieve the ConfigMap from the cluster.
	s := &corev1.ConfigMap{}
	if err := client.Get(ctx, req.NamespacedName, s); err != nil {
		if apierrors.IsNotFound(err) {
			// ConfigMap was deleted.
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to get configmap: %w", err)
	}

	log.Info("Reconciling ConfigMap", "name", s.Name, "uuid", s.UID)

	return ctrl.Result{}, nil
}
