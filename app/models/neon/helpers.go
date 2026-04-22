package neon

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ========== Type-Safe Join Helpers ==========

// GetTableName extracts the table name from a GORM model using reflection.
// It looks for the table name using GORM's Migrator or defaults to lowercase struct name.
//
// Parameters:
//   - model: a GORM model struct instance (e.g., &User{})
//
// Returns:
//   - string: the table name
//
// Example:
//
//	tableName := neon.GetTableName(&User{})  // Returns "users"
func GetTableName(model interface{}) string {
	if model == nil {
		return ""
	}

	var db *gorm.DB
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(model)
	return stmt.Table
}

// GetColumnName extracts the database column name from a struct field using GORM tags.
// Looks for the "column" tag in the struct field, defaults to lowercased field name.
//
// Parameters:
//   - model: a GORM model struct instance
//   - fieldName: the struct field name (e.g., "UserID")
//
// Returns:
//   - string: the database column name
//   - error: if field not found
//
// Example:
//
//	colName, err := neon.GetColumnName(&Order{}, "UserID")
//	// Returns "user_id"
func GetColumnName(model interface{}, fieldName string) (string, error) {
	if model == nil {
		return "", fmt.Errorf("model cannot be nil")
	}

	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	field, exists := modelType.FieldByName(fieldName)
	if !exists {
		return "", fmt.Errorf("field %s not found in %s", fieldName, modelType.Name())
	}

	// Check for gorm column tag
	if tag, ok := field.Tag.Lookup("gorm"); ok {
		for _, part := range strings.Split(tag, ";") {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "column:") {
				return strings.TrimPrefix(part, "column:"), nil
			}
		}
	}

	// Default to lowercased field name with snake_case conversion
	return toSnakeCase(fieldName), nil
}

// toSnakeCase converts camelCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// JoinBuilder provides a fluent interface for building type-safe joins
type JoinBuilder struct {
	joinType  string
	table     string
	condition string
	columns   []string
}

// NewInnerJoinBuilder creates a new join builder for INNER JOIN
func (s *Service) NewInnerJoinBuilder(toModel interface{}) *JoinBuilder {
	return &JoinBuilder{
		joinType: "INNER",
		table:    GetTableName(toModel),
	}
}

// NewLeftJoinBuilder creates a new join builder for LEFT JOIN
func (s *Service) NewLeftJoinBuilder(toModel interface{}) *JoinBuilder {
	return &JoinBuilder{
		joinType: "LEFT",
		table:    GetTableName(toModel),
	}
}

// NewRightJoinBuilder creates a new join builder for RIGHT JOIN
func (s *Service) NewRightJoinBuilder(toModel interface{}) *JoinBuilder {
	return &JoinBuilder{
		joinType: "RIGHT",
		table:    GetTableName(toModel),
	}
}

// On sets the join condition using struct fields instead of raw SQL
//
// Parameters:
//   - leftModel: the left side model
//   - leftField: field name in left model (e.g., "ID")
//   - rightField: field name in right model (e.g., "UserID")
//
// Returns:
//   - *JoinBuilder for method chaining
//
// Example:
//
//	builder := neonService.NewInnerJoinBuilder(&Order{}).
//	    On(&User{}, "ID", "UserID")
func (jb *JoinBuilder) On(leftModel interface{}, leftField, rightField string) *JoinBuilder {
	leftCol, err := GetColumnName(leftModel, leftField)
	if err != nil {
		leftCol = toSnakeCase(leftField)
	}

	rightCol, err := GetColumnName(leftModel, rightField)
	if err != nil {
		rightCol = toSnakeCase(rightField)
	}

	leftTable := GetTableName(leftModel)
	jb.condition = fmt.Sprintf("%s ON %s.%s = %s.%s", jb.table, leftTable, leftCol, jb.table, rightCol)
	return jb
}

// Select specifies which columns to retrieve using struct field names
//
// Parameters:
//   - selections: map of model struct to field names
//     Example: map[interface{}][]string{&User{}: {"ID", "Name"}}
//
// Returns:
//   - *JoinBuilder for method chaining
//
// Example:
//
//	builder.Select(map[interface{}][]string{
//	    &User{}: {"ID", "Name", "Email"},
//	    &Order{}: {"ID", "Amount", "Status"},
//	})
func (jb *JoinBuilder) Select(selections map[interface{}][]string) *JoinBuilder {
	for model, fields := range selections {
		tableName := GetTableName(model)
		for _, field := range fields {
			colName, err := GetColumnName(model, field)
			if err != nil {
				colName = toSnakeCase(field)
			}
			jb.columns = append(jb.columns, fmt.Sprintf("%s.%s", tableName, colName))
		}
	}
	return jb
}

