package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"time"

	monitoringCoreOsV1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	monitoringv1alpha1 "github.com/dacamposol/monitoring-operator/api/v1alpha1"
	"github.com/dacamposol/monitoring-operator/controllers"
	"github.com/dacamposol/monitoring-operator/pkg/silence"
)

type SilenceControllerSuite struct {
	suite.Suite

	kubernetesClient      client.Client
	kubernetesManager     ctrl.Manager
	kubernetesEnvironment *envtest.Environment

	cancel context.CancelFunc
}

func (suite *SilenceControllerSuite) SetupSuite() {
	suite.kubernetesEnvironment = &envtest.Environment{
		ErrorIfCRDPathMissing: true,
		CRDDirectoryPaths: []string{
			filepath.Join("crds"),
			filepath.Join("..", "chart", "crds"),
		},
	}

	cfg, err := suite.kubernetesEnvironment.Start()
	suite.Nil(err)

	utilruntime.Must(monitoringv1alpha1.AddToScheme(scheme.Scheme))
	utilruntime.Must(monitoringCoreOsV1.AddToScheme(scheme.Scheme))

	//+kubebuilder:scaffold:scheme

	suite.kubernetesClient, err = client.New(cfg, client.Options{
		Scheme: scheme.Scheme,
	})
	suite.Nil(err)

	suite.kubernetesManager, err = ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})

	err = (&controllers.SilenceReconciler{
		Client: suite.kubernetesManager.GetClient(),
		Scheme: scheme.Scheme,
	}).SetupWithManager(suite.kubernetesManager)

	go suite.startController()

	suite.initializeUniverse()
}

func (suite *SilenceControllerSuite) TearDownSuite() {
	suite.cancel()
	err := suite.kubernetesEnvironment.Stop()
	suite.Nil(err)
}

func (suite *SilenceControllerSuite) startController() {
	var controllerContext context.Context
	controllerContext, suite.cancel = context.WithCancel(context.Background())
	err := suite.kubernetesManager.Start(controllerContext)
	suite.Nil(err)
}

func (suite *SilenceControllerSuite) initializeUniverse() {
	initContext := context.Background()

	monitoringNamespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "monitoring",
		},
	}
	applicationNamespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "application-system",
		},
	}

	err := suite.kubernetesClient.Create(initContext, monitoringNamespace)
	suite.Nil(err)

	err = suite.kubernetesClient.Create(initContext, applicationNamespace)
	suite.Nil(err)

	suite.Require().Eventually(func() bool {
		testContext := context.Background()

		monitoringErr := suite.kubernetesClient.Get(testContext, types.NamespacedName{Name: "monitoring"}, monitoringNamespace)
		applicationErr := suite.kubernetesClient.Get(testContext, types.NamespacedName{Name: "application-system"}, applicationNamespace)

		return monitoringErr == nil && applicationErr == nil
	}, time.Second*10, time.Millisecond*10)
}

func (suite *SilenceControllerSuite) createSilence(testName string) {
	testContext := context.Background()

	resource := &monitoringv1alpha1.Silence{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName,
			Namespace: "application-system",
		},
		Spec: monitoringv1alpha1.SilenceSpec{
			Matchers: []monitoringv1alpha1.Matcher{
				{
					Name:    "alert-equal",
					Value:   "",
					IsRegex: false,
					IsEqual: true,
				},
			},
			CreatedBy: "Test User",
			Comment:   "This is a test for Silence resources",
		},
	}

	err := suite.kubernetesClient.Create(testContext, resource)
	suite.Nil(err)
}

func (suite *SilenceControllerSuite) deleteSilence(testName string) {
	testContext := context.Background()

	createdSilence := monitoringv1alpha1.Silence{}
	suite.Require().Eventually(func() bool {
		err := suite.kubernetesClient.Get(testContext, types.NamespacedName{
			Name:      testName,
			Namespace: "application-system",
		}, &createdSilence)

		return err == nil
	}, time.Second*10, time.Millisecond*10)

	err := suite.kubernetesClient.Delete(testContext, &createdSilence)
	suite.Nil(err)
}

func (suite *SilenceControllerSuite) createFakeAlertManagerServer() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.URL.Path == "/silences" {
				suite.Equal("POST", request.Method)

				writer.WriteHeader(200)
				body := &silence.Dao{
					ID:        "EXAMPLE_UID",
					Matchers:  nil,
					StartsAt:  time.Now(),
					EndsAt:    time.Now().Add(24 * time.Hour),
					CreatedBy: "Test User",
					Comment:   "This is a test for Silence resources",
				}
				bodyBytes, err := json.Marshal(body)
				suite.Nil(err)

				_, err = writer.Write(bodyBytes)
				suite.Nil(err)
				return
			}
			if request.Method == "DELETE" {
				suite.Equal("/silences/EXAMPLE_UID", request.URL.Path)
				return
			}

		}))
}
