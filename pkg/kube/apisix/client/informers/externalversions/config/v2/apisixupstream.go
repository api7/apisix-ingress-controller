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

// Code generated by informer-gen. DO NOT EDIT.

package v2

import (
	"context"
	time "time"

	configv2 "github.com/api7/apisix-ingress-controller/pkg/kube/apisix/apis/config/v2"
	versioned "github.com/api7/apisix-ingress-controller/pkg/kube/apisix/client/clientset/versioned"
	internalinterfaces "github.com/api7/apisix-ingress-controller/pkg/kube/apisix/client/informers/externalversions/internalinterfaces"
	v2 "github.com/api7/apisix-ingress-controller/pkg/kube/apisix/client/listers/config/v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ApisixUpstreamInformer provides access to a shared informer and lister for
// ApisixUpstreams.
type ApisixUpstreamInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v2.ApisixUpstreamLister
}

type apisixUpstreamInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewApisixUpstreamInformer constructs a new informer for ApisixUpstream type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewApisixUpstreamInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredApisixUpstreamInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredApisixUpstreamInformer constructs a new informer for ApisixUpstream type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredApisixUpstreamInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ApisixV2().ApisixUpstreams(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ApisixV2().ApisixUpstreams(namespace).Watch(context.TODO(), options)
			},
		},
		&configv2.ApisixUpstream{},
		resyncPeriod,
		indexers,
	)
}

func (f *apisixUpstreamInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredApisixUpstreamInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *apisixUpstreamInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&configv2.ApisixUpstream{}, f.defaultInformer)
}

func (f *apisixUpstreamInformer) Lister() v2.ApisixUpstreamLister {
	return v2.NewApisixUpstreamLister(f.Informer().GetIndexer())
}
