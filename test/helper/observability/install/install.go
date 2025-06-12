// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package install

import (
	"context"
	"flag"
	"fmt"
	"io"
	url "net/url"
	"os"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	clientgoretry "k8s.io/client-go/util/retry"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/observability"
)

var defaultBackOff = wait.Backoff{
	Duration: 2 * time.Second,
	Factor:   1.0,
	Steps:    30 * 5,
}

func InstallSnapshot(logger io.Writer) error {
	var (
		snapshotURL       string
		lokiSnapshotPath  string
		promSnapshotPath  string
		parcaSnapshotPath string
	)

	flag.StringVar(&snapshotURL, "snapshot-url", "", "The snapshot URL to download the snapshots from. If set, it has precedence over -loki-snapshot and -prom-snapshot.")
	flag.StringVar(&lokiSnapshotPath, "loki-snapshot", "", "The path to the loki snapshot .tar.gz file to use. If set, -prom-snapshot and -parca-snapshot must be provided.")
	flag.StringVar(&promSnapshotPath, "prom-snapshot", "", "The path to the prometheus snapshot .tar.gz file to use. If set, -loki-snapshot and -parca-snapshot must be provided.")
	flag.StringVar(&parcaSnapshotPath, "parca-snapshot", "", "The path to the parca snapshot .tar.gz file to use. If set, -loki-snapshot and -prom-snapshot must be provided.")
	err := flag.CommandLine.Parse(os.Args[2:])
	if err != nil {
		return fmt.Errorf("error parsing flags: %w", err)
	}

	switch {
	case snapshotURL != "":
		u, err := url.Parse(snapshotURL)
		if err != nil {
			return fmt.Errorf("failed to parse snapshot url: %w", err)
		}
		snapshotURL = u.String()
	case lokiSnapshotPath != "" && promSnapshotPath != "":
		var err error
		lokiSnapshotPath, err = filepath.Abs(lokiSnapshotPath)
		if err != nil {
			return fmt.Errorf("error getting absolute path for loki snapshot: %w", err)
		}
		promSnapshotPath, err = filepath.Abs(promSnapshotPath)
		if err != nil {
			return fmt.Errorf("error getting absolute path for prometheus snapshot: %w", err)
		}
		parcaSnapshotPath, err = filepath.Abs(parcaSnapshotPath)
		if err != nil {
			return fmt.Errorf("error getting absolute path for parca snapshot: %w", err)
		}
	default:
		return fmt.Errorf("either -snapshot-url or both -loki-snapshot and -prom-snapshot must be provided")
	}

	assetsDir, err := Unpack()
	if err != nil {
		return fmt.Errorf("error unpacking assets: %w", err)
	}
	defer os.RemoveAll(assetsDir)

	crdDir := filepath.Join(assetsDir, "config", "crd")
	assetsDir = filepath.Join(assetsDir, "test", "helper", "observability", "install", "assets")
	for _, cmdArgs := range [][]string{
		{"kind", "create", "cluster", `--name=dos`, fmt.Sprintf(`--config=%v/kind-config.yaml`, assetsDir)},
		{"kubectl", "apply", "-k", crdDir},
	} {
		if err := observability.ExecCommand(logger, cmdArgs...); err != nil {
			return err
		}
	}

	for _, cmdArgs := range [][]string{
		{"kubectl", "create", "ns", "parca"},
		{"kubectl", "apply", "--server-side", "-f", "https://github.com/parca-dev/parca/releases/download/v0.23.1/kubernetes-manifest.yaml"},
		{"kubectl", "apply", "--server-side", "-f", "https://github.com/parca-dev/parca-agent/releases/download/v0.39.0/kubernetes-manifest.yaml"},
		{"kubectl", "-n", "parca", "scale", "--replicas=0", "deployment/parca"},
		{"helm", "repo", "add", "prometheus-community", "https://prometheus-community.github.io/helm-charts"},
		{"helm", "repo", "add", "grafana", "https://grafana.github.io/helm-charts"},
		{"helm", "repo", "update"},
		{"helm", "upgrade", "--values", fmt.Sprintf("%v/kube-prometheus-helm.yaml", assetsDir), "--install", "kube-prometheus", "prometheus-community/kube-prometheus-stack", "-n", "monitoring", "--create-namespace"},
		{"helm", "upgrade", "--values", fmt.Sprintf("%v/loki-helm.yaml", assetsDir), "--install", "loki", "grafana/loki", "-n", "loki", "--create-namespace"},
		{"kubectl", "apply", "-f", fmt.Sprintf("%v/nodeports.yaml", assetsDir)},
		{"kubectl", "apply", "-f", fmt.Sprintf("%v/grafana-config.yaml", assetsDir)},
		{"kubectl", "-n", "loki", "scale", "--replicas=0", "sts/loki"},
	} {
		if err := observability.ExecCommand(logger, cmdArgs...); err != nil {
			return err
		}
	}

	if snapshotURL != "" {
		for _, cmdArgs := range [][]string{
			{"kubectl", "-n", "monitoring", "delete", "--ignore-not-found", "configmap", "artifact-urls"},
			{"kubectl", "-n", "monitoring", "create", "configmap", "artifact-urls",
				fmt.Sprintf(`--from-literal=PROMETHEUS_SNAPSHOT_URL='%v/prometheus.tar.gz'`, snapshotURL),
			},
			{"kubectl", "apply", "--server-side", "--force-conflicts", "-f", fmt.Sprintf("%v/prometheus-snapshot-url.yaml", assetsDir)},

			{"kubectl", "-n", "loki", "delete", "--ignore-not-found", "configmap", "artifact-urls"},
			{"kubectl", "-n", "loki", "create", "configmap", "artifact-urls",
				fmt.Sprintf(`--from-literal=LOKI_SNAPSHOT_URL='%v/loki.tar.gz'`, snapshotURL),
			},
			{"kubectl", "apply", "--server-side", "--force-conflicts", "-f", fmt.Sprintf("%v/loki-snapshot-url.yaml", assetsDir)},
		} {
			if err := observability.ExecCommand(logger, cmdArgs...); err != nil {
				return err
			}
		}
	} else {
		for _, cmdArgs := range [][]string{
			{"docker", "cp", lokiSnapshotPath, "dos-control-plane:/home/loki.tar.gz"},
			{"kubectl", "apply", "--server-side", "--force-conflicts", "-f", fmt.Sprintf("%v/loki-snapshot-file.yaml", assetsDir)},

			{"docker", "cp", promSnapshotPath, "dos-control-plane:/home/prometheus.tar.gz"},
			{"kubectl", "apply", "--server-side", "--force-conflicts", "-f", fmt.Sprintf("%v/prometheus-snapshot-file.yaml", assetsDir)},

			{"docker", "cp", parcaSnapshotPath, "dos-control-plane:/home/parca.tar.gz"},
			{"kubectl", "-n", "parca", "patch", "deployment/parca", "--patch-file", fmt.Sprintf("%v/parca-snapshot-file.yaml", assetsDir)},
		} {
			if err := observability.ExecCommand(logger, cmdArgs...); err != nil {
				return err
			}
		}
	}

	for _, cmdArgs := range [][]string{
		{"kubectl", "-n", "parca", "scale", "--replicas=1", "deployment/parca"},
		{"kubectl", "-n", "loki", "scale", "--replicas=1", "sts/loki"},
		{"kubectl", "-n", "loki", "wait", "pods", "-l", `app.kubernetes.io/name=loki`, "--for", "condition=Ready", "--timeout=600s"},
		{"kubectl", "-n", "monitoring", "wait", "pods", "-l", `app.kubernetes.io/instance=kube-prometheus-kube-prome-prometheus`, "--for", "condition=Ready", "--timeout=600s"},
		// flush loki, as it was disconnected.
		{"curl", "-XPOST", "-v", `http://localhost:30002/flush`},
	} {
		if err := observability.ExecCommand(logger, cmdArgs...); err != nil {
			return err
		}
	}

	return nil
}

