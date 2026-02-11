package services_test

import (
	"errors"
	"testing"

	"golang_starter_kit_2025/app/mocks"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/services"

	"go.uber.org/mock/gomock"
)

func TestUserService_List_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	userService := services.NewUserService(mockRepo)

	t.Run("success - returns users list with pagination", func(t *testing.T) {
		expectedUsers := []models.User{
			{ID: 1, Username: "user1", Email: "user1@example.com"},
			{ID: 2, Username: "user2", Email: "user2@example.com"},
			{ID: 3, Username: "user3", Email: "user3@example.com"},
		}
		var expectedTotal int64 = 3

		mockRepo.EXPECT().
			List(1, 10).
			Return(expectedUsers, expectedTotal, nil)

		users, total, err := userService.List(1, 10)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if total != expectedTotal {
			t.Errorf("expected total %d, got %d", expectedTotal, total)
		}
		if len(users) != len(expectedUsers) {
			t.Errorf("expected %d users, got %d", len(expectedUsers), len(users))
		}
		if users[0].Username != "user1" {
			t.Errorf("expected first user 'user1', got '%s'", users[0].Username)
		}
	})

	t.Run("error - repository returns error", func(t *testing.T) {
		expectedErr := errors.New("database connection failed")

		mockRepo.EXPECT().
			List(1, 10).
			Return(nil, int64(0), expectedErr)

		users, total, err := userService.List(1, 10)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if err.Error() != expectedErr.Error() {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if total != 0 {
			t.Errorf("expected total 0, got %d", total)
		}
		if users != nil {
			t.Error("expected nil users, got non-nil")
		}
	})

	t.Run("success - empty list when no users exist", func(t *testing.T) {
		expectedUsers := []models.User{}
		var expectedTotal int64 = 0

		mockRepo.EXPECT().
			List(1, 10).
			Return(expectedUsers, expectedTotal, nil)

		users, total, err := userService.List(1, 10)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if total != 0 {
			t.Errorf("expected total 0, got %d", total)
		}
		if len(users) != 0 {
			t.Errorf("expected 0 users, got %d", len(users))
		}
	})

	t.Run("success - pagination with different page sizes", func(t *testing.T) {
		// Test page 1 with size 20
		expectedUsers := make([]models.User, 20)
		for i := 0; i < 20; i++ {
			expectedUsers[i] = models.User{
				ID:       uint(i + 1),
				Username: "user" + string(rune(i+1)),
			}
		}
		var expectedTotal int64 = 50

		mockRepo.EXPECT().
			List(1, 20).
			Return(expectedUsers, expectedTotal, nil)

		users, total, err := userService.List(1, 20)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if total != 50 {
			t.Errorf("expected total 50, got %d", total)
		}
		if len(users) != 20 {
			t.Errorf("expected 20 users, got %d", len(users))
		}
	})
}

func TestUserService_FindByID_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	userService := services.NewUserService(mockRepo)

	t.Run("success - finds user by id", func(t *testing.T) {
		expectedUser := &models.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
		}

		mockRepo.EXPECT().
			FindByID(uint(1)).
			Return(expectedUser, nil)

		user, err := userService.FindByID(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user.ID != expectedUser.ID {
			t.Errorf("expected user ID %d, got %d", expectedUser.ID, user.ID)
		}
		if user.Username != expectedUser.Username {
			t.Errorf("expected username %s, got %s", expectedUser.Username, user.Username)
		}
		if user.Email != expectedUser.Email {
			t.Errorf("expected email %s, got %s", expectedUser.Email, user.Email)
		}
	})

	t.Run("error - user not found", func(t *testing.T) {
		mockRepo.EXPECT().
			FindByID(uint(999)).
			Return(nil, errors.New("user not found"))

		user, err := userService.FindByID(999)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if err.Error() != "user not found" {
			t.Errorf("expected 'user not found', got %v", err)
		}
		if user != nil {
			t.Error("expected nil user, got non-nil")
		}
	})

	t.Run("error - database error", func(t *testing.T) {
		mockRepo.EXPECT().
			FindByID(uint(1)).
			Return(nil, errors.New("database connection lost"))

		user, err := userService.FindByID(1)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if user != nil {
			t.Error("expected nil user, got non-nil")
		}
	})
}

