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
https://github.com/kubernetes/autoscaler/blob/cluster-autorepair-1.0.0/cluster-autoscaler/cloudprovider/builder/cloud_provider_builder.go
*/

package builder

import (
	"io"
	"os"

	"github.com/gardener/auto-node-repair/cloudprovider"
	"github.com/gardener/auto-node-repair/cloudprovider/aws"
	"github.com/gardener/auto-node-repair/cloudprovider/azure"

	"github.com/golang/glog"
)

// CloudProviderBuilder builds a cloud provider from all the necessary parameters including the name of a cloud provider e.g. aws, gce, azure
// and the path to a config file
type CloudProviderBuilder struct {
	cloudProviderFlag string
	cloudConfig       string
	clusterName       string
}

// NewCloudProviderBuilder builds a new builder from static settings
func NewCloudProviderBuilder(cloudProviderFlag string, cloudConfig string, clusterName string) CloudProviderBuilder {
	return CloudProviderBuilder{
		cloudProviderFlag: cloudProviderFlag,
		cloudConfig:       cloudConfig,
		clusterName:       clusterName,
	}
}

// Build a cloud provider from static settings contained in the builder and dynamic settings passed via args
func (b CloudProviderBuilder) Build(discoveryOpts cloudprovider.NodeGroupDiscoveryOptions, resourceLimiter *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
	var err error
	var cloudProvider cloudprovider.CloudProvider

	switch b.cloudProviderFlag {
	case "aws":
		var awsManager *aws.AwsManager
		var awsError error
		if b.cloudConfig != "" {
			config, fileErr := os.Open(b.cloudConfig)
			if fileErr != nil {
				glog.Fatalf("Couldn't open cloud provider configuration %s: %#v", b.cloudConfig, err)
			}
			defer config.Close()
			awsManager, awsError = aws.CreateAwsManager(config)
		} else {
			awsManager, awsError = aws.CreateAwsManager(nil)
		}
		if awsError != nil {
			glog.Fatalf("Failed to create AWS Manager: %v", err)
		}
		cloudProvider, err = aws.BuildAwsCloudProvider(awsManager)
		if err != nil {
			glog.Fatalf("Failed to create AWS cloud provider: %v", err)
		}
	case "azure":
		return b.buildAzure(discoveryOpts, resourceLimiter)
	}

	return cloudProvider
}

func (b CloudProviderBuilder) buildAzure(do cloudprovider.NodeGroupDiscoveryOptions, rl *cloudprovider.ResourceLimiter) cloudprovider.CloudProvider {
	var config io.ReadCloser
	if b.cloudConfig != "" {
		glog.Info("Creating Azure Manager using cloud-config file: %v", b.cloudConfig)
		var err error
		config, err := os.Open(b.cloudConfig)
		if err != nil {
			glog.Fatalf("Couldn't open cloud provider configuration %s: %#v", b.cloudConfig, err)
		}
		defer config.Close()
	} else {
		glog.Info("Creating Azure Manager with default configuration.")
	}
	manager, err := azure.CreateAzureManager(config, do)
	if err != nil {
		glog.Fatalf("Failed to create Azure Manager: %v", err)
	}
	provider, err := azure.BuildAzureCloudProvider(manager, rl)
	if err != nil {
		glog.Fatalf("Failed to create Azure cloud provider: %v", err)
	}
	return provider
}