func Teardown(logger io.Writer) error {
	if err := os.Chdir("test/helper/observability/install"); err != nil {
		return fmt.Errorf("error changing directory: %w", err)
	}

	for _, cmdArgs := range [][]string{
		{"kind", "delete", "cluster", `--name=dos`},
	} {
		if err := observability.ExecCommand(logger, cmdArgs...); err != nil {
			return err
		}
	}
	return nil
}

func Install(logger io.Writer) error {
	assetsDir, err := Unpack()
	if err != nil {
		return fmt.Errorf("error unpacking assets: %w", err)
	}
	defer os.RemoveAll(assetsDir)

	crdDir := filepath.Join(assetsDir, "config", "crd")
	assetsDir = filepath.Join(assetsDir, "test", "helper", "observability", "install", "assets")

	for _, cmdArgs := range [][]string{
		{"kind", "create", "cluster", `--name=dos`, fmt.Sprintf(`--config=%v/kind-config.yaml`, assetsDir)},
		{"kubectl", "apply", "-k", crdDir},
	} {
		if err := observability.ExecCommand(logger, cmdArgs...); err != nil {
			return err
		}
	}

	ctx := context.Background()
	for _, cmdArgs := range [][]string{
		{"kubectl", "create", "ns", "parca"},
		{"kubectl", "apply", "--server-side", "-f", "https://github.com/parca-dev/parca/releases/download/v0.23.1/kubernetes-manifest.yaml"},
		{"kubectl", "apply", "--server-side", "-f", "https://github.com/parca-dev/parca-agent/releases/download/v0.39.0/kubernetes-manifest.yaml"},
		{"kubectl", "-n", "parca", "scale", "--replicas=0", "deployment/parca"},
		{"kubectl", "-n", "parca", "delete", "configmap", "parca"},
		{"kubectl", "-n", "parca", "create", "configmap", "parca", fmt.Sprintf("--from-file=parca.yaml=%v/parca-config.yaml", assetsDir)},
		{"kubectl", "-n", "parca", "patch", "deployment/parca", "--patch-file", fmt.Sprintf("%v/parca-deployment.yaml", assetsDir)},
		{"kubectl", "-n", "parca", "scale", "--replicas=1", "deployment/parca"},
		{"helm", "repo", "add", "prometheus-community", "https://prometheus-community.github.io/helm-charts"},
		{"helm", "repo", "add", "grafana", "https://grafana.github.io/helm-charts"},
		{"helm", "repo", "update"},
		{"helm", "upgrade", "--values", fmt.Sprintf("%v/kube-prometheus-helm.yaml", assetsDir), "--install", "kube-prometheus", "prometheus-community/kube-prometheus-stack", "-n", "monitoring", "--create-namespace"},
		{"helm", "upgrade", "--values", fmt.Sprintf("%v/loki-helm.yaml", assetsDir), "--install", "loki", "grafana/loki", "-n", "loki", "--create-namespace"},
		{"helm", "upgrade", "--values", fmt.Sprintf("%v/promtail-helm.yaml", assetsDir), "--install", "promtail", "grafana/promtail", "-n", "promtail", "--create-namespace"},
		{"kubectl", "apply", "-f", fmt.Sprintf("%v/nodeports.yaml", assetsDir)},
		{"kubectl", "apply", "-f", fmt.Sprintf("%v/grafana-config.yaml", assetsDir)},
		{"kubectl", "-n", "monitoring", "create", "secret", "generic", "host-scrape-config", fmt.Sprintf("--from-file=%v/prometheus-host-scrape-config.yaml", assetsDir)},
		{"kubectl", "apply", "--server-side", "--force-conflicts", "-f", fmt.Sprintf("%v/prometheus.yaml", assetsDir)},
	} {
		if err := observability.ExecCommand(logger, cmdArgs...); err != nil {
			return err
		}
	}

	err = retry(ctx, func() error {
		return observability.ExecCommand(logger, "kubectl", "-n", "monitoring", "scale", "--replicas=0", "deployment/kube-prometheus-kube-state-metrics")
	})
	if err != nil {
		return fmt.Errorf("error executing command: %w", err)
	}

	for _, cmdArgs := range [][]string{
		{"kubectl", "-n", "monitoring", "create", "configmap", "kube-state-metrics-config", fmt.Sprintf("--from-file=%v/ksm-config.yaml", assetsDir)},
		{"kubectl", "apply", "--server-side", `--field-manager="dos"`, "--force-conflicts", "-f", fmt.Sprintf("%v/ksm-deployment.yaml", assetsDir)},
		{"kubectl", "apply", "--server-side", "-f", fmt.Sprintf("%v/ksm-cluster-role-binding.yaml", assetsDir)},
		{"kubectl", "-n", "monitoring", "scale", "--replicas=1", "deployment/kube-prometheus-kube-state-metrics"},
		{"kubectl", "-n", "loki", "rollout", "status", "--watch", "statefulset/loki"},
		{"kubectl", "-n", "promtail", "rollout", "status", "--watch", "deployment/promtail"},
		{"kubectl", "-n", "monitoring", "rollout", "status", "--watch", "deployment/kube-prometheus-kube-state-metrics"},
		{"kubectl", "-n", "monitoring", "rollout", "status", "--watch", "statefulset/prometheus-kube-prometheus-kube-prome-prometheus"},
	} {
		if err := observability.ExecCommand(logger, cmdArgs...); err != nil {
			return err
		}
	}
	return nil
}

func retry(ctx context.Context, f func() error) error {
	return clientgoretry.OnError(
		defaultBackOff, func(err error) bool { return true },
		func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			err := f()
			if err != nil {
				fmt.Println(err)
			}
			return err
		},
	)
}
