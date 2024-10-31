package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

func applyYAML(config *rest.Config, filePath string) error {
	// Read YAML file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}

	// Decode YAML into unstructured.Unstructured
	decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	_, gvk, err := decoder.Decode(data, nil, obj)
	if err != nil {
		return fmt.Errorf("error decoding YAML object: %v", err)
	}

	// Create a RESTMapper to find GVR
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return fmt.Errorf("error creating discovery client: %v", err)
	}
	resources, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		return fmt.Errorf("error getting API group resources: %v", err)
	}
	mapper := restmapper.NewDiscoveryRESTMapper(resources)

	// Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return fmt.Errorf("error finding REST mapping: %v", err)
	}

	// Create a dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("error creating dynamic client: %v", err)
	}

	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		if obj.GetNamespace() == "" {
			obj.SetNamespace("default")
		}
		dr = dynamicClient.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		dr = dynamicClient.Resource(mapping.Resource)
	}

	// Create or Update the resource
	ctx := context.TODO()
	_, err = dr.Create(ctx, obj, metav1.CreateOptions{})
	if errors.IsAlreadyExists(err) {
		fmt.Printf("Resource %s already exists, updating...\n", obj.GetName())
		_, err = dr.Update(ctx, obj, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("error updating resource: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("error creating resource: %v", err)
	}

	return nil
}
