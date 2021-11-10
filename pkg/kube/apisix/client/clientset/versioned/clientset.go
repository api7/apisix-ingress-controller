// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by client-gen. DO NOT EDIT.

package versioned

import (
	"fmt"

	apisixv2beta1 "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/client/clientset/versioned/typed/config/v2beta1"
	apisixv2beta2 "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/client/clientset/versioned/typed/config/v2beta2"
	apisixv2beta3 "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/client/clientset/versioned/typed/config/v2beta3"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	ApisixV2beta1() apisixv2beta1.ApisixV2beta1Interface
	ApisixV2beta2() apisixv2beta2.ApisixV2beta2Interface
	ApisixV2beta3() apisixv2beta3.ApisixV2beta3Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	apisixV2beta1 *apisixv2beta1.ApisixV2beta1Client
	apisixV2beta2 *apisixv2beta2.ApisixV2beta2Client
	apisixV2beta3 *apisixv2beta3.ApisixV2beta3Client
}

// ApisixV2beta1 retrieves the ApisixV2beta1Client
func (c *Clientset) ApisixV2beta1() apisixv2beta1.ApisixV2beta1Interface {
	return c.apisixV2beta1
}

// ApisixV2beta2 retrieves the ApisixV2beta2Client
func (c *Clientset) ApisixV2beta2() apisixv2beta2.ApisixV2beta2Interface {
	return c.apisixV2beta2
}

// ApisixV2beta3 retrieves the ApisixV2beta3Client
func (c *Clientset) ApisixV2beta3() apisixv2beta3.ApisixV2beta3Interface {
	return c.apisixV2beta3
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfig will generate a rate-limiter in configShallowCopy.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		if configShallowCopy.Burst <= 0 {
			return nil, fmt.Errorf("burst is required to be greater than 0 when RateLimiter is not set and QPS is set to greater than 0")
		}
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.apisixV2beta1, err = apisixv2beta1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.apisixV2beta2, err = apisixv2beta2.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.apisixV2beta3, err = apisixv2beta3.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.apisixV2beta1 = apisixv2beta1.NewForConfigOrDie(c)
	cs.apisixV2beta2 = apisixv2beta2.NewForConfigOrDie(c)
	cs.apisixV2beta3 = apisixv2beta3.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.apisixV2beta1 = apisixv2beta1.New(c)
	cs.apisixV2beta2 = apisixv2beta2.New(c)
	cs.apisixV2beta3 = apisixv2beta3.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