// SelectAll includes all columns from a model
//
// Example:
//
//	builder.SelectAll(&User{})
func (jb *JoinBuilder) SelectAll(model interface{}) *JoinBuilder {
	tableName := GetTableName(model)
	jb.columns = append(jb.columns, fmt.Sprintf("%s.*", tableName))
	return jb
}

// Execute runs the join query and populates the destination
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to result slice
//   - conds: optional WHERE conditions
//
// Returns:
//   - error if query fails
//
// Example:
//
//	var results []UserOrderResult
//	err := builder.Execute(ctx, &results, "orders.status = ?", "completed")
func (jb *JoinBuilder) Execute(s *Service, ctx context.Context, dest interface{}, conds ...interface{}) error {
	if s == nil || s.Client == nil {
		return gorm.ErrInvalidData
	}
	if jb.condition == "" {
		return fmt.Errorf("join condition not set, call On() first")
	}

	db := s.Client.WithContext(ctx).
		Joins(fmt.Sprintf("%s JOIN %s", jb.joinType, jb.condition))

	if len(jb.columns) > 0 {
		db = db.Select(jb.columns)
	}

	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	return db.Find(dest).Error
}

// ========== End Type-Safe Join Helpers ==========

// CreateWithValidation creates a record after running a validation function.
// If validation fails, the record is not created.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - value: the model/struct to insert
//   - validationFunc: function to validate the record before insert
//
// Returns:
//   - error: validation error or database error
//
// Example:
//
//	err := neonService.CreateWithValidation(ctx, &user, func() error {
//	    if user.Email == "" {
//	        return fmt.Errorf("email is required")
//	    }
//	    if len(user.Password) < 8 {
//	        return fmt.Errorf("password must be at least 8 characters")
//	    }
//	    return nil
//	})
func (s *Service) CreateWithValidation(ctx context.Context, value interface{}, validationFunc func() error) error {
	if err := validationFunc(); err != nil {
		return err
	}
	return s.Create(ctx, value)
}

// UpdateWithValidation updates a record after running a validation function.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - value: the model/struct with updated values
//   - validationFunc: function to validate the record before update
//   - conds: WHERE conditions
//
// Returns:
//   - *gorm.DB: for checking RowsAffected and Error
//
// Example:
//
//	result := neonService.UpdateWithValidation(ctx, &user,
//	    func() error {
//	        if user.Name == "" {
//	            return fmt.Errorf("name is required")
//	        }
//	        return nil
//	    },
//	    "id = ?", userID)
func (s *Service) UpdateWithValidation(ctx context.Context, value interface{}, validationFunc func() error, conds ...interface{}) *gorm.DB {
	if err := validationFunc(); err != nil {
		return &gorm.DB{Error: err}
	}
	return s.Update(ctx, value, conds...)
}

// FindOrCreate retrieves a record by conditions, or creates it if not found.
// Returns the found or newly created record and indicates if it was created.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to the model/struct to populate or create
//   - createFunc: function to populate default values for creation
//   - conds: WHERE conditions to search for existing record
//
// Returns:
//   - bool: true if record was created, false if found
//   - error: database or validation error
//
// Example:
//
//	var user User
//	created, err := neonService.FindOrCreate(ctx, &user,
//	    func() error {
//	        user.Status = "active"
//	        user.CreatedAt = time.Now()
//	        return nil
//	    },
//	    "email = ?", "john@example.com")
func (s *Service) FindOrCreate(ctx context.Context, dest interface{}, createFunc func() error, conds ...interface{}) (bool, error) {
	if s == nil || s.Client == nil {
		return false, gorm.ErrInvalidData
	}

	db := s.Client.WithContext(ctx)
	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	result := db.FirstOrCreate(dest)
	if result.Error != nil {
		return false, result.Error
	}

	// If RowsAffected is 1, the record was created
	if result.RowsAffected == 1 {
		if createFunc != nil {
			if err := createFunc(); err != nil {
				return true, err
			}
		}
		return true, nil
	}

	return false, nil
}

