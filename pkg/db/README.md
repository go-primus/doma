# DB Package

Database driver abstraction layer with auto-initialization support.

## Quick Start

```go
import "gitlab.etsus.net/meos/synthone/pkg/db"

// SQLite
gormDB, err := db.Open("sqlite", "data/index.db")

// MySQL
gormDB, err := db.Open("mysql", "user:pass@tcp(host:3306)/mydb?charset=utf8mb4")

// PostgreSQL
gormDB, err := db.Open("postgres", "user:pass@localhost:5432/mydb")
gormDB, err := db.Open("pgsql", "user:pass@localhost:5432/mydb")  // alias
```

## Supported Drivers

| Driver | Name | Aliases | DSN Format |
|--------|------|---------|------------|
| SQLite | `sqlite` | - | `path/to/file.db` |
| MySQL | `mysql` | - | `user:pass@tcp(host:port)/dbname?params` |
| PostgreSQL | `postgres` | `pgsql` | `user:pass@host:port/dbname?params` |

## API

### Open(dbType, dsn string) (*gorm.DB, error)

Convenience function that initializes database and returns GORM connection.

```go
gormDB, err := db.Open("mysql", "user:pass@host:3306/mydb")
if err != nil {
    log.Fatal(err)
}

// Use gormDB...
gormDB.AutoMigrate(&User{})
gormDB.Create(&User{Name: "test"})
```

### Get(name string) (Driver, bool)

Get driver instance by name.

```go
driver, ok := db.Get("postgres")
if !ok {
    log.Fatal("driver not found")
}
```

### Names() []string

List all registered driver names.

```go
names := db.Names()
// ["mysql", "postgres", "sqlite"]
```

## Driver Interface

Implement `Driver` interface for custom drivers:

```go
type Driver interface {
    Name() string
    Initialize(dsn string) error
    Connect(dsn string) (*gorm.DB, error)
}
```

## Directory Structure

```
pkg/db/
├── driver.go              # Core API
└── drivers/
    ├── mysql/driver.go
    ├── postgres/driver.go
    └── sqlite/driver.go
```
