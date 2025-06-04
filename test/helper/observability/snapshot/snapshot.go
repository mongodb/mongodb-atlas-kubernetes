package snapshot

import (
	"io"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/observability"
)

func Snapshot(logger io.Writer) error {
	for _, cmdArgs := range [][]string{
		// tell prometheus to take snapshot so WAL is flushed
		{"curl", "-XPOST", "-v", `http://localhost:30000/api/v1/admin/tsdb/snapshot`},
		{"sh", "-c", `kubectl exec -n monitoring prometheus-kube-prometheus-kube-prome-prometheus-0 -- tar cvzf - -C /prometheus . >prometheus.tar.gz`},
		{"sh", "-c", `kubectl exec -n loki -c loki loki-0 -- tar cvzf - -C /var/loki . >loki.tar.gz`},
	} {
		if err := observability.ExecCommand(logger, cmdArgs...); err != nil {
			return err
		}
	}
	return nil
}
