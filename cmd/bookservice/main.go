package main

import (
    "log"
    "net"

    "library_management/config"
    "library_management/internal/bookservice"
    pb "library_management/proto/book"

    "google.golang.org/grpc"
)

func main() {
    db, err := config.ConnectDB(config.BookServiceDB)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    redisClient := config.ConnectRedis()

    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    bookServiceServer := bookservice.NewBookServiceServer(db, redisClient)
    pb.RegisterBookServiceServer(grpcServer, bookServiceServer)

    log.Println("BookService is running on port 50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
