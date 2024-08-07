package main

import (
    "log"
    "net"

    "google.golang.org/grpc"
    pb "library_management/proto/category"
    "library_management/internal/categoryservice"
    "library_management/config"
)

func main() {
    db, err := config.ConnectDB(config.CategoryServiceDB)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    lis, err := net.Listen("tcp", ":50053")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterCategoryServiceServer(grpcServer, &categoryservice.CategoryServiceServer{DB: db})

    log.Println("Starting CategoryService server on :50053")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
