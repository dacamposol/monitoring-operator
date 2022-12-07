package silence

import (
	"context"
	"fmt"
	"net/http"

	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/dacamposol/monitoring-operator/api/v1alpha1"
	"github.com/dacamposol/monitoring-operator/pkg/alertmanager"
	"github.com/dacamposol/monitoring-operator/pkg/silence"
	"github.com/dacamposol/monitoring-operator/subroutines/common/conditions"
)

const ReadyForDeletion = "ReadyForDeletion"

type DeleteSubroutine struct {
	kubernetesClient alertmanager.Getter
	httpClient       silence.HttpClient
	conditionsSetter conditions.Setter
}

func NewDeletionSubroutine(kubernetes alertmanager.Getter, httpClient *http.Client) *DeleteSubroutine {
	return &DeleteSubroutine{
		kubernetesClient: kubernetes,
		httpClient:       httpClient,
		conditionsSetter: conditions.NewSetter(ReadyForDeletion),
	}
}

func (r *DeleteSubroutine) Run(ctx context.Context, silenceResource *v1alpha1.Silence) error {
	url, err := r.kubernetesClient.GetAlertManagerURL(ctx)
	if err != nil && !errors.IsNotFound(err) {
		r.conditionsSetter.SetFalse(
			silenceResource.ObjectMeta,
			&silenceResource.Status.Conditions,
			"KubernetesClientError",
			"Unable to retrieve AlertManager resources from Cluster",
		)
		return fmt.Errorf("unable to retrieve AlertManager instances from the current Kubernetes cluster: %w", err)
	}
	if url == "" {
		return nil
	}

	restClient := silence.NewService(r.httpClient, url)
	err = restClient.DeleteSilence(ctx, silenceResource.Status.ID)
	if err != nil {
		r.conditionsSetter.SetFalse(
			silenceResource.ObjectMeta,
			&silenceResource.Status.Conditions,
			"AlertManagerApiError",
			"Unable to delete Silence on AlertManager instance",
		)
		return fmt.Errorf("unable to delete Silence on AlertManager instance: %w", err)
	}

	return nil
}
