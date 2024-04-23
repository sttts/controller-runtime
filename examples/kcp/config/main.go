/*
Copyright 2024 The KCP Authors.

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

package main

import (
	kcpclienthelper "github.com/kcp-dev/apimachinery/v2/pkg/client"
	apisv1alpha1 "github.com/kcp-dev/kcp/sdk/apis/apis/v1alpha1"
	"github.com/kcp-dev/kcp/sdk/apis/core"
	corev1alpha1 "github.com/kcp-dev/kcp/sdk/apis/core/v1alpha1"
	tenancyv1alpha1 "github.com/kcp-dev/kcp/sdk/apis/tenancy/v1alpha1"
	"github.com/kcp-dev/logicalcluster/v3"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/kcp-dev/controller-runtime/examples/kcp/config/consumers"
	"github.com/kcp-dev/controller-runtime/examples/kcp/config/widgets"
	"github.com/kcp-dev/controller-runtime/examples/kcp/config/widgets/resources"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// config is bootstrap set of assets for the controller-runtime examples.
// It includes the following assets:
// - crds/* for the Widget type - autogenerated from the Widget type definition
// - widgets/resources/* - a set of Widget resources for KCP to manage. Automatically generated by kcp apigen
// see Makefile & hack/update-codegen-crds.sh for more details

// It is intended to be running with higher privileges than the examples themselves
// to ensure system (kcp) is bootstrapped. In real world scenarios, this would be
// done by the platform operator to enable service providers to deploy their
// controllers.

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(tenancyv1alpha1.AddToScheme(scheme))
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(corev1alpha1.AddToScheme(scheme))
	utilruntime.Must(apisv1alpha1.AddToScheme(scheme))

}

var (
	// clusterName is the workspace to host common APIs.
	clusterName = logicalcluster.NewPath("root:widgets")
)

func main() {
	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	ctx := ctrl.SetupSignalHandler()
	log := log.FromContext(ctx)

	restConfig, err := config.GetConfigWithContext("base")
	if err != nil {
		log.Error(err, "unable to get config")
	}

	restCopy := rest.CopyConfig(restConfig)
	restRoot := rest.AddUserAgent(kcpclienthelper.SetCluster(restCopy, core.RootCluster.Path()), "bootstrap-root")
	clientRoot, err := client.New(restRoot, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		log.Error(err, "unable to create client")
	}

	restCopy = rest.CopyConfig(restConfig)
	restWidgets := rest.AddUserAgent(kcpclienthelper.SetCluster(restCopy, clusterName), "bootstrap-widgets")
	clientWidgets, err := client.New(restWidgets, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		log.Error(err, "unable to create client")
	}

	err = widgets.Bootstrap(ctx, clientRoot)
	if err != nil {
		log.Error(err, "failed to bootstrap widgets")
	}

	err = resources.Bootstrap(ctx, clientWidgets)
	if err != nil {
		log.Error(err, "failed to bootstrap resources")
	}

	err = consumers.Bootstrap(ctx, clientRoot)
	if err != nil {
		log.Error(err, "failed to bootstrap consumers")
	}

}
