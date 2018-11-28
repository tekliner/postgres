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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/tekliner/apimachinery/apis/authorization/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// MySQLRoleLister helps list MySQLRoles.
type MySQLRoleLister interface {
	// List lists all MySQLRoles in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.MySQLRole, err error)
	// MySQLRoles returns an object that can list and get MySQLRoles.
	MySQLRoles(namespace string) MySQLRoleNamespaceLister
	MySQLRoleListerExpansion
}

// mySQLRoleLister implements the MySQLRoleLister interface.
type mySQLRoleLister struct {
	indexer cache.Indexer
}

// NewMySQLRoleLister returns a new MySQLRoleLister.
func NewMySQLRoleLister(indexer cache.Indexer) MySQLRoleLister {
	return &mySQLRoleLister{indexer: indexer}
}

// List lists all MySQLRoles in the indexer.
func (s *mySQLRoleLister) List(selector labels.Selector) (ret []*v1alpha1.MySQLRole, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.MySQLRole))
	})
	return ret, err
}

// MySQLRoles returns an object that can list and get MySQLRoles.
func (s *mySQLRoleLister) MySQLRoles(namespace string) MySQLRoleNamespaceLister {
	return mySQLRoleNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// MySQLRoleNamespaceLister helps list and get MySQLRoles.
type MySQLRoleNamespaceLister interface {
	// List lists all MySQLRoles in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.MySQLRole, err error)
	// Get retrieves the MySQLRole from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.MySQLRole, error)
	MySQLRoleNamespaceListerExpansion
}

// mySQLRoleNamespaceLister implements the MySQLRoleNamespaceLister
// interface.
type mySQLRoleNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all MySQLRoles in the indexer for a given namespace.
func (s mySQLRoleNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.MySQLRole, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.MySQLRole))
	})
	return ret, err
}

// Get retrieves the MySQLRole from the indexer for a given namespace and name.
func (s mySQLRoleNamespaceLister) Get(name string) (*v1alpha1.MySQLRole, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("mysqlrole"), name)
	}
	return obj.(*v1alpha1.MySQLRole), nil
}
