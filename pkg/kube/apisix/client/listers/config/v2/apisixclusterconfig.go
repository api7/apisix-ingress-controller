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

// ApisixClusterConfigLister helps list ApisixClusterConfigs.
// All objects returned here must be treated as read-only.
type ApisixClusterConfigLister interface {
	// List lists all ApisixClusterConfigs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v2.ApisixClusterConfig, err error)
	// Get retrieves the ApisixClusterConfig from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v2.ApisixClusterConfig, error)
	ApisixClusterConfigListerExpansion
}

// apisixClusterConfigLister implements the ApisixClusterConfigLister interface.
type apisixClusterConfigLister struct {
	indexer cache.Indexer
}

// NewApisixClusterConfigLister returns a new ApisixClusterConfigLister.
func NewApisixClusterConfigLister(indexer cache.Indexer) ApisixClusterConfigLister {
	return &apisixClusterConfigLister{indexer: indexer}
}

// List lists all ApisixClusterConfigs in the indexer.
func (s *apisixClusterConfigLister) List(selector labels.Selector) (ret []*v2.ApisixClusterConfig, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v2.ApisixClusterConfig))
	})
	return ret, err
}

// Get retrieves the ApisixClusterConfig from the index for a given name.
func (s *apisixClusterConfigLister) Get(name string) (*v2.ApisixClusterConfig, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v2.Resource("apisixclusterconfig"), name)
	}
	return obj.(*v2.ApisixClusterConfig), nil
}