// FindAndUpdate retrieves a record and updates it in one operation.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to the model/struct to populate and update
//   - updateFunc: function to modify the record before update
//   - conds: WHERE conditions
//
// Returns:
//   - error: database error or updateFunc error
//
// Example:
//
//	var user User
//	err := neonService.FindAndUpdate(ctx, &user,
//	    func() error {
//	        user.Status = "inactive"
//	        user.LastLogin = time.Now()
//	        return nil
//	    },
//	    "id = ?", userID)
func (s *Service) FindAndUpdate(ctx context.Context, dest interface{}, updateFunc func() error, conds ...interface{}) error {
	if s == nil || s.Client == nil {
		return gorm.ErrInvalidData
	}

	db := s.Client.WithContext(ctx)
	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	if err := db.First(dest).Error; err != nil {
		return err
	}

	if updateFunc != nil {
		if err := updateFunc(); err != nil {
			return err
		}
	}

	return db.Save(dest).Error
}

// BulkUpsert performs bulk insert or update (upsert) operations.
// Useful for syncing data from external sources.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - records: slice of records to upsert
//   - uniqueFields: field names that determine uniqueness
//
// Returns:
//   - error: if the operation fails
//
// Note: The actual behavior depends on the database constraints.
// Ensure your model has unique constraints on the specified fields.
//
// Example:
//
//	records := []User{
//	    {ID: 1, Name: "John", Email: "john@example.com"},
//	    {ID: 2, Name: "Jane", Email: "jane@example.com"},
//	}
//	err := neonService.BulkUpsert(ctx, records, "email")
func (s *Service) BulkUpsert(ctx context.Context, records interface{}, uniqueFields ...string) error {
	if s == nil || s.Client == nil {
		return gorm.ErrInvalidData
	}

	// Use GORM's Clauses for upsert support (PostgreSQL specific)
	db := s.Client.WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	})

	return db.CreateInBatches(records, 100).Error
}

// Soft delete support (if model has DeletedAt field)
// SoftDelete marks a record as deleted without removing it from database
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - value: the model/struct to soft delete
//   - conds: WHERE conditions
//
// Returns:
//   - *gorm.DB: for checking RowsAffected and Error
//
// Note: Model must have a DeletedAt field for soft delete to work
//
//	type User struct {
//	    ID uint
//	    DeletedAt gorm.DeletedAt `gorm:"index"`
//	}
//
// Example:
//
//	result := neonService.SoftDelete(ctx, &User{}, "id = ?", userID)
func (s *Service) SoftDelete(ctx context.Context, value interface{}, conds ...interface{}) *gorm.DB {
	if s == nil || s.Client == nil {
		return &gorm.DB{Error: gorm.ErrInvalidData}
	}

	db := s.Client.WithContext(ctx)
	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	return db.Delete(value)
}

// GetLastError returns a formatted error message for the last database operation
//
// Parameters:
//   - result: *gorm.DB from a previous operation
//   - operation: name of the operation (e.g., "create", "update")
//
// Returns:
//   - string: formatted error message or empty string if no error
//
// Example:
//
//	result := neonService.Create(ctx, &user)
//	if msg := neonService.GetLastError(result, "create"); msg != "" {
//	    logger.Error(msg)
//	}
func (s *Service) GetLastError(result *gorm.DB, operation string) string {
	if result == nil || result.Error == nil {
		return ""
	}

	return fmt.Sprintf("%s operation failed: %v", operation, result.Error)
}

// Stats returns statistics about the last database operation
//
// Parameters:
//   - result: *gorm.DB from a previous operation
//
// Returns:
//   - map[string]interface{}: map with RowsAffected and Error
//
// Example:
//
//	result := neonService.Update(ctx, updateData, "id = ?", userID)
//	stats := neonService.Stats(result)
//	fmt.Printf("Updated %d rows\n", stats["RowsAffected"])
func (s *Service) Stats(result *gorm.DB) map[string]interface{} {
	return map[string]interface{}{
		"RowsAffected": result.RowsAffected,
		"Error":        result.Error,
	}
}

// Limit adds LIMIT clause to a GORM query (utility for custom queries)
func (s *Service) Limit(db *gorm.DB, limit int) *gorm.DB {
	return db.Limit(limit)
}

