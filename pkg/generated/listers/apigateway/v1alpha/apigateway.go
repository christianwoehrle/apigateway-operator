/*
Copyright The Kubernetes Authors.

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

package v1alpha

import (
	v1alpha "github.com/christianwoehrle/apigateway-operator/pkg/apis/apigateway/v1alpha"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ApiGatewayLister helps list ApiGateways.
type ApiGatewayLister interface {
	// List lists all ApiGateways in the indexer.
	List(selector labels.Selector) (ret []*v1alpha.ApiGateway, err error)
	// ApiGateways returns an object that can list and get ApiGateways.
	ApiGateways(namespace string) ApiGatewayNamespaceLister
	ApiGatewayListerExpansion
}

// apiGatewayLister implements the ApiGatewayLister interface.
type apiGatewayLister struct {
	indexer cache.Indexer
}

// NewApiGatewayLister returns a new ApiGatewayLister.
func NewApiGatewayLister(indexer cache.Indexer) ApiGatewayLister {
	return &apiGatewayLister{indexer: indexer}
}

// List lists all ApiGateways in the indexer.
func (s *apiGatewayLister) List(selector labels.Selector) (ret []*v1alpha.ApiGateway, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha.ApiGateway))
	})
	return ret, err
}

// ApiGateways returns an object that can list and get ApiGateways.
func (s *apiGatewayLister) ApiGateways(namespace string) ApiGatewayNamespaceLister {
	return apiGatewayNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ApiGatewayNamespaceLister helps list and get ApiGateways.
type ApiGatewayNamespaceLister interface {
	// List lists all ApiGateways in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha.ApiGateway, err error)
	// Get retrieves the ApiGateway from the indexer for a given namespace and name.
	Get(name string) (*v1alpha.ApiGateway, error)
	ApiGatewayNamespaceListerExpansion
}

// apiGatewayNamespaceLister implements the ApiGatewayNamespaceLister
// interface.
type apiGatewayNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ApiGateways in the indexer for a given namespace.
func (s apiGatewayNamespaceLister) List(selector labels.Selector) (ret []*v1alpha.ApiGateway, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha.ApiGateway))
	})
	return ret, err
}

// Get retrieves the ApiGateway from the indexer for a given namespace and name.
func (s apiGatewayNamespaceLister) Get(name string) (*v1alpha.ApiGateway, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha.Resource("apigateway"), name)
	}
	return obj.(*v1alpha.ApiGateway), nil
}