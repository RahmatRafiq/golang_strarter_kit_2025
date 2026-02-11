package helpers

import (
	"fmt"
	"testing"
	"time"

	"golang_starter_kit_2025/app/models"
)

// TestUser wraps a user with its plain password for testing
type TestUser struct {
	User          *models.User
	PlainPassword string
}

// UserFactory creates test users (like Laravel's User::factory())
type UserFactory struct {
	db   *TestDB
	t    *testing.T
	user *models.User
}

// NewUserFactory creates a new user factory
func NewUserFactory(t *testing.T, db *TestDB) *UserFactory {
	t.Helper()

	return &UserFactory{
		db: db,
		t:  t,
		user: &models.User{
			Username: fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
			Email:    fmt.Sprintf("test_%d@example.com", time.Now().UnixNano()),
			Password: "password123", // Will be hashed
		},
	}
}

// WithEmail sets custom email
func (f *UserFactory) WithEmail(email string) *UserFactory {
	f.user.Email = email
	return f
}

// WithUsername sets custom username
func (f *UserFactory) WithUsername(username string) *UserFactory {
	f.user.Username = username
	return f
}

// WithPassword sets custom password (will be hashed)
func (f *UserFactory) WithPassword(password string) *UserFactory {
	f.user.Password = password
	return f
}

// Create saves user to database
func (f *UserFactory) Create() *models.User {
	f.t.Helper()

	// Store plain password for test reference (before it gets hashed by BeforeCreate hook)
	plainPassword := f.user.Password

	// Set default pin if empty (required by BeforeCreate hook)
	if f.user.Pin == "" {
		f.user.Pin = "123456" // Default test PIN
	}

	// GORM will automatically hash password via User.BeforeCreate hook
	if err := f.db.DB.Create(f.user).Error; err != nil {
		f.t.Fatalf("Failed to create user: %v", err)
	}

	// Store plain password in a temporary field for testing purposes
	// (so tests can use the plain password to login)
	f.t.Logf("Created test user: ID=%d, Email=%s, PlainPassword=%s", 
		f.user.ID, f.user.Email, plainPassword)
	
	return f.user
}

// CreateWithPlainPassword creates user and returns both user and plain password
// Useful for authentication tests
func (f *UserFactory) CreateWithPlainPassword() (*models.User, string) {
	f.t.Helper()
	plainPassword := f.user.Password
	user := f.Create()
	return user, plainPassword
}

// RoleFactory creates test roles
type RoleFactory struct {
	db   *TestDB
	t    *testing.T
	role *models.Role
}

// NewRoleFactory creates a new role factory
func NewRoleFactory(t *testing.T, db *TestDB) *RoleFactory {
	t.Helper()

	return &RoleFactory{
		db: db,
		t:  t,
		role: &models.Role{
			Name:  fmt.Sprintf("test_role_%d", time.Now().UnixNano()),
			Group: "test",
		},
	}
}

// WithName sets custom role name
func (f *RoleFactory) WithName(name string) *RoleFactory {
	f.role.Name = name
	return f
}

// WithGroup sets custom group
func (f *RoleFactory) WithGroup(group string) *RoleFactory {
	f.role.Group = group
	return f
}

// Create saves role to database
func (f *RoleFactory) Create() *models.Role {
	f.t.Helper()

	if err := f.db.DB.Create(f.role).Error; err != nil {
		f.t.Fatalf("Failed to create role: %v", err)
	}

	f.t.Logf("Created test role: ID=%d, Name=%s", f.role.ID, f.role.Name)
	return f.role
}

// PermissionFactory creates test permissions
type PermissionFactory struct {
	db         *TestDB
	t          *testing.T
	permission *models.Permission
}

// NewPermissionFactory creates a new permission factory
func NewPermissionFactory(t *testing.T, db *TestDB) *PermissionFactory {
	t.Helper()

	return &PermissionFactory{
		db: db,
		t:  t,
		permission: &models.Permission{
			Name:  fmt.Sprintf("test.permission.%d", time.Now().UnixNano()),
			Group: "test",
		},
	}
}

// WithName sets custom permission name
func (f *PermissionFactory) WithName(name string) *PermissionFactory {
	f.permission.Name = name
	return f
}

// WithGroup sets custom group
func (f *PermissionFactory) WithGroup(group string) *PermissionFactory {
	f.permission.Group = group
	return f
}

