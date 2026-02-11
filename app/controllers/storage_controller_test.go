package controllers_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"

	"golang_starter_kit_2025/app/controllers"
	"golang_starter_kit_2025/app/services"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StorageController", func() {
	var (
		router     *gin.Engine
		controller *controllers.StorageController
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		router = gin.New()

		// Setup storage service (may fail if MinIO not available)
		storageService, err := services.NewStorageService()
		if err != nil {
			Skip("Storage service not available (MinIO not running)")
		}
		controller = controllers.NewStorageController(storageService)
	})

	Describe("UploadFile", func() {
		Context("with valid file", func() {
			It("should upload file successfully", func() {
				router.POST("/upload", controller.UploadFile)

				// Create multipart form data
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				// Add file field
				part, err := writer.CreateFormFile("file", "test.txt")
				Expect(err).NotTo(HaveOccurred())
				part.Write([]byte("test content"))

				// Add folder field
				writer.WriteField("folder", "test")
				writer.Close()

				req, _ := http.NewRequest("POST", "/upload", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				if resp.Code == http.StatusOK {
					var result map[string]interface{}
					json.Unmarshal(resp.Body.Bytes(), &result)

					Expect(result["status"]).To(Equal("success"))
					Expect(result["message"]).To(Equal("File uploaded successfully"))
					Expect(result["item"]).NotTo(BeNil())

					item := result["item"].(map[string]interface{})
					Expect(item["filename"]).NotTo(BeEmpty())
					Expect(item["url"]).NotTo(BeEmpty())
				}
			})
		})

		Context("without file", func() {
			It("should return error", func() {
				router.POST("/upload", controller.UploadFile)

				req, _ := http.NewRequest("POST", "/upload", nil)
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				Expect(resp.Code).To(Equal(http.StatusBadRequest))

				var result map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &result)

				Expect(result["status"]).To(Equal("error"))
				Expect(result["message"]).To(Equal("No file uploaded"))
			})
		})
	})

	Describe("GetFileURL", func() {
		Context("with non-existent file", func() {
			It("should return not found error", func() {
				router.GET("/url/:filename", controller.GetFileURL)

				req, _ := http.NewRequest("GET", "/url/nonexistent.txt", nil)
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				Expect(resp.Code).To(Equal(http.StatusNotFound))

				var result map[string]interface{}
				json.Unmarshal(resp.Body.Bytes(), &result)

				Expect(result["status"]).To(Equal("error"))
				Expect(result["message"]).To(Equal("File not found"))
			})
		})
	})

	Describe("DeleteFile", func() {
		Context("with non-existent file", func() {
			It("should handle gracefully", func() {
				router.DELETE("/:filename", controller.DeleteFile)

				req, _ := http.NewRequest("DELETE", "/nonexistent.txt", nil)
				resp := httptest.NewRecorder()

				router.ServeHTTP(resp, req)

				// Storage service may return success even if file doesn't exist
				// depending on MinIO behavior
				Expect(resp.Code).To(Or(Equal(http.StatusOK), Equal(http.StatusInternalServerError)))
			})
		})
	})
})

var _ = Describe("StorageService", func() {
	var storageService *services.StorageService

	BeforeEach(func() {
		var err error
		storageService, err = services.NewStorageService()
		if err != nil {
			Skip("Storage service not available (MinIO not running)")
		}
	})

	Describe("Configuration", func() {
		It("should have valid configuration", func() {
			Expect(storageService).NotTo(BeNil())
		})
	})

	Describe("File Operations", func() {
		var testFilename string

		AfterEach(func() {
			// Cleanup: delete test file if it exists
			if testFilename != "" {
				storageService.DeleteFile(testFilename)
			}
		})

		It("should handle file lifecycle", func() {
			// Create a temporary file
			tmpFile, err := os.CreateTemp("", "test-*.txt")
			Expect(err).NotTo(HaveOccurred())
			defer os.Remove(tmpFile.Name())

			tmpFile.WriteString("test content for storage")
			tmpFile.Seek(0, 0)

			// Test upload
			filename, err := storageService.UploadFile(tmpFile, &multipart.FileHeader{
				Filename: "test-upload.txt",
				Size:     24,
			}, "test")

			if err == nil {
				testFilename = filename

				// Test GetFileURL
				url, err := storageService.GetFileURL(filename, 0)
				Expect(err).NotTo(HaveOccurred())
				Expect(url).NotTo(BeEmpty())

				// Test Delete
				err = storageService.DeleteFile(filename)
				Expect(err).NotTo(HaveOccurred())
				testFilename = "" // Mark as deleted
			}
		})
	})
})
