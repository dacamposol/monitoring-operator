/*
Copyright 2022.

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

package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	monitoringv1alpha1 "github.com/dacamposol/monitoring-operator/api/v1alpha1"
	"github.com/dacamposol/monitoring-operator/pkg/alertmanager"
	silence2 "github.com/dacamposol/monitoring-operator/subroutines/silence"
)

const (
	AlertManagerDependencyFinalizer = "monitoring.dacamposol.com/alertmanager"
)

// SilenceReconciler reconciles a Getter object
type SilenceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=monitoring.dacamposol.com,resources=silences,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.dacamposol.com,resources=silences/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.dacamposol.com,resources=silences/finalizers,verbs=update
//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=alertmanagers,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *SilenceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reconcileLogger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("component", "SilenceReconciler").
		Str("name", req.Name).
		Str("namespace", req.Namespace).
		Logger()

	reconcileLogger.Info().Msg("Reconciling Silence")
	silence := monitoringv1alpha1.Silence{}
	err := r.Client.Get(ctx, req.NamespacedName, &silence)
	if err != nil {
		if errors.IsNotFound(err) {
			reconcileLogger.Info().Msg("Silence not found")
			return ctrl.Result{}, nil
		}
		reconcileLogger.Error().Err(err).Msg("Cannot start reconciliation, Kubernetes client error")
	}

	kubernetesClient := alertmanager.NewService(r.Client)
	httpClient := http.DefaultClient

	if silence.GetDeletionTimestamp() == nil &&
		!controllerutil.ContainsFinalizer(&silence, AlertManagerDependencyFinalizer) {
		return r.addFinalizerToResource(ctx, silence)
	}

	if silence.GetDeletionTimestamp() != nil {
		deletionSubroutine := silence2.NewDeletionSubroutine(kubernetesClient, httpClient)
		reconcileLogger.Info().Msg("Starting Silence deletion readiness subroutine")
		err := deletionSubroutine.Run(ctx, &silence)
		if err == nil {
			return r.removeFinalizerFromResource(ctx, silence)
		}
		reconcileLogger.Error().Err(err).Msg("Reconciliation of Silence deletion readiness failed")
	}

	if silence.GetDeletionTimestamp() == nil {
		reconcileSubroutine := silence2.NewReconcileSubroutine(kubernetesClient, httpClient)
		reconcileLogger.Info().Msg("Starting Silence reconciliation subroutine")
		err = reconcileSubroutine.Run(ctx, &silence)
		if err != nil {
			reconcileLogger.Error().Err(err).Msg("Reconciliation of Silence against AlertManager failed")
			return ctrl.Result{}, err
		}
	}

	err = r.Status().Update(ctx, &silence)
	if err != nil {
		reconcileLogger.Error().Err(err).Msg("Cannot update Status in Silence")
		return ctrl.Result{}, err
	}

	if meta.IsStatusConditionTrue(silence.Status.Conditions, silence2.ReconciledWithAlertManager) {
		return ctrl.Result{
			RequeueAfter: 10 * time.Minute,
		}, nil
	}

	return ctrl.Result{
		RequeueAfter: 30 * time.Second,
	}, nil
}

func (r *SilenceReconciler) addFinalizerToResource(ctx context.Context, silence monitoringv1alpha1.Silence) (ctrl.Result, error) {
	controllerutil.AddFinalizer(&silence, AlertManagerDependencyFinalizer)
	err := r.Update(ctx, &silence)
	if err != nil {
		return ctrl.Result{}, err
	}
	// Schedule a new reconciliation cycle upon finalizer inclusion
	return ctrl.Result{Requeue: true}, nil
}

func (r *SilenceReconciler) removeFinalizerFromResource(ctx context.Context, silence monitoringv1alpha1.Silence) (ctrl.Result, error) {
	controllerutil.RemoveFinalizer(&silence, AlertManagerDependencyFinalizer)
	err := r.Update(ctx, &silence)
	if err != nil {
		return ctrl.Result{}, err
	}
	// When the finalizer is removed, the Kubernetes API will take care of deleting the resource
	return ctrl.Result{Requeue: true}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SilenceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("SilenceReconciler").
		For(&monitoringv1alpha1.Silence{}).
		Complete(r)
}
