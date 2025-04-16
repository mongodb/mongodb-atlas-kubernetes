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
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

func main() {
	if err := generateCert(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generateCert() error {
	cert, privatekey, publickey, err := utils.GenerateX509Cert()
	if err != nil {
		return err
	}

	parsedBasePath := flag.String("path", "tmp/x509/", "where to put the cert")
	flag.Parse()
	basePath := *parsedBasePath
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	if err := os.MkdirAll(basePath, fs.ModePerm); err != nil {
		return fmt.Errorf("failed to create path: %w", err)
	}

	// save private key
	pkey := x509.MarshalPKCS1PrivateKey(privatekey)
	pkeyPath := filepath.Join(basePath, "private.key")
	if err := os.WriteFile(pkeyPath, pkey, 0600); err != nil {
		return err
	}
	fmt.Println("private key saved to", pkeyPath)

	// save public key
	pubkey, _ := x509.MarshalPKIXPublicKey(publickey)
	pubkeyPath := filepath.Join(basePath, "public.key")
	if err := os.WriteFile(pubkeyPath, pubkey, 0600); err != nil {
		return err
	}
	fmt.Println("public key saved to", pubkeyPath)

	// this will create plain text PEM cert
	certPath := filepath.Join(basePath, "cert.pem")
	pemcert, err := os.Create(filepath.Clean(certPath))
	if err != nil {
		return err
	}
	var pemkey = &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	}
	if err := pem.Encode(pemcert, pemkey); err != nil {
		return err
	}
	if err := pemcert.Close(); err != nil {
		return err
	}
	fmt.Println("certificate saved to", certPath)

	return nil
}
