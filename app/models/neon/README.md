# Neon Database Wrapper

A lightweight Go wrapper around GORM for interacting with Neon (PostgreSQL) database. Provides a clean, consistent API for database operations with support for CRUD, validation, upserts, joins, and transactions.

## Overview

The Neon wrapper provides:
- **Initialization**: One-time database connection setup using environment variables
- **CRUD Operations**: Create, Read, Update, Delete with context support
- **Validation**: Built-in validation for create and update operations
- **Advanced Operations**: Upsert, FindOrCreate, FindAndUpdate, and more
- **Join Support**: Single, multiple, and ordered joins across tables
- **Soft Deletes**: Logical delete support for models with DeletedAt field
- **Query Utilities**: Limit, Offset, OrderBy, and Joins helpers
- **Error Handling**: Consistent error handling and statistics tracking
- **Migration Management**: Separate schema management from query operations

## Setup

### Prerequisites

1. Set the `GCP_NEON_DATABASE_URL` environment variable with your Neon connection string:

```bash
export GCP_NEON_DATABASE_URL="postgresql://user:password@host:port/database?sslmode=require"
```

2. Add the dependencies to your `go.mod`:

```bash
go get gorm.io/driver/postgres
go get gorm.io/gorm
```

### Initialization

Initialize the Neon service in your application startup (typically in `main.go`):

```go
package main

import (
	"context"
	"github.com/atharvyadav96k/gcp/app/models/neon"
)

func main() {
	// Initialize database connection
	neonService := neon.InitDB(&neon.Service{})
	defer neonService.Close()
	
	// Register all models and run migrations (once at startup)
	neonService.RegisterAndMigrate(&User{}, &Product{}, &Order{})
	
	// Now you can use neonService throughout your application
	ctx := context.Background()
	
	// Use the service for database operations
	var users []User
	neonService.FindAll(ctx, &users)
}
```

## Available Methods

### Initialization Methods

#### InitDB
Initializes the database connection using the `GCP_NEON_DATABASE_URL` environment variable.

```go
neonService := neon.InitDB(&neon.Service{})
```

#### RegisterAndMigrate
Performs auto-migration on the database for given models.

```go
neonService.RegisterAndMigrate(&User{}, &Product{}, &Order{})
```

**Important:** Call this once at application startup, not in every function.

#### Close
Closes the database connection and releases resources.

```go
defer neonService.Close()
```

---

### CRUD Operations

#### Create
Creates a single record.

```go
err := neonService.Create(ctx, &user)
```

#### CreateBatch
Creates multiple records in batches.

```go
err := neonService.CreateBatch(ctx, users, 100) // batch size of 100
```

#### CreateWithValidation
Creates a record after running validation.

```go
err := neonService.CreateWithValidation(ctx, &user, func() error {
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	return nil
})
```

#### FindByID
Retrieves a record by ID.

```go
var user User
err := neonService.FindByID(ctx, &user, userID)
```

#### FindOne
Retrieves a single record matching conditions.

```go
var user User
err := neonService.FindOne(ctx, &user, "email = ?", "john@example.com")
```

#### FindAll
Retrieves all records, optionally with conditions.

```go
var users []User
err := neonService.FindAll(ctx, &users)

// With conditions
var activeUsers []User
err := neonService.FindAll(ctx, &activeUsers, "status = ?", "active")
```

#### Update
Updates records matching conditions.

```go
result := neonService.Update(ctx, &user, "id = ?", userID)
if result.Error != nil {
	log.Printf("Update failed: %v", result.Error)
}
```

#### UpdateWithValidation
Updates a record after running validation.

```go
result := neonService.UpdateWithValidation(ctx, &user, func() error {
	if user.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}, "id = ?", userID)
```

#### Delete
Deletes records matching conditions.

```go
result := neonService.Delete(ctx, &User{}, "id = ?", userID)
```

#### SoftDelete
Marks records as deleted without removing from database.

