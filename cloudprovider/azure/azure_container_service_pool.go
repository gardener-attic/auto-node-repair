/*
Copyright 2018 The Kubernetes Authors.

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
https://github.com/kubernetes/autoscaler/blob/cluster-autoscaler-release-1.3/cluster-autoscaler/cloudprovider/azure/azure_container_service_pool.go
*/

package azure

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2017-09-30/containerservice"
	"github.com/golang/glog"

	"github.com/gardener/auto-node-repair/cloudprovider"
	"github.com/gardener/auto-node-repair/config/dynamic"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/plugin/pkg/scheduler/schedulercache"
)

//ContainerServiceAgentPool implements NodeGroup interface for agent pool deployed in ACS/AKS
type ContainerServiceAgentPool struct {
	azureRef
	manager *AzureManager
	util    *AzUtil

	minSize           int
	maxSize           int
	serviceType       string
	clusterName       string
	resourceGroup     string
	nodeResourceGroup string

	mutex       sync.Mutex
	lastRefresh time.Time
	curSize     int64
}

//NewContainerServiceAgentPool constructs ContainerServiceAgentPool from the --node param
//and azure manager
func NewContainerServiceAgentPool(spec *dynamic.NodeGroupSpec, am *AzureManager) (*ContainerServiceAgentPool, error) {
	asg := &ContainerServiceAgentPool{
		azureRef: azureRef{
			Name: spec.Name,
		},
		minSize: spec.MinSize,
		maxSize: spec.MaxSize,
		manager: am,
	}

	asg.util = &AzUtil{
		manager: am,
	}
	asg.serviceType = am.config.VMType
	asg.clusterName = am.config.ClusterName
	asg.resourceGroup = am.config.ResourceGroup

	// In case of AKS there is a different resource group for the worker nodes, where as for
	// ACS the vms are in the same group as that of the service.
	if am.config.VMType == vmTypeAKS {
		asg.nodeResourceGroup = am.config.NodeResourceGroup
	} else {
		asg.nodeResourceGroup = am.config.ResourceGroup
	}
	return asg, nil
}

//GetPool is an internal function which figures out agentPoolProfile from the list based on the pool name provided in the --node parameter passed
//to the autoscaler main
func (agentPool *ContainerServiceAgentPool) GetPool(agentProfiles *[]containerservice.AgentPoolProfile) (ret *containerservice.AgentPoolProfile) {
	for _, value := range *agentProfiles {
		profileName := *value.Name
		glog.V(5).Infof("AgentPool profile name: %s", profileName)
		if profileName == (agentPool.azureRef.Name) {
			return &value
		}
	}
	// Note: In some older ACS clusters, the name of the agentProfile can be different from kubernetes
	// label and vm pool tag. It can be have a "pool0" appended.
	// In the above loop we would check the normal case and if not returned yet, we will try matching
	// the node pool name with "pool0" and try a match as a workaround.
	if agentPool.serviceType == vmTypeACS {
		for _, value := range *agentProfiles {
			profileName := *value.Name
			poolName := agentPool.azureRef.Name + "pool0"
			glog.V(5).Infof("Workaround match check - Profile: %s <=> Poolname: %s", profileName, poolName)
			if profileName == poolName {
				return &value
			}
		}
	}
	return nil
}

//GetNodeCount returns the count of nodes from the agent pool profile
func (agentPool *ContainerServiceAgentPool) GetNodeCount(agentProfiles *[]containerservice.AgentPoolProfile) (count int, err error) {
	pool := agentPool.GetPool(agentProfiles)
	if pool != nil {
		count := (int)(*pool.Count)
		glog.V(5).Infof("Got node count: %d", count)
		return count, nil
	}
	return -1, fmt.Errorf("Could not find pool with name: %s", agentPool.azureRef)
}

//SetNodeCount sets the count of nodes in the in memory pool profile
func (agentPool *ContainerServiceAgentPool) SetNodeCount(agentProfiles *[]containerservice.AgentPoolProfile, count int) (err error) {
	pool := agentPool.GetPool(agentProfiles)
	if pool != nil {
		*(pool.Count) = int32(count)
		return nil
	}
	return fmt.Errorf("Could not find pool with name: %s", agentPool.azureRef)
}

