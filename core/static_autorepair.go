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
https://github.com/kubernetes/autoscaler/blob/cluster-autorepair-1.0.0/cluster-autoscaler/core/static_autoscaler.go
*/

package core

import (
	"time"
	
	"github.com/golang/glog"
	
	"github.com/gardener/auto-node-repair/clusterstate/utils"
	"github.com/gardener/auto-node-repair/utils/errors"
	
	kube_util "github.com/gardener/auto-node-repair/utils/kubernetes"
	kube_client "k8s.io/client-go/kubernetes"
	kube_record "k8s.io/client-go/tools/record"
)

// StaticAutorepair is an autorepair which has all the core functionality of a CA but without the reconfiguration feature
type StaticAutorepair struct {
	*AutorepairingContext
	kube_util.ListerRegistry
	lastScaleUpTime         time.Time
	lastScaleDownDeleteTime time.Time
	lastScaleDownFailTime   time.Time
}

// NewStaticAutorepair creates an instance of Autorepair filled with provided parameters
func NewStaticAutorepair(opts AutorepairingOptions,
	kubeClient kube_client.Interface, kubeEventRecorder kube_record.EventRecorder, listerRegistry kube_util.ListerRegistry) (*StaticAutorepair, errors.AutorepairError) {
	logRecorder, err := utils.NewStatusMapRecorder(kubeClient, opts.ConfigNamespace, kubeEventRecorder, opts.WriteStatusConfigMap)
	if err != nil {
		glog.Error("Failed to initialize status configmap, unable to write status events")
		logRecorder, _ = utils.NewStatusMapRecorder(kubeClient, opts.ConfigNamespace, kubeEventRecorder, false)
	}
	autorepairingContext, errctx := NewAutorepairContext(opts, kubeClient, kubeEventRecorder, logRecorder, listerRegistry)
	if errctx != nil {
		return nil, errctx
	}

	return &StaticAutorepair{
		AutorepairingContext:    autorepairingContext,
		ListerRegistry:          listerRegistry,
		lastScaleUpTime:         time.Now(),
		lastScaleDownDeleteTime: time.Now(),
		lastScaleDownFailTime:   time.Now(),
	}, nil
}

// RunOnce iterates over node groups and repairs them if necessary
func (a *StaticAutorepair) RunOnce(repairtime time.Duration) (errors.AutorepairError) {
	
	sleepTime := repairtime
	restartDelay := 1 * time.Minute

	notReadyNodeLister := a.NotReadyNodeLister()
	notReadyNodes, _ := notReadyNodeLister.List()
	if len(notReadyNodes) == 0 {
		//glog.Infof("All the nodes are healthy")
		return nil
	}
	
	glog.Infof("Number of not ready nodes detected: %v", len(notReadyNodes))
 	glog.Infof("Waiting for %v to confirm unhealthy nodes", repairtime)

	time.Sleep(sleepTime)

	stillNotReadyNodeLister := a.NotReadyNodeLister()
	stillNotReadyNodes, _ := stillNotReadyNodeLister.List()

	glog.Infof("Confirmed unhealthy nodes: %v", len(stillNotReadyNodes))

	toBeRepairedNodes := CommonNodes(notReadyNodes, stillNotReadyNodes)
	if len(toBeRepairedNodes) != 0 {
		glog.Infof("Auto-Node-Repair: Initializing auto repair. %v nodes unhealthy !", len(notReadyNodes))
		a.RepairNodes(toBeRepairedNodes)
		glog.Info("Auto-Node-Repair: Unhealthy nodes have been replaced")
		time.Sleep(restartDelay)
	} /*else{
		glog.Info("Not ready nodes were recovered")
	}*/
	
	return nil
}
