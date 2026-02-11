package services_test

import (
	"errors"
	"testing"

	"golang_starter_kit_2025/app/mocks"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/services"

	"go.uber.org/mock/gomock"
)

func TestUserService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	userService := services.NewUserService(mockRepo)

	t.Run("success - returns users list with pagination", func(t *testing.T) {
		expectedUsers := []models.User{
			{ID: 1, Username: "user1", Email: "user1@example.com"},
			{ID: 2, Username: "user2", Email: "user2@example.com"},
		}
		var expectedTotal int64 = 2

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
	})

	t.Run("error - repository returns error", func(t *testing.T) {
		expectedErr := errors.New("database error")

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

	t.Run("success - empty list", func(t *testing.T) {
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
}

func TestUserService_Find(t *testing.T) {
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

		user, err := userService.Find("1")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user.ID != expectedUser.ID {
			t.Errorf("expected user ID %d, got %d", expectedUser.ID, user.ID)
		}
		if user.Username != expectedUser.Username {
			t.Errorf("expected username %s, got %s", expectedUser.Username, user.Username)
		}
	})

	t.Run("error - invalid id format", func(t *testing.T) {
		_, err := userService.Find("invalid")

		if err == nil {
			t.Error("expected error for invalid ID, got nil")
		}
	})

	t.Run("error - user not found", func(t *testing.T) {
		mockRepo.EXPECT().
			FindByID(uint(999)).
			Return(nil, errors.New("user not found"))

		_, err := userService.Find("999")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestUserService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	userService := services.NewUserService(mockRepo)

	t.Run("success - deletes user", func(t *testing.T) {
		mockRepo.EXPECT().
			Delete(uint(1)).
			Return(nil)

		err := userService.Delete("1")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("error - invalid id format", func(t *testing.T) {
		err := userService.Delete("invalid")

		if err == nil {
			t.Error("expected error for invalid ID, got nil")
		}
	})

	t.Run("error - user not found", func(t *testing.T) {
		expectedErr := errors.New("user not found")

		mockRepo.EXPECT().
			Delete(uint(999)).
			Return(expectedErr)

		err := userService.Delete("999")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestUserService_AssignRolesToUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	userService := services.NewUserService(mockRepo)

	t.Run("success - assigns roles to user", func(t *testing.T) {
		roleIDs := []uint{1, 2, 3}

		mockRepo.EXPECT().
			AssignRoles(uint(1), roleIDs).
			Return(nil)

		err := userService.AssignRolesToUser("1", roleIDs)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("error - invalid user id", func(t *testing.T) {
		roleIDs := []uint{1, 2}

		err := userService.AssignRolesToUser("invalid", roleIDs)

		if err == nil {
			t.Error("expected error for invalid ID, got nil")
		}
	})

	t.Run("error - repository error", func(t *testing.T) {
		roleIDs := []uint{1, 2}
		expectedErr := errors.New("database error")

		mockRepo.EXPECT().
			AssignRoles(uint(1), roleIDs).
			Return(expectedErr)

		err := userService.AssignRolesToUser("1", roleIDs)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestUserService_GetRolesByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	userService := services.NewUserService(mockRepo)

	t.Run("success - gets user roles", func(t *testing.T) {
		expectedRoles := []models.Role{
			{ID: 1, Name: "admin"},
			{ID: 2, Name: "user"},
		}

		mockRepo.EXPECT().
			GetRoles(uint(1)).
			Return(expectedRoles, nil)

		roles, err := userService.GetRolesByUserID("1")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(roles) != len(expectedRoles) {
			t.Errorf("expected %d roles, got %d", len(expectedRoles), len(roles))
		}
	})

	t.Run("error - invalid user id", func(t *testing.T) {
		_, err := userService.GetRolesByUserID("invalid")

		if err == nil {
			t.Error("expected error for invalid ID, got nil")
		}
	})

	t.Run("error - user not found", func(t *testing.T) {
		mockRepo.EXPECT().
			GetRoles(uint(999)).
			Return(nil, errors.New("user not found"))

		_, err := userService.GetRolesByUserID("999")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