```go
result := neonService.SoftDelete(ctx, &User{}, "id = ?", userID)
```

**Note:** Model must have a `DeletedAt gorm.DeletedAt` field.

---

### Advanced Operations

#### FindOrCreate
Retrieves a record by conditions, or creates it if not found.

```go
created, err := neonService.FindOrCreate(ctx, &user,
	func() error {
		user.Status = "active"
		user.CreatedAt = time.Now()
		return nil
	},
	"email = ?", "john@example.com")

if created {
	log.Println("New user created")
} else {
	log.Println("User already existed")
}
```

#### FindAndUpdate
Retrieves and updates a record in one operation.

```go
err := neonService.FindAndUpdate(ctx, &user,
	func() error {
		user.Status = "inactive"
		user.LastLogin = time.Now()
		return nil
	},
	"id = ?", userID)
```

#### BulkUpsert
Performs bulk insert or update (upsert) operations. Useful for syncing data.

```go
records := []User{
	{ID: 1, Name: "John", Email: "john@example.com"},
	{ID: 2, Name: "Jane", Email: "jane@example.com"},
}
err := neonService.BulkUpsert(ctx, records, "email")
```

**Note:** Ensure your model has unique constraints on the specified fields.

---

### Join Operations

#### InnerJoin
Performs an INNER JOIN and returns results.

```go
var results []struct {
	UserID   uint
	UserName string
	OrderID  uint
	Amount   float64
}

err := neonService.InnerJoin(ctx, &results,
	"orders ON users.id = orders.user_id",
	[]string{"users.id", "users.name", "orders.id as order_id", "orders.amount"},
	"orders.status = ?", "completed")
```

#### LeftJoin
Performs a LEFT JOIN and returns results. Includes non-matching rows from left table.

```go
var userOrders []struct {
	UserName  string
	OrderID   sql.NullInt64
	Amount    sql.NullFloat64
}

err := neonService.LeftJoin(ctx, &userOrders,
	"orders ON users.id = orders.user_id",
	[]string{"users.name", "orders.id as order_id", "orders.amount"},
	"users.status = ?", "active")
```

**Note:** Use `sql.NullInt64`, `sql.NullFloat64`, etc. for optional columns from the joined table.

#### RightJoin
Performs a RIGHT JOIN and returns results.

```go
err := neonService.RightJoin(ctx, &results,
	"orders ON users.id = orders.user_id",
	[]string{"users.id", "users.name", "orders.amount"})
```

#### MultiJoin
Performs multiple JOINs in a single query.

```go
var results []struct {
	UserName    string
	ProductName string
	Quantity    int
}

joins := []neon.JoinInfo{
	{Type: "INNER", Condition: "orders ON users.id = orders.user_id"},
	{Type: "INNER", Condition: "order_items ON orders.id = order_items.order_id"},
	{Type: "INNER", Condition: "products ON order_items.product_id = products.id"},
}

err := neonService.MultiJoin(ctx, &results, joins,
	[]string{"users.name", "products.name", "order_items.quantity"},
	"orders.status = ?", "completed")
```

#### JoinWithOrder
Performs a JOIN with specific ordering.

```go
var userOrders []struct {
	UserName string
	OrderID  uint
	Amount   float64
}

err := neonService.JoinWithOrder(ctx, &userOrders,
	"LEFT",
	"orders ON users.id = orders.user_id",
	[]string{"users.name", "orders.id as order_id", "orders.amount"},
	"orders.created_at DESC",
	"users.status = ?", "active")
```

---

### Type-Safe Join Builder (Recommended)

The type-safe join builder uses struct fields instead of hardcoded strings. This provides compile-time safety and automatic column name resolution using reflection.

#### NewInnerJoinBuilder
Creates a join builder for INNER JOIN operations.

```go
builder := neonService.NewInnerJoinBuilder(&Order{})
```

#### NewLeftJoinBuilder
Creates a join builder for LEFT JOIN operations.

