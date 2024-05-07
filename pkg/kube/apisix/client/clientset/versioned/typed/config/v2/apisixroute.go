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
	"context"
	"time"

	v2 "github.com/api7/apisix-ingress-controller/pkg/kube/apisix/apis/config/v2"
	scheme "github.com/api7/apisix-ingress-controller/pkg/kube/apisix/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ApisixRoutesGetter has a method to return a ApisixRouteInterface.
// A group's client should implement this interface.
type ApisixRoutesGetter interface {
	ApisixRoutes(namespace string) ApisixRouteInterface
}

// ApisixRouteInterface has methods to work with ApisixRoute resources.
type ApisixRouteInterface interface {
	Create(ctx context.Context, apisixRoute *v2.ApisixRoute, opts v1.CreateOptions) (*v2.ApisixRoute, error)
	Update(ctx context.Context, apisixRoute *v2.ApisixRoute, opts v1.UpdateOptions) (*v2.ApisixRoute, error)
	UpdateStatus(ctx context.Context, apisixRoute *v2.ApisixRoute, opts v1.UpdateOptions) (*v2.ApisixRoute, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v2.ApisixRoute, error)
	List(ctx context.Context, opts v1.ListOptions) (*v2.ApisixRouteList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v2.ApisixRoute, err error)
	ApisixRouteExpansion
}

// apisixRoutes implements ApisixRouteInterface
type apisixRoutes struct {
	client rest.Interface
	ns     string
}

// newApisixRoutes returns a ApisixRoutes
func newApisixRoutes(c *ApisixV2Client, namespace string) *apisixRoutes {
	return &apisixRoutes{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the apisixRoute, and returns the corresponding apisixRoute object, and an error if there is any.
func (c *apisixRoutes) Get(ctx context.Context, name string, options v1.GetOptions) (result *v2.ApisixRoute, err error) {
	result = &v2.ApisixRoute{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("apisixroutes").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ApisixRoutes that match those selectors.
func (c *apisixRoutes) List(ctx context.Context, opts v1.ListOptions) (result *v2.ApisixRouteList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v2.ApisixRouteList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("apisixroutes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested apisixRoutes.
func (c *apisixRoutes) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("apisixroutes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a apisixRoute and creates it.  Returns the server's representation of the apisixRoute, and an error, if there is any.
func (c *apisixRoutes) Create(ctx context.Context, apisixRoute *v2.ApisixRoute, opts v1.CreateOptions) (result *v2.ApisixRoute, err error) {
	result = &v2.ApisixRoute{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("apisixroutes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(apisixRoute).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a apisixRoute and updates it. Returns the server's representation of the apisixRoute, and an error, if there is any.
func (c *apisixRoutes) Update(ctx context.Context, apisixRoute *v2.ApisixRoute, opts v1.UpdateOptions) (result *v2.ApisixRoute, err error) {
	result = &v2.ApisixRoute{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("apisixroutes").
		Name(apisixRoute.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(apisixRoute).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *apisixRoutes) UpdateStatus(ctx context.Context, apisixRoute *v2.ApisixRoute, opts v1.UpdateOptions) (result *v2.ApisixRoute, err error) {
	result = &v2.ApisixRoute{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("apisixroutes").
		Name(apisixRoute.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(apisixRoute).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the apisixRoute and deletes it. Returns an error if one occurs.
func (c *apisixRoutes) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("apisixroutes").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *apisixRoutes) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("apisixroutes").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched apisixRoute.
func (c *apisixRoutes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v2.ApisixRoute, err error) {
	result = &v2.ApisixRoute{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("apisixroutes").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
