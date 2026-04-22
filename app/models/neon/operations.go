package neon

import (
	"context"

	"gorm.io/gorm"
)

// Create inserts a single record into the database.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - value: the model/struct to insert
//
// Returns:
//   - error if the insert fails
//
// Example:
//
//	user := &User{Name: "John", Email: "john@example.com"}
//	err := neonService.Create(ctx, user)
func (s *Service) Create(ctx context.Context, value interface{}) error {
	if s == nil || s.Client == nil {
		return gorm.ErrInvalidData
	}

	return s.Client.WithContext(ctx).Create(value).Error
}

// CreateBatch inserts multiple records into the database in a single operation.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - records: slice of models/structs to insert
//   - batchSize: number of records to insert in each batch (recommended: 100-1000)
//
// Returns:
//   - error if the batch insert fails
//
// Example:
//
//	users := []User{{Name: "John"}, {Name: "Jane"}}
//	err := neonService.CreateBatch(ctx, users, 100)
func (s *Service) CreateBatch(ctx context.Context, records interface{}, batchSize int) error {
	if s == nil || s.Client == nil {
		return gorm.ErrInvalidData
	}

	return s.Client.WithContext(ctx).CreateInBatches(records, batchSize).Error
}

// FindByID retrieves a single record by primary key.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to the model/struct to populate
//   - conds: primary key value(s) or additional conditions
//
// Returns:
//   - error if the query fails or record not found (gorm.ErrRecordNotFound)
//
// Example:
//
//	var user User
//	err := neonService.FindByID(ctx, &user, userID)
func (s *Service) FindByID(ctx context.Context, dest interface{}, conds ...interface{}) error {
	if s == nil || s.Client == nil {
		return gorm.ErrInvalidData
	}

	return s.Client.WithContext(ctx).First(dest, conds...).Error
}

// FindOne retrieves a single record based on conditions.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to the model/struct to populate
//   - conds: WHERE conditions (field, value pairs or inline conditions)
//
// Returns:
//   - error if the query fails or record not found
//
// Example:
//
//	var user User
//	err := neonService.FindOne(ctx, &user, "email = ?", "john@example.com")
//	// or
//	err := neonService.FindOne(ctx, &user, User{Email: "john@example.com"})
func (s *Service) FindOne(ctx context.Context, dest interface{}, conds ...interface{}) error {
	if s == nil || s.Client == nil {
		return gorm.ErrInvalidData
	}

	db := s.Client.WithContext(ctx)
	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	return db.First(dest).Error
}

// FindAll retrieves all records matching the conditions.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to a slice of models/structs to populate
//   - conds: WHERE conditions (optional)
//
// Returns:
//   - error if the query fails
//
// Example:
//
//	var users []User
//	err := neonService.FindAll(ctx, &users)
//	// with conditions
//	err := neonService.FindAll(ctx, &users, "status = ?", "active")
func (s *Service) FindAll(ctx context.Context, dest interface{}, conds ...interface{}) error {
	if s == nil || s.Client == nil {
		return gorm.ErrInvalidData
	}

	db := s.Client.WithContext(ctx)
	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	return db.Find(dest).Error
}

// FindWithQuery retrieves records using a custom GORM query builder.
// This method provides flexibility for complex queries with joins, sorting, pagination, etc.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - queryFunc: function that builds the GORM query
//   - dest: pointer to a slice of models/structs to populate
//
// Returns:
//   - error if the query fails
//
// Example:
//
//	var users []User
//	err := neonService.FindWithQuery(ctx, func(db *gorm.DB) *gorm.DB {
//	    return db.Where("status = ?", "active").Order("created_at DESC").Limit(10)
//	}, &users)
func (s *Service) FindWithQuery(ctx context.Context, queryFunc func(*gorm.DB) *gorm.DB, dest interface{}) error {
	if s == nil || s.Client == nil {
		return gorm.ErrInvalidData
	}

	db := s.Client.WithContext(ctx)
	return queryFunc(db).Find(dest).Error
}

// Update updates records based on conditions.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - value: the model/struct with updated values (can use a struct or a map)
//   - conds: WHERE conditions (field, value pairs)
//
// Returns:
//   - *gorm.DB for chainable operations (check RowsAffected, Error)
//
// Example:
//
//	result := neonService.Update(ctx, User{Status: "inactive"}, "id = ?", userID)
//	if result.Error != nil { /* handle error */ }
//	fmt.Println("Updated rows:", result.RowsAffected)
func (s *Service) Update(ctx context.Context, value interface{}, conds ...interface{}) *gorm.DB {
	if s == nil || s.Client == nil {
		return &gorm.DB{Error: gorm.ErrInvalidData}
	}

	db := s.Client.WithContext(ctx)
	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	return db.Model(value).Updates(value)
}

// UpdateByID updates a record by its primary key.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: the model/struct with updated values
//   - id: primary key value
//
// Returns:
//   - *gorm.DB for chainable operations
//
// Example:
//
//	user := User{ID: 1, Name: "Updated Name"}
//	result := neonService.UpdateByID(ctx, user, user.ID)
func (s *Service) UpdateByID(ctx context.Context, dest interface{}, id interface{}) *gorm.DB {
	if s == nil || s.Client == nil {
		return &gorm.DB{Error: gorm.ErrInvalidData}
	}

	return s.Client.WithContext(ctx).Model(dest).Where("id = ?", id).Updates(dest)
}

