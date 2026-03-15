package models

import "time"

type User struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Gender    string    `json:"gender"`
    BirthDate time.Time `json:"birth_date"`
}

type PaginatedResponse struct {
    Data       []User `json:"data"`
    TotalCount int    `json:"totalCount"`
    Limit      int    `json:"limit"`
    Offset     int    `json:"offset"`
}

type UserFilters struct {
    ID        *int
    Name      string
    Email     string
    Gender    string
    BirthDate string
    Limit     int
    Offset    int
    OrderBy   string
}
