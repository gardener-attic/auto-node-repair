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
https://github.com/kubernetes/autoscaler/blob/cluster-autorepair-1.0.0/cluster-autoscaler/core/dynamic_autoscaler.go
*/

package core

import (
	"time"
	"github.com/gardener/auto-node-repair/utils/errors"
)

// DynamicAutorepair is a variant of autorepair which supports dynamic reconfiguration at runtime
type DynamicAutorepair struct {
	autorepair        Autorepair
	autorepairBuilder AutorepairBuilder
}

// NewDynamicAutorepair builds a DynamicAutorepair from required parameters
func NewDynamicAutorepair(autorepairBuilder AutorepairBuilder) (*DynamicAutorepair, errors.AutorepairError) {
	autorepair, err := autorepairBuilder.Build()
	if err != nil {
		return nil, err
	}
	return &DynamicAutorepair{
		autorepair:        autorepair,
		autorepairBuilder: autorepairBuilder,
	}, nil
}

// Dynamic method used to invoke the static RunOnce method to repair nodes
func (a *DynamicAutorepair) RunOnce(repairtime time.Duration) (errors.AutorepairError) {
	return a.autorepair.RunOnce(repairtime)
}
