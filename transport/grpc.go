package transport

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/eventflowdb/eventflowdb/api"
	"github.com/eventflowdb/eventflowdb/constants"
	"github.com/eventflowdb/eventflowdb/env"
	"github.com/eventflowdb/eventflowdb/store"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type EventStore struct {
	api.UnimplementedEventStoreServer
	eventstore store.EventStore
}

func (s *EventStore) AppendStream(ctx context.Context, req *api.AppendToStreamRequest) (*api.AppendToStreamResponse, error) {
	return s.eventstore.AppendToStream(req)
}

func (s *EventStore) GetStream(ctx context.Context, req *api.GetStreamRequest) (*api.GetStreamResponse, error) {
	return s.eventstore.GetStream(req)
}
func (s *EventStore) GetGlobalStream(ctx context.Context, req *api.GetGlobalStreamRequest) (*api.GetGlobalStreamResponse, error) {
	return s.eventstore.GetGlobalStream(req)
}

func (s *EventStore) EventCount(ctx context.Context, req *api.EventCountRequest) (*api.EventCountResponse, error) {
	return s.eventstore.EventCount(req)
}

func (s *EventStore) StreamCount(ctx context.Context, req *api.StreamCountRequest) (*api.StreamCountResponse, error) {
	return s.eventstore.StreamCount(req)
}

func (s *EventStore) ListStreams(ctx context.Context, req *api.ListStreamsRequest) (*api.ListStreamsReponse, error) {
	return s.eventstore.ListStreams(req)
}

func (s *EventStore) Size(ctx context.Context, req *api.SizeRequest) (*api.SizeResponse, error) {
	return s.eventstore.Size(req)
}

var (
	start = time.Now()
)

func (s *EventStore) Uptime(context.Context, *api.UptimeRequest) (*api.UptimeResponse, error) {
	uptime := time.Since(start)

	return &api.UptimeResponse{
		Uptime:      uptime.Milliseconds(),
		UptimeHuman: uptime.String(),
	}, nil
}

func (s *EventStore) Version(context.Context, *api.VersionRequest) (*api.VersionResponse, error) {
	return &api.VersionResponse{
		Version: constants.Version,
	}, nil
}

func RunGRPCServer(eventstore store.EventStore, logger *logrus.Logger) *grpc.Server {
	grpcPort := env.GetEnv("GRPC_PORT", "6543")
	tlsEnabled := env.GetEnv("TLS_ENABLED", "false") == "true"
	certFile := env.GetEnv("TLS_CERT_FILE", "certs/crt.pem")
	keyFile := env.GetEnv("TLS_KEY_FILE", "certs/key.pem")

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		logger.Fatal(err)
	}

	grpcOptions := []grpc.ServerOption{}

	if tlsEnabled {
		logger.Println("tls is enabled")

		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			logger.Fatal(err)
		}

		config := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.NoClientCert,
		}

		grpcOptions = append(grpcOptions, grpc.Creds(credentials.NewTLS(config)))
	}

	server := grpc.NewServer(grpcOptions...)

	api.RegisterEventStoreServer(server, &EventStore{eventstore: eventstore})

	go func() {
		logger.Printf("gRPC server listening on %s", grpcPort)

		server.Serve(lis)
	}()

	return server
}
