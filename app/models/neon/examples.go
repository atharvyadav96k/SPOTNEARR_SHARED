package neon

/*
EXAMPLES - Real-world usage patterns for the Neon wrapper

This file contains commented examples showing how to use the Neon wrapper
in typical application scenarios. Uncomment and adapt to your needs.
*/

// Example 1: User Registration with validation and transaction
/*
func RegisterUser(ctx context.Context, neon *neon.Service, userEmail, userName string) (uint, error) {
	// Validate input
	if userEmail == "" || userName == "" {
		return 0, fmt.Errorf("email and name are required")
	}

	// Check if user already exists
	exists, err := neon.Exists(ctx, &User{}, "email = ?", userEmail)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, fmt.Errorf("user with this email already exists")
	}

	// Create user within transaction
	var user User
	err = neon.Transaction(ctx, func(tx *gorm.DB) error {
		user = User{
			Email:     userEmail,
			Name:      userName,
			Status:    "active",
			CreatedAt: time.Now(),
		}
		return tx.Create(&user).Error
	})

	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
*/

// Example 2: Pagination with filtering
/*
func GetActiveUsers(ctx context.Context, neon *neon.Service, page int, pageSize int) ([]User, int64, error) {
	var users []User
	total, err := neon.Paginate(ctx, &users, page, pageSize, "status = ?", "active")
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
*/

// Example 3: Complex query with joins and ordering
/*
func GetUserOrdersWithDetails(ctx context.Context, neon *neon.Service, userID uint) ([]Order, error) {
	var orders []Order
	err := neon.FindWithQuery(ctx, func(db *gorm.DB) *gorm.DB {
		return db.
			Joins("LEFT JOIN users ON orders.user_id = users.id").
			Where("orders.user_id = ?", userID).
			Order("orders.created_at DESC").
			Preload("Items")  // If Order has Items relationship
	}, &orders)
	return orders, err
}
*/

// Example 4: Bulk import with validation
/*
func BulkImportUsers(ctx context.Context, neon *neon.Service, users []User) error {
	// Validate all records
	for i, user := range users {
		if user.Email == "" {
			return fmt.Errorf("user at index %d has empty email", i)
		}
	}

	// Insert in batches
	return neon.CreateBatch(ctx, users, 500)
}
*/

// Example 5: Find or create pattern
/*
func EnsureDefaultAdmin(ctx context.Context, neon *neon.Service) error {
	var admin User
	created, err := neon.FindOrCreate(ctx, &admin,
		func() error {
			admin.Email = "admin@example.com"
			admin.Name = "Admin"
			admin.Role = "admin"
			admin.Status = "active"
			admin.CreatedAt = time.Now()
			return nil
		},
		"email = ?", "admin@example.com")

	if err != nil {
		return err
	}

	if created {
		fmt.Println("Default admin user created")
	} else {
		fmt.Println("Default admin user already exists")
	}

	return nil
}
*/

// Example 6: Soft delete (requires DeletedAt field)
/*
func DeactivateUser(ctx context.Context, neon *neon.Service, userID uint) error {
	result := neon.SoftDelete(ctx, &User{}, "id = ?", userID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
*/

// Example 7: Update with condition checking
/*
func PromoteUser(ctx context.Context, neon *neon.Service, userID uint) error {
	result := neon.UpdateColumns(ctx,
		map[string]interface{}{
			"role":       "admin",
			"updated_at": time.Now(),
		},
		"id = ? AND status = ?", userID, "active")

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found or not active")
	}

	return nil
}
*/

// Example 8: Multiple updates in transaction
/*
func TransferBalance(ctx context.Context, neon *neon.Service, fromUserID, toUserID uint, amount decimal.Decimal) error {
	return neon.Transaction(ctx, func(tx *gorm.DB) error {
		// Deduct from source user
		if err := tx.Model(&Account{}).
			Where("user_id = ?", fromUserID).
			Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		// Add to destination user
		if err := tx.Model(&Account{}).
			Where("user_id = ?", toUserID).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		// Log transaction
		return tx.Create(&Transaction{
			FromUserID: fromUserID,
			ToUserID:   toUserID,
			Amount:     amount,
			CreatedAt:  time.Now(),
		}).Error
	})
}
*/

// Example 9: Count and aggregation
/*
func GetUserStats(ctx context.Context, neon *neon.Service) (map[string]int64, error) {
	stats := make(map[string]int64)

	// Count active users
	activeCount, err := neon.Count(ctx, &User{}, "status = ?", "active")
	if err != nil {
		return nil, err
	}
	stats["active_users"] = activeCount

	// Count inactive users
	inactiveCount, err := neon.Count(ctx, &User{}, "status = ?", "inactive")
	if err != nil {
		return nil, err
	}
	stats["inactive_users"] = inactiveCount

	return stats, nil
}
*/