```go
builder := neonService.NewLeftJoinBuilder(&Order{})
```

#### NewRightJoinBuilder
Creates a join builder for RIGHT JOIN operations.

```go
builder := neonService.NewRightJoinBuilder(&Order{})
```

#### On
Sets the join condition using struct field names instead of raw SQL.

**Signature:**
```go
func (jb *JoinBuilder) On(leftModel interface{}, leftField, rightField string) *JoinBuilder
```

**Parameters:**
- `leftModel`: the left side model (e.g., `&User{}`)
- `leftField`: field name in left model (e.g., `"ID"`)
- `rightField`: field name in right model (e.g., `"UserID"`)

**Returns:** `*JoinBuilder` for method chaining

**Example:**
```go
// Instead of: "orders ON users.id = orders.user_id"
// You write:
builder.On(&User{}, "ID", "UserID")
// Automatically converts field names to database column names
```

#### Select
Specifies which columns to retrieve using struct field names.

**Signature:**
```go
func (jb *JoinBuilder) Select(selections map[interface{}][]string) *JoinBuilder
```

**Parameters:**
- `selections`: map of model struct to field names (e.g., `map[interface{}][]string{&User{}: {"ID", "Name"}}`)

**Example:**
```go
builder.Select(map[interface{}][]string{
    &User{}: {"ID", "Name", "Email"},
    &Order{}: {"ID", "Amount", "Status"},
})
```

#### SelectAll
Includes all columns from a model.

**Example:**
```go
builder.SelectAll(&User{}).SelectAll(&Order{})
```

#### Execute
Runs the join query and populates the destination.

**Signature:**
```go
func (jb *JoinBuilder) Execute(s *Service, ctx context.Context, dest interface{}, conds ...interface{}) error
```

**Example:**
```go
var results []struct {
    UserID   uint
    UserName string
    OrderID  uint
    Amount   float64
}

err := builder.Execute(neonService, ctx, &results, "orders.status = ?", "completed")
```

#### Complete Type-Safe Join Example

```go
// Setup models
type User struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string
    Email     string
    Status    string
}

type Order struct {
    ID        uint
    UserID    uint
    Amount    float64
    Status    string
    CreatedAt time.Time
}

// Use type-safe join builder
var results []struct {
    UserID    uint
    UserName  string
    OrderID   uint
    Amount    float64
}

err := neonService.
    NewInnerJoinBuilder(&Order{}).
    On(&User{}, "ID", "UserID").
    Select(map[interface{}][]string{
        &User{}: {"ID", "Name"},
        &Order{}: {"ID", "Amount"},
    }).
    Execute(neonService, ctx, &results, 
        "orders.status = ?", "completed",
        "users.status = ?", "active")

if err != nil {
    log.Printf("Join failed: %v", err)
}

for _, result := range results {
    fmt.Printf("User %d (%s) - Order Amount: %.2f\n", 
        result.UserID, result.UserName, result.Amount)
}
```

#### Automatic Column Name Resolution

The type-safe join builder automatically:
1. **Extracts table names** from GORM models using reflection
2. **Converts field names** to database column names using:
   - GORM `column` tag if present
   - Automatic snake_case conversion (e.g., `UserID` → `user_id`)
3. **Validates field existence** at runtime and returns clear error messages

**Example:**
```go
// Given these models:
type User struct {
    ID        uint   `gorm:"primaryKey"`  // Column: id
    FirstName string                      // Column: first_name (auto snake_case)
    Email     string `gorm:"column:email_address"` // Custom column name
}

// The builder automatically resolves:
GetTableName(&User{})              // Returns: "users"
GetColumnName(&User{}, "ID")       // Returns: "id"
GetColumnName(&User{}, "FirstName") // Returns: "first_name"
GetColumnName(&User{}, "Email")    // Returns: "email_address"
```

#### Helper Functions

**GetTableName**
Extracts the table name from a GORM model.

```go
tableName := neon.GetTableName(&User{})  // Returns "users"
```

