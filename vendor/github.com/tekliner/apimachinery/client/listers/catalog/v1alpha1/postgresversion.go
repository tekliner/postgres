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
	v1alpha1 "github.com/tekliner/apimachinery/apis/catalog/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// PostgresVersionLister helps list PostgresVersions.
type PostgresVersionLister interface {
	// List lists all PostgresVersions in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.PostgresVersion, err error)
	// Get retrieves the PostgresVersion from the index for a given name.
	Get(name string) (*v1alpha1.PostgresVersion, error)
	PostgresVersionListerExpansion
}

// postgresVersionLister implements the PostgresVersionLister interface.
type postgresVersionLister struct {
	indexer cache.Indexer
}

// NewPostgresVersionLister returns a new PostgresVersionLister.
func NewPostgresVersionLister(indexer cache.Indexer) PostgresVersionLister {
	return &postgresVersionLister{indexer: indexer}
}

// List lists all PostgresVersions in the indexer.
func (s *postgresVersionLister) List(selector labels.Selector) (ret []*v1alpha1.PostgresVersion, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.PostgresVersion))
	})
	return ret, err
}

// Get retrieves the PostgresVersion from the index for a given name.
func (s *postgresVersionLister) Get(name string) (*v1alpha1.PostgresVersion, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("postgresversion"), name)
	}
	return obj.(*v1alpha1.PostgresVersion), nil
}
