package authorservice

import (
    "context"
    "database/sql"

    pb "library_management/proto/author"
    "library_management/internal/common"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type AuthorServiceServer struct {
    DB *sql.DB
    pb.UnimplementedAuthorServiceServer
}

func (s *AuthorServiceServer) CreateAuthor(ctx context.Context, req *pb.CreateAuthorRequest) (*pb.CreateAuthorResponse, error) {
    _, err := s.DB.Exec("INSERT INTO authors (name) VALUES ($1)", req.GetName())
    if err != nil {
        return nil, status.Error(codes.Internal, "Failed to create author")
    }
    return &pb.CreateAuthorResponse{Success: true}, nil
}

func (s *AuthorServiceServer) GetAuthor(ctx context.Context, req *pb.GetAuthorRequest) (*pb.GetAuthorResponse, error) {
    var author common.Author
    row := s.DB.QueryRow("SELECT id, name FROM authors WHERE id = $1", req.GetId())
    err := row.Scan(&author.ID, &author.Name)
    if err != nil {
        return nil, status.Error(codes.NotFound, "Author not found")
    }
    return &pb.GetAuthorResponse{
        Id:   int32(author.ID),
        Name: author.Name,
    }, nil
}

func (s *AuthorServiceServer) UpdateAuthor(ctx context.Context, req *pb.UpdateAuthorRequest) (*pb.UpdateAuthorResponse, error) {
    _, err := s.DB.Exec("UPDATE authors SET name=$1 WHERE id=$2", req.GetName(), req.GetId())
    if err != nil {
        return nil, status.Error(codes.Internal, "Failed to update author")
    }
    return &pb.UpdateAuthorResponse{Success: true}, nil
}

func (s *AuthorServiceServer) DeleteAuthor(ctx context.Context, req *pb.DeleteAuthorRequest) (*pb.DeleteAuthorResponse, error) {
    _, err := s.DB.Exec("DELETE FROM authors WHERE id = $1", req.GetId())
    if err != nil {
        return nil, status.Error(codes.Internal, "Failed to delete author")
    }
    return &pb.DeleteAuthorResponse{Success: true}, nil
}

func (s *AuthorServiceServer) ListAuthors(ctx context.Context, req *pb.ListAuthorsRequest) (*pb.ListAuthorsResponse, error) {
    rows, err := s.DB.Query("SELECT id, name FROM authors")
    if err != nil {
        return nil, status.Error(codes.Internal, "Failed to list authors")
    }
    defer rows.Close()

    var authors []*pb.GetAuthorResponse
    for rows.Next() {
        var author common.Author
        if err := rows.Scan(&author.ID, &author.Name); err != nil {
            return nil, status.Error(codes.Internal, "Failed to read author data")
        }
        authors = append(authors, &pb.GetAuthorResponse{
            Id:   int32(author.ID),
            Name: author.Name,
        })
    }

    return &pb.ListAuthorsResponse{Authors: authors}, nil
}
