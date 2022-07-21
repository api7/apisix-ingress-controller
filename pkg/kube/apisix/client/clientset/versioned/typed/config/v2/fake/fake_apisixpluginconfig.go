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

package fake

import (
	"context"

	v2 "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/apis/config/v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeApisixPluginConfigs implements ApisixPluginConfigInterface
type FakeApisixPluginConfigs struct {
	Fake *FakeApisixV2
	ns   string
}

var apisixpluginconfigsResource = schema.GroupVersionResource{Group: "apisix.apache.org", Version: "v2", Resource: "apisixpluginconfigs"}

var apisixpluginconfigsKind = schema.GroupVersionKind{Group: "apisix.apache.org", Version: "v2", Kind: "ApisixPluginConfig"}

// Get takes name of the apisixPluginConfig, and returns the corresponding apisixPluginConfig object, and an error if there is any.
func (c *FakeApisixPluginConfigs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v2.ApisixPluginConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(apisixpluginconfigsResource, c.ns, name), &v2.ApisixPluginConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ApisixPluginConfig), err
}

// List takes label and field selectors, and returns the list of ApisixPluginConfigs that match those selectors.
func (c *FakeApisixPluginConfigs) List(ctx context.Context, opts v1.ListOptions) (result *v2.ApisixPluginConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(apisixpluginconfigsResource, apisixpluginconfigsKind, c.ns, opts), &v2.ApisixPluginConfigList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v2.ApisixPluginConfigList{ListMeta: obj.(*v2.ApisixPluginConfigList).ListMeta}
	for _, item := range obj.(*v2.ApisixPluginConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested apisixPluginConfigs.
func (c *FakeApisixPluginConfigs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(apisixpluginconfigsResource, c.ns, opts))

}

// Create takes the representation of a apisixPluginConfig and creates it.  Returns the server's representation of the apisixPluginConfig, and an error, if there is any.
func (c *FakeApisixPluginConfigs) Create(ctx context.Context, apisixPluginConfig *v2.ApisixPluginConfig, opts v1.CreateOptions) (result *v2.ApisixPluginConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(apisixpluginconfigsResource, c.ns, apisixPluginConfig), &v2.ApisixPluginConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ApisixPluginConfig), err
}

// Update takes the representation of a apisixPluginConfig and updates it. Returns the server's representation of the apisixPluginConfig, and an error, if there is any.
func (c *FakeApisixPluginConfigs) Update(ctx context.Context, apisixPluginConfig *v2.ApisixPluginConfig, opts v1.UpdateOptions) (result *v2.ApisixPluginConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(apisixpluginconfigsResource, c.ns, apisixPluginConfig), &v2.ApisixPluginConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ApisixPluginConfig), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeApisixPluginConfigs) UpdateStatus(ctx context.Context, apisixPluginConfig *v2.ApisixPluginConfig, opts v1.UpdateOptions) (*v2.ApisixPluginConfig, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(apisixpluginconfigsResource, "status", c.ns, apisixPluginConfig), &v2.ApisixPluginConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ApisixPluginConfig), err
}

// Delete takes name of the apisixPluginConfig and deletes it. Returns an error if one occurs.
func (c *FakeApisixPluginConfigs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(apisixpluginconfigsResource, c.ns, name, opts), &v2.ApisixPluginConfig{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeApisixPluginConfigs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(apisixpluginconfigsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v2.ApisixPluginConfigList{})
	return err
}

// Patch applies the patch and returns the patched apisixPluginConfig.
func (c *FakeApisixPluginConfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v2.ApisixPluginConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(apisixpluginconfigsResource, c.ns, name, pt, data, subresources...), &v2.ApisixPluginConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ApisixPluginConfig), err
}
