package userservice

import (
    "context"
    "database/sql"
    "time"

    pb "library_management/proto/user"
    "library_management/internal/common"
    "library_management/config"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "github.com/golang-jwt/jwt"
)

type UserServiceServer struct {
    DB *sql.DB
    pb.UnimplementedUserServiceServer
}

func (s *UserServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
    hashedPassword, err := common.HashPassword(req.GetPassword())
    if err != nil {
        return nil, status.Error(codes.Internal, "Failed to hash password")
    }

    _, err = s.DB.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3)", req.GetUsername(), hashedPassword, req.GetRole())
    if err != nil {
        return nil, status.Error(codes.Internal, "Failed to register user")
    }
    return &pb.RegisterResponse{Success: true}, nil
}

func (s *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
    var user common.User
    row := s.DB.QueryRow("SELECT id, username, password, role FROM users WHERE username = $1", req.GetUsername())
    err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Role)
    if err != nil {
        return nil, status.Error(codes.NotFound, "User not found")
    }

    if !common.CheckPasswordHash(req.GetPassword(), user.Password) {
        return nil, status.Error(codes.Unauthenticated, "Invalid password")
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "id":    user.ID,
        "role":  user.Role,
        "exp":   time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString([]byte(config.JwtSecret))
    if err != nil {
        return nil, status.Error(codes.Internal, "Failed to generate token")
    }

    return &pb.LoginResponse{Token: tokenString}, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    var user common.User
    row := s.DB.QueryRow("SELECT id, username, role FROM users WHERE id = $1", req.GetId())
    err := row.Scan(&user.ID, &user.Username, &user.Role)
    if err != nil {
        return nil, status.Error(codes.NotFound, "User not found")
    }

    return &pb.GetUserResponse{
        Id:       int32(user.ID),
        Username: user.Username,
        Role:     user.Role,
    }, nil
}