func TestUserService_DeleteByID_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	userService := services.NewUserService(mockRepo)

	t.Run("success - deletes user", func(t *testing.T) {
		mockRepo.EXPECT().
			Delete(uint(1)).
			Return(nil)

		err := userService.DeleteByID(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("error - user not found", func(t *testing.T) {
		expectedErr := errors.New("user not found")

		mockRepo.EXPECT().
			Delete(uint(999)).
			Return(expectedErr)

		err := userService.DeleteByID(999)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if err.Error() != expectedErr.Error() {
			t.Errorf("expected '%v', got '%v'", expectedErr, err)
		}
	})

	t.Run("error - database error during deletion", func(t *testing.T) {
		mockRepo.EXPECT().
			Delete(uint(1)).
			Return(errors.New("foreign key constraint violation"))

		err := userService.DeleteByID(1)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestUserService_AssignRoles_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	userService := services.NewUserService(mockRepo)

	t.Run("success - assigns roles to user", func(t *testing.T) {
		roleIDs := []uint{1, 2, 3}

		mockRepo.EXPECT().
			AssignRoles(uint(1), roleIDs).
			Return(nil)

		err := userService.AssignRoles(1, roleIDs)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("success - assigns single role", func(t *testing.T) {
		roleIDs := []uint{1}

		mockRepo.EXPECT().
			AssignRoles(uint(1), roleIDs).
			Return(nil)

		err := userService.AssignRoles(1, roleIDs)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("error - invalid user id (zero)", func(t *testing.T) {
		roleIDs := []uint{1, 2}

		// Note: UserService does not validate ID=0, it passes to repository
		// Repository will return error for non-existent user
		mockRepo.EXPECT().
			AssignRoles(uint(0), roleIDs).
			Return(errors.New("invalid user id"))

		err := userService.AssignRoles(0, roleIDs)

		if err == nil {
			t.Error("expected error for invalid ID, got nil")
		}
	})

	t.Run("error - user not found", func(t *testing.T) {
		roleIDs := []uint{1, 2}

		mockRepo.EXPECT().
			AssignRoles(uint(999), roleIDs).
			Return(errors.New("user not found"))

		err := userService.AssignRoles(999, roleIDs)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("error - repository error", func(t *testing.T) {
		roleIDs := []uint{1, 2}
		expectedErr := errors.New("database transaction failed")

		mockRepo.EXPECT().
			AssignRoles(uint(1), roleIDs).
			Return(expectedErr)

		err := userService.AssignRoles(1, roleIDs)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if err.Error() != expectedErr.Error() {
			t.Errorf("expected '%v', got '%v'", expectedErr, err)
		}
	})

	t.Run("success - reassign roles (empty role list)", func(t *testing.T) {
		roleIDs := []uint{} // Empty list should clear all roles

		mockRepo.EXPECT().
			AssignRoles(uint(1), roleIDs).
			Return(nil)

		err := userService.AssignRoles(1, roleIDs)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}

func TestUserService_GetRoles_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	userService := services.NewUserService(mockRepo)

	t.Run("success - gets user roles", func(t *testing.T) {
		expectedRoles := []models.Role{
			{ID: 1, Name: "admin", Group: "system"},
			{ID: 2, Name: "user", Group: "default"},
		}

		mockRepo.EXPECT().
			GetRoles(uint(1)).
			Return(expectedRoles, nil)

		roles, err := userService.GetRoles(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(roles) != len(expectedRoles) {
			t.Errorf("expected %d roles, got %d", len(expectedRoles), len(roles))
		}
		if roles[0].Name != "admin" {
			t.Errorf("expected first role 'admin', got '%s'", roles[0].Name)
		}
	})

	t.Run("success - user has no roles", func(t *testing.T) {
		expectedRoles := []models.Role{}

		mockRepo.EXPECT().
			GetRoles(uint(1)).
			Return(expectedRoles, nil)

		roles, err := userService.GetRoles(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(roles) != 0 {
			t.Errorf("expected 0 roles, got %d", len(roles))
		}
	})

	t.Run("error - invalid user id (zero)", func(t *testing.T) {
		// Note: UserService does not validate ID=0, it passes to repository
		mockRepo.EXPECT().
			GetRoles(uint(0)).
			Return(nil, errors.New("invalid user id"))

		_, err := userService.GetRoles(0)

		if err == nil {
			t.Error("expected error for invalid ID, got nil")
		}
	})

	t.Run("error - user not found", func(t *testing.T) {
		mockRepo.EXPECT().
			GetRoles(uint(999)).
			Return(nil, errors.New("user not found"))

		roles, err := userService.GetRoles(999)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if roles != nil {
			t.Error("expected nil roles, got non-nil")
		}
	})

	t.Run("error - database error", func(t *testing.T) {
		mockRepo.EXPECT().
			GetRoles(uint(1)).
			Return(nil, errors.New("database query failed"))

		roles, err := userService.GetRoles(1)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if roles != nil {
			t.Error("expected nil roles, got non-nil")
		}
	})
}
