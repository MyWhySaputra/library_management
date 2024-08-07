package bookservice

import (
    "context"
    "database/sql"
    "time"

    pb "library_management/proto/book"
    "library_management/internal/common"

		"github.com/go-redis/redis/v8"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type Book struct {
    ID       string
    Title    string
    Author   string
    Category string
    Stock    int
}

type BookServiceServer struct {
    DB    *sql.DB
    Redis *redis.Client
    pb.UnimplementedBookServiceServer
}

func NewBookServiceServer(db *sql.DB, redis *redis.Client) *BookServiceServer {
    return &BookServiceServer{DB: db, Redis: redis}
}

func (s *BookServiceServer) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
    _, err := s.DB.Exec("INSERT INTO books (title, author_id, category_id, stock, created_at) VALUES ($1, $2, $3, $4, $5)", req.GetTitle(), req.GetAuthorId(), req.GetCategoryId(), req.GetStock(), time.Now())
    if err != nil {
        return nil, status.Error(codes.Internal, "Failed to create book")
    }
    return &pb.CreateBookResponse{Success: true}, nil
}

func (s *BookServiceServer) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.GetBookResponse, error) {
    var book common.Book
    row := s.DB.QueryRow("SELECT id, title, author_id, category_id, stock FROM books WHERE id = $1", req.GetId())
    err := row.Scan(&book.ID, &book.Title, &book.AuthorID, &book.CategoryID, &book.Stock)
    if err != nil {
        return nil, status.Error(codes.NotFound, "Book not found")
    }
    return &pb.GetBookResponse{
        Id:         int32(book.ID),
        Title:      book.Title,
        AuthorId:   int32(book.AuthorID),
        CategoryId: int32(book.CategoryID),
        Stock:      int32(book.Stock),
    }, nil
}

func (s *BookServiceServer) UpdateBook(ctx context.Context, req *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {
    _, err := s.DB.Exec("UPDATE books SET title=$1, author_id=$2, category_id=$3, stock=$4 WHERE id=$5", req.GetTitle(), req.GetAuthorId(), req.GetCategoryId(), req.GetStock(), req.GetId())
    if err != nil {
        return nil, status.Error(codes.Internal, "Failed to update book")
    }
    return &pb.UpdateBookResponse{Success: true}, nil
}

func (s *BookServiceServer) DeleteBook(ctx context.Context, req *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
    _, err := s.DB.Exec("DELETE FROM books WHERE id = $1", req.GetId())
    if err != nil {
        return nil, status.Error(codes.Internal, "Failed to delete book")
    }
    return &pb.DeleteBookResponse{Success: true}, nil
}

func (s *BookServiceServer) ListBooks(ctx context.Context, req *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
    rows, err := s.DB.Query("SELECT id, title, author_id, category_id, stock FROM books")
    if err != nil {
        return nil, status.Error(codes.Internal, "Failed to list books")
    }
    defer rows.Close()

    var books []*pb.GetBookResponse
    for rows.Next() {
        var book common.Book
        if err := rows.Scan(&book.ID, &book.Title, &book.AuthorID, &book.CategoryID, &book.Stock); err != nil {
            return nil, status.Error(codes.Internal, "Failed to read book data")
        }
        books = append(books, &pb.GetBookResponse{
            Id:         int32(book.ID),
            Title:      book.Title,
            AuthorId:   int32(book.AuthorID),
            CategoryId: int32(book.CategoryID),
            Stock:      int32(book.Stock),
        })
    }

    return &pb.ListBooksResponse{Books: books}, nil
}

func (s *BookServiceServer) SearchBooks(ctx context.Context, req *pb.SearchBooksRequest) (*pb.SearchBooksResponse, error) {
    rows, err := s.DB.Query("SELECT id, title, author, category, stock FROM books WHERE title ILIKE $1 OR author ILIKE $1 OR category ILIKE $1", "%"+req.Query+"%")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var books []*pb.GetBookResponse
    for rows.Next() {
        var book pb.GetBookResponse
        if err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.Category, &book.Stock); err != nil {
            return nil, err
        }
        books = append(books, &book)
    }

    return &pb.SearchBooksResponse{Books: books}, nil
}

func (s *BookServiceServer) RecommendBooks(ctx context.Context, req *pb.RecommendBooksRequest) (*pb.RecommendBooksResponse, error) {
    // Simple recommendation logic: recommend the most borrowed books
    rows, err := s.DB.Query("SELECT id, title, author, category, stock FROM books ORDER BY borrow_count DESC LIMIT 5")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var books []*pb.GetBookResponse
    for rows.Next() {
        var book pb.GetBookResponse
        if err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.Category, &book.Stock); err != nil {
            return nil, err
        }
        books = append(books, &book)
    }

    return &pb.RecommendBooksResponse{Books: books}, nil
}

func (s *BookServiceServer) BorrowBook(ctx context.Context, req *pb.BorrowBookRequest) (*pb.BorrowBookResponse, error) {
    // Start a transaction
    tx, err := s.DB.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    // Check stock
    var stock int
    err = tx.QueryRow("SELECT stock FROM books WHERE id = $1 FOR UPDATE", req.BookId).Scan(&stock)
    if err != nil {
        return nil, err
    }

    if stock <= 0 {
        return &pb.BorrowBookResponse{Success: false}, nil
    }

    // Update stock
    _, err = tx.Exec("UPDATE books SET stock = stock - 1, borrow_count = borrow_count + 1 WHERE id = $1", req.BookId)
    if err != nil {
        return nil, err
    }

    // Record borrowing
    _, err = tx.Exec("INSERT INTO borrowings (user_id, book_id, borrow_date) VALUES ($1, $2, NOW())", req.UserId, req.BookId)
    if err != nil {
        return nil, err
    }

    if err = tx.Commit(); err != nil {
        return nil, err
    }

    return &pb.BorrowBookResponse{Success: true}, nil
}

func (s *BookServiceServer) ReturnBook(ctx context.Context, req *pb.ReturnBookRequest) (*pb.ReturnBookResponse, error) {
    // Start a transaction
    tx, err := s.DB.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    // Check if the user has borrowed the book
    var borrowID int
    err = tx.QueryRow("SELECT id FROM borrowings WHERE user_id = $1 AND book_id = $2 AND return_date IS NULL", req.UserId, req.BookId).Scan(&borrowID)
    if err != nil {
        if err == sql.ErrNoRows {
            return &pb.ReturnBookResponse{Success: false}, nil
        }
        return nil, err
    }

    // Update stock
    _, err = tx.Exec("UPDATE books SET stock = stock + 1 WHERE id = $1", req.BookId)
    if err != nil {
        return nil, err
    }

    // Record returning
    _, err = tx.Exec("UPDATE borrowings SET return_date = NOW() WHERE id = $1", borrowID)
    if err != nil {
        return nil, err
    }

    if err = tx.Commit(); err != nil {
        return nil, err
    }

    return &pb.ReturnBookResponse{Success: true}, nil
}