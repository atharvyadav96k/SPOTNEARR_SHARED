package app

import (
	"context"
	"fmt"

	"github.com/atharvyadav96k/gcp/app/models/neon"
)

// RegisterModels performs auto-migration for the given models in Neon database.
// Call this once at application startup, not in every function.
//
// Parameters:
//   - models: variable number of model structs to migrate
//
// Example:
//
//	app.RegisterModels(&User{}, &Product{}, &Order{})
func (a *App) RegisterModels(models ...interface{}) error {
	if a.Neon == nil {
		fmt.Println("Warning: Neon service not initialized")
		return nil
	}
	return a.Neon.RegisterAndMigrate(models...)
}

// CreateRecord creates a single record in the database with validation.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - record: the model/struct to insert
//   - validationFunc: optional validation function (can be nil)
//
// Returns:
//   - error if creation fails
//
// Example:
//
//	user := &User{Name: "John", Email: "john@example.com"}
//	err := app.CreateRecord(ctx, user, nil)
func (a *App) CreateRecord(ctx context.Context, record interface{}, validationFunc func() error) error {
	if a.Neon == nil {
		return fmt.Errorf("Neon service not initialized")
	}

	if validationFunc != nil {
		return a.Neon.CreateWithValidation(ctx, record, validationFunc)
	}
	return a.Neon.Create(ctx, record)
}

// CreateRecords creates multiple records in batches.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - records: slice of records to create
//   - batchSize: number of records per batch (e.g., 100)
//
// Returns:
//   - error if creation fails
//
// Example:
//
//	users := []User{{Name: "John"}, {Name: "Jane"}}
//	err := app.CreateRecords(ctx, users, 100)
func (a *App) CreateRecords(ctx context.Context, records interface{}, batchSize int) error {
	if a.Neon == nil {
		return fmt.Errorf("Neon service not initialized")
	}
	return a.Neon.CreateBatch(ctx, records, batchSize)
}

// GetRecord retrieves a single record by ID.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to model to populate
//   - id: the ID value to search for
//
// Returns:
//   - error if retrieval fails
//
// Example:
//
//	var user User
//	err := app.GetRecord(ctx, &user, userID)
func (a *App) GetRecord(ctx context.Context, dest interface{}, id interface{}) error {
	if a.Neon == nil {
		return fmt.Errorf("Neon service not initialized")
	}
	return a.Neon.FindByID(ctx, dest, id)
}

// GetRecords retrieves all records, optionally with conditions.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to slice to populate
//   - conditions: optional WHERE conditions (cond, val, cond, val...)
//
// Returns:
//   - error if retrieval fails
//
// Example:
//
//	var users []User
//	err := app.GetRecords(ctx, &users, "status = ?", "active")
func (a *App) GetRecords(ctx context.Context, dest interface{}, conditions ...interface{}) error {
	if a.Neon == nil {
		return fmt.Errorf("Neon service not initialized")
	}
	return a.Neon.FindAll(ctx, dest, conditions...)
}

// GetOne retrieves a single record matching conditions.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to model to populate
//   - conditions: WHERE conditions (cond, val, cond, val...)
//
// Returns:
//   - error if retrieval fails
//
// Example:
//
//	var user User
//	err := app.GetOne(ctx, &user, "email = ?", "john@example.com")
func (a *App) GetOne(ctx context.Context, dest interface{}, conditions ...interface{}) error {
	if a.Neon == nil {
		return fmt.Errorf("Neon service not initialized")
	}
	return a.Neon.FindOne(ctx, dest, conditions...)
}

// UpdateRecord updates a record by ID.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - record: the model/struct with updated values
//   - id: the ID value to match
//
// Returns:
//   - error if update fails
//
// Example:
//
//	user := &User{ID: 1, Name: "Updated Name"}
//	err := app.UpdateRecord(ctx, user, 1)
func (a *App) UpdateRecord(ctx context.Context, record interface{}, id interface{}) error {
	if a.Neon == nil {
		return fmt.Errorf("Neon service not initialized")
	}
	return a.Neon.Update(ctx, record, "id = ?", id).Error
}

// UpdateRecords updates multiple records matching conditions.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - updates: map of field names to values
//   - conditions: WHERE conditions (cond, val, cond, val...)
//
// Returns:
//   - error if update fails
//
// Example:
//
//	err := app.UpdateRecords(ctx,
//	    map[string]interface{}{"status": "inactive"},
//	    "created_at < ?", time.Now().AddDate(0, -1, 0))
func (a *App) UpdateRecords(ctx context.Context, updates map[string]interface{}, conditions ...interface{}) error {
	if a.Neon == nil {
		return fmt.Errorf("Neon service not initialized")
	}

	// Build WHERE clause from conditions
	if len(conditions) == 0 {
		return fmt.Errorf("at least one condition is required for safety")
	}

	result := a.Neon.Client.WithContext(ctx).
		Where(conditions[0], conditions[1:]...).
		Updates(updates)

	return result.Error
}

// DeleteRecord deletes a record by ID.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - model: the model type to delete from
//   - id: the ID value to match
//
// Returns:
//   - error if deletion fails
//
// Example:
//
//	err := app.DeleteRecord(ctx, &User{}, userID)
func (a *App) DeleteRecord(ctx context.Context, model interface{}, id interface{}) error {
	if a.Neon == nil {
		return fmt.Errorf("Neon service not initialized")
	}
	return a.Neon.Delete(ctx, model, "id = ?", id).Error
}

