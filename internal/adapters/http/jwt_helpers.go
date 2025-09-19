package http

import (
	"os"
)

func getJWTSecret() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}
	return "super_secret_key"
}