// Create saves permission to database
func (f *PermissionFactory) Create() *models.Permission {
	f.t.Helper()

	if err := f.db.DB.Create(f.permission).Error; err != nil {
		f.t.Fatalf("Failed to create permission: %v", err)
	}

	f.t.Logf("Created test permission: ID=%d, Name=%s", f.permission.ID, f.permission.Name)
	return f.permission
}

// ProductFactory creates test products
type ProductFactory struct {
	db      *TestDB
	t       *testing.T
	product *models.Product
}

// NewProductFactory creates a new product factory
func NewProductFactory(t *testing.T, db *TestDB) *ProductFactory {
	t.Helper()

	return &ProductFactory{
		db: db,
		t:  t,
		product: &models.Product{
			Name:        fmt.Sprintf("Test Product %d", time.Now().UnixNano()),
			Description: "Test product created by factory",
			Price:       99.99,
			Stock:       100,
		},
	}
}

// WithName sets custom product name
func (f *ProductFactory) WithName(name string) *ProductFactory {
	f.product.Name = name
	return f
}

// WithPrice sets custom price
func (f *ProductFactory) WithPrice(price float64) *ProductFactory {
	f.product.Price = price
	return f
}

// WithStock sets custom stock
func (f *ProductFactory) WithStock(stock int) *ProductFactory {
	f.product.Stock = stock
	return f
}

// WithCategory sets category
func (f *ProductFactory) WithCategory(categoryID uint) *ProductFactory {
	f.product.CategoryID = categoryID
	return f
}

// Create saves product to database
func (f *ProductFactory) Create() *models.Product {
	f.t.Helper()

	if err := f.db.DB.Create(f.product).Error; err != nil {
		f.t.Fatalf("Failed to create product: %v", err)
	}

	f.t.Logf("Created test product: ID=%d, Name=%s", f.product.ID, f.product.Name)
	return f.product
}

// CategoryFactory creates test categories
type CategoryFactory struct {
	db       *TestDB
	t        *testing.T
	category *models.Category
}

// NewCategoryFactory creates a new category factory
func NewCategoryFactory(t *testing.T, db *TestDB) *CategoryFactory {
	t.Helper()

	return &CategoryFactory{
		db: db,
		t:  t,
		category: &models.Category{
			Category: fmt.Sprintf("Test Category %d", time.Now().UnixNano()),
		},
	}
}

// WithCategory sets custom category name
func (f *CategoryFactory) WithCategory(category string) *CategoryFactory {
	f.category.Category = category
	return f
}

// Create saves category to database
func (f *CategoryFactory) Create() *models.Category {
	f.t.Helper()

	if err := f.db.DB.Create(f.category).Error; err != nil {
		f.t.Fatalf("Failed to create category: %v", err)
	}

	f.t.Logf("Created test category: ID=%d, Category=%s", f.category.ID, f.category.Category)
	return f.category
}

// Helper functions for common test scenarios

// CreateUserWithRole creates a user and assigns a role
func CreateUserWithRole(t *testing.T, db *TestDB, roleName string) (*models.User, *models.Role) {
	t.Helper()

	role := NewRoleFactory(t, db).WithName(roleName).Create()
	user := NewUserFactory(t, db).Create()

	// Assign role to user
	if err := db.DB.Exec("INSERT INTO user_has_roles (user_id, role_id) VALUES (?, ?)", user.ID, role.ID).Error; err != nil {
		t.Fatalf("Failed to assign role to user: %v", err)
	}

	t.Logf("Assigned role '%s' to user ID=%d", roleName, user.ID)
	return user, role
}

// CreateRoleWithPermissions creates a role with specific permissions
func CreateRoleWithPermissions(t *testing.T, db *TestDB, roleName string, permissionNames []string) *models.Role {
	t.Helper()

	role := NewRoleFactory(t, db).WithName(roleName).Create()

	for _, permName := range permissionNames {
		perm := NewPermissionFactory(t, db).WithName(permName).Create()

		if err := db.DB.Exec("INSERT INTO role_has_permissions (role_id, permission_id) VALUES (?, ?)", role.ID, perm.ID).Error; err != nil {
			t.Fatalf("Failed to assign permission to role: %v", err)
		}
	}

	t.Logf("Created role '%s' with %d permissions", roleName, len(permissionNames))
	return role
}
