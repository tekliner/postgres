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

package v1alpha1

import (
	v1alpha1 "github.com/tekliner/apimachinery/apis/authorization/v1alpha1"
	scheme "github.com/tekliner/apimachinery/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MongoDBRolesGetter has a method to return a MongoDBRoleInterface.
// A group's client should implement this interface.
type MongoDBRolesGetter interface {
	MongoDBRoles(namespace string) MongoDBRoleInterface
}

// MongoDBRoleInterface has methods to work with MongoDBRole resources.
type MongoDBRoleInterface interface {
	Create(*v1alpha1.MongoDBRole) (*v1alpha1.MongoDBRole, error)
	Update(*v1alpha1.MongoDBRole) (*v1alpha1.MongoDBRole, error)
	UpdateStatus(*v1alpha1.MongoDBRole) (*v1alpha1.MongoDBRole, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.MongoDBRole, error)
	List(opts v1.ListOptions) (*v1alpha1.MongoDBRoleList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MongoDBRole, err error)
	MongoDBRoleExpansion
}

// mongoDBRoles implements MongoDBRoleInterface
type mongoDBRoles struct {
	client rest.Interface
	ns     string
}

// newMongoDBRoles returns a MongoDBRoles
func newMongoDBRoles(c *AuthorizationV1alpha1Client, namespace string) *mongoDBRoles {
	return &mongoDBRoles{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the mongoDBRole, and returns the corresponding mongoDBRole object, and an error if there is any.
func (c *mongoDBRoles) Get(name string, options v1.GetOptions) (result *v1alpha1.MongoDBRole, err error) {
	result = &v1alpha1.MongoDBRole{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mongodbroles").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MongoDBRoles that match those selectors.
func (c *mongoDBRoles) List(opts v1.ListOptions) (result *v1alpha1.MongoDBRoleList, err error) {
	result = &v1alpha1.MongoDBRoleList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mongodbroles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested mongoDBRoles.
func (c *mongoDBRoles) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("mongodbroles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a mongoDBRole and creates it.  Returns the server's representation of the mongoDBRole, and an error, if there is any.
func (c *mongoDBRoles) Create(mongoDBRole *v1alpha1.MongoDBRole) (result *v1alpha1.MongoDBRole, err error) {
	result = &v1alpha1.MongoDBRole{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("mongodbroles").
		Body(mongoDBRole).
		Do().
		Into(result)
	return
}

// Update takes the representation of a mongoDBRole and updates it. Returns the server's representation of the mongoDBRole, and an error, if there is any.
func (c *mongoDBRoles) Update(mongoDBRole *v1alpha1.MongoDBRole) (result *v1alpha1.MongoDBRole, err error) {
	result = &v1alpha1.MongoDBRole{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("mongodbroles").
		Name(mongoDBRole.Name).
		Body(mongoDBRole).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *mongoDBRoles) UpdateStatus(mongoDBRole *v1alpha1.MongoDBRole) (result *v1alpha1.MongoDBRole, err error) {
	result = &v1alpha1.MongoDBRole{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("mongodbroles").
		Name(mongoDBRole.Name).
		SubResource("status").
		Body(mongoDBRole).
		Do().
		Into(result)
	return
}

// Delete takes name of the mongoDBRole and deletes it. Returns an error if one occurs.
func (c *mongoDBRoles) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mongodbroles").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *mongoDBRoles) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mongodbroles").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched mongoDBRole.
func (c *mongoDBRoles) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MongoDBRole, err error) {
	result = &v1alpha1.MongoDBRole{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("mongodbroles").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
