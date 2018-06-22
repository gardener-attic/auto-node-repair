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
https://github.com/kubernetes/autoscaler/blob/cluster-autorepair-1.0.0/cluster-autorepair/main.go
*/

package main

import (
	"flag"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/pflag"

	"github.com/gardener/auto-node-repair/config"
	"github.com/gardener/auto-node-repair/core"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/kubernetes/pkg/apis/componentconfig"

	kube_util "github.com/gardener/auto-node-repair/utils/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kube_flag "k8s.io/apiserver/pkg/util/flag"
	kube_client "k8s.io/client-go/kubernetes"
	kube_leaderelection "k8s.io/client-go/tools/leaderelection"
)

// MultiStringFlag is a flag for passing multiple parameters using same flag
type MultiStringFlag []string

// String returns string representation of the node groups.
func (flag *MultiStringFlag) String() string {
	return "[" + strings.Join(*flag, " ") + "]"
}

// Set adds a new configuration.
func (flag *MultiStringFlag) Set(value string) error {
	*flag = append(*flag, value)
	return nil
}

const (
	defaultStartupDelay  = 10 * time.Minute
	defaultPollDelay     = 10 * time.Second
	defaultLeaseDuration = 15 * time.Second
	defaultRenewDeadline = 10 * time.Second
	defaultRetryPeriod   = 2 * time.Second
	defaultRepairPeriod  = 10 * time.Minute
)

var (
	nodeGroupsFlag    MultiStringFlag
	clusterName       = flag.String("clusterName", "", "Autorepairable cluster name, if available")
	address           = flag.String("address", ":8085", "The address to expose prometheus metrics.")
	kubernetes        = flag.String("kubernetes", "", "Kubernetes master location. Leave blank for default")
	kubeConfigFile    = flag.String("kubeconfig", "", "Path to kubeconfig file with authorization and master location information.")
	cloudConfig       = flag.String("cloud-config", "", "The path to the cloud provider configuration file.  Empty string for no configuration file.")
	namespace         = flag.String("namespace", "kube-system", "Namespace in which auto-node-repair run. If a --configmap flag is also provided, ensure that the configmap exists in this namespace before CA runs.")
	scanInterval      = flag.Duration("scan-interval", defaultPollDelay, "How often cluster is reevaluated for scale up or down")
	cloudProviderFlag = flag.String("cloud-provider", "", "Cloud provider type. Allowed values: aws, azure")
	repairPeriod      = flag.Duration("repair-period", defaultRepairPeriod, "How long controller should wait before attempting to repair unhealthy nodes")
	startUpDelay      = flag.Duration("initial-delay", defaultStartupDelay, "How long controller should wait while first boot up")
)

func createAutorepairOptions() core.AutorepairOptions {
	autorepairingOpts := core.AutorepairingOptions{
		CloudConfig:       *cloudConfig,
		CloudProviderName: *cloudProviderFlag,
		NodeGroups:        nodeGroupsFlag,
		ConfigNamespace:   *namespace,
		ClusterName:       *clusterName,
	}

	return core.AutorepairOptions{
		AutorepairingOptions: autorepairingOpts,
	}
}

func createKubeClient() kube_client.Interface {
	if *kubeConfigFile != "" {
		glog.V(1).Infof("Using kubeconfig file: %s", *kubeConfigFile)
		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeConfigFile)
		if err != nil {
			glog.Fatalf("Failed to build config: %v", err)
		}
		clientset, err := kube_client.NewForConfig(config)
		if err != nil {
			glog.Fatalf("Create clientset error: %v", err)
		}
		return clientset
	}
	url, err := url.Parse(*kubernetes)
	if err != nil {
		glog.Fatalf("Failed to parse Kubernetes url: %v", err)
	}

	kubeConfig, err := config.GetKubeClientConfig(url)
	if err != nil {
		glog.Fatalf("Failed to build Kubernetes client configuration: %v", err)
	}

	return kube_client.NewForConfigOrDie(kubeConfig)
}

func run() {
	glog.Infof("Auto-Node-Repair: Waiting %v for nodes to be ready", *startUpDelay)
	time.Sleep(*startUpDelay)

	kubeClient := createKubeClient()
	kubeEventRecorder := kube_util.CreateEventRecorder(kubeClient)
	opts := createAutorepairOptions()

	listerRegistryStopChannel := make(chan struct{})
	listerRegistry := kube_util.NewListerRegistryWithDefaultListers(kubeClient, listerRegistryStopChannel)
	autorepair, err := core.NewAutorepair(opts, kubeClient, kubeEventRecorder, listerRegistry)
	if err != nil {
		glog.Fatalf("Failed to create auto repair: %v", err)
	}

	glog.Info("Auto-Node-Repair: Started successfully")

	for {
		select {
		case <-time.After(*scanInterval):
			{
				err = autorepair.RunOnce(*repairPeriod)
				if err != nil {
					glog.Fatalf("Failed to run autorepair : %v", err)
				}
			}
		}
	}
}

