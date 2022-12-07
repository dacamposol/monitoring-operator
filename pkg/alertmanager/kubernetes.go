package alertmanager

import (
	"context"
	"fmt"

	monitoringCoreOsV1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type service struct {
	client.Client
}

func NewService(kubernetesClient client.Client) Getter {
	return &service{
		kubernetesClient,
	}
}

func (s *service) GetAlertManagerURL(ctx context.Context) (string, error) {
	alertManagerList := monitoringCoreOsV1.AlertmanagerList{}
	err := s.List(ctx, &alertManagerList)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve AlertManager resources in Kubernetes Cluster: %w", err)
	}
	if len(alertManagerList.Items) == 0 {
		return "", errors.NewNotFound(monitoringCoreOsV1.Resource("alertmanagers"), "AlertManager Instances")
	}

	return alertManagerList.Items[0].Spec.ExternalURL, nil
}
