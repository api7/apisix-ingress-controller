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

	v2 "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/apis/config/v2"
	scheme "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ApisixTlsesGetter has a method to return a ApisixTlsInterface.
// A group's client should implement this interface.
type ApisixTlsesGetter interface {
	ApisixTlses(namespace string) ApisixTlsInterface
}

// ApisixTlsInterface has methods to work with ApisixTls resources.
type ApisixTlsInterface interface {
	Create(ctx context.Context, apisixTls *v2.ApisixTls, opts v1.CreateOptions) (*v2.ApisixTls, error)
	Update(ctx context.Context, apisixTls *v2.ApisixTls, opts v1.UpdateOptions) (*v2.ApisixTls, error)
	UpdateStatus(ctx context.Context, apisixTls *v2.ApisixTls, opts v1.UpdateOptions) (*v2.ApisixTls, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v2.ApisixTls, error)
	List(ctx context.Context, opts v1.ListOptions) (*v2.ApisixTlsList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v2.ApisixTls, err error)
	ApisixTlsExpansion
}

// apisixTlses implements ApisixTlsInterface
type apisixTlses struct {
	client rest.Interface
	ns     string
}

// newApisixTlses returns a ApisixTlses
func newApisixTlses(c *ApisixV2Client, namespace string) *apisixTlses {
	return &apisixTlses{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the apisixTls, and returns the corresponding apisixTls object, and an error if there is any.
func (c *apisixTlses) Get(ctx context.Context, name string, options v1.GetOptions) (result *v2.ApisixTls, err error) {
	result = &v2.ApisixTls{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("apisixtlses").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ApisixTlses that match those selectors.
func (c *apisixTlses) List(ctx context.Context, opts v1.ListOptions) (result *v2.ApisixTlsList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v2.ApisixTlsList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("apisixtlses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested apisixTlses.
func (c *apisixTlses) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("apisixtlses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a apisixTls and creates it.  Returns the server's representation of the apisixTls, and an error, if there is any.
func (c *apisixTlses) Create(ctx context.Context, apisixTls *v2.ApisixTls, opts v1.CreateOptions) (result *v2.ApisixTls, err error) {
	result = &v2.ApisixTls{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("apisixtlses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(apisixTls).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a apisixTls and updates it. Returns the server's representation of the apisixTls, and an error, if there is any.
func (c *apisixTlses) Update(ctx context.Context, apisixTls *v2.ApisixTls, opts v1.UpdateOptions) (result *v2.ApisixTls, err error) {
	result = &v2.ApisixTls{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("apisixtlses").
		Name(apisixTls.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(apisixTls).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *apisixTlses) UpdateStatus(ctx context.Context, apisixTls *v2.ApisixTls, opts v1.UpdateOptions) (result *v2.ApisixTls, err error) {
	result = &v2.ApisixTls{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("apisixtlses").
		Name(apisixTls.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(apisixTls).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the apisixTls and deletes it. Returns an error if one occurs.
func (c *apisixTlses) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("apisixtlses").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *apisixTlses) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("apisixtlses").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched apisixTls.
func (c *apisixTlses) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v2.ApisixTls, err error) {
	result = &v2.ApisixTls{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("apisixtlses").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