**GetColumnName**
Extracts the database column name from a struct field.

```go
colName, err := neon.GetColumnName(&Order{}, "UserID")  // Returns "user_id"
```

---

### Query Utilities

#### Limit
Adds LIMIT clause to a GORM query.

```go
db := neonService.Limit(neonService.Client, 10)
```

#### Offset
Adds OFFSET clause to a GORM query.

```go
db := neonService.Offset(neonService.Client, 20)
```

#### OrderBy
Adds ORDER clause to a GORM query.

```go
db := neonService.OrderBy(neonService.Client, "created_at", "DESC")
```

#### Joins
Adds JOIN clause to a GORM query.

```go
db := neonService.Joins(neonService.Client, "orders ON users.id = orders.user_id")
```

---

### Error Handling & Utilities

#### GetLastError
Returns a formatted error message for the last database operation.

```go
result := neonService.Update(ctx, &user, "id = ?", userID)
if msg := neonService.GetLastError(result, "update"); msg != "" {
	logger.Error(msg)
}
```

#### Stats
Returns statistics about the last database operation.

```go
result := neonService.Update(ctx, updateData, "id = ?", userID)
stats := neonService.Stats(result)
fmt.Printf("Updated %d rows\n", stats["RowsAffected"])
```

---

## Usage Examples

### 1. Basic Create Operations

```go
// Create a single record
user := &User{
	Name:  "John Doe",
	Email: "john@example.com",
	Status: "active",
}
err := neonService.Create(ctx, user)
if err != nil {
	log.Printf("Failed to create user: %v", err)
}

// Create multiple records in batches
users := []User{
	{Name: "John", Email: "john@example.com"},
	{Name: "Jane", Email: "jane@example.com"},
	{Name: "Bob", Email: "bob@example.com"},
}
err := neonService.CreateBatch(ctx, users, 100) // batch size of 100
```

### 2. Read Operations

```go
// Find by ID
var user User
err := neonService.FindByID(ctx, &user, userID)
if err == gorm.ErrRecordNotFound {
	log.Println("User not found")
}

// Find one record with conditions
var user User
err := neonService.FindOne(ctx, &user, "email = ?", "john@example.com")

// Find all records
var users []User
err := neonService.FindAll(ctx, &users)

// Find with conditions
var activeUsers []User
err := neonService.FindAll(ctx, &activeUsers, "status = ?", "active")

// Complex queries with joins, ordering, etc.
var users []User
err := neonService.FindWithQuery(ctx, func(db *gorm.DB) *gorm.DB {
	return db.
		Joins("LEFT JOIN orders ON users.id = orders.user_id").
		Where("users.status = ?", "active").
		Order("users.created_at DESC").
		Limit(10)
}, &users)
```

### 3. Update Operations

```go
// Update records matching conditions
result := neonService.Update(ctx, User{Status: "inactive"}, "id = ?", userID)
if result.Error != nil {
	// Handle error
}
fmt.Printf("Updated %d rows\n", result.RowsAffected)

// Update by ID
user := User{ID: 1, Name: "Updated Name", Status: "active"}
result := neonService.UpdateByID(ctx, &user, user.ID)

// Update specific columns only
result := neonService.UpdateColumns(ctx,
	map[string]interface{}{
		"status":     "inactive",
		"updated_at": time.Now(),
	},
	"id = ?", userID)
```

### 4. Delete Operations

```go
// Delete records matching conditions
result := neonService.Delete(ctx, &User{}, "id = ?", userID)
if result.Error != nil {
	// Handle error
}

// Delete by ID
result := neonService.DeleteByID(ctx, &User{}, userID)

// Delete all matching conditions
result := neonService.Delete(ctx, &User{}, "status = ?", "inactive")
```

### 5. Query Operations

