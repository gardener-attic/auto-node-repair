/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http:  www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

This file was copied and modified from the kubernetes autoscaler project
https://github.com/kubernetes/autoscaler/blob/cluster-autorepair-1.0.0/cluster-autoscaler/core/autoscaler_builder.go
*/

package core

import (
	"github.com/gardener/auto-node-repair/utils/errors"
	kube_util "github.com/gardener/auto-node-repair/utils/kubernetes"
	kube_client "k8s.io/client-go/kubernetes"
	kube_record "k8s.io/client-go/tools/record"
)

// AutorepairBuilder builds an instance of Autorepair which is the core of CA
type AutorepairBuilder interface {
	//SetDynamicConfig(config dynamic.Config) AutorepairBuilder
	Build() (Autorepair, errors.AutorepairError)
}

// AutorepairBuilderImpl builds new autorepairs from its state including initial `AutoscalingOptions` given at startup and
// `dynamic.Config` read on demand from the configmap
type AutorepairBuilderImpl struct {
	autorepairingOptions AutorepairingOptions
	// dynamicConfig      *dynamic.Config TODO: Check if we really need dynamic repair based on configMap
	kubeClient         kube_client.Interface
	kubeEventRecorder  kube_record.EventRecorder
	listerRegistry     kube_util.ListerRegistry
}

// NewAutorepairBuilder builds an AutorepairBuilder from required parameters
func NewAutorepairBuilder(autorepairingOptions AutorepairingOptions,
	kubeClient kube_client.Interface, kubeEventRecorder kube_record.EventRecorder, listerRegistry kube_util.ListerRegistry) *AutorepairBuilderImpl {
	return &AutorepairBuilderImpl{
		autorepairingOptions: autorepairingOptions,
		kubeClient:         kubeClient,
		kubeEventRecorder:  kubeEventRecorder,
		listerRegistry:     listerRegistry,
	}
}

// Build an autorepair according to the builder's state
func (b *AutorepairBuilderImpl) Build() (Autorepair, errors.AutorepairError) {
	options := b.autorepairingOptions

	return NewStaticAutorepair(options, b.kubeClient, b.kubeEventRecorder, b.listerRegistry)
}
