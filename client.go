package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/marabunta/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Client marabunta (ant)
type Client struct {
	client pb.MarabuntaClient
	ctx    context.Context
}

// New create ant
func New(id, host string, port int) (*Client, error) {
	// TODO create client metadata
	md := metadata.Pairs(
		"ant", id,
	)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	return &Client{ctx: ctx}, nil
}

// Run ant
func (c *Client) Run(cert string) error {
	creds, err := credentials.NewClientTLSFromFile("server.crt", "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}
	// wait for 5 seconds
	connCtx, cancel := context.WithTimeout(c.ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(connCtx, "localhost:1415",
		grpc.WithTransportCredentials(creds),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("%s unable to connect", err)
	}
	defer conn.Close()

	c.client = pb.NewMarabuntaClient(conn)

	// TODO check how to register if then ... send / receive ...

	stream, err := c.client.Stream(c.ctx)
	if err != nil {
		return err
	}
	c.Send(stream)
	return c.Receive(stream)
}

func (c *Client) Send(stream pb.Marabunta_StreamClient) {
	msg := &pb.StreamRequest{
		Msg: fmt.Sprintf("date: %s", time.Now().Format(time.RFC3339Nano)),
	}
	err := stream.Send(msg)
	if s, ok := status.FromError(err); ok && s.Code() == codes.Canceled {
		log.Println("stream canceled")
		return
	} else if err == io.EOF {
		log.Println("stream closed by server")
		return
	} else if err != nil {
		log.Println("send", err)
		return
	}
}

func (c *Client) Receive(stream pb.Marabunta_StreamClient) error {
	for {
		res, err := stream.Recv()
		if s, ok := status.FromError(err); ok && s.Code() == codes.Canceled {
			return fmt.Errorf("%s, stream canceled", err)
		} else if err == io.EOF {
			return fmt.Errorf("%s, stream closed by server", err)
		} else if err != nil {
			return err
		}

		switch evt := res.Event.(type) {
		case *pb.StreamResponse_EPing:
			log.Printf("ping response = %+v\n", evt.EPing.Msg)
			fmt.Printf("update response = %+v\n\n", c.Update("foo"))
		case *pb.StreamResponse_EPulse:
			log.Printf("pulse response = %+v\n", evt.EPulse.Msg)
		default:
			log.Printf("event = %+v\n", evt)
		}
	}
}

func (c *Client) Update(name string) bool {
	r, err := c.client.Update(c.ctx, &pb.UpdateRequest{Name: name})
	if err != nil {
		log.Fatal(err)
	}
	return r.Ok
}
