package neon

import "gorm.io/gorm"

// GetDB returns the initialized GORM database client from the Service struct. It assumes that the database connection has already been established using the InitDB function. This method provides a convenient way to access the database client for performing queries and operations on the Neon database throughout the application.
// Example usage:
//	neonService := neon.InitNeon()
//	dbClient := neonService.GetDB()
//	// Use dbClient for database operations
func (s *Service) GetDB() *gorm.DB {
	return s.Client
}
