package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if err := generateCert(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generateCert() error {
	// see Certificate structure at
	// http://golang.org/pkg/crypto/x509/#Certificate
	template := &x509.Certificate{
		IsCA:                  true,
		BasicConstraintsValid: true,
		SubjectKeyId:          []byte{1, 2, 3},
		SerialNumber:          big.NewInt(1234),
		Issuer: pkix.Name{
			CommonName: "x509-testing-user",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(5, 5, 5),
		DNSNames:  []string{"x509-testing-user"},
		// see http://golang.org/pkg/crypto/x509/#KeyUsage
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	// generate private key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	publickey := &privatekey.PublicKey

	// create a self-signed certificate. template = parent
	var parent = template
	cert, err := x509.CreateCertificate(rand.Reader, template, parent, publickey, privatekey)
	if err != nil {
		return err
	}

	parsedBasePath := flag.String("path", "tmp/x509/", "where to put the cert")
	flag.Parse()
	basePath := *parsedBasePath
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create path: %w", err)
	}

	// save private key
	pkey := x509.MarshalPKCS1PrivateKey(privatekey)
	pkeyPath := filepath.Join(basePath, "private.key")
	if err := ioutil.WriteFile(pkeyPath, pkey, 0600); err != nil {
		return err
	}
	fmt.Println("private key saved to", pkeyPath)

	// save public key
	pubkey, _ := x509.MarshalPKIXPublicKey(publickey)
	pubkeyPath := filepath.Join(basePath, "public.key")
	if err := ioutil.WriteFile(pubkeyPath, pubkey, 0600); err != nil {
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