func main() {
	leaderElection := defaultLeaderElectionConfiguration()
	leaderElection.LeaderElect = false

	bindFlags(&leaderElection, pflag.CommandLine)
	flag.Var(&nodeGroupsFlag, "nodes", "sets min,max size and other configuration data for a node group in a format accepted by cloud provider."+
		"Can be used multiple times. Format: <min>:<max>:<other...>")
	kube_flag.InitFlags()

	if !leaderElection.LeaderElect {
		run()
	} else {
		id, err := os.Hostname()
		if err != nil {
			glog.Fatalf("Unable to get hostname: %v", err)
		}

		kubeClient := createKubeClient()

		// Validate that the client is ok.
		_, err = kubeClient.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			glog.Fatalf("Failed to get nodes from apiserver: %v", err)
		}

		lock, err := resourcelock.New(
			leaderElection.ResourceLock,
			*namespace,
			"auto-node-repair",
			kubeClient.CoreV1(),
			resourcelock.ResourceLockConfig{
				Identity:      id,
				EventRecorder: kube_util.CreateEventRecorder(kubeClient),
			},
		)
		if err != nil {
			glog.Fatalf("Unable to create leader election lock: %v", err)
		}

		kube_leaderelection.RunOrDie(kube_leaderelection.LeaderElectionConfig{
			Lock:          lock,
			LeaseDuration: leaderElection.LeaseDuration.Duration,
			RenewDeadline: leaderElection.RenewDeadline.Duration,
			RetryPeriod:   leaderElection.RetryPeriod.Duration,
			Callbacks: kube_leaderelection.LeaderCallbacks{
				OnStartedLeading: func(_ <-chan struct{}) {
					// Since we are committing a suicide after losing
					// mastership, we can safely ignore the argument.
					run()
				},
				OnStoppedLeading: func() {
					glog.Fatalf("lost master")
				},
			},
		})
	}
}

func defaultLeaderElectionConfiguration() componentconfig.LeaderElectionConfiguration {
	return componentconfig.LeaderElectionConfiguration{
		LeaderElect:   false,
		LeaseDuration: metav1.Duration{Duration: defaultLeaseDuration},
		RenewDeadline: metav1.Duration{Duration: defaultRenewDeadline},
		RetryPeriod:   metav1.Duration{Duration: defaultRetryPeriod},
		ResourceLock:  resourcelock.EndpointsResourceLock,
	}
}

func bindFlags(l *componentconfig.LeaderElectionConfiguration, fs *pflag.FlagSet) {
	fs.BoolVar(&l.LeaderElect, "leader-elect", l.LeaderElect, ""+
		"Start a leader election client and gain leadership before "+
		"executing the main loop. Enable this when running replicated "+
		"components for high availability.")
	fs.DurationVar(&l.LeaseDuration.Duration, "leader-elect-lease-duration", l.LeaseDuration.Duration, ""+
		"The duration that non-leader candidates will wait after observing a leadership "+
		"renewal until attempting to acquire leadership of a led but unrenewed leader "+
		"slot. This is effectively the maximum duration that a leader can be stopped "+
		"before it is replaced by another candidate. This is only applicable if leader "+
		"election is enabled.")
	fs.DurationVar(&l.RenewDeadline.Duration, "leader-elect-renew-deadline", l.RenewDeadline.Duration, ""+
		"The interval between attempts by the acting master to renew a leadership slot "+
		"before it stops leading. This must be less than or equal to the lease duration. "+
		"This is only applicable if leader election is enabled.")
	fs.DurationVar(&l.RetryPeriod.Duration, "leader-elect-retry-period", l.RetryPeriod.Duration, ""+
		"The duration the clients should wait between attempting acquisition and renewal "+
		"of a leadership. This is only applicable if leader election is enabled.")
	fs.StringVar(&l.ResourceLock, "leader-elect-resource-lock", l.ResourceLock, ""+
		"The type of resource resource object that is used for locking during"+
		"leader election. Supported options are `endpoints` (default) and `configmap`.")
}
