package categoryservice

import (
    "context"
    "database/sql"

    pb "library_management/proto/category"
    "library_management/internal/common"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type CategoryServiceServer struct {
    DB *sql.DB
    pb.UnimplementedCategoryServiceServer
}

func (s *CategoryServiceServer) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
    _, err := s.DB.Exec("INSERT INTO categories (name) VALUES ($1)", req.GetName())
    if err != nil {
        return nil, status.Error(codes.Internal, "Failed to create category")
    }
    return &pb.CreateCategoryResponse{Success: true}, nil
}

func (s *CategoryServiceServer) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.GetCategoryResponse, error) {
    var category common.Category
    row := s.DB.QueryRow("SELECT id, name FROM categories WHERE id = $1", req.GetId())
    err := row.Scan(&category.ID, &category.Name)
    if err != nil {
        return nil, status.Error(codes.NotFound, "Category not found")
    }
    return &pb.GetCategoryResponse{
        Id:   int32(category.ID),
        Name: category.Name,
    }, nil
}

func (s *CategoryServiceServer) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.UpdateCategoryResponse, error) {
    _, err := s.DB.Exec("UPDATE categories SET name=$1 WHERE id=$2", req.GetName(), req.GetId())
    if err != nil {
        return nil, status.Error(codes.Internal, "Failed to update category")
    }
    return &pb.UpdateCategoryResponse{Success: true}, nil
}

func (s *CategoryServiceServer) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.DeleteCategoryResponse, error) {
    _, err := s.DB.Exec("DELETE FROM categories WHERE id = $1", req.GetId())
    if err != nil {
        return nil, status.Error(codes.Internal, "Failed to delete category")
    }
    return &pb.DeleteCategoryResponse{Success: true}, nil
}

func (s *CategoryServiceServer) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
    rows, err := s.DB.Query("SELECT id, name FROM categories")
    if err != nil {
        return nil, status.Error(codes.Internal, "Failed to list categories")
    }
    defer rows.Close()

    var categories []*pb.GetCategoryResponse
    for rows.Next() {
        var category common.Category
        if err := rows.Scan(&category.ID, &category.Name); err != nil {
            return nil, status.Error(codes.Internal, "Failed to read category data")
        }
        categories = append(categories, &pb.GetCategoryResponse{
            Id:   int32(category.ID),
            Name: category.Name,
        })
    }

    return &pb.ListCategoriesResponse{Categories: categories}, nil
}
