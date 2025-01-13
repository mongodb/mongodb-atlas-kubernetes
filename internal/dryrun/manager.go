package dryrun

import (
	"context"
	"net/http"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

type Manager struct {
	origin    manager.Manager
	transport *DryRunTransport
	k8sClient client.Client
}

func NewManager(origin manager.Manager, transport *DryRunTransport, k8sClient client.Client) *Manager {
	return &Manager{
		origin:    origin,
		transport: transport,
		k8sClient: k8sClient,
	}
}

func (m *Manager) GetHTTPClient() *http.Client {
	return m.origin.GetHTTPClient()
}

func (m *Manager) GetConfig() *rest.Config {
	return m.origin.GetConfig()
}

func (m *Manager) GetCache() cache.Cache {
	return m.origin.GetCache()
}

func (m *Manager) GetScheme() *runtime.Scheme {
	return m.origin.GetScheme()
}

func (m *Manager) GetClient() client.Client {
	return m.k8sClient
}

func (m *Manager) GetFieldIndexer() client.FieldIndexer {
	return m.origin.GetFieldIndexer()
}

func (m *Manager) GetEventRecorderFor(name string) record.EventRecorder {
	return m.origin.GetEventRecorderFor(name)
}

func (m *Manager) GetRESTMapper() meta.RESTMapper {
	return m.origin.GetRESTMapper()
}

func (m *Manager) GetAPIReader() client.Reader {
	return m.origin.GetAPIReader()
}

func (m *Manager) Start(ctx context.Context) error {
	return m.origin.Start(ctx)
}

func (m *Manager) Add(runnable manager.Runnable) error {
	return m.origin.Add(runnable)
}

func (m *Manager) Elected() <-chan struct{} {
	return m.origin.Elected()
}

func (m *Manager) AddMetricsServerExtraHandler(path string, handler http.Handler) error {
	return m.origin.AddMetricsServerExtraHandler(path, handler)
}

func (m *Manager) AddHealthzCheck(name string, check healthz.Checker) error {
	return m.origin.AddHealthzCheck(name, check)
}

func (m *Manager) AddReadyzCheck(name string, check healthz.Checker) error {
	return m.origin.AddReadyzCheck(name, check)
}

func (m *Manager) GetWebhookServer() webhook.Server {
	return m.origin.GetWebhookServer()
}

func (m *Manager) GetLogger() logr.Logger {
	return m.origin.GetLogger()
}

func (m *Manager) GetControllerOptions() config.Controller {
	return m.origin.GetControllerOptions()
}
