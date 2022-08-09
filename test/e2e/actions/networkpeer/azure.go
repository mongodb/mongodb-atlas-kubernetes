package networkpeer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func AzureLogin() error {
	cmd := exec.Command("az", "account", "show")
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return err
	}
	scanner := bufio.NewScanner(stderr)
	errString := ""
	for scanner.Scan() {
		errString += scanner.Text()
	}
	if strings.Contains(errString, "Please run 'az login' to setup account.") {
		return nil
	}
	if strings.Contains(errString, "ERROR:") {
		return fmt.Errorf("cant see account info %s", errString)
	}

	clientID := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	tenantID := os.Getenv("AZURE_TENANT_ID")
	cmd = exec.Command("az", "login", "--service-principal", "-u", clientID, "-p", clientSecret, "--tenant", tenantID)
	stderr, _ = cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return err
	}

	scanner = bufio.NewScanner(stderr)
	errString = ""
	for scanner.Scan() {
		errString += scanner.Text()
	}
	if strings.Contains(errString, "ERROR:") {
		return fmt.Errorf("cant login %s", errString)
	}
	return nil
}

type azVersion struct {
	AzureCli          string `json:"azure-cli"`
	AzureCliCore      string `json:"azure-cli-core"`
	AzureCliTelemetry string `json:"azure-cli-telemetry"`
}

func CheckAZ() (bool, error) {
	cmd := exec.Command("az", "version")
	stderr, _ := cmd.StderrPipe()
	stdoud, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return false, err
	}

	scanner := bufio.NewScanner(stderr)
	errString := ""
	for scanner.Scan() {
		errString += scanner.Text()
	}
	if errString != "" {
		return false, fmt.Errorf(errString)
	}

	scanner = bufio.NewScanner(stdoud)
	response := []byte{}
	for scanner.Scan() {
		response = append(response, scanner.Bytes()...)
	}
	version := azVersion{}
	err := json.Unmarshal(response, &version)
	if err != nil {
		return false, err
	}

	return true, nil
}
