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

package v2beta3

import (
	"context"
	time "time"

	configv2beta3 "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/apis/config/v2beta3"
	versioned "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/client/clientset/versioned"
	internalinterfaces "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/client/informers/externalversions/internalinterfaces"
	v2beta3 "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/client/listers/config/v2beta3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ApisixRouteInformer provides access to a shared informer and lister for
// ApisixRoutes.
type ApisixRouteInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v2beta3.ApisixRouteLister
}

type apisixRouteInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewApisixRouteInformer constructs a new informer for ApisixRoute type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewApisixRouteInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredApisixRouteInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredApisixRouteInformer constructs a new informer for ApisixRoute type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredApisixRouteInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ApisixV2beta3().ApisixRoutes(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ApisixV2beta3().ApisixRoutes(namespace).Watch(context.TODO(), options)
			},
		},
		&configv2beta3.ApisixRoute{},
		resyncPeriod,
		indexers,
	)
}

func (f *apisixRouteInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredApisixRouteInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *apisixRouteInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&configv2beta3.ApisixRoute{}, f.defaultInformer)
}

func (f *apisixRouteInformer) Lister() v2beta3.ApisixRouteLister {
	return v2beta3.NewApisixRouteLister(f.Informer().GetIndexer())
}
