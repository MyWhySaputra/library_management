package userservice

import (
    "context"
    "library_management/config"

    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "github.com/golang-jwt/jwt"
)

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, status.Errorf(codes.Unauthenticated, "missing context metadata")
    }

    token := md["authorization"]
    if len(token) == 0 {
        return nil, status.Errorf(codes.Unauthenticated, "missing token")
    }

    claims := &jwt.MapClaims{}
    _, err := jwt.ParseWithClaims(token[0], claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(config.JwtSecret), nil
    })

    if err != nil {
        return nil, status.Errorf(codes.Unauthenticated, "invalid token")
    }

    return handler(ctx, req)
}
