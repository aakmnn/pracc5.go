# Practice 5 Ready

## Run locally

1. Start PostgreSQL:
```bash
docker compose up -d
```

2. Install dependencies:
```bash
go mod tidy
```

3. Create tables:
```bash
docker exec -i practice5-postgres psql -U postgres -d practice5 < sql/schema.sql
```

4. Seed data:
```bash
docker exec -i practice5-postgres psql -U postgres -d practice5 < sql/seed.sql
```

5. Run API:
```bash
go run ./cmd/api
```

Server runs on:
```text
http://localhost:8080
```

## Postman demo requests

### Health
`GET http://localhost:8080/health`

### Pagination + sorting
`GET http://localhost:8080/users?limit=5&offset=0&order_by=name`

### Filtering by name
`GET http://localhost:8080/users?limit=5&offset=0&name=alice`

### Filtering by email
`GET http://localhost:8080/users?limit=5&offset=0&email=example.com`

### Filtering by 3 fields + pagination + sorting
`GET http://localhost:8080/users?limit=3&offset=0&name=a&gender=female&birth_date=1999-03-14&order_by=email`

### Common friends
`GET http://localhost:8080/friends/common?user1=1&user2=2`
