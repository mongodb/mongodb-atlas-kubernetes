package kube

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func ObjectKey(namespace, name string) client.ObjectKey {
	return types.NamespacedName{Name: name, Namespace: namespace}
}

func ObjectKeyFromObject(obj metav1.Object) client.ObjectKey {
	return ObjectKey(obj.GetNamespace(), obj.GetName())
}

var (
	ResultRetry   = reconcile.Result{RequeueAfter: time.Second * 10}
	ResultSuccess = reconcile.Result{}
)
