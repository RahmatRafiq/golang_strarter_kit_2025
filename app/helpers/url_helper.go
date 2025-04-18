package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

func GenerateSignedURL(path string) string {
	baseURL := os.Getenv("APP_URL") + path
	appKey := os.Getenv("APP_KEY")

	expirationTime := time.Now().Add(15 * time.Minute).Unix() // Set expiration time
	data := fmt.Sprintf("%s?expires=%d", baseURL, expirationTime)

	// Generate the HMAC signature
	signature := generateHMACSignature(data, appKey)

	// Create the full signed URL
	signedURL := fmt.Sprintf("%s&signature=%s", data, signature)

	return signedURL
}

func generateHMACSignature(data string, secret string) string {
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(data))
    return hex.EncodeToString(h.Sum(nil))
}