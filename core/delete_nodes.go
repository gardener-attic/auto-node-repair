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
https://github.com/kubernetes/autoscaler/blob/cluster-autorepair-1.0.0/cluster-autoscaler/core/scale_down.go
*/

package core

import (
	"fmt"
	"reflect"
	"time"

	"github.com/gardener/auto-node-repair/cloudprovider"
	"github.com/gardener/auto-node-repair/clusterstate"
	"github.com/gardener/auto-node-repair/utils/deletetaint"
	"github.com/gardener/auto-node-repair/utils/errors"
	
	apiv1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1beta1"
	kube_errors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kube_client "k8s.io/client-go/kubernetes"
	kube_record "k8s.io/client-go/tools/record"
	
	"github.com/golang/glog"
)

const (
	// MaxKubernetesEmptyNodeDeletionTime is the maximum time needed by Kubernetes to delete an empty node.
	MaxKubernetesEmptyNodeDeletionTime = 3 * time.Minute
	// MaxCloudProviderNodeDeletionTime is the maximum time needed by cloud provider to delete a node.
	MaxCloudProviderNodeDeletionTime = 5 * time.Minute
	// MaxPodEvictionTime is the maximum time CA tries to evict a pod before giving up.
	MaxPodEvictionTime = 2 * time.Minute
	// EvictionRetryTime is the time after CA retries failed pod eviction.
	EvictionRetryTime = 10 * time.Second
	// PodEvictionHeadroom is the extra time we wait to catch situations when the pod is ignoring SIGTERM and
	// is killed with SIGKILL after MaxGracefulTerminationTime
	PodEvictionHeadroom = 30 * time.Second
	// UnremovableNodeRecheckTimeout is the timeout before we check again a node that couldn't be removed before
	UnremovableNodeRecheckTimeout = 5 * time.Minute
)

// Delete a node by first cordoning, draining and then deleting from relavent ASG
func deleteNode(context *AutorepairingContext, node *apiv1.Node, pods []*apiv1.Pod) errors.AutorepairError {
	deleteSuccessful := false
	drainSuccessful := false

	if err := deletetaint.MarkToBeDeleted(node, context.ClientSet); err != nil {
		context.Recorder.Eventf(node, apiv1.EventTypeWarning, "AutoRepairFailed", "failed to mark the node as toBeDeleted/unschedulable: %v", err)
		return errors.ToAutorepairError(errors.ApiCallError, err)
	}

	// If we fail to evict all the pods from the node we want to remove delete taint
	defer func() {
		if !deleteSuccessful {
			deletetaint.CleanToBeDeleted(node, context.ClientSet)
			if !drainSuccessful {
				context.Recorder.Eventf(node, apiv1.EventTypeWarning, "AutoRepairFailed", "failed to drain the node, aborting ScaleDown")
			} else {
				context.Recorder.Eventf(node, apiv1.EventTypeWarning, "AutoRepairFailed", "failed to delete the node")
			}
		}
	}()

	context.Recorder.Eventf(node, apiv1.EventTypeNormal, "AutoRepairFailed", "marked the node as toBeDeleted/unschedulable")

	// attempt drain
	if err := drainNode(node, pods, context.ClientSet, context.Recorder, context.MaxGracefulTerminationSec, MaxPodEvictionTime, EvictionRetryTime); err != nil {
		return err
	}
	drainSuccessful = true
	
	// attempt delete from cloud provider
	
	err := deleteNodeFromCloudProvider(node, context.CloudProvider, context.Recorder, context.ClusterStateRegistry)
	if err != nil {
		return err
	}
	
	deleteSuccessful = true // Let the deferred function know there is no need to cleanup
	return nil
}

// Evict pods
func evictPod(podToEvict *apiv1.Pod, client kube_client.Interface, recorder kube_record.EventRecorder,
	maxGracefulTerminationSec int, retryUntil time.Time, waitBetweenRetries time.Duration) error {
	recorder.Eventf(podToEvict, apiv1.EventTypeNormal, "AutoRepairFailed", "deleting pod for node scale down")
	maxTermination := int64(0)

	var lastError error
	for first := true; first || time.Now().Before(retryUntil); time.Sleep(waitBetweenRetries) {
		first = false
		eviction := &policyv1.Eviction{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: podToEvict.Namespace,
				Name:      podToEvict.Name,
			},
			DeleteOptions: &metav1.DeleteOptions{
				GracePeriodSeconds: &maxTermination,
			},
		}
		lastError = client.CoreV1().Pods(podToEvict.Namespace).Evict(eviction)
		if lastError == nil || kube_errors.IsNotFound(lastError) {
			return nil
		}
	}
	glog.Errorf("Failed to evict pod %s, error: %v", podToEvict.Name, lastError)
	recorder.Eventf(podToEvict, apiv1.EventTypeWarning, "AutoRepairFailed", "failed to delete pod for AutoRepairFailed")
	return fmt.Errorf("Failed to evict pod %s/%s within allowed timeout (last error: %v)", podToEvict.Namespace, podToEvict.Name, lastError)
}

