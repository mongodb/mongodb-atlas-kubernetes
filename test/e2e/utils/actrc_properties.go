package utils

import (
	"bufio"
	"os"
	"regexp"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
)

type Secrets map[string]string

// GetSecretsFromActrc act use .actrc file for define the secrets, we can use the same file
func GetSecretsFromActrc() (Secrets, error) {
	s := make(map[string]string)
	file, err := os.Open(config.ActrcPath)
	if err != nil {
		return s, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		re := regexp.MustCompile("^-s (?P<key>.+)=(?P<value>.+)$")
		if re.MatchString(line) {
			matchers := re.FindStringSubmatch(line)
			s[matchers[1]] = matchers[2]
		}
	}
	return s, err
}

func GetSecretEnvOrActrc(keys []string) Secrets {
	actrcSecrets, _ := GetSecretsFromActrc()
	finalSecrets := make(map[string]string)
	for _, key := range keys {
		value := os.Getenv(key)
		if value == "" {
			if actValue, ok := actrcSecrets[key]; ok {
				finalSecrets[key] = actValue
			}
		} else {
			finalSecrets[key] = value
		}
	}
	return finalSecrets
}