//GetProviderID converts the name of a node into the form that kubernetes cloud
//provider id is presented in.
func (agentPool *ContainerServiceAgentPool) GetProviderID(name string) string {
	//TODO: come with a generic way to make it work with provider id formats
	// in different version of k8s.
	return "azure://" + name
}

//GetName extracts the name of the node (a format which underlying cloud service understands)
//from the cloud providerID (format which kubernetes understands)
func (agentPool *ContainerServiceAgentPool) GetName(providerID string) (string, error) {
	// Remove the "azure://" string from it
	providerID = strings.TrimPrefix(providerID, "azure://")
	ctx, cancel := getContextWithCancel()
	defer cancel()
	vms, err := agentPool.manager.azClient.virtualMachinesClient.List(ctx, agentPool.nodeResourceGroup)
	if err != nil {
		return "", err
	}
	for _, vm := range vms {
		if strings.Compare(*vm.ID, providerID) == 0 {
			return *vm.Name, nil
		}
	}
	return "", fmt.Errorf("VM list empty")
}

//MaxSize returns the maximum size scale limit provided by --node
//parameter to the autoscaler main
func (agentPool *ContainerServiceAgentPool) MaxSize() int {
	return agentPool.maxSize
}

//MinSize returns the minimum size the cluster is allowed to scaled down
//to as provided by the node spec in --node parameter.
func (agentPool *ContainerServiceAgentPool) MinSize() int {
	return agentPool.minSize
}

//TargetSize gathers the target node count set for the cluster by
//querying the underlying service.
func (agentPool *ContainerServiceAgentPool) TargetSize() (int, error) {
	ctx, cancel := getContextWithCancel()
	defer cancel()

	// AKS
	if agentPool.serviceType == vmTypeAKS {
		aksContainerService, err := agentPool.manager.azClient.managedContainerServicesClient.Get(ctx,
			agentPool.resourceGroup,
			agentPool.clusterName)
		if err != nil {
			glog.Error(err)
			return -1, err
		}
		return agentPool.GetNodeCount(aksContainerService.AgentPoolProfiles)
	}

	// ACS
	acsContainerService, err := agentPool.manager.azClient.containerServicesClient.Get(ctx,
		agentPool.resourceGroup,
		agentPool.clusterName)
	if err != nil {
		glog.Error(err)
		return -1, err
	}
	return agentPool.GetNodeCount(acsContainerService.AgentPoolProfiles)
}

