package integration

import (
	"testing"

	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/repositories"
	"golang_starter_kit_2025/app/services"
	"golang_starter_kit_2025/tests/helpers"
)

// TestUserService_CRUD_Integration tests complete user CRUD operations
func TestUserService_CRUD_Integration(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	userRepo := repositories.NewUserRepository(testDB.DB)
	userService := services.NewUserService(userRepo)

	t.Run("success - list users with pagination", func(t *testing.T) {
		// Arrange: Create test users
		helpers.NewUserFactory(t, testDB).WithEmail("user1@test.com").Create()
		helpers.NewUserFactory(t, testDB).WithEmail("user2@test.com").Create()
		helpers.NewUserFactory(t, testDB).WithEmail("user3@test.com").Create()

		// Act: List users
		users, total, err := userService.List(1, 10)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if total != 3 {
			t.Errorf("Expected total 3, got %d", total)
		}
		if len(users) != 3 {
			t.Errorf("Expected 3 users, got %d", len(users))
		}
		t.Logf("✓ Listed %d users successfully", len(users))
	})

	t.Run("success - find user by id", func(t *testing.T) {
		// Arrange: Create test user
		createdUser := helpers.NewUserFactory(t, testDB).
			WithEmail("find@test.com").
			WithUsername("findme").
			Create()

		// Act: Find user
		foundUser, err := userService.FindByID(createdUser.ID)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if foundUser.Email != createdUser.Email {
			t.Errorf("Expected email %s, got %s", createdUser.Email, foundUser.Email)
		}
		if foundUser.Username != createdUser.Username {
			t.Errorf("Expected username %s, got %s", createdUser.Username, foundUser.Username)
		}
		t.Logf("✓ Found user: %s", foundUser.Email)
	})

	t.Run("success - delete user", func(t *testing.T) {
		// Arrange: Create user to delete
		user := helpers.NewUserFactory(t, testDB).
			WithEmail("delete@test.com").
			Create()

		// Act: Delete user
		err := userService.DeleteByID(user.ID)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Verify user soft deleted (DeletedAt should be set)
		var deletedUser models.User
		testDB.DB.Unscoped().Where("id = ?", user.ID).First(&deletedUser)
		if deletedUser.DeletedAt.Time.IsZero() {
			t.Error("Expected user to be soft deleted (DeletedAt set), but DeletedAt is null")
		}

		// Verify user not found in normal queries (soft delete hides it)
		var count int64
		testDB.DB.Model(&models.User{}).Where("id = ?", user.ID).Count(&count)
		if count != 0 {
			t.Error("Expected soft deleted user to be hidden from normal queries")
		}
		t.Log("✓ User deleted successfully")
	})
}

// TestUserService_RoleAssignment_Integration tests user-role relationships
func TestUserService_RoleAssignment_Integration(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	userRepo := repositories.NewUserRepository(testDB.DB)
	userService := services.NewUserService(userRepo)

	t.Run("success - assign roles to user", func(t *testing.T) {
		// Arrange: Create user and roles
		user := helpers.NewUserFactory(t, testDB).WithEmail("roles@test.com").Create()
		adminRole := helpers.NewRoleFactory(t, testDB).WithName("admin").Create()
		editorRole := helpers.NewRoleFactory(t, testDB).WithName("editor").Create()

		roleIDs := []uint{adminRole.ID, editorRole.ID}

		// Act: Assign roles
		err := userService.AssignRoles(user.ID, roleIDs)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Verify roles assigned
		var count int64
		testDB.DB.Table("user_has_roles").Where("user_id = ?", user.ID).Count(&count)
		if count != 2 {
			t.Errorf("Expected 2 roles assigned, got %d", count)
		}
		t.Logf("✓ Assigned %d roles to user", len(roleIDs))
	})

	t.Run("success - get user roles", func(t *testing.T) {
		// Arrange: Create user with roles
		user, role := helpers.CreateUserWithRole(t, testDB, "viewer")

		// Act: Get user roles
		roles, err := userService.GetRoles(user.ID)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if len(roles) != 1 {
			t.Errorf("Expected 1 role, got %d", len(roles))
		}
		if len(roles) > 0 && roles[0].Name != role.Name {
			t.Errorf("Expected role name %s, got %s", role.Name, roles[0].Name)
		}
		t.Logf("✓ Retrieved %d roles for user", len(roles))
	})

	t.Run("success - reassign roles (replace existing)", func(t *testing.T) {
		// Arrange: Create user with initial role
		user := helpers.NewUserFactory(t, testDB).WithEmail("reassign@test.com").Create()
		oldRole := helpers.NewRoleFactory(t, testDB).WithName("old_role").Create()
		userService.AssignRoles(user.ID, []uint{oldRole.ID})

		// Create new roles
		newRole1 := helpers.NewRoleFactory(t, testDB).WithName("new_role_1").Create()
		newRole2 := helpers.NewRoleFactory(t, testDB).WithName("new_role_2").Create()

		// Act: Reassign with new roles
		err := userService.AssignRoles(user.ID, []uint{newRole1.ID, newRole2.ID})

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Verify new roles assigned (old role should be replaced)
		roles, _ := userService.GetRoles(user.ID)
		if len(roles) != 2 {
			t.Errorf("Expected 2 roles after reassignment, got %d", len(roles))
		}

		// Verify old role not present
		hasOldRole := false
		for _, r := range roles {
			if r.ID == oldRole.ID {
				hasOldRole = true
			}
		}
		if hasOldRole {
			t.Error("Old role should not be present after reassignment")
		}

		t.Log("✓ Roles reassigned successfully")
	})
}

