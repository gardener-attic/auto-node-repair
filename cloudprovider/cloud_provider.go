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
https://github.com/kubernetes/autoscaler/blob/cluster-autorepair-1.0.0/cluster-autoscaler/cloudprovider/cloud_provider.go
*/

package cloudprovider

import (
	"bytes"
	"fmt"
	"math"
	"time"

	"github.com/gardener/auto-node-repair/utils/errors"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/kubernetes/plugin/pkg/scheduler/schedulercache"
)

// CloudProvider contains configuration info and functions for interacting with
// cloud provider (GCE, AWS, etc).
type CloudProvider interface {
	// Name returns name of the cloud provider.
	Name() string

	// NodeGroups returns all node groups configured for this cloud provider.
	NodeGroups() []NodeGroup

	// NodeGroupForNode returns the node group for the given node, nil if the node
	// should not be processed by Cluster autorepair, or non-nil error if such
	// occurred. Must be implemented.
	NodeGroupForNode(*apiv1.Node) (NodeGroup, error)

	// Pricing returns pricing model for this cloud provider or error if not available.
	// Implementation optional.
	Pricing() (PricingModel, errors.AutorepairError)

	// GetAvailableMachineTypes get all machine types that can be requested from the cloud provider.
	// Implementation optional.
	GetAvailableMachineTypes() ([]string, error)

	// NewNodeGroup builds a theoretical node group based on the node definition provided. The node group is not automatically
	// created on the cloud provider side. The node group is not returned by NodeGroups() until it is created.
	// Implementation optional.
	NewNodeGroup(machineType string, labels map[string]string, extraResources map[string]resource.Quantity) (NodeGroup, error)
}

// ErrNotImplemented is returned if a method is not implemented.
var ErrNotImplemented errors.AutorepairError = errors.NewAutorepairError(errors.InternalError, "Not implemented")

// ErrAlreadyExist is returned if a method is not implemented.
var ErrAlreadyExist errors.AutorepairError = errors.NewAutorepairError(errors.InternalError, "Already exist")

// NodeGroup contains configuration info and functions to control a set
// of nodes that have the same capacity and set of labels.
type NodeGroup interface {
	// MaxSize returns maximum size of the node group.
	MaxSize() int

	// MinSize returns minimum size of the node group.
	MinSize() int

	// TargetSize returns the current target size of the node group. It is possible that the
	// number of nodes in Kubernetes is different at the moment but should be equal
	// to Size() once everything stabilizes (new nodes finish startup and registration or
	// removed nodes are deleted completely). Implementation required.
	TargetSize() (int, error)

	// IncreaseSize increases the size of the node group. To delete a node you need
	// to explicitly name it and use DeleteNode. This function should wait until
	// node group size is updated. Implementation required.
	IncreaseSize(delta int) error

	// DeleteNodes deletes nodes from this node group. Error is returned either on
	// failure or if the given node doesn't belong to this node group. This function
	// should wait until node group size is updated. Implementation required.
	DeleteNodes([]*apiv1.Node) error

	// DecreaseTargetSize decreases the target size of the node group. This function
	// doesn't permit to delete any existing node and can be used only to reduce the
	// request for new nodes that have not been yet fulfilled. Delta should be negative.
	// It is assumed that cloud provider will not delete the existing nodes when there
	// is an option to just decrease the target. Implementation required.
	DecreaseTargetSize(delta int) error

	// Id returns an unique identifier of the node group.
	Id() string

	// Debug returns a string containing all information regarding this node group.
	Debug() string

	// Nodes returns a list of all nodes that belong to this node group.
	Nodes() ([]string, error)

	// TemplateNodeInfo returns a schedulercache.NodeInfo structure of an empty
	// (as if just started) node. This will be used in scale-up simulations to
	// predict what would a new node look like if a node group was expanded. The returned
	// NodeInfo is expected to have a fully populated Node object, with all of the labels,
	// capacity and allocatable information as well as all pods that are started on
	// the node by default, using manifest (most likely only kube-proxy). Implementation optional.
	TemplateNodeInfo() (*schedulercache.NodeInfo, error)

	// Exist checks if the node group really exists on the cloud provider side. Allows to tell the
	// theoretical node group from the real one. Implementation required.
	Exist() bool

	// Create creates the node group on the cloud provider side. Implementation optional.
	Create() error

	// Delete deletes the node group on the cloud provider side.
	// This will be executed only for autoprovisioned node groups, once their size drops to 0.
	// Implementation optional.
	Delete() error

	// Autoprovisioned returns true if the node group is autoprovisioned. An autoprovisioned group
	// was created by CA and can be deleted when scaled to 0.
	Autoprovisioned() bool
}

// PricingModel contains information about the node price and how it changes in time.
type PricingModel interface {
	// NodePrice returns a price of running the given node for a given period of time.
	// All prices returned by the structure should be in the same currency.
	NodePrice(node *apiv1.Node, startTime time.Time, endTime time.Time) (float64, error)

	// PodePrice returns a theoretical minimum priece of running a pod for a given
	// period of time on a perfectly matching machine.
	PodPrice(pod *apiv1.Pod, startTime time.Time, endTime time.Time) (float64, error)
}

const (
	// ResourceNameCores is string name for cores. It's used by ResourceLimiter.
	ResourceNameCores = "cpu"
	// ResourceNameMemory is string name for memory. It's used by ResourceLimiter.
	// Memory should always be provided in megabytes.
	ResourceNameMemory = "memory"
)

// ResourceLimiter contains limits (max, min) for resources (cores, memory etc.).
type ResourceLimiter struct {
	minLimits map[string]int64
	maxLimits map[string]int64
}

// NewResourceLimiter creates new ResourceLimiter for map. Maps are deep copied.
func NewResourceLimiter(minLimits map[string]int64, maxLimits map[string]int64) *ResourceLimiter {
	minLimitsCopy := make(map[string]int64)
	maxLimitsCopy := make(map[string]int64)
	for key, value := range minLimits {
		minLimitsCopy[key] = value
	}
	for key, value := range maxLimits {
		maxLimitsCopy[key] = value
	}
	return &ResourceLimiter{minLimitsCopy, maxLimitsCopy}
}

// GetMin returns minimal number of resources for a given resource type.
func (r *ResourceLimiter) GetMin(resourceName string) int64 {
	result, found := r.minLimits[resourceName]
	if found {
		return result
	}
	return 0
}

// GetMax returns maximal number of resources for a given resource type.
func (r *ResourceLimiter) GetMax(resourceName string) int64 {
	result, found := r.maxLimits[resourceName]
	if found {
		return result
	}
	return math.MaxInt64
}

func (r *ResourceLimiter) String() string {
	var buffer bytes.Buffer
	for name, maxLimit := range r.maxLimits {
		if buffer.Len() > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(fmt.Sprintf("{%s : %d - %d}", name, r.minLimits[name], maxLimit))
	}
	return buffer.String()
}
