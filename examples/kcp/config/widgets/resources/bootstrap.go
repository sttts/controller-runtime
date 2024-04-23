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

package resources

import (
	"context"
	"embed"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	confighelpers "github.com/kcp-dev/controller-runtime/examples/kcp/config/helpers"
)

//go:embed *.yaml
var fs embed.FS

// Bootstrap creates resources in this package by continuously retrying the list.
// This is blocking, i.e. it only returns (with error) when the context is closed or with nil when
// the bootstrapping is successfully completed.
func Bootstrap(
	ctx context.Context,
	client client.Client,
) error {
	log := log.FromContext(ctx)

	log.Info("Bootstrapping widgets resources")
	if err := confighelpers.Bootstrap(ctx, client, fs); err != nil {
		return err
	}

	return nil
}
