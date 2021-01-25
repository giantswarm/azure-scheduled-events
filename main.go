package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/k8sclient/v4/pkg/k8srestconfig"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"k8s.io/client-go/rest"

	"github.com/giantswarm/azure-scheduled-events/pkg/azuremetadataclient"
	"github.com/giantswarm/azure-scheduled-events/pkg/eventhandler"
	"github.com/giantswarm/azure-scheduled-events/pkg/eventhandler/drainer"
)

var (
	k8sAddress       string
	cafile           string
	crtfile          string
	keyfile          string
	kubeconfigPath   string
	inCluster        bool
	mainLoopInterval int
)

const k8sNodeNameEnvVarName = "K8S_NODE_NAME"

func main() {
	flag.StringVar(&k8sAddress, "k8saddress", "", "k8s address.")
	flag.StringVar(&cafile, "cafile", "", "TLS ca file.")
	flag.StringVar(&crtfile, "crtfile", "", "TLS crt file.")
	flag.StringVar(&keyfile, "keyfile", "", "TLS key file.")
	flag.StringVar(&kubeconfigPath, "kubeconfigpath", "", "kubeconfig path.")
	flag.BoolVar(&inCluster, "incluster", true, "whether it runs in k8s cluster or not.")
	flag.IntVar(&mainLoopInterval, "check-interval", 5, "The interval in seconds between two checks of the events")

	flag.Parse()

	ctx := context.Background()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		done <- true
	}()

	logger, err := micrologger.New(micrologger.Config{})
	if err != nil {
		log.Fatal(err)
	}

	k8sclients, err := getK8sClient(logger, k8sAddress, cafile, crtfile, keyfile, kubeconfigPath, inCluster)
	if err != nil {
		log.Fatal(err)
	}

	azureMetadata, err := azuremetadataclient.New(azuremetadataclient.Config{})
	if err != nil {
		log.Fatal(err)
	}

	localNodeName := os.Getenv(k8sNodeNameEnvVarName)
	if localNodeName == "" {
		log.Fatalf("Missing required environment variable %s", k8sNodeNameEnvVarName)
	}

	logger.LogCtx(ctx, "level", "info", "message", fmt.Sprintf("Local Kubernetes Node Name is %q", localNodeName))

	eventHandlers := []eventhandler.EventHandler{
		drainer.NewDrainEventHandler(logger, azureMetadata, k8sclients.K8sClient(), localNodeName),
	}

	interval := 5
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	go func() {
		logger.LogCtx(ctx, "level", "info", "message", fmt.Sprintf("Starting main loop with %d seconds interval", interval))
		for range ticker.C {
			events, err := azureMetadata.FetchEvents()
			if err != nil {
				logger.Errorf(ctx, err, "Error fetching events from azure metadata service")
				break
			}

			for _, event := range events {
				for _, handler := range eventHandlers {
					err = handler.HandleEvent(ctx, event)
					if err != nil {
						logger.Errorf(ctx, err, "Error handling event")
					}
				}
			}
		}
	}()

	<-done
	ticker.Stop()
	logger.LogCtx(ctx, "message", "Exiting")
}

func getK8sClient(logger micrologger.Logger, k8sAddress, cafile, crtfile, keyfile, kubeconfigPath string, incluster bool) (*k8sclient.Clients, error) {
	var err error
	var k8sClient *k8sclient.Clients
	{
		defined := 0
		if k8sAddress != "" {
			defined++
		}
		if incluster {
			defined++
		}
		if kubeconfigPath != "" {
			defined++
		}

		if defined == 0 {
			return nil, microerror.Maskf(invalidConfigError, "address or inCluster or kubeConfigPath must be defined")
		}
		if defined > 1 {
			return nil, microerror.Maskf(invalidConfigError, "address and inCluster and kubeConfigPath must not be defined at the same time")
		}

		var restConfig *rest.Config
		if kubeconfigPath == "" {
			restConfig, err = buildK8sRestConfig(logger, k8sAddress, cafile, crtfile, keyfile, kubeconfigPath, incluster)
			if err != nil {
				return nil, microerror.Mask(err)
			}
		}

		c := k8sclient.ClientsConfig{
			Logger:         logger,
			SchemeBuilder:  k8sclient.SchemeBuilder{},
			KubeConfigPath: kubeconfigPath,
			RestConfig:     restConfig,
		}

		k8sClient, err = k8sclient.NewClients(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return k8sClient, nil
}

func buildK8sRestConfig(logger micrologger.Logger, k8sAddress, cafile, crtfile, keyfile, kubeconfigPath string, incluster bool) (*rest.Config, error) {
	c := k8srestconfig.Config{
		Logger: logger,

		Address:    k8sAddress,
		InCluster:  incluster,
		KubeConfig: kubeconfigPath,
		TLS: k8srestconfig.ConfigTLS{
			CAFile:  cafile,
			CrtFile: crtfile,
			KeyFile: keyfile,
		},
	}

	restConfig, err := k8srestconfig.New(c)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return restConfig, nil
}

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}