// Example 10: Raw SQL for complex analytics
/*
func GetUserActivitySummary(ctx context.Context, neon *neon.Service, days int) ([]ActivitySummary, error) {
	var summary []ActivitySummary
	err := neon.Raw(ctx, &summary, `
		SELECT
			u.id,
			u.name,
			COUNT(o.id) as total_orders,
			SUM(o.amount) as total_spent,
			MAX(o.created_at) as last_order_date
		FROM users u
		LEFT JOIN orders o ON u.id = o.user_id
		WHERE o.created_at >= NOW() - INTERVAL '?' DAY
		GROUP BY u.id, u.name
		ORDER BY total_spent DESC
	`, days)

	return summary, err
}
*/

// Example 11: Error handling pattern
/*
func GetUserWithErrorHandling(ctx context.Context, neon *neon.Service, userID uint) (*User, error) {
	var user User
	err := neon.FindByID(ctx, &user, userID)

	if err != nil {
		if neon.IsRecordNotFound(err) {
			return nil, fmt.Errorf("user %d not found", userID)
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}
*/

// Example 12: Validation before create/update
/*
func CreateProductWithValidation(ctx context.Context, neon *neon.Service, product *Product) error {
	err := neon.CreateWithValidation(ctx, product, func() error {
		// Validate product data
		if product.Name == "" {
			return fmt.Errorf("product name is required")
		}
		if product.Price <= 0 {
			return fmt.Errorf("product price must be positive")
		}
		if product.SKU == "" {
			return fmt.Errorf("product SKU is required")
		}

		// Check if SKU already exists
		exists, err := neon.Exists(ctx, &Product{}, "sku = ?", product.SKU)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("product with SKU %s already exists", product.SKU)
		}

		return nil
	})

	return err
}
*/

// Example 13: Batch operations with error handling
/*
func UpdateUserStatuses(ctx context.Context, neon *neon.Service, userIDs []uint, status string) (int64, error) {
	result := neon.Update(ctx, User{Status: status}, "id IN ?", userIDs)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
*/

// Example 14: Pagination helper
/*
type PaginationResult struct {
	Data      interface{} `json:"data"`
	Total     int64       `json:"total"`
	Page      int         `json:"page"`
	PageSize  int         `json:"page_size"`
	TotalPages int64      `json:"total_pages"`
}

func GetPaginatedUsers(ctx context.Context, neon *neon.Service, page, pageSize int) (*PaginationResult, error) {
	var users []User
	total, err := neon.Paginate(ctx, &users, page, pageSize, "status = ?", "active")
	if err != nil {
		return nil, err
	}

	return &PaginationResult{
		Data:       users,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (total + int64(pageSize) - 1) / int64(pageSize),
	}, nil
}
*/

// Example 15: Initialization in main.go
/*
func main() {
	// Initialize database service
	neonService := neon.InitDB(&neon.Service{})
	defer neonService.Close()

	// Register all models (migrations) once at startup
	neonService.RegisterAndMigrate(
		&User{},
		&Product{},
		&Order{},
		&Transaction{},
	)

	// Create app with all services
	app := &App{
		Neon: neonService,
		// ... other services
	}

	// Now ready to use in handlers
	router.GET("/api/users/:id", func(c *gin.Context) {
		userID := c.Param("id")
		var user User
		err := app.Neon.FindByID(c.Request.Context(), &user, userID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, user)
	})

	router.Run(":8080")
}
*/

// Model definitions for examples (uncomment and adapt)
/*
type User struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string
	Email     string         `gorm:"unique"`
	Password  string
	Status    string         `gorm:"default:active"`
	Role      string         `gorm:"default:user"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Order struct {
	ID        uint
	UserID    uint
	Amount    decimal.Decimal
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Items     []OrderItem    `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID      uint
	OrderID uint
	Product Product
}

type Product struct {
	ID        uint
	Name      string
	SKU       string `gorm:"unique"`
	Price     decimal.Decimal
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Transaction struct {
	ID        uint
	FromUserID uint
	ToUserID   uint
	Amount    decimal.Decimal
	CreatedAt time.Time
}

type ActivitySummary struct {
	ID           uint
	Name         string
	TotalOrders  int64
	TotalSpent   decimal.Decimal
	LastOrderDate *time.Time
}
*/
