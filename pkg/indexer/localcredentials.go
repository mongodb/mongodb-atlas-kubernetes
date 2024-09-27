package indexer

import (
	"context"
	"reflect"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
)

type LocalCredentialIndexer struct {
	obj    client.Object
	name   string
	logger *zap.SugaredLogger
}

func NewLocalCredentialsIndexer(name string, obj client.Object, logger *zap.Logger) *LocalCredentialIndexer {
	return &LocalCredentialIndexer{
		obj:    obj,
		name:   name,
		logger: logger.Named(name).Sugar(),
	}
}

func (lc *LocalCredentialIndexer) Object() client.Object {
	return lc.obj
}

func (lc *LocalCredentialIndexer) Name() string {
	return lc.name
}

func (lc *LocalCredentialIndexer) Keys(object client.Object) []string {
	if reflect.TypeOf(object) != reflect.TypeOf(lc.obj) {
		lc.logger.Errorf("expected %T but got %T", lc.obj, object)
		return nil
	}

	credentialProvider, ok := (object).(api.CredentialsProvider)
	if !ok {
		lc.logger.Errorf("expected %T to implement api.CredentialProvider", object)
		return nil
	}

	if localRef := credentialProvider.Credentials(); localRef != nil && localRef.Name != "" {
		return []string{types.NamespacedName{Namespace: object.GetNamespace(), Name: localRef.Name}.String()}
	}
	return []string{}
}

type requestsFunc[L client.ObjectList] func(L) []reconcile.Request

func CredentialsIndexMapperFunc[L client.ObjectList](indexerName string, list L, reqsFn requestsFunc[L], kubeClient client.Client, logger *zap.SugaredLogger) handler.MapFunc {
	return func(ctx context.Context, obj client.Object) []reconcile.Request {
		secret, ok := obj.(*corev1.Secret)
		if !ok {
			logger.Warnf("watching Secret but got %T", obj)
			return nil
		}

		listOpts := &client.ListOptions{
			FieldSelector: fields.OneTermEqualSelector(
				indexerName,
				client.ObjectKeyFromObject(secret).String(),
			),
		}
		err := kubeClient.List(ctx, list, listOpts)
		if err != nil {
			logger.Errorf("failed to list from indexer %s: %v", indexerName, err)
			return nil
		}
		return reqsFn(list)
	}
}

// ToRequest is a helper to turns CRD objects into reconcile requests.
// Most Reconciliable implementations may leverage it.
func toRequest(obj client.Object) reconcile.Request {
	return reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      obj.GetName(),
			Namespace: obj.GetNamespace(),
		},
	}
}
