package services_test

import (
	"errors"
	"testing"

	"golang_starter_kit_2025/app/mocks"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/services"

	"go.uber.org/mock/gomock"
)

func TestRoleService_List_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepositoryInterface(ctrl)
	mockPermRepo := mocks.NewMockPermissionRepositoryInterface(ctrl)
	roleService := services.NewRoleService(mockRoleRepo, mockPermRepo)

	t.Run("success - returns paginated roles", func(t *testing.T) {
		expectedRoles := []models.Role{
			{ID: 1, Name: "admin"},
			{ID: 2, Name: "user"},
		}
		var expectedTotal int64 = 10

		mockRoleRepo.EXPECT().
			List(1, 10).
			Return(expectedRoles, expectedTotal, nil)

		roles, total, err := roleService.List(1, 10)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if total != 10 {
			t.Errorf("expected total 10, got %d", total)
		}
		if len(roles) != 2 {
			t.Errorf("expected 2 roles, got %d", len(roles))
		}
	})

	t.Run("error - invalid pagination parameters", func(t *testing.T) {
		mockRoleRepo.EXPECT().
			List(0, -1).
			Return(nil, int64(0), errors.New("invalid pagination parameters"))

		roles, total, err := roleService.List(0, -1)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if total != 0 {
			t.Errorf("expected total 0, got %d", total)
		}
		if roles != nil {
			t.Error("expected nil roles, got non-nil")
		}
	})
}

func TestRoleService_Create_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepositoryInterface(ctrl)
	mockPermRepo := mocks.NewMockPermissionRepositoryInterface(ctrl)
	roleService := services.NewRoleService(mockRoleRepo, mockPermRepo)

	t.Run("success - creates new role", func(t *testing.T) {
		newRole := models.Role{
			Name:  "editor",
			Group: "content",
		}

		mockRoleRepo.EXPECT().
			Create(gomock.Any()).
			DoAndReturn(func(role *models.Role) error {
				role.ID = 1
				return nil
			})

		result, err := roleService.Create(newRole)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result.ID == 0 {
			t.Error("expected role ID to be set")
		}
	})

	t.Run("error - duplicate role name", func(t *testing.T) {
		newRole := models.Role{
			Name:  "admin",
			Group: "system",
		}

		mockRoleRepo.EXPECT().
			Create(gomock.Any()).
			Return(errors.New("duplicate entry"))

		_, err := roleService.Create(newRole)

		if err == nil {
			t.Error("expected error for duplicate role, got nil")
		}
	})
}