// UpdateColumns updates specific columns of records matching conditions.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - updates: map of column names to new values (map[string]interface{})
//   - conds: WHERE conditions
//
// Returns:
//   - *gorm.DB for chainable operations
//
// Example:
//
//	result := neonService.UpdateColumns(ctx,
//	    map[string]interface{}{"status": "inactive", "updated_at": time.Now()},
//	    "id = ?", userID)
func (s *Service) UpdateColumns(ctx context.Context, updates map[string]interface{}, conds ...interface{}) *gorm.DB {
	if s == nil || s.Client == nil {
		return &gorm.DB{Error: gorm.ErrInvalidData}
	}

	db := s.Client.WithContext(ctx)
	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	return db.Model(nil).Updates(updates)
}

// Delete deletes records based on conditions.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - value: the model/struct (or its type)
//   - conds: WHERE conditions
//
// Returns:
//   - *gorm.DB for chainable operations (check RowsAffected, Error)
//
// Example:
//
//	result := neonService.Delete(ctx, &User{}, "id = ?", userID)
//	if result.Error != nil { /* handle error */ }
func (s *Service) Delete(ctx context.Context, value interface{}, conds ...interface{}) *gorm.DB {
	if s == nil || s.Client == nil {
		return &gorm.DB{Error: gorm.ErrInvalidData}
	}

	db := s.Client.WithContext(ctx)
	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	return db.Delete(value)
}

// DeleteByID deletes a record by its primary key.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: the model/struct type
//   - id: primary key value
//
// Returns:
//   - *gorm.DB for chainable operations
//
// Example:
//
//	result := neonService.DeleteByID(ctx, &User{}, userID)
func (s *Service) DeleteByID(ctx context.Context, dest interface{}, id interface{}) *gorm.DB {
	if s == nil || s.Client == nil {
		return &gorm.DB{Error: gorm.ErrInvalidData}
	}

	return s.Client.WithContext(ctx).Delete(dest, id)
}

// Count returns the number of records matching the conditions.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - model: the model/struct type
//   - conds: WHERE conditions (optional)
//
// Returns:
//   - int64: the count of matching records
//   - error: if the query fails
//
// Example:
//
//	count, err := neonService.Count(ctx, &User{}, "status = ?", "active")
func (s *Service) Count(ctx context.Context, model interface{}, conds ...interface{}) (int64, error) {
	if s == nil || s.Client == nil {
		return 0, gorm.ErrInvalidData
	}

	var count int64
	db := s.Client.WithContext(ctx).Model(model)
	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	err := db.Count(&count).Error
	return count, err
}

// Exists checks if any record exists matching the conditions.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - model: the model/struct type
//   - conds: WHERE conditions
//
// Returns:
//   - bool: true if record exists, false otherwise
//   - error: if the query fails
//
// Example:
//
//	exists, err := neonService.Exists(ctx, &User{}, "email = ?", "john@example.com")
func (s *Service) Exists(ctx context.Context, model interface{}, conds ...interface{}) (bool, error) {
	if s == nil || s.Client == nil {
		return false, gorm.ErrInvalidData
	}

	var count int64
	db := s.Client.WithContext(ctx).Model(model).Limit(1)
	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	err := db.Count(&count).Error
	return count > 0, err
}

// Paginate applies pagination to a query.
// Returns the results, total count, and error.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to a slice of models/structs to populate
//   - page: page number (starting from 1)
//   - pageSize: number of records per page
//   - conds: WHERE conditions (optional)
//
// Returns:
//   - int64: total count of matching records
//   - error: if the query fails
//
// Example:
//
//	var users []User
//	total, err := neonService.Paginate(ctx, &users, 1, 10, "status = ?", "active")
func (s *Service) Paginate(ctx context.Context, dest interface{}, page int, pageSize int, conds ...interface{}) (int64, error) {
	if s == nil || s.Client == nil {
		return 0, gorm.ErrInvalidData
	}

	var count int64
	db := s.Client.WithContext(ctx)
	if len(conds) > 0 {
		db = db.Where(conds[0], conds[1:]...)
	}

	// Get total count
	if err := db.Model(dest).Count(&count).Error; err != nil {
		return 0, err
	}

	// Calculate offset and apply pagination
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(dest).Error; err != nil {
		return count, err
	}

	return count, nil
}

// Raw executes a raw SQL query and scans the results.
// Use this for complex queries not easily expressed with GORM methods.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to the result (struct or slice)
//   - sql: raw SQL query string
//   - values: query parameter values
//
// Returns:
//   - error: if the query fails
//
// Example:
//
//	var user User
//	err := neonService.Raw(ctx, &user, "SELECT * FROM users WHERE id = ?", userID)
func (s *Service) Raw(ctx context.Context, dest interface{}, sql string, values ...interface{}) error {
	if s == nil || s.Client == nil {
		return gorm.ErrInvalidData
	}

	return s.Client.WithContext(ctx).Raw(sql, values...).Scan(dest).Error
}

// Exec executes raw SQL without expecting results.
// Use this for INSERT, UPDATE, DELETE operations with custom SQL.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - sql: raw SQL query string
//   - values: query parameter values
//
// Returns:
//   - *gorm.DB for checking RowsAffected and Error
//
// Example:
//
//	result := neonService.Exec(ctx, "UPDATE users SET status = ? WHERE status = ?", "active", "pending")
//	if result.Error != nil { /* handle error */ }
//	fmt.Println("Updated rows:", result.RowsAffected)
func (s *Service) Exec(ctx context.Context, sql string, values ...interface{}) *gorm.DB {
	if s == nil || s.Client == nil {
		return &gorm.DB{Error: gorm.ErrInvalidData}
	}

	return s.Client.WithContext(ctx).Exec(sql, values...)
}