// Offset adds OFFSET clause to a GORM query (utility for custom queries)
func (s *Service) Offset(db *gorm.DB, offset int) *gorm.DB {
	return db.Offset(offset)
}

// OrderBy adds ORDER clause to a GORM query (utility for custom queries)
func (s *Service) OrderBy(db *gorm.DB, column string, direction string) *gorm.DB {
	if direction != "" {
		return db.Order(fmt.Sprintf("%s %s", column, direction))
	}
	return db.Order(column)
}

// Joins adds JOIN clause to a GORM query (utility for custom queries)
func (s *Service) Joins(db *gorm.DB, query string) *gorm.DB {
	return db.Joins(query)
}

// JoinInfo holds information about a single join operation
type JoinInfo struct {
	Type      string // "INNER", "LEFT", "RIGHT", "FULL"
	Condition string // e.g., "orders ON users.id = orders.user_id"
}

// InnerJoin performs an INNER JOIN and returns results
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to result slice
//   - joinTable: table to join (e.g., "orders ON users.id = orders.user_id")
//   - columns: columns to select (e.g., "users.id", "users.name", "orders.amount")
//   - conds: WHERE conditions (optional)
//
// Returns:
//   - error if query fails
//
// Example:
//
//	var results []struct {
//	    UserID   uint
//	    UserName string
//	    OrderID  uint
//	    Amount   float64
//	}
//	err := neonService.InnerJoin(ctx, &results,
//	    "orders ON users.id = orders.user_id",
//	    []string{"users.id", "users.name", "orders.id as order_id", "orders.amount"},
//	    "orders.status = ?", "completed")
func (s *Service) InnerJoin(ctx context.Context, dest interface{}, joinTable string, columns []string, conds ...interface{}) error {
	if s == nil || s.Client == nil {
		return gorm.ErrInvalidData
	}

	db := s.Client.WithContext(ctx).Joins("INNER JOIN " + joinTable)

	if len(columns) > 0 {
		db = db.Select(columns)
	}

	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	return db.Find(dest).Error
}

// LeftJoin performs a LEFT JOIN and returns results
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to result slice
//   - joinTable: table to join (e.g., "orders ON users.id = orders.user_id")
//   - columns: columns to select
//   - conds: WHERE conditions (optional)
//
// Returns:
//   - error if query fails
//
// Example:
//
//	var userOrders []struct {
//	    UserName  string
//	    OrderID   sql.NullInt64
//	    Amount    sql.NullFloat64
//	}
//	err := neonService.LeftJoin(ctx, &userOrders,
//	    "orders ON users.id = orders.user_id",
//	    []string{"users.name", "orders.id as order_id", "orders.amount"},
//	    "users.status = ?", "active")
func (s *Service) LeftJoin(ctx context.Context, dest interface{}, joinTable string, columns []string, conds ...interface{}) error {
	if s == nil || s.Client == nil {
		return gorm.ErrInvalidData
	}

	db := s.Client.WithContext(ctx).Joins("LEFT JOIN " + joinTable)

	if len(columns) > 0 {
		db = db.Select(columns)
	}

	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	return db.Find(dest).Error
}

// RightJoin performs a RIGHT JOIN and returns results
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to result slice
//   - joinTable: table to join
//   - columns: columns to select
//   - conds: WHERE conditions (optional)
//
// Returns:
//   - error if query fails
func (s *Service) RightJoin(ctx context.Context, dest interface{}, joinTable string, columns []string, conds ...interface{}) error {
	if s == nil || s.Client == nil {
		return gorm.ErrInvalidData
	}

	db := s.Client.WithContext(ctx).Joins("RIGHT JOIN " + joinTable)

	if len(columns) > 0 {
		db = db.Select(columns)
	}

	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	return db.Find(dest).Error
}

