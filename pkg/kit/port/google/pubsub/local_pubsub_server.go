package pubsub

import (
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"context"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

type LocalPubSubServer struct {
	server *pstest.Server
}

func NewLocalPubSubServer() *LocalPubSubServer {
	return &LocalPubSubServer{
		server: pstest.NewServer(),
	}
}

func (s *LocalPubSubServer) CreateClient() (*pubsub.Client, error) {
	conn, err := grpc.Dial(s.server.Addr, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to dial %s", s.server.Addr)
	}
	client, err := pubsub.NewClient(context.Background(), "project", option.WithGRPCConn(conn))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create a client")
	}
	return client, nil
}

func (s *LocalPubSubServer) TearDown() error {
	if err := s.server.Close(); err != nil {
		return errors.Wrap(err, "failed to close the server")
	}
	return nil
}
