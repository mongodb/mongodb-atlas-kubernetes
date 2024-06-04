package launcher

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

const (
	ExpectedContext     = "kind-kind"
	LauncherTestInstall = "test-ako"
	HelmRepoURL         = "https://mongodb.github.io/helm-charts"
	RepoRef             = "mongodb"
	OperatorChart       = "mongodb-atlas-operator"
	AtlasURI            = "https://cloud-qa.mongodb.com"
	// #nosec G101 -- This is just a name
	AtlasSecretName = "mongodb-atlas-operator-api-key"
)

type Launcher struct {
	credentials   AtlasCredentials
	version       string
	cmdOutput     bool
	clearKind     bool
	clearOperator bool
	clearSecret   bool
	appliedYAMLs  []string
}

// NewLauncher creates a new operator Launcher
func NewLauncher(creds AtlasCredentials, version string, cmdOutput bool) *Launcher {
	return &Launcher{credentials: creds, version: version, cmdOutput: cmdOutput, appliedYAMLs: []string{}}
}

// NewFromEnv creates an operator Launcher using defaults and credentials from env vars
func NewFromEnv(version string) *Launcher {
	return NewLauncher(credentialsFromEnv(), version, true)
}

// MustLookupEnv sets the env var value given, or panics if the env var is not set
func MustLookupEnv(envvar string) string {
	value, ok := os.LookupEnv(envvar)
	if !ok {
		panic(fmt.Errorf("environment variable %s not set", envvar))
	}
	return value
}

// Launch will try to launch the operator and apply the given YAML for it to handle
func (l *Launcher) Launch(yml string, waitCfg *WaitConfig) error {
	if err := l.ensureK8sCluster(); err != nil {
		return fmt.Errorf("failed to setup Kubernetes cluster: %w", err)
	}
	if err := l.ensureOperator(); err != nil {
		return fmt.Errorf("failed to setup Atlas Kubernetes Operator: %w", err)
	}
	if err := l.ensureAtlasSecret(); err != nil {
		return fmt.Errorf("failed to setup Atlas Secrets: %w", err)
	}
	if err := l.kubeApply(yml); err != nil {
		return fmt.Errorf("failed to setup Kubernetes resources: %w", err)
	}
	l.appliedYAMLs = append(l.appliedYAMLs, yml)
	return l.kubeWait(waitCfg)
}

func (l *Launcher) Kubectl(args ...string) (string, error) {
	return l.run("kubectl", args...)
}

// Cleanup related Launcher test resources, merely wipe the kind cluster
func (l *Launcher) Cleanup() error {
	for _, yml := range l.appliedYAMLs {
		l.kubeDelete(yml)
	}
	if l.clearOperator {
		if err := l.uninstall(); err != nil {
			return fmt.Errorf("failed to uninstall Atlas operator: %w", err)
		}
	}
	if l.clearSecret {
		if err := l.kubeDeleteAtlasSecret(); err != nil {
			return fmt.Errorf("failed to delete Atlas secrets: %w", err)
		}
	}
	if l.clearKind {
		return l.stopKind()
	}
	return nil
}

func (l *Launcher) ensureK8sCluster() error {
	if !l.isKubeConfigAvailable() {
		return l.startKind()
	}
	return nil
}

func (l *Launcher) isKubeConfigAvailable() bool {
	out, err := l.run("kubectl", "config", "current-context")
	return err == nil && strings.Contains(out, ExpectedContext)
}

func (l *Launcher) startKind() error {
	err := l.silentRun("kind", "create", "cluster")
	if err != nil {
		return fmt.Errorf("failed to start kind: %w", err)
	}
	l.clearKind = true
	return nil
}

func (l *Launcher) stopKind() error {
	return l.silentRun("kind", "delete", "cluster")
}

func (l *Launcher) ensureOperator() error {
	if !l.isInstalled() {
		return l.install()
	}
	return nil
}

