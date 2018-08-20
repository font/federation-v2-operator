package main

import (
	"context"
	"runtime"

	"github.com/marun/federation-v2-operator/pkg/stub"
	sdk "github.com/operator-framework/operator-sdk/pkg/sdk"
	k8sutil "github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
	sdkVersion "github.com/operator-framework/operator-sdk/version"

	"github.com/sirupsen/logrus"
)

func printVersion() {
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

func main() {
	printVersion()

	sdk.ExposeMetricsPort()

	resource := "operator.federation.k8s.io/v1alpha1"
	kinds := []string{"ClusterRegistry", "FederationV2"}
	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		logrus.Fatalf("Failed to get watch namespace: %v", err)
	}
	resyncPeriod := 0 // Disable periodic events i.e. only send events when updated.
	for _, kind := range kinds {
		logrus.Infof("Watching %s, %s, %s, %d", resource, kind, namespace, resyncPeriod)
		sdk.Watch(resource, kind, namespace, resyncPeriod)
	}
	sdk.Handle(stub.NewHandler())
	sdk.Run(context.TODO())
}