func TestRoleService_CreateAndUpdate_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepositoryInterface(ctrl)
	mockPermRepo := mocks.NewMockPermissionRepositoryInterface(ctrl)
	roleService := services.NewRoleService(mockRoleRepo, mockPermRepo)

	t.Run("success - creates new role", func(t *testing.T) {
		newRole := models.Role{
			Name:  "editor",
			Group: "content",
		}

		mockRoleRepo.EXPECT().
			Create(gomock.Any()).
			DoAndReturn(func(role *models.Role) error {
				role.ID = 1
				return nil
			})

		result, err := roleService.Create(newRole)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result.ID == 0 {
			t.Error("expected role ID to be set")
		}
	})

	t.Run("success - updates existing role", func(t *testing.T) {
		updatedRole := models.Role{
			Name:  "updated_admin",
			Group: "system",
		}

		mockRoleRepo.EXPECT().
			Update(gomock.Any()).
			Return(nil)

		mockRoleRepo.EXPECT().
			FindByID(uint(1)).
			Return(&updatedRole, nil)

		result, err := roleService.Update(1, updatedRole)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result.Name != "updated_admin" {
			t.Errorf("expected name 'updated_admin', got '%s'", result.Name)
		}
	})

	t.Run("error - update fails", func(t *testing.T) {
		updatedRole := models.Role{
			Name: "admin",
		}

		mockRoleRepo.EXPECT().
			Update(gomock.Any()).
			Return(errors.New("database error"))

		_, err := roleService.Update(1, updatedRole)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestRoleService_DeleteByID_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepositoryInterface(ctrl)
	mockPermRepo := mocks.NewMockPermissionRepositoryInterface(ctrl)
	roleService := services.NewRoleService(mockRoleRepo, mockPermRepo)

	t.Run("success - deletes role", func(t *testing.T) {
		mockRoleRepo.EXPECT().
			Delete(uint(1)).
			Return(nil)

		err := roleService.DeleteByID(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("error - role not found", func(t *testing.T) {
		mockRoleRepo.EXPECT().
			Delete(uint(999)).
			Return(errors.New("role not found"))

		err := roleService.DeleteByID(999)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("error - role has users assigned", func(t *testing.T) {
		mockRoleRepo.EXPECT().
			Delete(uint(1)).
			Return(errors.New("cannot delete role with assigned users"))

		err := roleService.DeleteByID(1)

		if err == nil {
			t.Error("expected error when role has users, got nil")
		}
	})
}

func TestRoleService_FindByID_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepositoryInterface(ctrl)
	mockPermRepo := mocks.NewMockPermissionRepositoryInterface(ctrl)
	roleService := services.NewRoleService(mockRoleRepo, mockPermRepo)

	t.Run("success - finds role by id", func(t *testing.T) {
		expectedRole := &models.Role{
			ID:    1,
			Name:  "admin",
			Group: "system",
		}

		mockRoleRepo.EXPECT().
			FindByID(uint(1)).
			Return(expectedRole, nil)

		role, err := roleService.FindByID(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if role.Name != "admin" {
			t.Errorf("expected role 'admin', got '%s'", role.Name)
		}
	})

	t.Run("error - role not found", func(t *testing.T) {
		mockRoleRepo.EXPECT().
			FindByID(uint(999)).
			Return(nil, errors.New("role not found"))

		role, err := roleService.FindByID(999)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if role != nil {
			t.Error("expected nil role, got non-nil")
		}
	})
}

func TestRoleService_FindByName_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepositoryInterface(ctrl)
	mockPermRepo := mocks.NewMockPermissionRepositoryInterface(ctrl)
	roleService := services.NewRoleService(mockRoleRepo, mockPermRepo)

	t.Run("success - finds role by name", func(t *testing.T) {
		expectedRole := &models.Role{
			ID:    1,
			Name:  "admin",
			Group: "system",
		}

		mockRoleRepo.EXPECT().
			FindByName("admin").
			Return(expectedRole, nil)

		role, err := roleService.FindByName("admin")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if role.ID != 1 {
			t.Errorf("expected role ID 1, got %d", role.ID)
		}
	})

	t.Run("error - role not found", func(t *testing.T) {
		mockRoleRepo.EXPECT().
			FindByName("nonexistent").
			Return(nil, errors.New("role not found"))

		role, err := roleService.FindByName("nonexistent")

		if err == nil {
			t.Error("expected error, got nil")
		}
		if role != nil {
			t.Error("expected nil role, got non-nil")
		}
	})
}

func TestRoleService_ExistsByName_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepositoryInterface(ctrl)
	mockPermRepo := mocks.NewMockPermissionRepositoryInterface(ctrl)
	roleService := services.NewRoleService(mockRoleRepo, mockPermRepo)

	t.Run("success - role exists", func(t *testing.T) {
		mockRoleRepo.EXPECT().
			ExistsByName("admin").
			Return(true, nil)

		exists, err := roleService.ExistsByName("admin")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !exists {
			t.Error("expected role to exist")
		}
	})

	t.Run("success - role does not exist", func(t *testing.T) {
		mockRoleRepo.EXPECT().
			ExistsByName("nonexistent").
			Return(false, nil)

		exists, err := roleService.ExistsByName("nonexistent")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if exists {
			t.Error("expected role to not exist")
		}
	})

	t.Run("error - database error", func(t *testing.T) {
		mockRoleRepo.EXPECT().
			ExistsByName("admin").
			Return(false, errors.New("database error"))

		exists, err := roleService.ExistsByName("admin")

		if err == nil {
			t.Error("expected error, got nil")
		}
		if exists {
			t.Error("expected exists to be false on error")
		}
	})
}

func TestRoleService_Count_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepositoryInterface(ctrl)
	mockPermRepo := mocks.NewMockPermissionRepositoryInterface(ctrl)
	roleService := services.NewRoleService(mockRoleRepo, mockPermRepo)

	t.Run("success - returns role count", func(t *testing.T) {
		var expectedCount int64 = 5

		mockRoleRepo.EXPECT().
			Count().
			Return(expectedCount, nil)

		count, err := roleService.Count()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if count != 5 {
			t.Errorf("expected count 5, got %d", count)
		}
	})

	t.Run("success - zero roles", func(t *testing.T) {
		var expectedCount int64 = 0

		mockRoleRepo.EXPECT().
			Count().
			Return(expectedCount, nil)

		count, err := roleService.Count()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if count != 0 {
			t.Errorf("expected count 0, got %d", count)
		}
	})
}

func TestRoleService_AssignPermissions_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepositoryInterface(ctrl)
	mockPermRepo := mocks.NewMockPermissionRepositoryInterface(ctrl)
	roleService := services.NewRoleService(mockRoleRepo, mockPermRepo)

	t.Run("success - assigns permissions to role", func(t *testing.T) {
		permissionIDs := []uint{1, 2, 3}

		mockRoleRepo.EXPECT().
			AssignPermissions(uint(1), permissionIDs).
			Return(nil)

		err := roleService.AssignPermissions(1, permissionIDs)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("success - empty permission list clears permissions", func(t *testing.T) {
		permissionIDs := []uint{}

		mockRoleRepo.EXPECT().
			AssignPermissions(uint(1), permissionIDs).
			Return(nil)

		err := roleService.AssignPermissions(1, permissionIDs)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("error - role not found", func(t *testing.T) {
		permissionIDs := []uint{1, 2}

		mockRoleRepo.EXPECT().
			AssignPermissions(uint(999), permissionIDs).
			Return(errors.New("role not found"))

		err := roleService.AssignPermissions(999, permissionIDs)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestRoleService_GetPermissions_Unit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRoleRepo := mocks.NewMockRoleRepositoryInterface(ctrl)
	mockPermRepo := mocks.NewMockPermissionRepositoryInterface(ctrl)
	roleService := services.NewRoleService(mockRoleRepo, mockPermRepo)

	t.Run("success - gets role permissions", func(t *testing.T) {
		expectedPerms := []models.Permission{
			{ID: 1, Name: "create:users"},
			{ID: 2, Name: "delete:users"},
		}

		mockRoleRepo.EXPECT().
			GetPermissions(uint(1)).
			Return(expectedPerms, nil)

		perms, err := roleService.GetPermissions(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(perms) != 2 {
			t.Errorf("expected 2 permissions, got %d", len(perms))
		}
	})

	t.Run("success - role has no permissions", func(t *testing.T) {
		expectedPerms := []models.Permission{}

		mockRoleRepo.EXPECT().
			GetPermissions(uint(1)).
			Return(expectedPerms, nil)

		perms, err := roleService.GetPermissions(1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(perms) != 0 {
			t.Errorf("expected 0 permissions, got %d", len(perms))
		}
	})

	t.Run("error - role not found", func(t *testing.T) {
		mockRoleRepo.EXPECT().
			GetPermissions(uint(999)).
			Return(nil, errors.New("role not found"))

		perms, err := roleService.GetPermissions(999)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if perms != nil {
			t.Error("expected nil permissions, got non-nil")
		}
	})
}
