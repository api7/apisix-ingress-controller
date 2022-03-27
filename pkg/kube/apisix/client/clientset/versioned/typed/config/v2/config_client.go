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

package v2

import (
	v2 "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/apis/config/v2"
	"github.com/apache/apisix-ingress-controller/pkg/kube/apisix/client/clientset/versioned/scheme"
	rest "k8s.io/client-go/rest"
)

type ApisixV2Interface interface {
	RESTClient() rest.Interface
	ApisixClusterConfigsGetter
	ApisixConsumersGetter
	ApisixPluginConfigsGetter
	ApisixRoutesGetter
	ApisixTlsesGetter
	ApisixUpstreamsGetter
}

// ApisixV2Client is used to interact with features provided by the apisix.apache.org group.
type ApisixV2Client struct {
	restClient rest.Interface
}

func (c *ApisixV2Client) ApisixClusterConfigs() ApisixClusterConfigInterface {
	return newApisixClusterConfigs(c)
}

func (c *ApisixV2Client) ApisixConsumers(namespace string) ApisixConsumerInterface {
	return newApisixConsumers(c, namespace)
}

func (c *ApisixV2Client) ApisixPluginConfigs(namespace string) ApisixPluginConfigInterface {
	return newApisixPluginConfigs(c, namespace)
}

func (c *ApisixV2Client) ApisixRoutes(namespace string) ApisixRouteInterface {
	return newApisixRoutes(c, namespace)
}

func (c *ApisixV2Client) ApisixTlses(namespace string) ApisixTlsInterface {
	return newApisixTlses(c, namespace)
}

func (c *ApisixV2Client) ApisixUpstreams(namespace string) ApisixUpstreamInterface {
	return newApisixUpstreams(c, namespace)
}

// NewForConfig creates a new ApisixV2Client for the given config.
func NewForConfig(c *rest.Config) (*ApisixV2Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &ApisixV2Client{client}, nil
}

// NewForConfigOrDie creates a new ApisixV2Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *ApisixV2Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new ApisixV2Client for the given RESTClient.
func New(c rest.Interface) *ApisixV2Client {
	return &ApisixV2Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v2.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *ApisixV2Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
