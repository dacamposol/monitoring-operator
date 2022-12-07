package test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	monitoringCoreOsV1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/suite"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	monitoringv1alpha1 "github.com/dacamposol/monitoring-operator/api/v1alpha1"
	"github.com/dacamposol/monitoring-operator/controllers"
	"github.com/dacamposol/monitoring-operator/subroutines/silence"
)

func TestSilenceControllerSuite(t *testing.T) {
	suite.Run(t, new(SilenceControllerSuite))
}

func (suite *SilenceControllerSuite) TestFinalizer_NoAlertManagers() {
	testContext := context.Background()
	testSilenceUUID := uuid.New().String()

	suite.createSilence(testSilenceUUID)

	createdSilence := monitoringv1alpha1.Silence{}
	suite.Assert().Eventually(func() bool {
		err := suite.kubernetesClient.Get(testContext, types.NamespacedName{
			Name:      testSilenceUUID,
			Namespace: "application-system",
		}, &createdSilence)

		if err == nil {
			return controllerutil.ContainsFinalizer(&createdSilence, controllers.AlertManagerDependencyFinalizer)
		}

		return false
	}, time.Second*10, time.Millisecond*10)

	suite.deleteSilence(testSilenceUUID)
}

func (suite *SilenceControllerSuite) TestReconciliation_AlertManager() {
	testContext := context.Background()
	testSilenceUUID := uuid.New().String()
	alertManagerUUID := uuid.New().String()

	alertManagerServer := suite.createFakeAlertManagerServer()

	alertmanager := &monitoringCoreOsV1.Alertmanager{
		ObjectMeta: metav1.ObjectMeta{
			Name:      alertManagerUUID,
			Namespace: "monitoring",
		},
		Spec: monitoringCoreOsV1.AlertmanagerSpec{
			ExternalURL: alertManagerServer.URL,
		},
	}
	err := suite.kubernetesClient.Create(testContext, alertmanager)
	suite.Nil(err)

	suite.Require().Eventually(func() bool {
		err := suite.kubernetesClient.Get(testContext, types.NamespacedName{
			Name:      alertManagerUUID,
			Namespace: "monitoring",
		}, alertmanager)

		return err == nil
	}, time.Second*10, time.Millisecond*10)

	suite.createSilence(testSilenceUUID)

	createdSilence := monitoringv1alpha1.Silence{}
	suite.Require().Eventually(func() bool {
		err := suite.kubernetesClient.Get(testContext, types.NamespacedName{
			Name:      testSilenceUUID,
			Namespace: "application-system",
		}, &createdSilence)

		return err == nil &&
			meta.IsStatusConditionTrue(createdSilence.Status.Conditions, silence.ReconciledWithAlertManager) &&
			createdSilence.Status.ID == "EXAMPLE_UID"
	}, time.Second*10, time.Millisecond*10)

	suite.deleteSilence(testSilenceUUID)
}