// TestUserService_MultipleUsers_Integration tests handling multiple users
func TestUserService_MultipleUsers_Integration(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	userRepo := repositories.NewUserRepository(testDB.DB)
	userService := services.NewUserService(userRepo)

	t.Run("success - create and list 50 users", func(t *testing.T) {
		// Arrange: Create 50 users
		for i := 0; i < 50; i++ {
			helpers.NewUserFactory(t, testDB).Create()
		}

		// Act: List with pagination
		page1Users, total, err := userService.List(1, 20)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if total != 50 {
			t.Errorf("Expected total 50, got %d", total)
		}
		if len(page1Users) != 20 {
			t.Errorf("Expected 20 users on page 1, got %d", len(page1Users))
		}

		// Act: Get page 2
		page2Users, _, err := userService.List(2, 20)
		if err != nil {
			t.Errorf("Expected no error on page 2, got: %v", err)
		}
		if len(page2Users) != 20 {
			t.Errorf("Expected 20 users on page 2, got %d", len(page2Users))
		}

		// Act: Get page 3
		page3Users, _, err := userService.List(3, 20)
		if err != nil {
			t.Errorf("Expected no error on page 3, got: %v", err)
		}
		if len(page3Users) != 10 {
			t.Errorf("Expected 10 users on page 3, got %d", len(page3Users))
		}

		t.Logf("✓ Pagination working: Page1=%d, Page2=%d, Page3=%d",
			len(page1Users), len(page2Users), len(page3Users))
	})
}

// TestUserService_ErrorCases_Integration tests error handling scenarios
func TestUserService_ErrorCases_Integration(t *testing.T) {
	testDB := helpers.SetupTestDB(t)
	defer helpers.CleanupTestDB(t, testDB)

	userRepo := repositories.NewUserRepository(testDB.DB)
	userService := services.NewUserService(userRepo)

	t.Run("error - find non-existent user", func(t *testing.T) {
		// Act: Try to find user that doesn't exist
		_, err := userService.FindByID(99999)

		// Assert
		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}
		t.Logf("✓ Correctly returned error for non-existent user: %v", err)
	})

	t.Run("success - delete non-existent user (GORM behavior)", func(t *testing.T) {
		// Act: Try to delete user that doesn't exist
		// Note: GORM Delete() does not return error for non-existent records
		err := userService.DeleteByID(99999)

		// Assert: GORM behavior - no error for deleting non-existent record
		if err != nil {
			t.Errorf("Expected no error (GORM behavior), got: %v", err)
		}
		t.Logf("✓ Delete non-existent user completed (GORM does not error)")
	})

	t.Run("error - assign roles to non-existent user", func(t *testing.T) {
		// Act: Try to assign roles to non-existent user
		err := userService.AssignRoles(99999, []uint{1, 2})

		// Assert
		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}
		t.Logf("✓ Correctly returned error when assigning roles to non-existent user: %v", err)
	})

	t.Run("error - get roles of non-existent user", func(t *testing.T) {
		// Act: Try to get roles of non-existent user
		_, err := userService.GetRoles(99999)

		// Assert
		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}
		t.Logf("✓ Correctly returned error when getting roles of non-existent user: %v", err)
	})

	t.Run("error - list users returns empty when no users exist", func(t *testing.T) {
		// Act: List users in fresh DB (no users created)
		users, total, err := userService.List(1, 10)

		// Assert
		if err != nil {
			t.Errorf("Expected no error for empty list, got: %v", err)
		}
		if total != 0 {
			t.Errorf("Expected total 0, got %d", total)
		}
		if len(users) != 0 {
			t.Errorf("Expected empty users slice, got %d users", len(users))
		}
		t.Logf("✓ Correctly returned empty list when no users exist")
	})
}
