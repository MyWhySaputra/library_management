package main

import (
    "log"
    "net"

    "google.golang.org/grpc"
    pb "library_management/proto/user"
    "library_management/internal/userservice"
    "library_management/config"
)

func main() {
    db, err := config.ConnectDB(config.UserServiceDB)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    lis, err := net.Listen("tcp", ":50054")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(userservice.AuthInterceptor),
    )
    pb.RegisterUserServiceServer(grpcServer, &userservice.UserServiceServer{DB: db})

    log.Println("Starting UserService server on :50054")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