```go
// Count records
count, err := neonService.Count(ctx, &User{}, "status = ?", "active")

// Check if exists
exists, err := neonService.Exists(ctx, &User{}, "email = ?", "john@example.com")
if exists {
	// User with this email already exists
}

// Pagination
var users []User
total, err := neonService.Paginate(ctx, &users, 
	1,    // page number
	10,   // page size
	"status = ?", "active")
fmt.Printf("Page 1: %d/%d users\n", len(users), total)
```

### 6. Joins & Multiple Tables

```go
// ========== RECOMMENDED: Type-Safe Join Builder ==========

// Simple INNER JOIN with type safety
var results []struct {
	UserID   uint
	UserName string
	OrderID  uint
	Amount   float64
}

err := neonService.
	NewInnerJoinBuilder(&Order{}).
	On(&User{}, "ID", "UserID").
	Select(map[interface{}][]string{
		&User{}: {"ID", "Name"},
		&Order{}: {"ID", "Amount"},
	}).
	Execute(neonService, ctx, &results, "orders.status = ?", "completed")

// LEFT JOIN with type safety - includes users with no orders
var userOrders []struct {
	UserName  string
	OrderID   sql.NullInt64
	Amount    sql.NullFloat64
}

err := neonService.
	NewLeftJoinBuilder(&Order{}).
	On(&User{}, "ID", "UserID").
	Select(map[interface{}][]string{
		&User{}: {"Name"},
		&Order{}: {"ID", "Amount"},
	}).
	Execute(neonService, ctx, &userOrders, "users.status = ?", "active")

// RIGHT JOIN with type safety
var orderUsers []struct {
	UserName string
	OrderID  uint
	Amount   float64
}

err := neonService.
	NewRightJoinBuilder(&User{}).
	On(&Order{}, "UserID", "ID").
	Select(map[interface{}][]string{
		&User{}: {"Name"},
		&Order{}: {"ID", "Amount"},
	}).
	Execute(neonService, ctx, &orderUsers)

// ========== ALTERNATIVE: Raw String Joins (For Complex Multi-Table Scenarios) ==========

// Multiple JOINs across 3+ tables (use raw strings for complex scenarios)
var complexResults []struct {
	UserName    string
	ProductName string
	Quantity    int
	OrderDate   time.Time
}

joins := []neon.JoinInfo{
	{Type: "INNER", Condition: "orders ON users.id = orders.user_id"},
	{Type: "INNER", Condition: "order_items ON orders.id = order_items.order_id"},
	{Type: "INNER", Condition: "products ON order_items.product_id = products.id"},
}

err := neonService.MultiJoin(ctx, &complexResults, joins,
	[]string{"users.name", "products.name", "order_items.quantity", "orders.created_at as order_date"},
	"orders.status = ?", "completed")

// JOIN with ordering
var recentOrders []struct {
	UserName string
	OrderID  uint
	Amount   float64
}

err := neonService.JoinWithOrder(ctx, &recentOrders,
	"LEFT",
	"orders ON users.id = orders.user_id",
	[]string{"users.name", "orders.id as order_id", "orders.amount"},
	"orders.created_at DESC",
	"users.status = ?", "active")
```

### 7. Advanced Operations

```go
// Find or Create - get existing user or create new one
var user User
created, err := neonService.FindOrCreate(ctx, &user,
	func() error {
		user.Status = "active"
		user.CreatedAt = time.Now()
		return nil
	},
	"email = ?", "john@example.com")

if created {
	log.Println("New user created")
}

// Find and Update - retrieve and update in one call
var user User
err := neonService.FindAndUpdate(ctx, &user,
	func() error {
		user.Status = "inactive"
		user.LastLogin = time.Now()
		return nil
	},
	"id = ?", userID)

// Bulk Upsert - useful for syncing data from external sources
records := []User{
	{ID: 1, Name: "John", Email: "john@example.com", Status: "active"},
	{ID: 2, Name: "Jane", Email: "jane@example.com", Status: "active"},
	{ID: 3, Name: "Bob", Email: "bob@example.com", Status: "active"},
}
err := neonService.BulkUpsert(ctx, records, "email")

// Soft Delete - mark as deleted but keep data
result := neonService.SoftDelete(ctx, &User{}, "id = ?", userID)
if result.Error != nil {
	log.Printf("Soft delete failed: %v", result.Error)
}

// Create with validation - validate before insert
err := neonService.CreateWithValidation(ctx, &user, func() error {
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	if len(user.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	return nil
})

// Update with validation - validate before update
result := neonService.UpdateWithValidation(ctx, &user, 
	func() error {
		if user.Name == "" {
			return fmt.Errorf("name is required")
		}
		return nil
	},
	"id = ?", userID)
```

