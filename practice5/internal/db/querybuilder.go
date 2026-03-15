package db

import (
    "fmt"
    "strings"
)

type SQLBuilder struct {
    clauses []string
    args    []any
}

func NewSQLBuilder() *SQLBuilder {
    return &SQLBuilder{clauses: []string{"1=1"}}
}

func (b *SQLBuilder) AddExact(field string, value any) {
    b.args = append(b.args, value)
    b.clauses = append(b.clauses, fmt.Sprintf("%s = $%d", field, len(b.args)))
}

func (b *SQLBuilder) AddILike(field string, value string) {
    b.args = append(b.args, "%"+value+"%")
    b.clauses = append(b.clauses, fmt.Sprintf("%s ILIKE $%d", field, len(b.args)))
}

func (b *SQLBuilder) Where() string {
    return strings.Join(b.clauses, " AND ")
}

func (b *SQLBuilder) Args() []any {
    return b.args
}
