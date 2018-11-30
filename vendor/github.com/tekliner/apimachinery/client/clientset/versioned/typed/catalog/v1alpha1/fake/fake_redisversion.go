/*
Copyright 2018 The KubeDB Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/tekliner/apimachinery/apis/catalog/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeRedisVersions implements RedisVersionInterface
type FakeRedisVersions struct {
	Fake *FakeCatalogV1alpha1
}

var redisversionsResource = schema.GroupVersionResource{Group: "catalog.kubedb.com", Version: "v1alpha1", Resource: "redisversions"}

var redisversionsKind = schema.GroupVersionKind{Group: "catalog.kubedb.com", Version: "v1alpha1", Kind: "RedisVersion"}

// Get takes name of the redisVersion, and returns the corresponding redisVersion object, and an error if there is any.
func (c *FakeRedisVersions) Get(name string, options v1.GetOptions) (result *v1alpha1.RedisVersion, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(redisversionsResource, name), &v1alpha1.RedisVersion{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RedisVersion), err
}

// List takes label and field selectors, and returns the list of RedisVersions that match those selectors.
func (c *FakeRedisVersions) List(opts v1.ListOptions) (result *v1alpha1.RedisVersionList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(redisversionsResource, redisversionsKind, opts), &v1alpha1.RedisVersionList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.RedisVersionList{ListMeta: obj.(*v1alpha1.RedisVersionList).ListMeta}
	for _, item := range obj.(*v1alpha1.RedisVersionList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested redisVersions.
func (c *FakeRedisVersions) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(redisversionsResource, opts))
}

// Create takes the representation of a redisVersion and creates it.  Returns the server's representation of the redisVersion, and an error, if there is any.
func (c *FakeRedisVersions) Create(redisVersion *v1alpha1.RedisVersion) (result *v1alpha1.RedisVersion, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(redisversionsResource, redisVersion), &v1alpha1.RedisVersion{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RedisVersion), err
}

// Update takes the representation of a redisVersion and updates it. Returns the server's representation of the redisVersion, and an error, if there is any.
func (c *FakeRedisVersions) Update(redisVersion *v1alpha1.RedisVersion) (result *v1alpha1.RedisVersion, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(redisversionsResource, redisVersion), &v1alpha1.RedisVersion{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RedisVersion), err
}

// Delete takes name of the redisVersion and deletes it. Returns an error if one occurs.
func (c *FakeRedisVersions) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(redisversionsResource, name), &v1alpha1.RedisVersion{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRedisVersions) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(redisversionsResource, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.RedisVersionList{})
	return err
}

// Patch applies the patch and returns the patched redisVersion.
func (c *FakeRedisVersions) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.RedisVersion, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(redisversionsResource, name, data, subresources...), &v1alpha1.RedisVersion{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RedisVersion), err
}