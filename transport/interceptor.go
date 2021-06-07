package transport

import (
	"context"
	"errors"
	"log"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AuthenticationInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("invalid metadata")
	}

	if len(md.Get("token")) == 0 {
		return nil, errors.New("unauthorized")
	}

	token := md.Get("token")[0]

	log.Println(token)

	return handler(ctx, req)
}
