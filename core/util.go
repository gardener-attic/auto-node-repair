/*
Copyright (c) 2017 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http:  www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package core

import (
	"time"
	"flag"

	"github.com/golang/glog"
	
	"k8s.io/kubernetes/plugin/pkg/scheduler/schedulercache"
	"github.com/gardener/auto-node-repair/simulator"
	"github.com/gardener/auto-node-repair/cloudprovider"
	
	apiv1 "k8s.io/api/core/v1"
)

var (
	skipNodesWithSystemPods = flag.Bool("skip-nodes-with-system-pods", false,
		"If true Cluster autorepair will never delete nodes with pods from kube-system (except for DaemonSet "+
			"or mirror pods)")
	skipNodesWithLocalStorage = flag.Bool("skip-nodes-with-local-storage", false,
		"If true Cluster autorepair will never delete nodes with pods with local storage, e.g. EmptyDir or HostPath")

	minReplicaCount = flag.Int("min-replica-count", 0,
		"Minimum number or replicas that a replica set or replication controller should have to allow their pods deletion in auto repair")
)

// Get common nodes between two nodelists
func CommonNodes(nodelist1 []*apiv1.Node, nodelist2 []*apiv1.Node) []*apiv1.Node {
	var commonElems [] *apiv1.Node
	for _, elem := range nodelist1 {
		if ContainsNode(elem, nodelist2) {
			commonElems = append(commonElems, elem)
		}
	}
	return commonElems
}

// Check if nodelist contains passed node
func ContainsNode(node *apiv1.Node, nodelist []*apiv1.Node) bool {
	for _, elem := range nodelist {
		if node.Name == elem.Name {
			return true
		}
	}
	return false	
}

// Calculate possible expansion size for a ASG based on allowable maximum size
func (a *StaticAutorepair) CalcExpandSize(asg cloudprovider.NodeGroup, toBeRepairedNodes []* apiv1.Node) int {
	desiredSize, _ := asg.TargetSize()
	maxSize := asg.MaxSize()

	desiredIncreaseSize := len(toBeRepairedNodes)
    allowableIncreaseSize := maxSize - desiredSize
	effectiveIncreaseSize := 0

    if allowableIncreaseSize >= desiredIncreaseSize {
    	effectiveIncreaseSize = desiredIncreaseSize
    } else {
    	effectiveIncreaseSize = allowableIncreaseSize
    }

    return effectiveIncreaseSize
}

// Method used to extract not ready nodes for a particular cluster, notReadyNodes passed (maybe stale node status)
func (a *StaticAutorepair) AsgNotReadyNodes(desiredAsg cloudprovider.NodeGroup, notReadyNodes []* apiv1.Node) ([]*apiv1.Node, error) {
	belongsToAsg := make([]*apiv1.Node, 0, len(notReadyNodes))
	
	for  _, node := range notReadyNodes {	
		asg, err := a.AutorepairingContext.CloudProvider.NodeGroupForNode(node)
		if err != nil {
			glog.Errorf("Failed to get node group: %v", err)
			return []*apiv1.Node{}, err
		}

		if asg == desiredAsg {
			belongsToAsg = append(belongsToAsg, node)
		}
	}

	return belongsToAsg, nil
}

// Ready ASG nodes based on real-time, used to poll on newly created nodes (Latest node status)
func (a *StaticAutorepair) AsgReadyNodes(desiredAsg cloudprovider.NodeGroup) ([]*apiv1.Node, error) {
	ReadyNodeLister := a.ReadyNodeLister()
	readyNodes, _ := ReadyNodeLister.List()
	
	belongsToAsg := make([]*apiv1.Node, 0, len(readyNodes))
	for  _, node := range readyNodes {
		
		asg, err := a.AutorepairingContext.CloudProvider.NodeGroupForNode(node)
		if err != nil {
			glog.Errorf("Failed to get node group: %v", err)
			return []*apiv1.Node{}, err
		}

		if asg == desiredAsg {
			belongsToAsg = append(belongsToAsg, node)
		}
	}

	return belongsToAsg, nil
}

// Increases size of ASG, Blocking call untill all nodes are ready
func (a *StaticAutorepair) IncreaseSize(asg cloudprovider.NodeGroup, toBeRepairedNodes []* apiv1.Node, incSize int) {
	sleepTime := 30 * time.Second
	timeOut := 10
	iter := 0
	
	currReadyNodes, _ := a.AsgReadyNodes(asg)
	
	prevSize := len(currReadyNodes)
	currSize := len(currReadyNodes)

	err := asg.IncreaseSize(incSize)
	if err != nil {
		glog.Errorf("Failed to increase size of ASG: %v", err)
		return
	}

	// Waiting for nodes to become ready
    for (currSize != prevSize + incSize) && (iter < timeOut)  {
    	//glog.Infof("Waiting for %d node(s) to become ready", prevSize + incSize - currSize)
    	time.Sleep(sleepTime)
    	
    	currReadyNodes, _ = a.AsgReadyNodes(asg)
    	currSize = len(currReadyNodes)
    	iter += 1
    }
    //glog.Info("Newly created nodes are now ready !\n")
    return 
}

// Delete nodes from ASG, Blocking call
func (a *StaticAutorepair) DeleteNodes(asg cloudprovider.NodeGroup, nodeList []*apiv1.Node) {
	sleepTime := 30 * time.Second
	timeOut := 4
	iter := 0

	decSize := len(nodeList)
	prevSize, _ := asg.TargetSize()
	currSize, _ := asg.TargetSize()

	pdbs, err := a.PodDisruptionBudgetLister().List()
	allScheduledPods, err := a.ScheduledPodLister().List()

	if err != nil {
		glog.Errorf("Failed to list scheduled pods: %v", err)	
	}

	nodeNameToNodeInfo := schedulercache.CreateNodeNameToInfoMap(allScheduledPods, nodeList)
	
	for _, node := range nodeList {
		//glog.Infof("Deleting node: %v", node.ObjectMeta.Name)

		var podsToRemove []*apiv1.Pod
		
		if nodeInfo, found := nodeNameToNodeInfo[node.Name]; found {
			podsToRemove, err = simulator.DetailedGetPodsForMove(nodeInfo, *skipNodesWithSystemPods, *skipNodesWithLocalStorage, a.AutorepairingContext.ClientSet, int32(*minReplicaCount),
					pdbs)
			if err != nil {
				glog.Errorf("node %s cannot be removed: %v", node.Name, err)
				continue
			}
		} else {
			glog.Errorf("nodeInfo for %s not found",node.Name)
			continue
		}


		err := deleteNode(a.AutorepairingContext, node, podsToRemove)
		if err != nil {
			glog.Errorf("Failed to delete %s: %v", node.Name, err)
			return
		}
	}

	// Waiting for node to be deleted
	for (currSize != prevSize - decSize) && (iter < timeOut)  {		
		//glog.Infof("Waiting for %d node(s) to be deleted", currSize + decSize - prevSize)
		time.Sleep(sleepTime)
		
		currSize, _ = asg.TargetSize()
		iter += 1
	}
	//glog.Info("Nodes are successfully deleted !\n")

}

// Repair nodes in a particular ASG
func (a *StaticAutorepair) AsgRepairNodes(asg cloudprovider.NodeGroup, toBeRepairedNodes []*apiv1.Node) {
	desiredSize, _ := asg.TargetSize()
	maxSize := asg.MaxSize()
	
	if len(toBeRepairedNodes) == 0 {
		return
	} else {
		incSize := 1
		if desiredSize >= maxSize {
			// Delete first as ASG is already at max size 
			a.DeleteNodes(asg, toBeRepairedNodes[0:incSize])	
			a.IncreaseSize(asg, toBeRepairedNodes[0:incSize], incSize)
		} else{
			incSize = a.CalcExpandSize(asg, toBeRepairedNodes)
			a.IncreaseSize(asg, toBeRepairedNodes[0:incSize], incSize)
			a.DeleteNodes(asg, toBeRepairedNodes[0:incSize])
		}
		// Recurrisvely call ASG until its healthy
		a.AsgRepairNodes(asg, toBeRepairedNodes[incSize:])
	}
}

// Repair nodes for a particular kubernetes cluster
func (a *StaticAutorepair) RepairNodes(toBeRepairedNodes []*apiv1.Node) {
	asgs := a.AutorepairingContext.CloudProvider.NodeGroups()

	for _, asg := range asgs {
		asgToBeRepariedNodes, _ := a.AsgNotReadyNodes(asg, toBeRepairedNodes)
		//glog.Info("\n\n\n")
		//glog.Infof("--- Repairing %d nodes for asg-%d ---", len(asgToBeRepariedNodes), (i+1))
		a.AsgRepairNodes(asg, asgToBeRepariedNodes)
	}
}
