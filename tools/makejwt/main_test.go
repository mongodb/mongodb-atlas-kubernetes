package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"encoding/base64"
	"testing"

	"github.com/golang-jwt/jwt"
)

func generateRandomRSAKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	return privateKey, nil
}

func asPEM(keyType string, key *rsa.PrivateKey) []byte {
	// Encode private key to PEM format
	privateKeyPEMBlock := &pem.Block{
		Type:  keyType,
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	return pem.EncodeToMemory(privateKeyPEMBlock)
}

func testSpec(appId, b64Key string) *JWTSpec {
	return &JWTSpec{
		AppID:          appId,
		Base64PEMBytes: b64Key,
		Raw:            true,
		Duration:       10 * time.Minute,
	}
}

func TestPrintJWT(t *testing.T) {
	key, err := generateRandomRSAKey()
	if err != nil {
		t.Fatal(err)
	}
	b64Key := base64.StdEncoding.EncodeToString(asPEM("RSA PRIVATE KEY", key))
	buf := bytes.NewBufferString("")
	if err := printJWT(buf, testSpec("123456", b64Key)); err != nil {
		t.Fatal(err)
	}
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return key.Public(), nil
	}
	if _, err := jwt.Parse(buf.String(), keyFunc); err != nil {
		t.Fatal(err)
	}
}

func TestParseJWTSpecArgsErrors(t *testing.T) {
	testCases := []struct {
		title         string
		args          []string
		expectedError error
	}{
		{
			title:         "No args",
			args:          []string{},
			expectedError: ErrorEmptyPEMKey,
		},
		{
			title:         "Missing PEM key",
			args:          []string{"-appId=123456"},
			expectedError: ErrorEmptyPEMKey,
		},
		{
			title:         "Missing App Id",
			args:          []string{"-key=fake"},
			expectedError: ErrorEmptyAppId,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			_, err := parseJWTSpecArgs(tc.args)
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("got %v want %v", err, tc.expectedError)
			}
		})
	}
}
