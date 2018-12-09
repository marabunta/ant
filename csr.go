package ant

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func createCertificate(cfg *Config) error {
	var do bool
	if !isFile(cfg.TLS.Crt) {
		do = true
	}
	if !isFile(cfg.TLS.Key) {
		do = true
	}
	if !do {
		return nil
	}

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	id, err := GetID(filepath.Join(cfg.Home, "ant.id"))
	if err != nil {
		return err
	}

	subj := pkix.Name{
		CommonName:         id,
		Organization:       []string{"marabunta"},
		OrganizationalUnit: []string{"ant"},
	}

	asn1Subj, err := asn1.Marshal(subj.ToRDNSequence())
	if err != nil {
		return fmt.Errorf("unable to marshal asn1: %v", err)
	}

	template := x509.CertificateRequest{
		RawSubject:         asn1Subj,
		SignatureAlgorithm: x509.ECDSAWithSHA256,
	}

	csrCertificate, err := x509.CreateCertificateRequest(rand.Reader, &template, key)
	if err != nil {
		return fmt.Errorf("could not create CSR, %s", err)
	}

	x509Encoded, _ := x509.MarshalECPrivateKey(key)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	err = ioutil.WriteFile(cfg.TLS.Key, pemEncoded, 0600)
	if err != nil {
		return err
	}

	csr := pem.EncodeToMemory(&pem.Block{
		Type: "CERTIFICATE REQUEST", Bytes: csrCertificate,
	})

	err = ioutil.WriteFile(filepath.Join(cfg.Home, "ant.csr"), csr, 0644)
	if err != nil {
		return err
	}

	return RequestCertificate(
		fmt.Sprintf("https://%s:%d/register", cfg.Marabunta, cfg.HTTPPort),
		cfg.Home,
		id,
		csr)
}