// SoftDeleteRecord marks a record as deleted without removing it from database.
// Model must have a DeletedAt field.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - model: the model type to soft delete
//   - id: the ID value to match
//
// Returns:
//   - error if deletion fails
//
// Example:
//
//	err := app.SoftDeleteRecord(ctx, &User{}, userID)
func (a *App) SoftDeleteRecord(ctx context.Context, model interface{}, id interface{}) error {
	if a.Neon == nil {
		return fmt.Errorf("Neon service not initialized")
	}
	return a.Neon.SoftDelete(ctx, model, "id = ?", id).Error
}

// UpsertRecords performs bulk insert or update operations.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - records: slice of records to upsert
//
// Returns:
//   - error if operation fails
//
// Example:
//
//	records := []User{{ID: 1, Name: "John"}, {ID: 2, Name: "Jane"}}
//	err := app.UpsertRecords(ctx, records)
func (a *App) UpsertRecords(ctx context.Context, records interface{}) error {
	if a.Neon == nil {
		return fmt.Errorf("Neon service not initialized")
	}
	return a.Neon.BulkUpsert(ctx, records)
}

// FindOrCreateRecord retrieves a record or creates it if not found.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to model to populate or create
//   - createFunc: function to set default values for creation
//   - conditions: WHERE conditions to search
//
// Returns:
//   - bool: true if record was created, false if found
//   - error if operation fails
//
// Example:
//
//	var user User
//	created, err := app.FindOrCreateRecord(ctx, &user,
//	    func() error {
//	        user.Status = "active"
//	        return nil
//	    },
//	    "email = ?", "john@example.com")
func (a *App) FindOrCreateRecord(ctx context.Context, dest interface{}, createFunc func() error, conditions ...interface{}) (bool, error) {
	if a.Neon == nil {
		return false, fmt.Errorf("Neon service not initialized")
	}
	return a.Neon.FindOrCreate(ctx, dest, createFunc, conditions...)
}

// FindAndUpdateRecord retrieves and updates a record in one operation.
//
// Parameters:
//   - ctx: context for timeout and cancellation
//   - dest: pointer to model to populate and update
//   - updateFunc: function to modify the record before update
//   - conditions: WHERE conditions
//
// Returns:
//   - error if operation fails
//
// Example:
//
//	var user User
//	err := app.FindAndUpdateRecord(ctx, &user,
//	    func() error {
//	        user.Status = "inactive"
//	        return nil
//	    },
//	    "id = ?", userID)
func (a *App) FindAndUpdateRecord(ctx context.Context, dest interface{}, updateFunc func() error, conditions ...interface{}) error {
	if a.Neon == nil {
		return fmt.Errorf("Neon service not initialized")
	}
	return a.Neon.FindAndUpdate(ctx, dest, updateFunc, conditions...)
}

// TypeSafeInnerJoin creates a type-safe INNER JOIN builder.
//
// Parameters:
//   - toModel: the model to join with (e.g., &Order{})
//
// Returns:
//   - *neon.JoinBuilder for method chaining
//
// Example:
//
//	var results []UserOrder
//	err := app.TypeSafeInnerJoin(&Order{}).
//	    On(&User{}, "ID", "UserID").
//	    Select(map[interface{}][]string{
//	        &User{}: {"ID", "Name"},
//	        &Order{}: {"ID", "Amount"},
//	    }).
//	    Execute(app.Neon, ctx, &results)
func (a *App) TypeSafeInnerJoin(toModel interface{}) *neon.JoinBuilder {
	if a.Neon == nil {
		return nil
	}
	return a.Neon.NewInnerJoinBuilder(toModel)
}

// TypeSafeLeftJoin creates a type-safe LEFT JOIN builder.
//
// Parameters:
//   - toModel: the model to join with
//
// Returns:
//   - *neon.JoinBuilder for method chaining
//
// Example:
//
//	builder := app.TypeSafeLeftJoin(&Order{})
func (a *App) TypeSafeLeftJoin(toModel interface{}) *neon.JoinBuilder {
	if a.Neon == nil {
		return nil
	}
	return a.Neon.NewLeftJoinBuilder(toModel)
}

// TypeSafeRightJoin creates a type-safe RIGHT JOIN builder.
//
// Parameters:
//   - toModel: the model to join with
//
// Returns:
//   - *neon.JoinBuilder for method chaining
//
// Example:
//
//	builder := app.TypeSafeRightJoin(&Order{})
func (a *App) TypeSafeRightJoin(toModel interface{}) *neon.JoinBuilder {
	if a.Neon == nil {
		return nil
	}
	return a.Neon.NewRightJoinBuilder(toModel)
}

// GetNeonService returns the Neon database service for direct access.
// Use this when you need advanced operations not covered by the wrapper functions.
//
// Returns:
//   - *neon.Service: the Neon service instance
//
// Example:
//
//	neonService := app.GetNeonService()
//	err := neonService.FindAll(ctx, &users)
func (a *App) GetNeonService() *neon.Service {
	return a.Neon
}
