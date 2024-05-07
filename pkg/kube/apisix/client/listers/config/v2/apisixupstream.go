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

// Code generated by lister-gen. DO NOT EDIT.

package v2

import (
	v2 "github.com/api7/apisix-ingress-controller/pkg/kube/apisix/apis/config/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ApisixUpstreamLister helps list ApisixUpstreams.
// All objects returned here must be treated as read-only.
type ApisixUpstreamLister interface {
	// List lists all ApisixUpstreams in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v2.ApisixUpstream, err error)
	// ApisixUpstreams returns an object that can list and get ApisixUpstreams.
	ApisixUpstreams(namespace string) ApisixUpstreamNamespaceLister
	ApisixUpstreamListerExpansion
}

// apisixUpstreamLister implements the ApisixUpstreamLister interface.
type apisixUpstreamLister struct {
	indexer cache.Indexer
}

// NewApisixUpstreamLister returns a new ApisixUpstreamLister.
func NewApisixUpstreamLister(indexer cache.Indexer) ApisixUpstreamLister {
	return &apisixUpstreamLister{indexer: indexer}
}

// List lists all ApisixUpstreams in the indexer.
func (s *apisixUpstreamLister) List(selector labels.Selector) (ret []*v2.ApisixUpstream, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v2.ApisixUpstream))
	})
	return ret, err
}

// ApisixUpstreams returns an object that can list and get ApisixUpstreams.
func (s *apisixUpstreamLister) ApisixUpstreams(namespace string) ApisixUpstreamNamespaceLister {
	return apisixUpstreamNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ApisixUpstreamNamespaceLister helps list and get ApisixUpstreams.
// All objects returned here must be treated as read-only.
type ApisixUpstreamNamespaceLister interface {
	// List lists all ApisixUpstreams in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v2.ApisixUpstream, err error)
	// Get retrieves the ApisixUpstream from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v2.ApisixUpstream, error)
	ApisixUpstreamNamespaceListerExpansion
}

// apisixUpstreamNamespaceLister implements the ApisixUpstreamNamespaceLister
// interface.
type apisixUpstreamNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ApisixUpstreams in the indexer for a given namespace.
func (s apisixUpstreamNamespaceLister) List(selector labels.Selector) (ret []*v2.ApisixUpstream, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v2.ApisixUpstream))
	})
	return ret, err
}

// Get retrieves the ApisixUpstream from the indexer for a given namespace and name.
func (s apisixUpstreamNamespaceLister) Get(name string) (*v2.ApisixUpstream, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v2.Resource("apisixupstream"), name)
	}
	return obj.(*v2.ApisixUpstream), nil
}
