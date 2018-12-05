package ant

import (
	"bytes"
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"

	uuid "github.com/satori/go.uuid"
)

// isDir return true if path is a dir
func isDir(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	if m := f.Mode(); m.IsDir() && m&400 != 0 {
		return true
	}
	return false
}

// isFile return true if path is a regular file
func isFile(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	if m := f.Mode(); !m.IsDir() && m.IsRegular() && m&400 != 0 {
		return true
	}
	return false
}

func GetID(file string) (string, error) {
	if isFile(file) {
		id, err := ioutil.ReadFile(file)
		if err == nil {
			id = bytes.TrimSpace(id)
			if len(id) > 0 {
				if len(id) > 36 {
					return string(id[:36]), nil
				}
				return string(id), nil
			}
		}
	}
	uuid1, err := uuid.NewV1()
	if err != nil {
		return "", fmt.Errorf("could not create UUID, %s", err)
	}
	err = ioutil.WriteFile(file, []byte(uuid1.String()), 0644)
	if err != nil {
		return "", err
	}
	return uuid1.String(), nil
}

// GetHome returns the $HOME/.marabunta
func GetHome() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("error getting user home: %s", err)
		}
		home = usr.HomeDir
	}
	home = filepath.Join(home, ".marabunta")
	if err := os.MkdirAll(home, os.ModePerm); err != nil {
		return "", err
	}
	return home, nil
}

func RequestCertificate(url, home, id string, data []byte) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", fmt.Sprintf("ant-%s", id))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: tr,
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(io.LimitReader(res.Body, 4096))
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("%s", b)
	}

	block, _ := pem.Decode(b)
	if block == nil {
		return fmt.Errorf("failed to parse certificate PEM")
	}
	//cert, err := x509.ParseCertificate(block.Bytes)
	//if err != nil {
	//return fmt.Errorf("failed to parse certificate: %v", err.Error())
	//}
	crt := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: block.Bytes})

	fmt.Printf("crt = %s\n", crt)
	return nil
}
