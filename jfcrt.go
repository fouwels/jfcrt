package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
)

// Just fork me a certificate
// Just generate a vanilla, self signed, non CA, x509 certificate, without extensions and all the openssl-is
// ...you probably need the extensions for your use case though
// Do not use self signed CA=TRUE (ala OpenSSL) certificates for mTLS in the browser... will cause you trust any external certificate signed with your "malign" test cert, bad idea.
// You can do this with openSSL, but it's a massive hassle.

func main() {
	err := run()
	log.Printf("exit with err: %v", err)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
func run() error {

	const length = 3072
	const days = 36500
	const signature = x509.SHA256WithRSA

	subject := flag.String("s", "", "Subject for Certificate")
	flag.Parse()

	if *subject == "" {
		return fmt.Errorf("-s subject not set")
	}

	key, err := rsa.GenerateKey(rand.Reader, length)
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	cert := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: *subject,
		},
		SignatureAlgorithm: signature,
		NotBefore:          time.Now(),
		NotAfter:           time.Now().Add(time.Hour * 24 * days),
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &cert, &cert, &key.PublicKey, key)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	out := &bytes.Buffer{}
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	f, err := os.Create(fmt.Sprintf("%v.crt", *subject))
	if err != nil {
		return fmt.Errorf("failed to create output certificate file: %w", err)
	}

	defer f.Close()
	_, err = f.Write(out.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write certificate to file: %w", err)
	}

	fp, err := os.Create(fmt.Sprintf("%v.pem", *subject))
	if err != nil {
		return fmt.Errorf("failed to create output key file: %w", err)
	}

	keypem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	defer fp.Close()
	_, err = fp.Write(keypem)
	if err != nil {
		return fmt.Errorf("failed to write key to file: %w", err)
	}

	return nil
}
