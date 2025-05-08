package controllers_test

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"

	// . "github.com/onsi/gomega"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/facades"
)

var _ = Describe("Login", Ordered, func() {
	BeforeAll(func() {
		// Get the current working directory
		wd, err := os.Getwd()
		if err != nil {
			Fail(fmt.Sprintf("Error getting current working directory: %v", err))
		}

		// Load the .env.test file from the current working directory
		envPath := fmt.Sprintf("%s/../../.env.test", wd)
		if err := godotenv.Load(envPath); err != nil {
			Fail("Error loading .env.test file")
		}

		facades.ConnectDB(envPath)
		facades.DB.AutoMigrate(
			&models.User{},
			&models.Role{},
			&models.UserHasRole{},
		)
	})

	AfterAll(func() {
		facades.DB.Migrator().DropTable(&models.UserHasRole{})
		facades.DB.Migrator().DropTable(&models.Role{})
		facades.DB.Migrator().DropTable(&models.User{})

		_, err := facades.DB.DB()
		if err != nil {
			Fail(fmt.Sprintf("Error getting facades connection: %v", err))
		}
	})

	// Describe("Login", func() {
	// 	Context("Success", func() {
	// 		It("should return token", func() {
	// 			// Add test data
	// 			user := models.User{
	// 				Email:    "test@mail.com",
	// 				Password: "12345678",
	// 			}

	// 			if err := facades.DB.Create(&user).Error; err != nil {
	// 				Fail(fmt.Sprintf("Error creating user: %v", err))
	// 			}

	// 			router := bootstrap.Router()

	// 			w := httptest.NewRecorder()
	// 			payload, _ := json.Marshal(user)
	// 			req, _ := http.NewRequest("PUT", "http://localhost:8080/auth/login", strings.NewReader(string(payload)))
	// 			router.ServeHTTP(w, req)

	// 			Expect(w.Code).To(Equal(http.StatusOK))
	// 			Expect(w.Body.String()).To(ContainSubstring("token"))
	// 		})
	// 	})
	// })
})