## Model Definition Examples

### Basic User Model

```go
package models

import "gorm.io"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Password  string    `json:"-"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // For soft deletes
}
```

### Model with Relations

```go
type Order struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `json:"order_id"`
	ProductID uint      `json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"-"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
}

type Product struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}
```

## Error Handling

### Common Errors

- `gorm.ErrRecordNotFound`: Record doesn't exist
- `gorm.ErrInvalidData`: Invalid data passed to operation
- Validation errors: Custom validation function returned error
- Database errors: Connection issues, constraint violations, etc.

### Error Handling Pattern

```go
result := neonService.Update(ctx, &user, "id = ?", userID)
if result.Error != nil {
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("User not found")
	} else {
		log.Printf("Update failed: %v", result.Error)
	}
	return result.Error
}

// Check rows affected
if result.RowsAffected == 0 {
	log.Println("No records were updated")
}

// Get stats
stats := neonService.Stats(result)
log.Printf("Operation stats: %+v", stats)
```

## Best Practices

1. **Call RegisterAndMigrate Once**: Do this at application startup, not in every function
2. **Use Context with Timeouts**: Always pass context with appropriate timeouts
3. **Validate Before Database Operations**: Use validation functions to catch issues early
4. **Use Transactions for Multi-Step Operations**: Chain operations logically
5. **Use Joins Instead of Multiple Queries**: Reduce database round trips
6. **Defer Close()**: Always defer `neonService.Close()` to ensure proper cleanup
7. **Use Soft Deletes**: For audit trails and data recovery
8. **Batch Operations**: Use `CreateBatch` and `BulkUpsert` for large datasets
9. **Check RowsAffected**: Verify operations affected expected number of rows
10. **Use Proper Indexes**: Define database indexes in your models for queries

## Service Structure

```go
type Service struct {
	// Client is the underlying GORM database client
	Client *gorm.DB

	// Once ensures the client is initialized only once
	Once sync.Once
}

type JoinInfo struct {
	Type      string // "INNER", "LEFT", "RIGHT", "FULL"
	Condition string // e.g., "orders ON users.id = orders.user_id"
}
```

## See Also

