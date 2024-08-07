package common

import "time"

type Book struct {
    ID        int
    Title     string
    AuthorID  int
    CategoryID int
    Stock     int
    CreatedAt time.Time
}

type Author struct {
    ID   int
    Name string
}

type Category struct {
    ID   int
    Name string
}

type User struct {
    ID       int
    Username string
    Password string
    Role     string
}
