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

	"github.com/satori/go.uuid"
)

func createCertificate(home string) error {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	// create UUID1
	uuid1, err := uuid.NewV1()
	if err != nil {
		return fmt.Errorf("could not create UUID, %s", err)
	}

	// save UUID1 in ant.id
	err = ioutil.WriteFile(filepath.Join(home, "ant.id"), []byte(uuid1.String()), 0644)
	if err != nil {
		return err
	}

	subj := pkix.Name{
		CommonName:         uuid1.String(),
		Country:            []string{"-"},
		Province:           []string{"-"},
		Locality:           []string{"-"},
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

	csr := pem.EncodeToMemory(&pem.Block{
		Type: "CERTIFICATE REQUEST", Bytes: csrCertificate,
	})

	err = ioutil.WriteFile(filepath.Join(home, "ant.csr"), csr, 0644)
	if err != nil {
		return err
	}

	x509Encoded, _ := x509.MarshalECPrivateKey(key)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	err = ioutil.WriteFile(filepath.Join(home, "ant.key"), pemEncoded, 0600)
	if err != nil {
		return err
	}

	// TODO REQUEST HTTP to marabunta the sign the cert
	return nil
}
