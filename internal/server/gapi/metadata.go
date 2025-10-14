package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

func (s *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// log.Printf("metadata: %+v\n", md)

		if values := md.Get(grpcGatewayUserAgentHeader); len(values) > 0 {
			mtdt.UserAgent = values[0]
		}

		if values := md.Get(userAgentHeader); len(values) > 0 {
			mtdt.UserAgent = values[0]
		}

		if values := md.Get(xForwardedForHeader); len(values) > 0 {
			mtdt.ClientIP = values[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