// MultiJoin performs multiple JOINs and returns results
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to result slice
//   - joins: slice of JoinInfo for multiple joins
//   - columns: columns to select
//   - conds: WHERE conditions (optional)
//
// Returns:
//   - error if query fails
//
// Example:
//
//	var results []struct {
//	    UserName    string
//	    ProductName string
//	    Quantity    int
//	}
//	joins := []JoinInfo{
//	    {Type: "INNER", Condition: "orders ON users.id = orders.user_id"},
//	    {Type: "INNER", Condition: "order_items ON orders.id = order_items.order_id"},
//	    {Type: "INNER", Condition: "products ON order_items.product_id = products.id"},
//	}
//	err := neonService.MultiJoin(ctx, &results, joins,
//	    []string{"users.name", "products.name", "order_items.quantity"},
//	    "orders.status = ?", "completed")
func (s *Service) MultiJoin(ctx context.Context, dest interface{}, joins []JoinInfo, columns []string, conds ...interface{}) error {
	if s == nil || s.Client == nil {
		return gorm.ErrInvalidData
	}

	db := s.Client.WithContext(ctx)

	// Apply all joins
	for _, join := range joins {
		db = db.Joins(join.Type + " JOIN " + join.Condition)
	}

	if len(columns) > 0 {
		db = db.Select(columns)
	}

	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	return db.Find(dest).Error
}

// JoinWithOrder performs a JOIN with specific ordering
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to result slice
//   - joinType: "INNER", "LEFT", "RIGHT"
//   - joinTable: table to join
//   - columns: columns to select
//   - orderBy: ORDER BY clause (e.g., "users.created_at DESC")
//   - conds: WHERE conditions (optional)
//
// Returns:
//   - error if query fails
//
// Example:
//
//	var userOrders []struct {
//	    UserName string
//	    OrderID  uint
//	    Amount   float64
//	}
//	err := neonService.JoinWithOrder(ctx, &userOrders,
//	    "LEFT",
//	    "orders ON users.id = orders.user_id",
//	    []string{"users.name", "orders.id as order_id", "orders.amount"},
//	    "orders.created_at DESC",
//	    "users.status = ?", "active")
func (s *Service) JoinWithOrder(ctx context.Context, dest interface{}, joinType string, joinTable string, columns []string, orderBy string, conds ...interface{}) error {
	if s == nil || s.Client == nil {
		return gorm.ErrInvalidData
	}

	db := s.Client.WithContext(ctx).Joins(joinType + " JOIN " + joinTable)

	if len(columns) > 0 {
		db = db.Select(columns)
	}

	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	if orderBy != "" {
		db = db.Order(orderBy)
	}

	return db.Find(dest).Error
}

// JoinWithPagination performs a JOIN with pagination
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to result slice
//   - joinType: "INNER", "LEFT", "RIGHT"
//   - joinTable: table to join
//   - columns: columns to select
//   - page: page number (starting from 1)
//   - pageSize: number of records per page
//   - conds: WHERE conditions (optional)
//
// Returns:
//   - int64: total count of matching records
//   - error if query fails
//
// Example:
//
//	var userOrders []struct {
//	    UserName string
//	    OrderID  uint
//	}
//	total, err := neonService.JoinWithPagination(ctx, &userOrders,
//	    "LEFT",
//	    "orders ON users.id = orders.user_id",
//	    []string{"users.name", "orders.id"},
//	    1, 10,
//	    "users.status = ?", "active")
func (s *Service) JoinWithPagination(ctx context.Context, dest interface{}, joinType string, joinTable string, columns []string, page int, pageSize int, conds ...interface{}) (int64, error) {
	if s == nil || s.Client == nil {
		return 0, gorm.ErrInvalidData
	}

	var count int64
	db := s.Client.WithContext(ctx).Joins(joinType + " JOIN " + joinTable)

	// Count total before pagination
	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}
	if err := db.Model(dest).Count(&count).Error; err != nil {
		return 0, err
	}

	// Apply pagination and columns
	offset := (page - 1) * pageSize
	db = s.Client.WithContext(ctx).Joins(joinType + " JOIN " + joinTable)
	if len(columns) > 0 {
		db = db.Select(columns)
	}
	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	if err := db.Offset(offset).Limit(pageSize).Find(dest).Error; err != nil {
		return count, err
	}

	return count, nil
}

// ========== Type-Safe Join Examples ==========
// NOTE: See examples below for recommended usage patterns with the JoinBuilder

