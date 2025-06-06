package snapshot

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/observability"
)

func Snapshot(logger io.Writer) error {
	kubeClient, err := kubecli.CreateNewClient()
	if err != nil {
		return fmt.Errorf("Failed to create client: %w", err)
	}

	n := &corev1.Node{}
	err = kubeClient.Get(context.Background(), client.ObjectKey{Name: "dos-control-plane"}, n)
	if err != nil {
		return fmt.Errorf("Failed to get dos-control-plane node: %w", err)
	}
	start := n.GetObjectMeta().GetCreationTimestamp()
	now := time.Now()

	promLayout := "2006-01-02+15:04:05"
	endString := strings.ReplaceAll(url.QueryEscape(now.UTC().Format(promLayout)), "%2B", "+")
	durationString := fmt.Sprintf("%dm", int(now.Sub(start.Time).Minutes()))
	promQuery := `http://localhost:30000/query?g0.expr=%7Bjob%3D%22ako%22%7D&g0.show_tree=0&g0.tab=graph&g0.end_input=` + endString + `&g0.range_input=` + durationString + `&g0.res_type=auto&g0.res_density=medium&g0.display_mode=lines&g0.show_exemplars=0`
	parcaQuery := `http://localhost:30003/?step_count=113&query_browser_mode=advanced&expression_a=goroutine%3Agoroutine%3Acount%3Agoroutine%3Acount%7Bjob%3D%22ako%22%7D&time_selection_a=absolute:` + fmt.Sprintf("%d", start.UnixMilli()) + `-` + fmt.Sprintf("%d", now.UnixMilli()) + `&sum_by_a=__none__`
	lokiQuery := `http://localhost:30001/explore?schemaVersion=1&panes=%7B%22ooz%22%3A%7B%22datasource%22%3A%22loki%22%2C%22queries%22%3A%5B%7B%22refId%22%3A%22A%22%2C%22expr%22%3A%22%7Bjob%3D%5C%22ako%5C%22%7D+%7C+json+%7C+line_format+%5C%22%7B%7B+.ts+%7D%7D+%5C%5C033%5B1%3B37m%7B%7B+.level+%7D%7D%5C%5C033%5B0m+%5C%5C033%5B1%3B32m%7B%7B+.logger+%7D%7D%5C%5C033%5B0m+%7B%7B+.msg+%7D%7D%5C%22%22%2C%22queryType%22%3A%22range%22%2C%22datasource%22%3A%7B%22type%22%3A%22loki%22%2C%22uid%22%3A%22loki%22%7D%2C%22editorMode%22%3A%22code%22%2C%22direction%22%3A%22backward%22%7D%5D%2C%22range%22%3A%7B%22from%22%3A%22` + fmt.Sprintf("%d", start.UnixMilli()) + `%22%2C%22to%22%3A%22` + fmt.Sprintf("%d", now.UnixMilli()) + `%22%7D%7D%7D&orgId=1`

	out := fmt.Sprintf(`Prometheus query:
%s

Loki query:
%s

Parca query:
%s
`, promQuery, lokiQuery, parcaQuery)
	//nolint:gosec
	err = os.WriteFile("urls.txt", []byte(out), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	for _, cmdArgs := range [][]string{
		// tell prometheus to take snapshot so WAL is flushed
		{"curl", "-XPOST", "-v", `http://localhost:30000/api/v1/admin/tsdb/snapshot`},
		{"sh", "-c", `kubectl exec -n monitoring prometheus-kube-prometheus-kube-prome-prometheus-0 -- tar cvzf - -C /prometheus . >prometheus.tar.gz`},
		{"sh", "-c", `kubectl exec -n loki -c loki loki-0 -- tar cvzf - -C /var/loki . >loki.tar.gz`},
		{"sh", "-c", `kubectl exec -n parca deployment/parca -- tar cvzf - -C /var/lib/parca . >parca.tar.gz`},
	} {
		if err := observability.ExecCommand(logger, cmdArgs...); err != nil {
			return err
		}
	}
	return nil
}