func (l *Launcher) isInstalled() bool {
	result, err := l.run("helm", "ls", "-a", "-A")
	if err != nil || !strings.Contains(result, LauncherTestInstall) {
		return false
	}
	scanner := bufio.NewScanner(strings.NewReader(result))
	for scanner.Scan() {
		str := scanner.Text()
		if strings.Contains(str, LauncherTestInstall) {
			parts := strings.Split(str, "\t")
			installedVersion := strings.TrimSpace(parts[len(parts)-1])
			return l.version == installedVersion
		}
	}
	return false
}

func (l *Launcher) install() error {
	err := l.silentRun("helm", "repo", "add", RepoRef, HelmRepoURL)
	update := false
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			update = true
		} else {
			return fmt.Errorf("failed to set up mongodb helm chart repo: %w", err)
		}
	}
	if update {
		if err := l.silentRun("helm", "repo", "update", RepoRef); err != nil {
			return fmt.Errorf("failed to update mongodb helm chart repo: %w", err)
		}
	}
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return fmt.Errorf("failed to set up mongodb helm chart repo: %w", err)
	}
	err = l.silentRun("helm", "install", LauncherTestInstall, path.Join(RepoRef, OperatorChart),
		"--version", l.version, "--atomic", "--set", fmt.Sprintf("atlasURI=%s", AtlasURI))
	if err != nil {
		return fmt.Errorf("failed to install operator: %w", err)
	}
	l.clearOperator = true
	return nil
}

func (l *Launcher) uninstall() error {
	return l.silentRun("helm", "uninstall", LauncherTestInstall)
}

func (l *Launcher) ensureAtlasSecret() error {
	if !l.atlasSecretAvailable() {
		l.clearSecret = true
		return l.kubeCreateAtlasSecret()
	}
	return nil
}

func (l *Launcher) atlasSecretAvailable() bool {
	return l.silentRun("kubectl", "get", "secret", AtlasSecretName) == nil
}

func (l *Launcher) kubeCreateAtlasSecret() error {
	if err := l.silentRun("kubectl", "create", "secret", "generic", AtlasSecretName,
		fmt.Sprintf("--from-literal=orgId=%s", l.credentials.OrgID),
		fmt.Sprintf("--from-literal=publicApiKey=%s", l.credentials.PublicKey),
		fmt.Sprintf("--from-literal=privateApiKey=%s", l.credentials.PrivateKey)); err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}
	return l.silentRun("kubectl", "label", "secret", AtlasSecretName, "atlas.mongodb.com/type=credentials")
}

func (l *Launcher) kubeDeleteAtlasSecret() error {
	return l.silentRun("kubectl", "delete", "secret", AtlasSecretName)
}

func (l *Launcher) kubeApply(yml string) error {
	return l.pipeRun(yml, "kubectl", "apply", "-f", "-")
}

func (l *Launcher) kubeDelete(yml string) error {
	return l.pipeRun(yml, "kubectl", "delete", "-f", "-")
}

func (l *Launcher) kubeWait(cfg *WaitConfig) error {
	if cfg == NoWait {
		return nil
	}
	return l.silentRun("kubectl", cfg.waitArgs()...)
}

func (l *Launcher) run(cmd string, args ...string) (string, error) {
	buf := bytes.NewBufferString("")
	err := run(nil, l.stdout(buf), l.stderr(buf), cmd, args...)
	return buf.String(), err
}

func (l *Launcher) silentRun(cmd string, args ...string) error {
	return silenceReturn(l.run(cmd, args...))
}

func (l *Launcher) pipeRun(stdin, cmd string, args ...string) error {
	buf := bytes.NewBufferString("")
	input := bytes.NewBufferString(stdin)
	err := run(input, l.stdout(buf), l.stderr(buf), cmd, args...)
	return silenceReturn(buf.String(), err)
}

func (l *Launcher) stdout(w io.Writer) io.Writer {
	if l.cmdOutput {
		return io.MultiWriter(w, os.Stdout)
	}
	return w
}

func (l *Launcher) stderr(w io.Writer) io.Writer {
	if l.cmdOutput {
		return io.MultiWriter(w, os.Stderr)
	}
	return w
}

func silenceReturn(msg string, err error) error {
	if err != nil {
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}
