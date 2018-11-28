package ant

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func createCertificate(home string) error {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	id, err := GetID(filepath.Join(home, "ant.id"))
	if err != nil {
		return err
	}

	subj := pkix.Name{
		CommonName:         id,
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

	x509Encoded, _ := x509.MarshalECPrivateKey(key)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	err = ioutil.WriteFile(filepath.Join(home, "ant.key"), pemEncoded, 0600)
	if err != nil {
		return err
	}

	csr := pem.EncodeToMemory(&pem.Block{
		Type: "CERTIFICATE REQUEST", Bytes: csrCertificate,
	})

	err = ioutil.WriteFile(filepath.Join(home, "ant.csr"), csr, 0644)
	if err != nil {
		return err
	}

	// REQUEST HTTP to marabunta the sign the cert
	req, err := http.NewRequest("POST", "https://httpbin.org/post", bytes.NewBuffer(csr))
	req.Header.Set("User-Agent", fmt.Sprintf("ant-%s", id))
	if err != nil {
		return err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	fmt.Printf("res = %+v\n", res)
	if res.StatusCode == 200 {
		crt, err := os.Create(filepath.Join(home, "ant.crt"))
		if err != nil {
			return err
		}
		_, err = io.Copy(crt, res.Body)
		if err != nil {
			return err
		}
	}
	// TODO
	return nil
}
