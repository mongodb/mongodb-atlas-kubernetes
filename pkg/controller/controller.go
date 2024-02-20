package controller

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MultiNamespacedCacheBuilder returns a manager cache builder for a list of namespaces
func MultiNamespacedCacheBuilder(namespaces []string) cache.NewCacheFunc {
	return func(config *rest.Config, opts cache.Options) (cache.Cache, error) {
		opts.DefaultNamespaces = map[string]cache.Config{}
		for _, ns := range namespaces {
			opts.DefaultNamespaces[ns] = cache.Config{}
		}
		return cache.New(config, opts)
	}
}

// CustomLabelSelectorCacheBuilder returns a manager cache builder for a custom label selector
func CustomLabelSelectorCacheBuilder(obj client.Object, labelsSelector labels.Selector) cache.NewCacheFunc {
	return func(config *rest.Config, opts cache.Options) (cache.Cache, error) {
		if opts.ByObject == nil {
			opts.ByObject = map[client.Object]cache.ByObject{}
		}
		opts.ByObject[obj] = cache.ByObject{Label: labelsSelector}
		return cache.New(config, opts)
	}
}
