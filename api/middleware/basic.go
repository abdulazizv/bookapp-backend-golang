package middleware

import (
	"github.com/golanguzb70/middleware/gin/basicauth"
)

func BasicAuth() basicauth.Config {
	// This configuration checks for all incoming requests for authentication
	return basicauth.Config{
		Users: []basicauth.User{
			{
				UserName: "bookapp",
				Password: "54321",
			},
		},
		RestrictedUrls: []string{
			"/v1/swagger/*",
		},
	}
}
