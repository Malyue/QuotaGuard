package server

import (
	quotaController "github.com/Malyue/quotaguard/pkg/controller/quota"
	clientset "github.com/Malyue/quotaguard/pkg/generated/clientset/versioned"
	informers "github.com/Malyue/quotaguard/pkg/generated/informers/externalversions"
	"github.com/Malyue/quotaguard/pkg/quota"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sync"
	"time"
)

type Server struct {
	QuotaManager *quota.QuotaManager
	Controller   *quotaController.Controller
	stopCh       chan struct{}
	wg           sync.WaitGroup
}

func NewServer(kubeconfigPath string) (*Server, error) {
	// init k8s client
	config, err := buildKubeConfig(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	crdClient := clientset.NewForConfigOrDie(config)

	informerFactory := informers.NewSharedInformerFactoryWithOptions(
		crdClient,
		10*time.Minute,
		informers.WithNamespace(""))

	qm := quota.NewQuotaManager()

	controller := quotaController.NewController(
		crdClient,
		informerFactory,
		qm,
		5,
	)

	controller.Start()

	svr := &Server{
		QuotaManager: qm,
		Controller:   controller,
	}

	return svr, nil
}

func buildKubeConfig(kubeconfigPath string) (*rest.Config, error) {
	if kubeconfigPath == "" {
		return rest.InClusterConfig()
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}
