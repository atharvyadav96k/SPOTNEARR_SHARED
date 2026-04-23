package neon

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewClient() *Service {
	return &Service{}
}

// InitDB initializes the database connection using the GCP_NEON_DATABASE_URL environment variable.
// It ensures the connection is created only once using sync.Once.
// Set the GCP_NEON_DATABASE_URL environment variable to the connection string for your Neon database before calling InitDB.
//
// Example usage:
//
//	neonService := neon.InitDB(&neon.Service{})
func InitDB(s *Service) *Service {
	s.Once.Do(func() {
		dsn := os.Getenv("GCP_NEON_DATABASE_URL")
		var err error
		s.Client, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to Neon: %v", err)
		}
	})
	return s
}

// RegisterAndMigrate takes a variable number of model structs, performs auto-migration on the Neon database, and logs the process.
// This method is designed to be called after the database connection has been initialized.
// Migrations should be registered once at application startup, not in every function.
//
// Example usage:
//
//	neonService := neon.InitDB(&neon.Service{})
//	neonService.RegisterAndMigrate(&User{}, &Product{})
func (s *Service) RegisterAndMigrate(models ...interface{}) error {
	db := s.GetDB()

	log.Printf("Starting migration for %d models...", len(models))
	return db.AutoMigrate(models...)
}

// Close shuts down the database connection and releases underlying resources.
//
// It is safe to call multiple times. If the client is not initialized,
// the function simply returns nil.
//
// Returns:
//   - error if closing the connection fails
func (s *Service) Close() error {
	if s == nil || s.Client == nil {
		return nil
	}

	sqlDB, err := s.Client.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
