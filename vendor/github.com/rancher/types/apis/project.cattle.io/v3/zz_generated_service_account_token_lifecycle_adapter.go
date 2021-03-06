package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type ServiceAccountTokenLifecycle interface {
	Create(obj *ServiceAccountToken) (runtime.Object, error)
	Remove(obj *ServiceAccountToken) (runtime.Object, error)
	Updated(obj *ServiceAccountToken) (runtime.Object, error)
}

type serviceAccountTokenLifecycleAdapter struct {
	lifecycle ServiceAccountTokenLifecycle
}

func (w *serviceAccountTokenLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*ServiceAccountToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *serviceAccountTokenLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*ServiceAccountToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *serviceAccountTokenLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*ServiceAccountToken))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewServiceAccountTokenLifecycleAdapter(name string, clusterScoped bool, client ServiceAccountTokenInterface, l ServiceAccountTokenLifecycle) ServiceAccountTokenHandlerFunc {
	adapter := &serviceAccountTokenLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *ServiceAccountToken) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
