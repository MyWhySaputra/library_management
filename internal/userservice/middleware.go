package userservice

import (
    "context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func AdminInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    // Extract user role from context
    role, ok := ctx.Value("role").(string)
    if !ok || role != "admin" {
        return nil, status.Errorf(codes.PermissionDenied, "permission denied")
    }

    // Proceed to handler
    return handler(ctx, req)
}