- [GORM Documentation](https://gorm.io)
- [PostgreSQL Documentation](https://www.postgresql.org/docs)
- [Neon Documentation](https://neon.tech/docs)
// Simple INNER JOIN
var results []struct {
	UserID   uint
	UserName string
	OrderID  uint
	Amount   float64
}
err := neonService.InnerJoin(ctx, &results,
	"orders ON users.id = orders.user_id",
	[]string{"users.id", "users.name", "orders.id as order_id", "orders.amount"},
	"orders.status = ?", "completed")

// LEFT JOIN with conditions
var userOrders []struct {
	UserName string
	OrderID  uint
	Amount   float64
}
err := neonService.LeftJoin(ctx, &userOrders,
	"orders ON users.id = orders.user_id",
	[]string{"users.name", "orders.id as order_id", "orders.amount"},
	"users.status = ?", "active")

// Multiple JOINs
var userOrderItems []struct {
	UserName    string
	OrderID     uint
	ProductName string
	Quantity    int
}
joins := []neon.JoinInfo{
	{Type: "INNER", Condition: "orders ON users.id = orders.user_id"},
	{Type: "INNER", Condition: "order_items ON orders.id = order_items.order_id"},
	{Type: "INNER", Condition: "products ON order_items.product_id = products.id"},
}
err := neonService.MultiJoin(ctx, &userOrderItems, joins,
	[]string{"users.name", "orders.id", "products.name", "order_items.quantity"},
	"orders.status = ?", "completed")

// JOIN with ORDER BY
var sortedOrders []struct {
	UserName  string
	OrderID   uint
	CreatedAt time.Time
}
err := neonService.JoinWithOrder(ctx, &sortedOrders,
	"LEFT",
	"orders ON users.id = orders.user_id",
	[]string{"users.name", "orders.id", "orders.created_at"},
	"orders.created_at DESC",
	"users.status = ?", "active")

// JOIN with PAGINATION
var paginatedResults []struct {
	UserName string
	OrderID  uint
}
total, err := neonService.JoinWithPagination(ctx, &paginatedResults,
	"LEFT",
	"orders ON users.id = orders.user_id",
	[]string{"users.name", "orders.id"},
	1, 10,
	"users.status = ?", "active")
fmt.Printf("Page 1: %d/%d results\n", len(paginatedResults), total)
```

### 7. Raw SQL Queries

```go
// Execute raw query and scan results
var user User
err := neonService.Raw(ctx, &user,
	"SELECT * FROM users WHERE email = ? AND status = ?",
	"john@example.com", "active")

// Execute raw command (INSERT, UPDATE, DELETE)
result := neonService.Exec(ctx,
	"UPDATE users SET status = ? WHERE status = ?",
	"active", "pending")
if result.Error != nil {
	// Handle error
}
fmt.Printf("Updated %d rows\n", result.RowsAffected)
```

## API Reference

### Initialization Functions

- `NewClient()` - Create a new Service instance
- `InitDB(s *Service) *Service` - Initialize database connection
- `RegisterAndMigrate(models ...interface{})` - Register and migrate models (call once at startup)
- `Close() error` - Close database connection

### CRUD Operations

- `Create(ctx, value) error` - Insert single record
- `CreateBatch(ctx, records, batchSize) error` - Insert multiple records
- `FindByID(ctx, dest, id) error` - Get record by primary key
- `FindOne(ctx, dest, conditions...) error` - Get single record by conditions
- `FindAll(ctx, dest, conditions...) error` - Get multiple records
- `FindWithQuery(ctx, queryFunc, dest) error` - Get records with custom query builder
- `Update(ctx, value, conditions...) *gorm.DB` - Update records
- `UpdateByID(ctx, dest, id) *gorm.DB` - Update record by ID
- `UpdateColumns(ctx, updates, conditions...) *gorm.DB` - Update specific columns
- `Delete(ctx, value, conditions...) *gorm.DB` - Delete records
- `DeleteByID(ctx, value, id) *gorm.DB` - Delete record by ID

### Query Operations

- `Count(ctx, model, conditions...) (int64, error)` - Count matching records
- `Exists(ctx, model, conditions...) (bool, error)` - Check if record exists
- `Paginate(ctx, dest, page, pageSize, conditions...) (int64, error)` - Get paginated results
- `Raw(ctx, dest, sql, values...) error` - Execute raw query
- `Exec(ctx, sql, values...) *gorm.DB` - Execute raw command

### Join Operations

- `InnerJoin(ctx, dest, joinTable, columns, conditions...) error` - Perform INNER JOIN
- `LeftJoin(ctx, dest, joinTable, columns, conditions...) error` - Perform LEFT JOIN
- `RightJoin(ctx, dest, joinTable, columns, conditions...) error` - Perform RIGHT JOIN
- `MultiJoin(ctx, dest, joins, columns, conditions...) error` - Perform multiple JOINs
- `JoinWithOrder(ctx, dest, joinType, joinTable, columns, orderBy, conditions...) error` - JOIN with ordering
- `JoinWithPagination(ctx, dest, joinType, joinTable, columns, page, pageSize, conditions...) (int64, error)` - JOIN with pagination

### Error Handling

- `IsRecordNotFound(err) bool` - Check if record not found
- `IsInvalidData(err) bool` - Check if service uninitialized
- `HandleQueryError(err, entityType, identifier) string` - Get user-friendly error message

## Best Practices

1. **Always pass context**: All operations accept `context.Context` for timeout and cancellation control.

```go
// Good
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
err := neonService.FindByID(ctx, &user, userID)

// Avoid
err := neonService.FindByID(context.Background(), &user, userID)
```

2. **Separate migrations from queries**: Run migrations once at startup, never in request handlers.

```go
// In main.go (startup)
neonService.RegisterAndMigrate(&User{}, &Product{})

// In handlers (runtime)
err := neonService.FindByID(ctx, &user, userID) // Just query
```

3. **Handle "not found" errors properly**:

```go
var user User
err := neonService.FindByID(ctx, &user, userID)
if neon.IsRecordNotFound(err) {
	return nil, fmt.Errorf("user not found")
}
if err != nil {
	return nil, fmt.Errorf("database error: %w", err)
}
```

4. **Use batch operations for bulk inserts**:

```go
// Good for large datasets
err := neonService.CreateBatch(ctx, largeUserList, 500)

// Avoid for large datasets
for _, user := range largeUserList {
	neonService.Create(ctx, &user)
}
```

## Differences from Direct GORM Usage

| Aspect | GORM | Neon Wrapper |
|--------|------|-------------|
| **Context** | Optional | Required (enforced) |
| **Lifecycle** | Manual | Managed singleton |
| **Consistency** | Varies | Standardized |
| **Error handling** | Direct error types | Helper functions included |
| **Initialization** | Multiple places | Once at startup |

## Integration with App Structure

```go
package app

import (
	"github.com/atharvyadav96k/gcp/app/models/neon"
	"github.com/atharvyadav96k/gcp/app/models/firestore"
	"github.com/atharvyadav96k/gcp/app/models/pubsub"
)

type App struct {
	Env       secrets.Env
	FireStore *firestore.Service
	PubSub    *pubsub.Service
	Neon      *neon.Service  // Add Neon service
}

// In main.go
func init() {
	app.Neon = neon.InitDB(&neon.Service{})
	app.Neon.RegisterAndMigrate(&User{}, &Order{})
}
```

## Common Patterns

### Pattern 1: Create with validation and notification

```go
err := neonService.Transaction(ctx, func(tx *gorm.DB) error {
	// Create user
	if err := tx.Create(&user).Error; err != nil {
		return err
	}
	
	// Publish notification event
	return publishEvent(ctx, "user.created", user.ID)
})
```

### Pattern 2: Find or create

```go
var user User
err := neonService.FindOne(ctx, &user, "email = ?", email)
if neon.IsRecordNotFound(err) {
	// Create new user
	user = User{Email: email, Name: name}
	err = neonService.Create(ctx, &user)
}
```

### Pattern 3: Bulk update with validation

```go
// Update all inactive users
updates := map[string]interface{}{
	"status": "archived",
	"updated_at": time.Now(),
}
result := neonService.UpdateColumns(ctx, updates, "status = ?", "inactive")
if result.RowsAffected > 0 {
	// Log the bulk update
}
```

## Migration and Schema Management

Migrations are handled separately from query operations:

```go
// 1. Define models with GORM tags
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"index"`
	Email     string    `gorm:"unique"`
	Status    string    `gorm:"default:active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 2. Register migrations once (in main.go or init)
neonService.RegisterAndMigrate(&User{})

// 3. Use service for queries (no migration involved)
err := neonService.Create(ctx, &user) // Just creates a record
```

For schema changes, update your model and call `RegisterAndMigrate` again. GORM will handle the migration.

## Closing Connection

Always close the connection when shutting down:

```go
defer neonService.Close()
```

This ensures proper cleanup of database resources.
