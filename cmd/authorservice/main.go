package main

import (
    "log"
    "net"

    "google.golang.org/grpc"
    pb "library_management/proto/author"
    "library_management/internal/authorservice"
    "library_management/config"
)

func main() {
    db, err := config.ConnectDB(config.AuthorServiceDB)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    lis, err := net.Listen("tcp", ":50052")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterAuthorServiceServer(grpcServer, &authorservice.AuthorServiceServer{DB: db})

    log.Println("Starting AuthorService server on :50052")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
