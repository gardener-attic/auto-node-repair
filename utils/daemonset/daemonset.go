/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

This file was copied from the kubernetes autoscaler project
https://github.com/kubernetes/autoscaler/blob/cluster-autorepair-1.0.0/cluster-autoscaler/util/daemonset/daemonset.go
*/

package daemonset

import (
	"fmt"
	"math/rand"

	"github.com/gardener/auto-node-repair/simulator"

	apiv1 "k8s.io/api/core/v1"
	extensionsv1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/kubernetes/plugin/pkg/scheduler/schedulercache"
)

// GetDaemonSetPodsForNode returns daemonset nodes for the given pod.
func GetDaemonSetPodsForNode(nodeInfo *schedulercache.NodeInfo, daemonsets []*extensionsv1.DaemonSet, predicateChecker *simulator.PredicateChecker) []*apiv1.Pod {
	result := make([]*apiv1.Pod, 0)
	for _, ds := range daemonsets {
		pod := newPod(ds, nodeInfo.Node().Name)
		if err := predicateChecker.CheckPredicates(pod, nil, nodeInfo, simulator.ReturnSimpleError); err == nil {
			result = append(result, pod)
		}
	}
	return result
}

func newPod(ds *extensionsv1.DaemonSet, nodeName string) *apiv1.Pod {
	newPod := &apiv1.Pod{Spec: ds.Spec.Template.Spec, ObjectMeta: ds.Spec.Template.ObjectMeta}
	newPod.Namespace = ds.Namespace
	newPod.Name = fmt.Sprintf("%s-pod-%d", ds.Name, rand.Int63())
	newPod.Spec.NodeName = nodeName
	return newPod
}
