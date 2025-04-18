package helpers

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetFileURL(key string, path string) string {
	jwtKey := []byte(GetEnv("APP_KEY", "your_secret_key"))
	expires := GetEnvInt("IMAGE_EXPIRE_MINUTES", 2)
	expiredAt := time.Now().Add(time.Minute * time.Duration(expires)).Unix()
	signature := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"key":       key,
		"path":      path,
		"expired_at": expiredAt,
	})
	token, _ := signature.SignedString(jwtKey)

	mainUrl := GetEnv("APP_URL", "http://localhost:8080")
	return fmt.Sprintf("%s/file/%s/%s?signature=%s", mainUrl, path, key, token)
}
