package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	pb "github.com/marabunta/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//var version string

func main() {
	//var (
	//v        = flag.Bool("v", false, fmt.Sprintf("Print version: %s", version))
	//id       = flag.String("id", "", "ant ID")
	//host     = flag.String("host", "", "Connect to `host` default (marabunta)")
	//port     = flag.Int("port", 1415, "Port number to use for the connection")
	//certFile = flag.String("cert", "server.crt", "TLS cert")
	//)

	//flag.Parse()
	//if *v {
	//fmt.Printf("%s\n", version)
	//os.Exit(0)
	//}

	certificate, err := tls.LoadX509KeyPair(
		"./certs/client.crt",
		"./certs/client.key",
	)

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile("./certs/CA.crt")
	if err != nil {
		log.Fatalf("failed to read ca cert: %s", err)
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		log.Fatal("failed to append certs")
	}

	transportCreds := credentials.NewTLS(&tls.Config{
		ServerName:   "marabunta.host",
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	dialOption := grpc.WithTransportCredentials(transportCreds)
	conn, err := grpc.Dial("localhost:1415", dialOption)
	if err != nil {
		log.Fatalf("failed to dial server: %s", err)
	}
	defer conn.Close()

	client := pb.NewMarabuntaClient(conn)

	stream, err := client.Stream(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("stream = %+v\n", stream)
	for {

	}
}
