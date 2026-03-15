package repository

import (
    "database/sql"
    "fmt"
    "strings"

    qbuilder "practice5_ready/internal/db"
    "practice5_ready/internal/models"
)

type Repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
    return &Repository{db: db}
}

var allowedOrderBy = map[string]string{
    "id":         "u.id",
    "name":       "u.name",
    "email":      "u.email",
    "gender":     "u.gender",
    "birth_date": "u.birth_date",
}

func normalizeOrderBy(orderBy string) string {
    orderBy = strings.TrimSpace(strings.ToLower(orderBy))
    if orderBy == "" {
        return "u.id ASC"
    }

    direction := "ASC"
    if strings.HasPrefix(orderBy, "-") {
        direction = "DESC"
        orderBy = strings.TrimPrefix(orderBy, "-")
    }

    column, ok := allowedOrderBy[orderBy]
    if !ok {
        return "u.id ASC"
    }

    return fmt.Sprintf("%s %s", column, direction)
}

func (r *Repository) GetPaginatedUsers(filters models.UserFilters) (models.PaginatedResponse, error) {
    builder := qbuilder.NewSQLBuilder()

    if filters.ID != nil {
        builder.AddExact("u.id", *filters.ID)
    }
    if filters.Name != "" {
        builder.AddILike("u.name", filters.Name)
    }
    if filters.Email != "" {
        builder.AddILike("u.email", filters.Email)
    }
    if filters.Gender != "" {
        builder.AddExact("u.gender", filters.Gender)
    }
    if filters.BirthDate != "" {
        builder.AddExact("u.birth_date", filters.BirthDate)
    }

    whereClause := builder.Where()
    baseArgs := builder.Args()

    countQuery := fmt.Sprintf(`
        SELECT COUNT(*)
        FROM users u
        WHERE %s
    `, whereClause)

    var totalCount int
    if err := r.db.QueryRow(countQuery, baseArgs...).Scan(&totalCount); err != nil {
        return models.PaginatedResponse{}, err
    }

    args := append([]any{}, baseArgs...)
    args = append(args, filters.Limit, filters.Offset)

    query := fmt.Sprintf(`
        SELECT u.id, u.name, u.email, u.gender, u.birth_date
        FROM users u
        WHERE %s
        ORDER BY %s
        LIMIT $%d OFFSET $%d
    `, whereClause, normalizeOrderBy(filters.OrderBy), len(baseArgs)+1, len(baseArgs)+2)

    rows, err := r.db.Query(query, args...)
    if err != nil {
        return models.PaginatedResponse{}, err
    }
    defer rows.Close()

    users := make([]models.User, 0)
    for rows.Next() {
        var u models.User
        if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
            return models.PaginatedResponse{}, err
        }
        users = append(users, u)
    }

    if err := rows.Err(); err != nil {
        return models.PaginatedResponse{}, err
    }

    return models.PaginatedResponse{
        Data:       users,
        TotalCount: totalCount,
        Limit:      filters.Limit,
        Offset:     filters.Offset,
    }, nil
}

func (r *Repository) GetCommonFriends(user1, user2 int) ([]models.User, error) {
    query := `
        SELECT DISTINCT u.id, u.name, u.email, u.gender, u.birth_date
        FROM user_friends uf1
        JOIN user_friends uf2 ON uf1.friend_id = uf2.friend_id
        JOIN users u ON u.id = uf1.friend_id
        WHERE uf1.user_id = $1
          AND uf2.user_id = $2
          AND $1 <> $2
        ORDER BY u.id
    `

    rows, err := r.db.Query(query, user1, user2)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    commonFriends := make([]models.User, 0)
    for rows.Next() {
        var u models.User
        if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
            return nil, err
        }
        commonFriends = append(commonFriends, u)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return commonFriends, nil
}
