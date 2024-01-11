package manager

import (
	"flag"
	"os"
	"strings"
	"time"
)

const (
	LeaderElectionDefault         = true
	MetricsBindAddressDefault     = ":8080"
	HealthProbeBindAddressDefault = ":8081"
	SyncPeriodDefault             = 3 * time.Hour
)

type Config struct {
	Namespaces             []string
	EnableLeaderElection   bool
	MetricsBindAddress     string
	HealthProbeBindAddress string
	SyncPeriod             time.Duration
}

func (c *Config) ParseWatchedNamespaces() {
	// dev note: we pass the watched namespace as the env variable to use the Kubernetes Downward API. Unfortunately
	// there is no way to use it for container arguments
	watchedNamespace := os.Getenv("WATCH_NAMESPACE")

	c.Namespaces = make([]string, 0, 1)

	for _, namespace := range strings.Split(watchedNamespace, ",") {
		c.Namespaces = append(c.Namespaces, strings.TrimSpace(namespace))
	}
}

func DefaultConfig() Config {
	return Config{
		EnableLeaderElection:   LeaderElectionDefault,
		MetricsBindAddress:     MetricsBindAddressDefault,
		HealthProbeBindAddress: HealthProbeBindAddressDefault,
		SyncPeriod:             SyncPeriodDefault,
	}
}

func RegisterFlags(conf *Config, fs *flag.FlagSet) {
	fs.StringVar(&conf.MetricsBindAddress, "metrics-bind-address", MetricsBindAddressDefault, "The address the metric endpoint binds to.")
	fs.StringVar(&conf.HealthProbeBindAddress, "health-probe-bind-address", HealthProbeBindAddressDefault, "The address the probe endpoint binds to.")
	fs.BoolVar(&conf.EnableLeaderElection, "leader-elect", LeaderElectionDefault,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
}