// Performs drain logic on the node. Marks the node as unschedulable and later removes all pods, giving
// them up to MaxGracefulTerminationTime to finish.
func drainNode(node *apiv1.Node, pods []*apiv1.Pod, client kube_client.Interface, recorder kube_record.EventRecorder,
	maxGracefulTerminationSec int, maxPodEvictionTime time.Duration, waitBetweenRetries time.Duration) errors.AutorepairError {

	toEvict := len(pods)
	retryUntil := time.Now().Add(maxPodEvictionTime)
	confirmations := make(chan error, toEvict)
	for _, pod := range pods {
		go func(podToEvict *apiv1.Pod) {
			confirmations <- evictPod(podToEvict, client, recorder, maxGracefulTerminationSec, retryUntil, waitBetweenRetries)
		}(pod)
	}

	evictionErrs := make([]error, 0)

	for range pods {
		select {
		case err := <-confirmations:
			if err != nil {
				evictionErrs = append(evictionErrs, err)
			} 
		case <-time.After(retryUntil.Sub(time.Now()) + 5*time.Second):
			return errors.NewAutorepairError(
				errors.ApiCallError, "Failed to drain node %s/%s: timeout when waiting for creating evictions", node.Namespace, node.Name)
		}
	}
	if len(evictionErrs) != 0 {
		return errors.NewAutorepairError(
			errors.ApiCallError, "Failed to drain node %s/%s, due to following errors: %v", node.Namespace, node.Name, evictionErrs)
	}

	// Evictions created successfully, wait maxGracefulTerminationSec + PodEvictionHeadroom to see if pods really disappeared.
	allGone := true
	for start := time.Now(); time.Now().Sub(start) < time.Duration(maxGracefulTerminationSec)*time.Second+PodEvictionHeadroom; time.Sleep(5 * time.Second) {
		allGone = true
		for _, pod := range pods {
			podreturned, err := client.CoreV1().Pods(pod.Namespace).Get(pod.Name, metav1.GetOptions{})
			if err == nil {
				glog.Errorf("Not deleted yet %v", podreturned)
				allGone = false
				break
			}
			if !kube_errors.IsNotFound(err) {
				glog.Errorf("Failed to check pod %s/%s: %v", pod.Namespace, pod.Name, err)
				allGone = false
			}
		}
		if allGone {
			glog.V(1).Infof("All pods removed from %s", node.Name)
			// Let the deferred function know there is no need for cleanup
			return nil
		}
	}
	return errors.NewAutorepairError(
		errors.TransientError, "Failed to drain node %s/%s: pods remaining after timeout", node.Namespace, node.Name)
}

// cleanToBeDeleted cleans ToBeDeleted taints.
func cleanToBeDeleted(nodes []*apiv1.Node, client kube_client.Interface, recorder kube_record.EventRecorder) {
	for _, node := range nodes {
		cleaned, err := deletetaint.CleanToBeDeleted(node, client)
		if err != nil {
			glog.Warningf("Error while releasing taints on node %v: %v", node.Name, err)
			recorder.Eventf(node, apiv1.EventTypeWarning, "ClusterAutorepairCleanup",
				"failed to clean toBeDeletedTaint: %v", err)
		} else if cleaned {
			glog.V(1).Infof("Successfully released toBeDeletedTaint on node %v", node.Name)
			recorder.Eventf(node, apiv1.EventTypeNormal, "ClusterAutorepairCleanup", "marking the node as schedulable")
		}
	}
}

// Removes the given node from cloud provider. No extra pre-deletion actions are executed on
// the Kubernetes side.
func deleteNodeFromCloudProvider(node *apiv1.Node, cloudProvider cloudprovider.CloudProvider,
	recorder kube_record.EventRecorder, registry *clusterstate.ClusterStateRegistry) errors.AutorepairError {
	nodeGroup, err := cloudProvider.NodeGroupForNode(node)
	if err != nil {
		return errors.NewAutorepairError(
			errors.CloudProviderError, "failed to find node group for %s: %v", node.Name, err)
	}
	if nodeGroup == nil || reflect.ValueOf(nodeGroup).IsNil() {
		return errors.NewAutorepairError(errors.InternalError, "picked node that doesn't belong to a node group: %s", node.Name)
	}
	if err = nodeGroup.DeleteNodes([]*apiv1.Node{node}); err != nil {
		return errors.NewAutorepairError(errors.CloudProviderError, "failed to delete %s: %v", node.Name, err)
	}
	recorder.Eventf(node, apiv1.EventTypeNormal, "AutoRepairFailed", "node removed by Cluster autorepair")
	registry.RegisterScaleDown(&clusterstate.ScaleDownRequest{
		NodeGroupName:      nodeGroup.Id(),
		NodeName:           node.Name,
		Time:               time.Now(),
		ExpectedDeleteTime: time.Now().Add(MaxCloudProviderNodeDeletionTime),
	})
	return nil
}
