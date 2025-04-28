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

package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrorEmptyPEMKey when the PEM key was not passed in
	ErrorEmptyPEMKey = errors.New("empty PEM key or filename")

	// ErrorNoPEMData when the PEM data was not found
	ErrorNoPEMData = errors.New("no PEM data found")

	// ErrorEmptyAppId when the App ID was not passed in
	ErrorEmptyAppId = errors.New("empty App Id")
)

type JWTSpec struct {
	// Base64PEMBytes is the base64 encoding of the key in PEM format
	Base64PEMBytes string
	// AppID is the OAuth Application id
	AppID string
	// Duration is how long we want the JWT to be valid for
	Duration time.Duration
	// Raw means the CLI outputs the raw JWT and no other debug info, such as the public key
	Raw bool
}

type jsonReply struct {
	PublicKey string `json:"publicKey"`
	JWT       string `json:"jwt"`
}

func makeJWT(spec *JWTSpec) (string, error) {
	// decode base64 of the key to PEM
	pemBytes, err := base64.StdEncoding.DecodeString(string(spec.Base64PEMBytes))
	if err != nil {
		return "", fmt.Errorf("error decoding base64 of the PEM key: %w", err)
	}

	// Parse PEM block
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return "", ErrorNoPEMData
	}

	// Parse RSA private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("error parsing private key: %w", err)
	}

	// Create JWT claims
	now := time.Now()
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(spec.Duration).Unix(),
		"iss": spec.AppID,
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the token with the private key
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}
	if spec.Raw {
		return tokenString, nil
	}
	pubKey, ok := (privateKey.Public()).(*rsa.PublicKey)
	if !ok {
		return "", errors.New("error expected RSA public key")
	}
	return buildJsonReply(pubKey, tokenString)
}

func buildJsonReply(publicKey *rsa.PublicKey, jwt string) (string, error) {
	pubKeyDER, err := asn1.Marshal(*publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key into ASN1: %w", err)
	}
	pubKeyPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubKeyDER})
	jsonBytes, err := json.Marshal(jsonReply{PublicKey: string(pubKeyPem), JWT: jwt})
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key into PEM: %w", err)
	}
	return string(jsonBytes), nil
}

func printJWT(w io.Writer, spec *JWTSpec) error {
	newJWT, err := makeJWT(spec)
	if err != nil {
		return fmt.Errorf("failed to create JWT from spec: %w", err)
	}
	_, err = fmt.Fprintln(w, newJWT)
	return err
}

func parseJWTSpecArgs(args []string) (*JWTSpec, error) {
	var spec JWTSpec
	fs := flag.NewFlagSet("jwtSpecFlags", flag.ContinueOnError)
	fs.StringVar(&spec.Base64PEMBytes, "key", "", "Base64 of the PEM key contents")
	fs.StringVar(&spec.AppID, "appId", "", "Application ID for the JWT token")
	fs.DurationVar(&spec.Duration, "duration", 10*time.Minute, "Duration the JWT token will be valid for")
	fs.BoolVar(&spec.Raw, "raw", true, "Emit the raw JWT or a JSON with the jwt and the public key certificate")
	if err := fs.Parse(args); err != nil {
		fs.PrintDefaults()
		return nil, fmt.Errorf("error parsing command line arguments: %w", err)
	}
	if spec.Base64PEMBytes == "" {
		return nil, ErrorEmptyPEMKey
	}
	if spec.AppID == "" {
		return nil, ErrorEmptyAppId
	}
	return &spec, nil
}

func main() {
	spec, err := parseJWTSpecArgs(os.Args[1:])
	if err != nil {
		log.Fatalf("Failed to parse input arguments: %v", err)
	}
	if err := printJWT(os.Stdout, spec); err != nil {
		log.Fatalf("Failed to create JWT token: %v", err)
	}
}
