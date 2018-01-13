/*
Copyright 2016 The Kubernetes Authors.

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
https://github.com/kubernetes/autoscaler/blob/cluster-autorepair-1.0.0/cluster-autoscaler/core/autoscaler.go
*/

package core

import (
	"time"
	
	"github.com/gardener/auto-node-repair/utils/errors"
	kube_util "github.com/gardener/auto-node-repair/utils/kubernetes"
	kube_client "k8s.io/client-go/kubernetes"
	kube_record "k8s.io/client-go/tools/record"
)

// AutorepairOptions is the whole set of options for configuring an autorepair
type AutorepairOptions struct {
	AutorepairingOptions
}

// Autorepair is the main component of CA which scales up/down node groups according to its configuration
// The configuration can be injected at the creation of an autorepair
type Autorepair interface {
	// RunOnce represents an iteration in the control-loop of CA
	RunOnce(time.Duration) (errors.AutorepairError) 
}

// NewAutorepair creates an autorepair of an appropriate type according to the parameters
func NewAutorepair(opts AutorepairOptions, kubeClient kube_client.Interface,
	kubeEventRecorder kube_record.EventRecorder, listerRegistry kube_util.ListerRegistry) (Autorepair, errors.AutorepairError) {

	autorepairBuilder := NewAutorepairBuilder(opts.AutorepairingOptions, kubeClient, kubeEventRecorder, listerRegistry )
	return NewDynamicAutorepair(autorepairBuilder)
}