// ExampleTypeSafeInnerJoin demonstrates using the type-safe join builder for INNER JOIN
//
// Example models:
//
//	type User struct {
//	    ID        uint   `gorm:"primaryKey"`
//	    Name      string
//	    Email     string
//	}
//
//	type Order struct {
//	    ID     uint
//	    UserID uint
//	    Amount float64
//	    Status string
//	}
//
// Usage:
//
//	var results []struct {
//	    UserID   uint
//	    UserName string
//	    OrderID  uint
//	    Amount   float64
//	}
//
//	builder := neonService.NewInnerJoinBuilder(&Order{}).
//	    On(&User{}, "ID", "UserID").
//	    Select(map[interface{}][]string{
//	        &User{}: {"ID", "Name"},
//	        &Order{}: {"ID", "Amount"},
//	    })
//
//	err := builder.Execute(neonService, ctx, &results, "orders.status = ?", "completed")
//	if err != nil {
//	    log.Printf("Join failed: %v", err)
//	}
//
//	for _, result := range results {
//	    fmt.Printf("User %d (%s) - Order Amount: %.2f\n", result.UserID, result.UserName, result.Amount)
//	}

// ExampleTypeSafeLeftJoin demonstrates using the type-safe join builder for LEFT JOIN
//
// Usage:
//
//	var userOrders []struct {
//	    UserID   uint
//	    UserName string
//	    OrderID  sql.NullInt64
//	    Amount   sql.NullFloat64
//	}
//
//	builder := neonService.NewLeftJoinBuilder(&Order{}).
//	    On(&User{}, "ID", "UserID").
//	    Select(map[interface{}][]string{
//	        &User{}: {"ID", "Name"},
//	        &Order{}: {"ID", "Amount"},
//	    })
//
//	// Note: Use sql.Null* types for optional columns from joined table
//	err := builder.Execute(neonService, ctx, &userOrders, "users.status = ?", "active")
//	if err != nil {
//	    log.Printf("Join failed: %v", err)
//	}

// ExampleMultipleTableJoins demonstrates joining 3+ tables type-safely
//
// Example models:
//
//	type User struct {
//	    ID   uint `gorm:"primaryKey"`
//	    Name string
//	}
//
//	type Order struct {
//	    ID     uint
//	    UserID uint
//	}
//
//	type OrderItem struct {
//	    ID      uint
//	    OrderID uint
//	}
//
//	type Product struct {
//	    ID   uint
//	    Name string
//	}
//
// Usage (using the old method, but could be extended):
//
//	var complexResults []struct {
//	    UserName    string
//	    ProductName string
//	    Quantity    int
//	}
//
//	joins := []JoinInfo{
//	    {Type: "INNER", Condition: "orders ON users.id = orders.user_id"},
//	    {Type: "INNER", Condition: "order_items ON orders.id = order_items.order_id"},
//	    {Type: "INNER", Condition: "products ON order_items.product_id = products.id"},
//	}
//
//	err := neonService.MultiJoin(ctx, &complexResults, joins,
//	    []string{"users.name", "products.name", "order_items.quantity"},
//	    "orders.status = ?", "completed")

// ExampleSelectAllColumnsFromTable demonstrates selecting all columns from a table
//
// Usage:
//
//	var userOrderPairs []struct {
//	    User  User   `gorm:"embedded"`
//	    Order Order  `gorm:"embedded"`
//	}
//
//	builder := neonService.NewInnerJoinBuilder(&Order{}).
//	    On(&User{}, "ID", "UserID").
//	    SelectAll(&User{}).
//	    SelectAll(&Order{})
//
//	err := builder.Execute(neonService, ctx, &userOrderPairs)

// ExampleChainableJoinBuilder demonstrates method chaining with the join builder
//
// Usage:
//
//	var results []UserOrderResult
//
//	err := neonService.
//	    NewInnerJoinBuilder(&Order{}).
//	    On(&User{}, "ID", "UserID").
//	    Select(map[interface{}][]string{
//	        &User{}: {"ID", "Name", "Email"},
//	        &Order{}: {"ID", "Amount", "Status"},
//	    }).
//	    Execute(neonService, ctx, &results,
//	        "orders.status = ?", "completed",
//	        "users.status = ?", "active")
//
//	if err != nil {
//	    log.Printf("Error: %v", err)
//	    return err
//	}
//
//	return nil

// Helper functions for common column selections

// SelectUserColumns returns common columns to select from User table
func SelectUserColumns() []string {
	return []string{"ID", "Name", "Email", "Status"}
}

// SelectOrderColumns returns common columns to select from Order table
func SelectOrderColumns() []string {
	return []string{"ID", "UserID", "Amount", "Status", "CreatedAt"}
}

// SelectOrderItemColumns returns common columns to select from OrderItem table
func SelectOrderItemColumns() []string {
	return []string{"ID", "OrderID", "ProductID", "Quantity", "Price"}
}
