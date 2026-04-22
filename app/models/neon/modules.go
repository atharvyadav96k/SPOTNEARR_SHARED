package neon

import (
	"sync"

	"gorm.io/gorm"
)

// Service manages the Neon (PostgreSQL) database connection and provides a wrapper
// around GORM for performing database operations.
//
// The Service uses sync.Once to ensure the database connection is initialized only once,
// even when accessed from multiple goroutines. This implements the singleton pattern for
// the database connection.
//
// Key characteristics:
//   - One-time initialization using InitDB()
//   - Context-aware query operations (all methods require context.Context)
//   - Migration management separated from query operations (RegisterAndMigrate is called once at startup)
//   - CRUD operations: Create, Read, Update, Delete
//   - Advanced queries: Pagination, Transactions, Raw SQL
//   - Consistent error handling with helper functions
//
// Setup:
//  1. Set GCP_NEON_DATABASE_URL environment variable with your Neon connection string
//  2. Call InitDB(&neon.Service{}) once at application startup
//  3. Call RegisterAndMigrate() to set up schema
//  4. Use the service methods for queries (no migrations needed in handlers)
//
// Example:
//
//	neonService := neon.InitDB(&neon.Service{})
//	defer neonService.Close()
//	neonService.RegisterAndMigrate(&User{}, &Order{})
//
//	// Later in handlers
//	err := neonService.FindByID(ctx, &user, userID)
//	err = neonService.Create(ctx, &newOrder)
type Service struct {
	// Client is the underlying GORM database client.
	// It is initialized once by InitDB and reused for all database operations.
	Client *gorm.DB

	// Once ensures the database connection is established only once,
	// even if InitDB is called multiple times from different goroutines.
	Once sync.Once
}
