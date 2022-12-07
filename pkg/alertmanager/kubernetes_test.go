package alertmanager

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"sigs.k8s.io/controller-runtime/pkg/client"

	monitoringCoreOsV1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"

	"github.com/dacamposol/monitoring-operator/pkg/alertmanager/mocks"
)

type ServiceTestSuite struct {
	suite.Suite

	testObj *service

	clientMock *mocks.Client
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.clientMock = new(mocks.Client)
	suite.testObj = &service{
		suite.clientMock,
	}
}

func (suite *ServiceTestSuite) TestGetAlertManagerURL_NoResources_NOK() {
	testCtx := context.Background()

	suite.clientMock.EXPECT().
		List(testCtx, mock.Anything).
		Run(func(ctx context.Context, list client.ObjectList, opts ...client.ListOption) {
			_, ok := list.(*monitoringCoreOsV1.AlertmanagerList)
			suite.True(ok)
		}).
		Return(nil)
	url, err := suite.testObj.GetAlertManagerURL(context.Background())
	suite.EqualError(err, "alertmanagers.monitoring.coreos.com \"AlertManager Instances\" not found")
	suite.Equal("", url)
}

func (suite *ServiceTestSuite) TestGetAlertManagerURL_KubernetesError_NOK() {
	testCtx := context.Background()

	suite.clientMock.EXPECT().
		List(testCtx, mock.Anything).
		Run(func(ctx context.Context, list client.ObjectList, opts ...client.ListOption) {
			_, ok := list.(*monitoringCoreOsV1.AlertmanagerList)
			suite.True(ok)
		}).
		Return(errors.New("list error"))
	url, err := suite.testObj.GetAlertManagerURL(context.Background())
	suite.EqualError(err, "unable to retrieve AlertManager resources in Kubernetes Cluster: list error")
	suite.Equal("", url)
}

func (suite *ServiceTestSuite) TestGetAlertManagerURL_AvailableResources_OK() {
	testCtx := context.Background()

	suite.clientMock.EXPECT().
		List(testCtx, mock.Anything).
		Run(func(ctx context.Context, list client.ObjectList, opts ...client.ListOption) {
			alertList, ok := list.(*monitoringCoreOsV1.AlertmanagerList)
			suite.True(ok)

			alertList.Items = []monitoringCoreOsV1.Alertmanager{
				{
					Spec: monitoringCoreOsV1.AlertmanagerSpec{
						ExternalURL: "https://alertmanager.local",
					},
				},
			}
		}).
		Return(nil)
	url, err := suite.testObj.GetAlertManagerURL(context.Background())
	suite.Nil(err)
	suite.Equal("https://alertmanager.local", url)
}