//SetSize contacts the underlying service and sets the size of the pool.
//This will be called when a scale up occurs and will be called just after
//a delete is performed from a scale down.
func (agentPool *ContainerServiceAgentPool) SetSize(targetSize int) error {
	var err error
	var poolProfiles *[]containerservice.AgentPoolProfile
	var containerService containerservice.ContainerService
	var managedCluster containerservice.ManagedCluster

	glog.Infof("Set size request: %d", targetSize)
	if targetSize > agentPool.MaxSize() || targetSize < agentPool.MinSize() {
		glog.Errorf("Target size %d requested outside Max: %d, Min: %d", targetSize, agentPool.MaxSize(), agentPool.MaxSize)
		return fmt.Errorf("Target size %d requested outside Max: %d, Min: %d", targetSize, agentPool.MaxSize(), agentPool.MinSize())
	}
	ctx, cancel := getContextWithCancel()
	defer cancel()

	if agentPool.serviceType == vmTypeAKS {
		managedCluster, err = agentPool.manager.azClient.managedContainerServicesClient.Get(ctx, agentPool.resourceGroup,
			agentPool.clusterName)
		poolProfiles = managedCluster.AgentPoolProfiles
		if err != nil {
			return err
		}
		glog.V(5).Infof("AKS: %+v", managedCluster)
	} else {
		containerService, err = agentPool.manager.azClient.containerServicesClient.Get(ctx, agentPool.resourceGroup,
			agentPool.clusterName)
		if err != nil {
			return err
		}
		poolProfiles = containerService.AgentPoolProfiles
		glog.V(5).Infof("ACS: %+v", containerService)
	}
	if err != nil {
		glog.Error(err)
		return err
	}
	currentSize, err := agentPool.GetNodeCount(poolProfiles)
	if err != nil {
		glog.Error(err)
		return err
	}

	glog.Infof("Current size: %d, Target size requested: %d", currentSize, targetSize)

	// Set the value in the volatile structure.
	err = agentPool.SetNodeCount(poolProfiles, targetSize)
	if err != nil {
		glog.Error(err)
		return err
	}

	var end time.Time
	start := time.Now()

	updateCtx, updateCancel := getContextWithCancel()
	defer updateCancel()

	//Update the service with the new value.
	if agentPool.serviceType == vmTypeAKS {
		aksClient := agentPool.manager.azClient.managedContainerServicesClient
		updatedVal, err := aksClient.CreateOrUpdate(updateCtx, agentPool.resourceGroup,
			agentPool.clusterName, managedCluster)
		if err != nil {
			glog.Error(err)
			return err
		}
		err = updatedVal.WaitForCompletion(updateCtx, aksClient.Client)
		if err != nil {
			glog.Error(err)
			return err
		}
		newVal, err := updatedVal.Result(aksClient)
		if err != nil {
			glog.Error(err)
			return err
		}
		end = time.Now()
		glog.Infof("Target size set done, AKS. Value: %+v\n", newVal)
	} else {
		acsClient := agentPool.manager.azClient.containerServicesClient
		updatedVal, err := acsClient.CreateOrUpdate(updateCtx, agentPool.resourceGroup,
			agentPool.clusterName, containerService)
		if err != nil {
			glog.Error(err)
			return err
		}
		err = updatedVal.WaitForCompletion(updateCtx, acsClient.Client)
		if err != nil {
			glog.Error(err)
			return err
		}
		newVal, err := updatedVal.Result(acsClient)
		if err != nil {
			glog.Error(err)
			return err
		}
		end = time.Now()
		glog.Infof("Target size set done, ACS. Val: %+v", newVal)
	}

	glog.Infof("Got Updated value. Time taken: %s", end.Sub(start).String())
	return nil
}

//IncreaseSize calls in the underlying SetSize to increase the size in response
//to a scale up. It calculates the expected size based on a delta provided as
//parameter
func (agentPool *ContainerServiceAgentPool) IncreaseSize(delta int) error {
	if delta <= 0 {
		return fmt.Errorf("Size increase must be +ve")
	}
	currentSize, err := agentPool.TargetSize()
	if err != nil {
		return err
	}
	targetSize := int(currentSize) + delta
	if targetSize > agentPool.MaxSize() {
		return fmt.Errorf("Size increase request of %d more than max size %d set", targetSize, agentPool.MaxSize())
	}
	return agentPool.SetSize(targetSize)
}

//DeleteNodesInternal calls the underlying vm service to delete the node.
//Additionally it also call into the container service to reflect the node count resulting
//from the deletion.
func (agentPool *ContainerServiceAgentPool) DeleteNodesInternal(providerIDs []string) error {
	currentSize, err := agentPool.TargetSize()
	if err != nil {
		return err
	}
	// Set the size to current size. Reduce as we are able to go through.
	targetSize := currentSize

	for _, providerID := range providerIDs {
		glog.Infof("ProviderID got to delete: %s", providerID)
		nodeName, err := agentPool.GetName(providerID)
		if err != nil {
			return err
		}
		glog.Infof("VM name got to delete: %s", nodeName)

		err = agentPool.util.DeleteVirtualMachine(agentPool.nodeResourceGroup, nodeName)
		if err != nil {
			glog.Error(err)
			return err
		}
		targetSize--
	}

	// TODO: handle the errors from delete operation.
	if currentSize != targetSize {
		agentPool.SetSize(targetSize)
	}
	return nil
}

