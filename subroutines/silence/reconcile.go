package silence

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/dacamposol/monitoring-operator/api/v1alpha1"
	"github.com/dacamposol/monitoring-operator/pkg/alertmanager"
	"github.com/dacamposol/monitoring-operator/pkg/silence"
	"github.com/dacamposol/monitoring-operator/subroutines/common/conditions"
)

const ReconciledWithAlertManager = "AlertManagerReconciled"

type ReconcileSubroutine struct {
	kubernetesClient alertmanager.Getter
	httpClient       silence.HttpClient
	conditionsSetter conditions.Setter
}

func NewReconcileSubroutine(kubernetes alertmanager.Getter, httpClient *http.Client) *ReconcileSubroutine {
	return &ReconcileSubroutine{
		kubernetesClient: kubernetes,
		httpClient:       httpClient,
		conditionsSetter: conditions.NewSetter(ReconciledWithAlertManager),
	}
}

func (r *ReconcileSubroutine) Run(ctx context.Context, silenceResource *v1alpha1.Silence) error {
	url, err := r.kubernetesClient.GetAlertManagerURL(ctx)
	if err != nil {
		r.conditionsSetter.SetFalse(
			silenceResource.ObjectMeta,
			&silenceResource.Status.Conditions,
			"KubernetesClientError",
			"Unable to retrieve AlertManager resources from Cluster",
		)
		return fmt.Errorf("unable to retrieve any AlertManager instance in the current Kubernetes cluster: %w", err)
	}
	if url == "" {
		r.conditionsSetter.SetFalse(
			silenceResource.ObjectMeta,
			&silenceResource.Status.Conditions,
			"AlertManagerConfigurationError",
			"Unable to retrieve external URL from existing AlertManager",
		)
		return errors.New("there isn't any configured external URL for AlertManager")
	}

	restClient := silence.NewService(r.httpClient, url)

	dao := &silence.Dao{
		ID:        uuid.New().String(),
		Matchers:  silenceResource.Spec.Matchers,
		StartsAt:  time.Now(),
		EndsAt:    time.Now().Add(time.Hour * 24),
		CreatedBy: silenceResource.Spec.CreatedBy,
		Comment:   silenceResource.Spec.Comment,
	}
	err = restClient.UpdateSilence(ctx, dao)
	if err != nil {
		r.conditionsSetter.SetFalse(
			silenceResource.ObjectMeta,
			&silenceResource.Status.Conditions,
			"AlertManagerApiError",
			"Unable to update Silence on AlertManager instance",
		)
		return fmt.Errorf("unable to update Silence on AlertManager instance: %w", err)
	}

	silenceResource.Status = v1alpha1.SilenceStatus{
		ID:       dao.ID,
		StartsAt: dao.StartsAt.GoString(),
		EndsAt:   dao.EndsAt.GoString(),
	}
	r.conditionsSetter.SetTrue(
		silenceResource.ObjectMeta,
		&silenceResource.Status.Conditions,
		"SilenceReconciled",
		"Writer successfully reconcile against AlertManager instance",
	)

	return nil
}
