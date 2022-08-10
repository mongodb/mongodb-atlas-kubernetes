package networkpeer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const fileName = "role13.json"

type RoleAssignment struct {
	Id string `json:"id"`
}

type AzureRole struct {
	SubscriptionId    string
	ResourceGroupName string
	VnetName          string
}

func RoleName(subID, resGroupName, vnetName string) string {
	return fmt.Sprintf("AtlasPeering/%s/%s/%s", subID, resGroupName, vnetName)
}

type AzureRoleDefinition struct {
	Name             string   `json:"Name"`
	IsCustom         bool     `json:"IsCustom"`
	Description      string   `json:"Description"`
	Actions          []string `json:"Actions"`
	AssignableScopes []string `json:"AssignableScopes"`
}

func NewAzureRole(subID, resGroupName, vnetName string) AzureRoleDefinition {
	DescriptionTemplate := "Grants MongoDB access to manage peering connections on network " +
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network" +
		"/virtualNetworks/%s"
	return AzureRoleDefinition{
		Name:        RoleName(subID, resGroupName, vnetName),
		IsCustom:    true,
		Description: fmt.Sprintf(DescriptionTemplate, subID, resGroupName, vnetName),
		Actions: []string{
			"Microsoft.Network/virtualNetworks/virtualNetworkPeerings/read",
			"Microsoft.Network/virtualNetworks/virtualNetworkPeerings/write",
			"Microsoft.Network/virtualNetworks/virtualNetworkPeerings/delete",
			"Microsoft.Network/virtualNetworks/peer/action",
		},
		AssignableScopes: []string{
			fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s",
				subID, resGroupName, vnetName),
		},
	}
}

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

func CreateServicePrincipal(id string) error {
	cmd := exec.Command("az", "ad", "sp", "create", "--id", id)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "ERROR:") {
			return fmt.Errorf(scanner.Text())
		}
	}
	return nil
}

func AzureCreateRole(subID, resGroupName, vnetName string) error {
	role := NewAzureRole(subID, resGroupName, vnetName)
	json, err := json.Marshal(role)
	if err != nil {
		return err
	}
	roleFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer roleFile.Close()
	_, err = roleFile.Write(json)
	if err != nil {
		return err
	}

	cmd := exec.Command("az", "role", "definition", "create", "--role-definition", fileName)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "ERROR:") {
			return fmt.Errorf(scanner.Text())
		}
	}
	return nil
}

func AzureAssignRole(principalID string, role AzureRole) (string, error) {
	roleArgTemplate := "AtlasPeering/%s/%s/%s"
	roleArg := fmt.Sprintf(roleArgTemplate, role.SubscriptionId, role.ResourceGroupName, role.VnetName)
	assigneeTemplate := "%s"
	assigneeArg := fmt.Sprintf(assigneeTemplate, principalID)
	scopeTemplate := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s"
	scopeArg := fmt.Sprintf(scopeTemplate, role.SubscriptionId, role.ResourceGroupName, role.VnetName)

	cmd := exec.Command("az", "role", "assignment", "create",
		"--role", roleArg,
		"--assignee", assigneeArg,
		"--scope", scopeArg)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	roleAssignment := RoleAssignment{}
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "ERROR:") {
			return "", fmt.Errorf(scanner.Text())
		}
	}

	scanner = bufio.NewScanner(stdout)
	var response []byte
	for scanner.Scan() {
		response = append(response, scanner.Bytes()...)
	}
	err = json.Unmarshal(response, &roleAssignment)
	if err != nil {
		return "", err
	}

	return roleAssignment.Id, nil
}

func DeleteServicePrincipal(id string) error {
	cmd := exec.Command("az", "ad", "sp", "delete", "--id", id)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "ERROR:") {
			return fmt.Errorf(scanner.Text())
		}
	}
	return nil
}

func DeleteRole(role AzureRole) error {
	roleName := RoleName(role.SubscriptionId, role.ResourceGroupName, role.VnetName)
	cmd := exec.Command("az", "role", "definition", "delete", "--name", roleName)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "ERROR:") {
			return fmt.Errorf(scanner.Text())
		}
	}
	return nil
}

func DeleteRoleAssignment(principalID string, roleAssigmentId string) error {
	cmd := exec.Command("az", "role", "assignment", "delete", "--ids", roleAssigmentId)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "ERROR:") {
			return fmt.Errorf(scanner.Text())
		}
	}
	return nil
}