//DeleteNodes extracts the providerIDs from the node spec and calls into the internal
//delete method.
func (agentPool *ContainerServiceAgentPool) DeleteNodes(nodes []*apiv1.Node) error {
	var providerIDs []string
	for _, node := range nodes {
		glog.Infof("Node: %s", node.Spec.ProviderID)
		providerIDs = append(providerIDs, node.Spec.ProviderID)
	}
	for _, p := range providerIDs {
		glog.Infof("ProviderID before calling acsmgr: %s", p)
	}
	return agentPool.DeleteNodesInternal(providerIDs)
}

//IsContainerServiceNode checks if the tag from the vm matches the agentPool name
func (agentPool *ContainerServiceAgentPool) IsContainerServiceNode(tags map[string]*string) bool {
	poolName := tags["poolName"]
	if poolName != nil {
		glog.V(5).Infof("Matching agentPool name: %s with tag name: %s", agentPool.azureRef.Name, *poolName)
		if *poolName == agentPool.azureRef.Name {
			return true
		}
	}
	return false
}

//GetNodes extracts the node list from the underlying vm service and returns back
//equivalent providerIDs  as list.
func (agentPool *ContainerServiceAgentPool) GetNodes() ([]string, error) {
	ctx, cancel := getContextWithCancel()
	defer cancel()
	vmList, err := agentPool.manager.azClient.virtualMachinesClient.List(ctx, agentPool.nodeResourceGroup)
	if err != nil {
		glog.Error("Error", err)
		return nil, err
	}
	var nodeArray []string
	for _, node := range vmList {
		glog.V(5).Infof("Node Name: %s, ID: %s", *node.Name, *node.ID)
		if agentPool.IsContainerServiceNode(node.Tags) {
			providerID := agentPool.GetProviderID(*node.ID)
			glog.V(5).Infof("Returning back the providerID: %s", providerID)
			nodeArray = append(nodeArray, providerID)
		}
	}
	return nodeArray, nil
}

//DecreaseTargetSize requests the underlying service to decrease the node count.
func (agentPool *ContainerServiceAgentPool) DecreaseTargetSize(delta int) error {
	if delta >= 0 {
		glog.Errorf("Size decrease error: %d", delta)
		return fmt.Errorf("Size decrease must be negative")
	}
	currentSize, err := agentPool.TargetSize()
	if err != nil {
		glog.Error(err)
		return err
	}
	// Get the current nodes in the list
	nodes, err := agentPool.GetNodes()
	if err != nil {
		glog.Error(err)
		return err
	}
	targetSize := int(currentSize) + delta
	if targetSize < len(nodes) {
		return fmt.Errorf("attempt to delete existing nodes targetSize:%d delta:%d existingNodes: %d",
			currentSize, delta, len(nodes))
	}
	return agentPool.SetSize(targetSize)
}

//Id returns the name of the agentPool
func (agentPool *ContainerServiceAgentPool) Id() string {
	return agentPool.azureRef.Name
}

//Debug returns a string with basic details of the agentPool
func (agentPool *ContainerServiceAgentPool) Debug() string {
	return fmt.Sprintf("%s (%d:%d)", agentPool.Id(), agentPool.MinSize(), agentPool.MaxSize())
}

//Nodes returns the list of nodes in the agentPool.
func (agentPool *ContainerServiceAgentPool) Nodes() ([]string, error) {
	return agentPool.GetNodes()
}

//TemplateNodeInfo is not implemented.
func (agentPool *ContainerServiceAgentPool) TemplateNodeInfo() (*schedulercache.NodeInfo, error) {
	return nil, cloudprovider.ErrNotImplemented
}

//Exist is always true since we are initialized with an existing agentpool
func (agentPool *ContainerServiceAgentPool) Exist() bool {
	return true
}

//Create is returns already exists since we don't support the
//agent pool creation.
func (agentPool *ContainerServiceAgentPool) Create() error {
	return cloudprovider.ErrAlreadyExist
}

//Delete is not implemented since we don't support agent pool
//deletion.
func (agentPool *ContainerServiceAgentPool) Delete() error {
	return cloudprovider.ErrNotImplemented
}

//Autoprovisioned is set to false to indicate that this code
//does not create agentPools by itself.
func (agentPool *ContainerServiceAgentPool) Autoprovisioned() bool {
	return false
}
